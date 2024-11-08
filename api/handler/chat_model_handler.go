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
	session := h.DB.Session(&gorm.Session{}).Where("enabled", true)
	t := c.Query("type")
	if t != "" {
		session = session.Where("type", t)
	}

	session = session.Where("open", true)
	if h.IsLogin(c) {
		user, _ := h.GetLoginUser(c)
		var models []int
		err := utils.JsonDecode(user.ChatModels, &models)
		// 查询用户有权限访问的模型以及所有开放的模型
		if err == nil {
			session = session.Or("id IN ?", models)
		}

	}

	res := session.Order("sort_num ASC").Find(&items)
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
