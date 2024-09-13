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
	tid := h.GetInt(c, "tid", 0)
	var roles []model.ChatRole
	session := h.DB.Where("enable", true)
	if tid > 0 {
		session = session.Where("tid", tid)
	}
	err := session.Order("sort_num ASC").Find(&roles).Error
	if err != nil {
		resp.ERROR(c, err.Error())
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

// ListByUser 获取用户添加的角色列表
func (h *ChatRoleHandler) ListByUser(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var roles []model.ChatRole
	session := h.DB.Where("enable", true)
	// 如果用户没登录，则获取所有角色
	if userId > 0 {
		var user model.User
		h.DB.First(&user, userId)
		var roleKeys []string
		err := utils.JsonDecode(user.ChatRoles, &roleKeys)
		if err != nil {
			resp.ERROR(c, "角色解析失败！")
			return
		}
		// 保证用户至少有一个角色可用
		if len(roleKeys) > 0 {
			session = session.Where("marker IN ?", roleKeys)
		}
	}

	if id > 0 {
		session = session.Or("id", id)
	}
	res := session.Order("sort_num ASC").Find(&roles)
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
