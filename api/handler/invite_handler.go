package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// InviteHandler 用户邀请
type InviteHandler struct {
	BaseHandler
}

func NewInviteHandler(app *core.AppServer, db *gorm.DB) *InviteHandler {
	return &InviteHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// Code 获取当前用户邀请码
func (h *InviteHandler) Code(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	var inviteCode model.InviteCode
	res := h.DB.Where("user_id = ?", userId).First(&inviteCode)
	// 如果邀请码不存在，则创建一个
	if res.Error != nil {
		code := strings.ToUpper(utils.RandString(8))
		for {
			res = h.DB.Where("code = ?", code).First(&inviteCode)
			if res.Error != nil { // 不存在相同的邀请码则退出
				break
			}
		}
		inviteCode.UserId = userId
		inviteCode.Code = code
		h.DB.Create(&inviteCode)
	}

	var codeVo vo.InviteCode
	err := utils.CopyObject(inviteCode, &codeVo)
	if err != nil {
		resp.ERROR(c, "拷贝对象失败")
		return
	}

	resp.SUCCESS(c, codeVo)
}

// List Log 用户邀请记录
func (h *InviteHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	userId := h.GetLoginUserId(c)
	session := h.DB.Session(&gorm.Session{}).Where("inviter_id = ?", userId)
	var total int64
	session.Model(&model.InviteLog{}).Count(&total)
	var items []model.InviteLog
	var list = make([]vo.InviteLog, 0)
	offset := (page - 1) * pageSize
	res := session.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var v vo.InviteLog
			err := utils.CopyObject(item, &v)
			if err == nil {
				v.Id = item.Id
				v.CreatedAt = item.CreatedAt.Unix()
				list = append(list, v)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, list))
}

// Hits 访问邀请码
func (h *InviteHandler) Hits(c *gin.Context) {
	code := c.Query("code")
	h.DB.Model(&model.InviteCode{}).Where("code = ?", code).UpdateColumn("hits", gorm.Expr("hits + ?", 1))
	resp.SUCCESS(c)
}
