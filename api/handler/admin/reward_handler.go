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

type RewardHandler struct {
	handler.BaseHandler
}

func NewRewardHandler(app *core.AppServer, db *gorm.DB) *RewardHandler {
	return &RewardHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

func (h *RewardHandler) List(c *gin.Context) {
	var items []model.Reward
	res := h.DB.Order("id DESC").Find(&items)
	var rewards = make([]vo.Reward, 0)
	if res.Error == nil {
		userIds := make([]uint, 0)
		for _, v := range items {
			userIds = append(userIds, v.UserId)
		}
		var users []model.User
		h.DB.Where("id IN ?", userIds).Find(&users)
		var userMap = make(map[uint]model.User)
		for _, u := range users {
			userMap[u.Id] = u
		}

		for _, v := range items {
			var r vo.Reward
			err := utils.CopyObject(v, &r)
			if err != nil {
				continue
			}

			r.Id = v.Id
			r.Username = userMap[v.UserId].Username
			r.CreatedAt = v.CreatedAt.Unix()
			r.UpdatedAt = v.UpdatedAt.Unix()
			rewards = append(rewards, r)
		}
	}

	resp.SUCCESS(c, rewards)
}

func (h *RewardHandler) Remove(c *gin.Context) {
	var data struct {
		Id uint
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if data.Id > 0 {
		res := h.DB.Where("id = ?", data.Id).Delete(&model.Reward{})
		if res.Error != nil {
			logger.Error("error with update database：", res.Error)
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}
