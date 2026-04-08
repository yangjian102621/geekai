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
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PowerLogHandler struct {
	BaseHandler
}

func NewPowerLogHandler(app *core.AppServer, db *gorm.DB) *PowerLogHandler {
	return &PowerLogHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *PowerLogHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/powerLog/")

	// 需要用户授权的接口
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.POST("list", h.List)
		group.GET("stats", h.Stats)
	}
}

func (h *PowerLogHandler) List(c *gin.Context) {
	var data struct {
		Model    string   `json:"model"`
		Date     []string `json:"date"`
		Page     int      `json:"page"`
		PageSize int      `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	userId := h.GetLoginUserId(c)
	session = session.Where("user_id", userId)
	if data.Model != "" {
		session = session.Where("model", data.Model)
	}
	if len(data.Date) == 2 {
		start := data.Date[0] + " 00:00:00"
		end := data.Date[1] + " 00:00:00"
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}

	var total int64
	session.Model(&model.PowerLog{}).Count(&total)
	var items []model.PowerLog
	var list = make([]vo.PowerLog, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var log vo.PowerLog
			err := utils.CopyObject(item, &log)
			if err != nil {
				continue
			}
			log.Id = item.Id
			log.CreatedAt = item.CreatedAt.Unix()
			log.TypeStr = item.Type.String()
			list = append(list, log)
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}

// Stats 获取用户算力统计
func (h *PowerLogHandler) Stats(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.NotAuth(c)
		return
	}

	// 获取用户信息（包含余额）
	var user model.User
	if err := h.DB.Where("id", userId).First(&user).Error; err != nil {
		resp.ERROR(c, "用户不存在")
		return
	}

	// 计算总消费（所有支出记录）
	var totalConsume int64
	h.DB.Model(&model.PowerLog{}).
		Where("user_id", userId).
		Where("mark", types.PowerSub).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalConsume)

	// 计算今日消费
	today := time.Now().Format("2006-01-02")
	var todayConsume int64
	h.DB.Model(&model.PowerLog{}).
		Where("user_id", userId).
		Where("mark", types.PowerSub).
		Where("DATE(created_at) = ?", today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayConsume)

	stats := map[string]interface{}{
		"total":   totalConsume,
		"today":   todayConsume,
		"balance": user.Power,
	}

	resp.SUCCESS(c, stats)
}
