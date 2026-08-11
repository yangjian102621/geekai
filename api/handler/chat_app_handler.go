package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatAppHandler struct {
	BaseHandler
}

func NewChatAppHandler(app *core.AppServer, db *gorm.DB) *ChatAppHandler {
	return &ChatAppHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *ChatAppHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/app/")
	group.GET("list", h.List)

	// 需要用户授权的接口
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.GET("list/user", h.ListByUser)
		group.POST("create", h.Create)
		group.POST("copy", h.Copy)
		group.POST("update", h.UpdateApp)
		group.POST("workspace", h.UpdateWorkArea)
		group.POST("remove", h.Remove)
	}
}

// List 获取用户聊天应用列表
func (h *ChatAppHandler) List(c *gin.Context) {
	tid := h.GetInt(c, "tid", 0)
	var roles []model.ChatApp
	session := h.DB.Where("enable = ? AND user_id = 0", true)
	if tid > 0 {
		session = session.Where("tid", tid)
	}
	err := session.Order("sort_num ASC").Find(&roles).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var roleVos = make([]vo.ChatApp, 0)
	for _, r := range roles {
		var v vo.ChatApp
		err := utils.CopyObject(r, &v)
		if err == nil {
			v.Id = r.Id
			if r.UserId == 0 {
				v.SystemPrompt = ""
			}
			roleVos = append(roleVos, v)
		}
	}
	resp.SUCCESS(c, roleVos)
}

// ListByUser 获取用户添加的角色列表
func (h *ChatAppHandler) ListByUser(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var roles []model.ChatApp
	session := h.DB.Where("enable = ?", true)
	if userId > 0 {
		session = session.Where("(user_id = 0 OR user_id = ?)", userId)
	} else {
		session = session.Where("user_id = 0")
	}

	if id > 0 {
		session = session.Or("id", id)
	}
	res := session.Order("sort_num ASC").Find(&roles)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var roleVos = make([]vo.ChatApp, 0)
	for _, r := range roles {
		var v vo.ChatApp
		err := utils.CopyObject(r, &v)
		if err == nil {
			v.Id = r.Id
			if r.UserId == 0 {
				v.SystemPrompt = ""
			}
			roleVos = append(roleVos, v)
		}
	}
	resp.SUCCESS(c, roleVos)
}

// Create 用户创建智能体
func (h *ChatAppHandler) Create(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.NotAuth(c)
		return
	}
	var data vo.ChatApp
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	role := model.ChatApp{
		Name:         data.Name,
		Tid:          data.Tid,
		UserId:       userId,
		SystemPrompt: data.SystemPrompt,
		HelloMsg:     data.HelloMsg,
		Icon:         data.Icon,
		Enable:       true,
		SortNum:      int(data.SortNum),
		ModelId:      data.ModelId,
	}
	if role.Icon == "" {
		role.Icon = "/images/avatar/gpt.png"
	}
	if err := h.DB.Create(&role).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	data.Id = role.Id
	data.UserId = role.UserId
	resp.SUCCESS(c, data)
}

// Copy 用户复制智能体（复制为当前用户名下）
func (h *ChatAppHandler) Copy(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.NotAuth(c)
		return
	}
	var body struct {
		SourceId uint `json:"source_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.SourceId == 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var src model.ChatApp
	if err := h.DB.First(&src, body.SourceId).Error; err != nil {
		resp.ERROR(c, "智能体不存在")
		return
	}
	role := model.ChatApp{
		Name:         src.Name,
		Tid:          src.Tid,
		UserId:       userId,
		SystemPrompt: src.SystemPrompt,
		HelloMsg:     src.HelloMsg,
		Icon:         src.Icon,
		Enable:       true,
		SortNum:      src.SortNum,
		ModelId:      src.ModelId,
	}
	if err := h.DB.Create(&role).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, gin.H{"id": role.Id})
}

// UpdateApp 更新用户聊天应用（仅允许更新自己创建的）
func (h *ChatAppHandler) UpdateApp(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.NotAuth(c)
		return
	}
	var data vo.ChatApp
	if err := c.ShouldBindJSON(&data); err != nil || data.Id == 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var role model.ChatApp
	if err := h.DB.First(&role, data.Id).Error; err != nil {
		resp.ERROR(c, "智能体不存在")
		return
	}
	if role.UserId != userId {
		resp.ERROR(c, "无权限修改该智能体")
		return
	}
	updates := map[string]interface{}{
		"name":          data.Name,
		"hello_msg":     data.HelloMsg,
		"icon":          data.Icon,
		"model_id":      data.ModelId,
		"system_prompt": data.SystemPrompt,
	}
	if err := h.DB.Model(&role).Updates(updates).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, nil)
}

// UpdateWorkArea 更新用户工作区应用列表（存为应用 id 数组）
func (h *ChatAppHandler) UpdateWorkArea(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.NotAuth(c)
		return
	}
	var body struct {
		Ids []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if err := h.DB.Model(&model.User{}).Where("id = ?", userId).Update("chat_roles_json", utils.JsonEncode(body.Ids)).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, nil)
}

// Remove 删除用户智能体（仅允许删除自己创建的）
func (h *ChatAppHandler) Remove(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.NotAuth(c)
		return
	}
	var body struct {
		Id uint `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.Id == 0 {
		body.Id = uint(h.GetInt(c, "id", 0))
	}
	if body.Id == 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var role model.ChatApp
	if err := h.DB.First(&role, body.Id).Error; err != nil {
		resp.ERROR(c, "智能体不存在")
		return
	}
	if role.UserId != userId {
		resp.ERROR(c, "无权限删除该智能体")
		return
	}
	if err := h.DB.Delete(&role).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, nil)
}
