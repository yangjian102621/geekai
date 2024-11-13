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
)

type ChatModelHandler struct {
	BaseHandler
}

func NewChatModelHandler(app *core.AppServer, db *gorm.DB) *ChatModelHandler {
	return &ChatModelHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// List 模型列表
func (h *ChatModelHandler) List(c *gin.Context) {
	var items []model.ChatModel
	var chatModels = make([]vo.ChatModel, 0)
	var res *gorm.DB
	session := h.DB.Session(&gorm.Session{}).Where("enabled", true)
	t := c.Query("type")
	if t != "" {
		session = session.Where("type", t)
	}
	// 如果用户没有登录，则加载所有开放模型
	if !h.IsLogin(c) {
		res = session.Where("open", true).Order("sort_num ASC").Find(&items)
	} else {
		user, _ := h.GetLoginUser(c)
		var models []int
		err := utils.JsonDecode(user.ChatModels, &models)
		if err != nil {
			resp.ERROR(c, "当前用户没有订阅任何模型")
			return
		}
		// 查询用户有权限访问的模型以及所有开放的模型
		res = h.DB.Where("enabled = ?", true).Where(
			h.DB.Where("id IN ?", models).Or("open", true),
		).Order("sort_num ASC").Find(&items)
	}

	if res.Error == nil {
		for _, item := range items {
			var cm vo.ChatModel
			err := utils.CopyObject(item, &cm)
			if err == nil {
				cm.Id = item.Id
				cm.CreatedAt = item.CreatedAt.Unix()
				cm.UpdatedAt = item.UpdatedAt.Unix()
				chatModels = append(chatModels, cm)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, chatModels)
}
