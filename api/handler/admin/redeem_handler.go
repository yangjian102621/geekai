package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RedeemHandler struct {
	handler.BaseHandler
}

func NewRedeemHandler(app *core.AppServer, db *gorm.DB) *RedeemHandler {
	return &RedeemHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

func (h *RedeemHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	code := c.Query("code")
	status := h.GetInt(c, "status", -1)

	session := h.DB.Session(&gorm.Session{})
	if code != "" {
		session.Where("code LIKE ?", "%"+code+"%")
	}
	if status == 0 {
		session.Where("redeem_at = ?", 0)
	} else if status == 1 {
		session.Where("redeem_at > ?", 0)
	}

	var total int64
	session.Model(&model.Redeem{}).Count(&total)
	var redeems []model.Redeem
	offset := (page - 1) * pageSize
	err := session.Order("id DESC").Offset(offset).Limit(pageSize).Find(&redeems).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	var items = make([]vo.Redeem, 0)
	userIds := make([]uint, 0)
	for _, v := range redeems {
		userIds = append(userIds, v.UserId)
	}
	var users []model.User
	h.DB.Where("id IN ?", userIds).Find(&users)
	var userMap = make(map[uint]model.User)
	for _, u := range users {
		userMap[u.Id] = u
	}

	for _, v := range redeems {
		var r vo.Redeem
		err = utils.CopyObject(v, &r)
		if err != nil {
			continue
		}

		r.Id = v.Id
		r.Username = userMap[v.UserId].Username
		r.CreatedAt = v.CreatedAt.Unix()
		items = append(items, r)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, items))
}

func (h *RedeemHandler) Create(c *gin.Context) {
	var data struct {
		Name  string `json:"name"`
		Power int    `json:"power"`
		Num   int    `json:"num"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	counter := 0
	codes := make([]string, 0)
	var errMsg = ""
	if data.Num > 0 {
		for i := 0; i < data.Num; i++ {
			code, err := utils.GenRedeemCode(32)
			if err != nil {
				errMsg = err.Error()
				continue
			}
			err = h.DB.Create(&model.Redeem{
				Code:    code,
				Name:    data.Name,
				Power:   data.Power,
				Enabled: true,
			}).Error
			if err != nil {
				errMsg = err.Error()
				continue
			}
			codes = append(codes, code)
			counter++
		}
	}
	if counter == 0 {
		resp.ERROR(c, errMsg)
		return
	}

	resp.SUCCESS(c, gin.H{
		"counter": counter,
	})
}

func (h *RedeemHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Model(&model.Redeem{}).Where("id = ?", data.Id).Update(data.Filed, data.Value).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *RedeemHandler) Remove(c *gin.Context) {
	var data struct {
		Id uint
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if data.Id > 0 {
		err := h.DB.Where("id", data.Id).Delete(&model.Redeem{}).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	resp.SUCCESS(c)
}
