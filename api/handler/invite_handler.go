package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InviteHandler 用户邀请
type InviteHandler struct {
	BaseHandler
}

func NewInviteHandler(app *core.AppServer, db *gorm.DB) *InviteHandler {
	return &InviteHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *InviteHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/invite/")
	group.GET("code", h.Code)
	group.GET("list", h.List)
	group.GET("hits", h.Hits)
	group.GET("stats", h.Stats)
	group.GET("rules", h.Rules)
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
	offset := (page - 1) * pageSize
	err := session.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	userIds := make([]uint, 0)
	for _, item := range items {
		userIds = append(userIds, item.UserId)
	}
	userMap := make(map[uint]model.User)
	var users []model.User
	h.DB.Model(&model.User{}).Where("id IN (?)", userIds).Find(&users)
	for _, user := range users {
		userMap[user.Id] = user
	}

	var list = make([]vo.InviteLog, 0)
	for _, item := range items {
		var v vo.InviteLog
		err := utils.CopyObject(item, &v)
		if err != nil {
			continue
		}
		v.CreatedAt = item.CreatedAt.Unix()
		v.Avatar = userMap[item.UserId].Avatar
		list = append(list, v)
	}
	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, list))
}

// Hits 访问邀请码
func (h *InviteHandler) Hits(c *gin.Context) {
	code := c.Query("code")
	h.DB.Model(&model.InviteCode{}).Where("code = ?", code).UpdateColumn("hits", gorm.Expr("hits + ?", 1))
	resp.SUCCESS(c)
}

// Stats 获取邀请统计
func (h *InviteHandler) Stats(c *gin.Context) {
	userId := h.GetLoginUserId(c)

	// 获取邀请码
	var inviteCode model.InviteCode
	res := h.DB.Where("user_id = ?", userId).First(&inviteCode)
	if res.Error != nil {
		resp.ERROR(c, "邀请码不存在")
		return
	}

	// 统计累计邀请数
	var totalInvite int64
	h.DB.Model(&model.InviteLog{}).Where("inviter_id = ?", userId).Count(&totalInvite)

	// 统计今日邀请数
	today := time.Now().Format("2006-01-02")
	var todayInvite int64
	h.DB.Model(&model.InviteLog{}).Where("inviter_id = ? AND DATE(created_at) = ?", userId, today).Count(&todayInvite)

	// 获取系统配置中的邀请奖励
	var config model.Config
	var invitePower int = 200 // 默认值
	if h.DB.Where("name = ?", "system").First(&config).Error == nil {
		var configMap map[string]any
		if utils.JsonDecode(config.Value, &configMap) == nil {
			if power, ok := configMap["invite_power"].(float64); ok {
				invitePower = int(power)
			}
		}
	}

	// 计算获得奖励总数
	rewardTotal := int(totalInvite) * invitePower

	// 构建邀请链接
	inviteLink := fmt.Sprintf("%s/register?invite=%s", h.App.Config.StaticUrl, inviteCode.Code)

	stats := vo.InviteStats{
		InviteCount: int(totalInvite),
		RewardTotal: rewardTotal,
		TodayInvite: int(todayInvite),
		InviteCode:  inviteCode.Code,
		InviteLink:  inviteLink,
	}

	resp.SUCCESS(c, stats)
}

// Rules 获取奖励规则
func (h *InviteHandler) Rules(c *gin.Context) {
	// 获取系统配置中的邀请奖励
	var config model.Config
	var invitePower int = 200 // 默认值
	if h.DB.Where("name = ?", "system").First(&config).Error == nil {
		var configMap map[string]interface{}
		if utils.JsonDecode(config.Value, &configMap) == nil {
			if power, ok := configMap["invite_power"].(float64); ok {
				invitePower = int(power)
			}
		}
	}

	rules := []vo.RewardRule{
		{
			Id:     1,
			Title:  "好友注册",
			Desc:   "好友通过邀请链接成功注册",
			Icon:   "icon-user-fill",
			Color:  "#1989fa",
			Reward: invitePower,
		},
		{
			Id:     2,
			Title:  "好友首次充值",
			Desc:   "好友首次充值任意金额",
			Icon:   "icon-money",
			Color:  "#07c160",
			Reward: invitePower * 2, // 假设首次充值奖励是注册奖励的2倍
		},
	}

	resp.SUCCESS(c, rules)
}
