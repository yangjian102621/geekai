package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatRoleHandler struct {
	BaseHandler
}

func NewChatRoleHandler(app *core.AppServer, db *gorm.DB) *ChatRoleHandler {
	return &ChatRoleHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// List 获取用户聊天应用列表
func (h *ChatRoleHandler) List(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var roles []model.ChatRole
	query := h.DB.Where("enable", true)
	if userId > 0 {
		var user model.User
		h.DB.First(&user, userId)
		var roleKeys []string
		err := utils.JsonDecode(user.ChatRoles, &roleKeys)
		if err != nil {
			resp.ERROR(c, "角色解析失败！")
			return
		}
		query = query.Where("marker IN ?", roleKeys)
	}
	if id > 0 {
		query = query.Or("id", id)
	}
	res := h.DB.Where("enable", true).Order("sort_num ASC").Find(&roles)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var roleVos = make([]vo.ChatRole, 0)
	for _, r := range roles {
		var v vo.ChatRole
		err := utils.CopyObject(r, &v)
		if err == nil {
			v.Id = r.Id
			roleVos = append(roleVos, v)
		}
	}
	resp.SUCCESS(c, roleVos)
}

// UpdateRole 更新用户聊天角色
func (h *ChatRoleHandler) UpdateRole(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	var data struct {
		Keys []string `json:"keys"`
	}
	if err = c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err = h.DB.Model(&model.User{}).Where("id = ?", user.Id).UpdateColumn("chat_roles_json", utils.JsonEncode(data.Keys)).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
