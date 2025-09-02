package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service/moderation"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModerationHandler struct {
	handler.BaseHandler
	sysConfig         *types.SystemConfig
	moderationManager *moderation.ServiceManager
}

func NewModerationHandler(app *core.AppServer, db *gorm.DB, sysConfig *types.SystemConfig, moderationManager *moderation.ServiceManager) *ModerationHandler {
	return &ModerationHandler{BaseHandler: handler.BaseHandler{DB: db, App: app}, sysConfig: sysConfig, moderationManager: moderationManager}
}

// RegisterRoutes 注册路由
func (h *ModerationHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/admin/moderation/")

	// 需要管理员授权的接口
	group.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		group.POST("list", h.List)
		group.GET("remove", h.Remove)
		group.POST("batch-remove", h.BatchRemove)
		group.GET("source-list", h.GetSourceList)
		group.POST("config", h.UpdateModeration)
		group.POST("test", h.TestModeration)
	}
}

// List 获取文本审核记录列表
func (h *ModerationHandler) List(c *gin.Context) {
	var data struct {
		Username  string `json:"username"`
		Source    string `json:"source"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Page      int    `json:"page"`
		PageSize  int    `json:"page_size"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})

	// 构建查询条件
	if data.Username != "" {
		// 通过用户名查找用户ID
		var user model.User
		if err := h.DB.Where("username LIKE ?", "%"+data.Username+"%").First(&user).Error; err == nil {
			session = session.Where("user_id", user.Id)
		}
	}

	if data.Source != "" {
		session = session.Where("source", data.Source)
	}

	if data.StartDate != "" && data.EndDate != "" {
		startTime := data.StartDate + " 00:00:00"
		endTime := data.EndDate + " 23:59:59"
		session = session.Where("created_at >= ? AND created_at <= ?", startTime, endTime)
	}

	// 统计总数
	var total int64
	session.Model(&model.Moderation{}).Count(&total)

	// 分页
	page := data.Page
	pageSize := data.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	session = session.Offset(offset).Limit(pageSize)

	// 查询数据
	var items []model.Moderation
	err := session.Order("id DESC").Find(&items).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 获取用户信息
	userIds := make([]uint, 0)
	for _, item := range items {
		userIds = append(userIds, item.UserId)
	}

	var users []model.User
	if len(userIds) > 0 {
		h.DB.Where("id IN ?", userIds).Find(&users)
	}

	userMap := make(map[uint]string)
	for _, user := range users {
		userMap[user.Id] = user.Username
	}

	// 转换为响应数据
	list := make([]map[string]any, 0)
	for _, item := range items {
		var moderation types.ModerationResult
		err := utils.JsonDecode(item.Result, &moderation)
		if err != nil {
			continue
		}
		var result []string
		for value, label := range types.ModerationCategories {
			if moderation.Categories[value] {
				result = append(result, label)
			}
		}
		list = append(list, map[string]any{
			"id":         item.Id,
			"user_id":    item.UserId,
			"username":   userMap[item.UserId],
			"source":     item.Source,
			"input":      item.Input,
			"output":     item.Output,
			"result":     result,
			"created_at": item.CreatedAt.Unix(),
		})
	}

	resp.SUCCESS(c, map[string]any{
		"items":     list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *ModerationHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Where("id", id).Delete(&model.Moderation{}).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

// BatchRemove 批量删除文本审核记录
func (h *ModerationHandler) BatchRemove(c *gin.Context) {
	var data struct {
		Ids []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if len(data.Ids) == 0 {
		resp.ERROR(c, "请选择要删除的记录")
		return
	}

	err := h.DB.Where("id IN ?", data.Ids).Delete(&model.Moderation{}).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// 获取 source 列表
func (h *ModerationHandler) GetSourceList(c *gin.Context) {
	sources := []gin.H{
		{
			"id":   types.ModerationSourceChat,
			"name": "AI对话",
		},
		{
			"id":   types.ModerationSourceMJ,
			"name": "Midjourney 绘图",
		},
		{
			"id":   types.ModerationSourceDalle,
			"name": "Dalle 绘图",
		},
		{
			"id":   types.ModerationSourceSD,
			"name": "StableDiffusion 绘图",
		},
		{
			"id":   types.ModerationSourceSuno,
			"name": "Suno 音乐",
		},
		{
			"id":   types.ModerationSourceVideo,
			"name": "视频生成",
		},
		{
			"id":   types.ModerationSourceJiMeng,
			"name": "即梦AI",
		},
	}

	resp.SUCCESS(c, sources)
}

// UpdateModeration 更新文本审查配置
func (h *ModerationHandler) UpdateModeration(c *gin.Context) {
	var data types.ModerationConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Where("name", types.ConfigKeyModeration).FirstOrCreate(&model.Config{Name: types.ConfigKeyModeration, Value: utils.JsonEncode(data)}).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	h.moderationManager.UpdateConfig(data)
	h.sysConfig.Moderation = data

	resp.SUCCESS(c, data)
}

// 测试结果类型，用于前端显示
type ModerationTestResult struct {
	IsAbnormal bool                   `json:"isAbnormal"`
	Details    []ModerationTestDetail `json:"details"`
}

type ModerationTestDetail struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Confidence  string `json:"confidence"`
	IsCategory  bool   `json:"isCategory"`
}

// TestModeration 测试文本审查服务
func (h *ModerationHandler) TestModeration(c *gin.Context) {
	var data struct {
		Text    string `json:"text"`
		Service string `json:"service"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Text == "" {
		resp.ERROR(c, "测试文本不能为空")
		return
	}

	// 检查是否启用了文本审查
	if !h.sysConfig.Moderation.Enable {
		resp.ERROR(c, "文本审查服务未启用")
		return
	}

	// 获取当前激活的审核服务
	service := h.moderationManager.GetService()
	// 执行文本审核
	result, err := service.Moderate(data.Text)
	if err != nil {
		resp.ERROR(c, "审核服务调用失败: "+err.Error())
		return
	}

	// 转换为前端需要的格式
	testResult := ModerationTestResult{
		IsAbnormal: result.Flagged,
		Details:    make([]ModerationTestDetail, 0),
	}

	// 构建详细信息
	for category, description := range types.ModerationCategories {
		score := result.CategoryScores[category]
		isCategory := result.Categories[category]

		testResult.Details = append(testResult.Details, ModerationTestDetail{
			Category:    category,
			Description: description,
			Confidence:  fmt.Sprintf("%.2f", score),
			IsCategory:  isCategory,
		})
	}

	resp.SUCCESS(c, testResult)
}
