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
	"geekai/utils/resp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	handler.BaseHandler
}

func NewDashboardHandler(app *core.AppServer, db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *DashboardHandler) RegisterRoutes(r *core.Router) {
	// 需要管理员授权的接口（改造前漏挂管理员鉴权中间件，本次修复）
	group := r.AdminGroup("/api/admin/dashboard/")
	group.GET("stats", h.Stats)
}

// statsVo 增加 recentOrders、recentUsers 字段
// 最近订单
type OrderBrief struct {
	OrderNo   string    `json:"order_no"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// 最近用户
type UserBrief struct {
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	LastActive time.Time `json:"last_active"`
}

type statsVo struct {
	Users          int64                         `json:"users"`
	Chats          int64                         `json:"chats"`
	Tokens         int                           `json:"tokens"`
	Income         float64                       `json:"income"`
	Chart          map[string]map[string]float64 `json:"chart"`
	TodayUsers     int64                         `json:"todayUsers"`
	TodayChats     int64                         `json:"todayChats"`
	TodayTokens    int                           `json:"todayTokens"`
	TodayIncome    float64                       `json:"todayIncome"`
	TodayOrders    int64                         `json:"todayOrders"`
	TodayImageJobs int64                         `json:"todayImageJobs"`
	TodayVideoJobs int64                         `json:"todayVideoJobs"`
	TodayMusicJobs int64                         `json:"todayMusicJobs"`
	Orders         int64                         `json:"orders"`
	ImageJobs      int64                         `json:"imageJobs"`
	VideoJobs      int64                         `json:"videoJobs"`
	MusicJobs      int64                         `json:"musicJobs"`
	RecentOrders   []OrderBrief                  `json:"recentOrders"`
	RecentUsers    []UserBrief                   `json:"recentUsers"`
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	stats := statsVo{}
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 总用户数
	h.DB.Model(&model.User{}).Count(&stats.Users)

	// 今日新增用户
	h.DB.Model(&model.User{}).Where("created_at > ?", zeroTime).Count(&stats.TodayUsers)

	// 总对话数
	h.DB.Model(&model.ChatItem{}).Count(&stats.Chats)

	// 今日新增对话
	h.DB.Model(&model.ChatItem{}).Where("created_at > ?", zeroTime).Count(&stats.TodayChats)

	// 总积分消耗
	var powerLogs []model.PowerLog
	h.DB.Where("mark = ?", types.PowerSub).Find(&powerLogs)
	for _, item := range powerLogs {
		stats.Tokens += item.Amount
	}

	// 今日积分消耗
	var todayPowerLogs []model.PowerLog
	h.DB.Where("mark = ?", types.PowerSub).Where("created_at > ?", zeroTime).Find(&todayPowerLogs)
	for _, item := range todayPowerLogs {
		stats.TodayTokens += item.Amount
	}

	// 总收入
	var allOrders []model.Order
	h.DB.Where("status = ?", types.OrderPaidSuccess).Find(&allOrders)
	for _, item := range allOrders {
		stats.Income += item.Amount
	}

	// 今日收入
	var todayOrders []model.Order
	h.DB.Where("status = ?", types.OrderPaidSuccess).Where("created_at > ?", zeroTime).Find(&todayOrders)
	for _, item := range todayOrders {
		stats.TodayIncome += item.Amount
	}

	// 订单总数
	h.DB.Model(&model.Order{}).Where("status = ?", types.OrderPaidSuccess).Count(&stats.Orders)

	// 今日订单数
	h.DB.Model(&model.Order{}).Where("status = ?", types.OrderPaidSuccess).Where("created_at > ?", zeroTime).Count(&stats.TodayOrders)

	// 图片生成任务统计
	var mjJobs, imageJobs, jimengImageJobs int64
	h.DB.Model(&model.MidJourneyJob{}).Count(&mjJobs)
	h.DB.Model(&model.ImageJob{}).Count(&imageJobs)
	h.DB.Model(&model.JimengJob{}).Where("type IN ?", []string{"text_to_image", "image_to_image", "image_edit", "image_effects"}).Count(&jimengImageJobs)
	stats.ImageJobs = mjJobs + imageJobs + jimengImageJobs

	logger.Info("stats.ImageJobs", stats.ImageJobs)

	// 今日图片生成任务统计
	var todayMjJobs, todayImageJobs, todayJimengImageJobs int64
	h.DB.Model(&model.MidJourneyJob{}).Where("created_at > ?", zeroTime).Count(&todayMjJobs)
	h.DB.Model(&model.ImageJob{}).Where("created_at > ?", zeroTime).Count(&todayImageJobs)
	h.DB.Model(&model.JimengJob{}).Where("type IN ?", []string{"text_to_image", "image_to_image", "image_edit", "image_effects"}).Where("created_at > ?", zeroTime).Count(&todayJimengImageJobs)
	stats.TodayImageJobs = todayMjJobs + todayImageJobs + todayJimengImageJobs

	// 视频生成任务统计
	var videoJobs, jimengVideoJobs int64
	h.DB.Model(&model.VideoJob{}).Count(&videoJobs)
	h.DB.Model(&model.JimengJob{}).Where("type IN ?", []string{"text_to_video", "image_to_video"}).Count(&jimengVideoJobs)
	stats.VideoJobs = videoJobs + jimengVideoJobs

	// 今日视频生成任务统计
	var todayVideoJobs, todayJimengVideoJobs int64
	h.DB.Model(&model.VideoJob{}).Where("created_at > ?", zeroTime).Count(&todayVideoJobs)
	h.DB.Model(&model.JimengJob{}).Where("type IN ?", []string{"text_to_video", "image_to_video"}).Where("created_at > ?", zeroTime).Count(&todayJimengVideoJobs)
	stats.TodayVideoJobs = todayVideoJobs + todayJimengVideoJobs

	// 音乐生成任务统计
	h.DB.Model(&model.SunoJob{}).Count(&stats.MusicJobs)

	// 今日音乐生成任务统计
	h.DB.Model(&model.SunoJob{}).Where("created_at > ?", zeroTime).Count(&stats.TodayMusicJobs)

	// recentOrders: 最近10条已支付订单
	var orderList []model.Order
	h.DB.Model(&model.Order{}).Where("status = ?", types.OrderPaidSuccess).Order("created_at desc").Limit(10).Find(&orderList)
	for _, o := range orderList {
		stats.RecentOrders = append(stats.RecentOrders, OrderBrief{
			OrderNo:   o.OrderNo,
			Amount:    o.Amount,
			CreatedAt: o.CreatedAt,
		})
	}
	// recentUsers: 最近10个注册用户
	var userList []model.User
	h.DB.Model(&model.User{}).Order("created_at desc").Limit(10).Find(&userList)
	for _, u := range userList {
		lastActive := u.UpdatedAt
		if lastActive.IsZero() {
			lastActive = u.CreatedAt
		}
		stats.RecentUsers = append(stats.RecentUsers, UserBrief{
			Nickname:   u.Nickname,
			Avatar:     u.Avatar,
			LastActive: lastActive,
		})
	}

	// 趋势图表筛选：默认最近 7 天，支持 days，或 start/end + granularity + 用户名称
	days := 7
	if v, err := strconv.Atoi(c.Query("days")); err == nil && v > 0 && v <= 365 {
		days = v
	}
	granularity := c.DefaultQuery("granularity", "day")
	var startTime, endTime time.Time
	if st, err := time.ParseInLocation("2006-01-02 15:04:05", c.Query("start"), now.Location()); err == nil {
		if et, err2 := time.ParseInLocation("2006-01-02 15:04:05", c.Query("end"), now.Location()); err2 == nil && et.After(st) {
			startTime, endTime = st, et
		}
	}
	if startTime.IsZero() || endTime.IsZero() {
		granularity = "day"
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		startTime = time.Date(now.Year(), now.Month(), now.Day()-(days-1), 0, 0, 0, 0, now.Location())
	}
	// 按小时统计最多支持 31 天，超出则退化为按天
	if granularity == "hour" && endTime.Sub(startTime) > 31*24*time.Hour {
		granularity = "day"
	}
	if granularity != "hour" {
		granularity = "day"
	}

	// 按用户名称过滤（仅作用于趋势图）
	var filterUserId uint
	filterByUser := false
	if name := c.Query("username"); name != "" {
		filterByUser = true
		var filterUser model.User
		if err := h.DB.Where("username = ? OR nickname = ?", name, name).First(&filterUser).Error; err == nil {
			filterUserId = filterUser.Id
		}
	}

	bucketKey := func(t time.Time) string {
		if granularity == "hour" {
			return t.Format("2006-01-02 15:00")
		}
		return t.Format("2006-01-02")
	}

	var statsChart = make(map[string]map[string]float64)
	chartKeys := []string{"users", "historyMessage", "orders", "orderCount", "chats", "imageJobs", "videoJobs", "musicJobs"}
	for _, k := range chartKeys {
		statsChart[k] = make(map[string]float64)
	}
	//// 初始化
	if granularity == "hour" {
		for t := startTime.Truncate(time.Hour); !t.After(endTime); t = t.Add(time.Hour) {
			key := bucketKey(t)
			for _, k := range chartKeys {
				statsChart[k][key] = 0
			}
		}
	} else {
		startDay := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		for t := startDay; !t.After(endTime); t = t.AddDate(0, 0, 1) {
			key := bucketKey(t)
			for _, k := range chartKeys {
				statsChart[k][key] = 0
			}
		}
	}

	// 统计区间内新增用户曲线
	userQuery := h.DB.Model(&model.User{}).Where("created_at BETWEEN ? AND ?", startTime, endTime)
	if filterByUser {
		userQuery = userQuery.Where("id = ?", filterUserId)
	}
	var users []model.User
	if err := userQuery.Find(&users).Error; err == nil {
		for _, item := range users {
			statsChart["users"][bucketKey(item.CreatedAt)] += 1
		}
	}

	// 统计积分消耗曲线
	powerQuery := h.DB.Where("mark = ?", types.PowerSub).Where("created_at BETWEEN ? AND ?", startTime, endTime)
	if filterByUser {
		powerQuery = powerQuery.Where("user_id = ?", filterUserId)
	}
	var chartPowerLogs []model.PowerLog
	if err := powerQuery.Find(&chartPowerLogs).Error; err == nil {
		for _, item := range chartPowerLogs {
			statsChart["historyMessage"][bucketKey(item.CreatedAt)] += float64(item.Amount)
		}
	}

	// 统计订单收入与订单数曲线
	orderQuery := h.DB.Where("status = ?", types.OrderPaidSuccess).Where("created_at BETWEEN ? AND ?", startTime, endTime)
	if filterByUser {
		orderQuery = orderQuery.Where("user_id = ?", filterUserId)
	}
	var chartOrders []model.Order
	if err := orderQuery.Find(&chartOrders).Error; err == nil {
		for _, item := range chartOrders {
			key := bucketKey(item.CreatedAt)
			statsChart["orderCount"][key] += 1
			statsChart["orders"][key], _ = decimal.NewFromFloat(statsChart["orders"][key]).Add(decimal.NewFromFloat(item.Amount)).Float64()
		}
	}

	// 统计对话数曲线
	var chartChats []model.ChatItem
	if err := h.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&chartChats).Error; err == nil {
		for _, item := range chartChats {
			statsChart["chats"][bucketKey(item.CreatedAt)] += 1
		}
	}

	// 统计图片生成曲线
	var chartMjJobs []model.MidJourneyJob
	if err := h.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&chartMjJobs).Error; err == nil {
		for _, item := range chartMjJobs {
			statsChart["imageJobs"][bucketKey(item.CreatedAt)] += 1
		}
	}
	var chartImageJobs []model.ImageJob
	if err := h.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&chartImageJobs).Error; err == nil {
		for _, item := range chartImageJobs {
			statsChart["imageJobs"][bucketKey(item.CreatedAt)] += 1
		}
	}

	// 统计视频生成曲线
	var chartVideoJobs []model.VideoJob
	if err := h.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&chartVideoJobs).Error; err == nil {
		for _, item := range chartVideoJobs {
			statsChart["videoJobs"][bucketKey(item.CreatedAt)] += 1
		}
	}

	// 即梦任务按类型归入图片/视频曲线
	var chartJimengJobs []model.JimengJob
	if err := h.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&chartJimengJobs).Error; err == nil {
		for _, item := range chartJimengJobs {
			switch item.Type {
			case "text_to_image", "image_to_image", "image_edit", "image_effects":
				statsChart["imageJobs"][bucketKey(item.CreatedAt)] += 1
			case "text_to_video", "image_to_video":
				statsChart["videoJobs"][bucketKey(item.CreatedAt)] += 1
			}
		}
	}

	// 统计音乐生成曲线
	var chartSunoJobs []model.SunoJob
	if err := h.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime).Find(&chartSunoJobs).Error; err == nil {
		for _, item := range chartSunoJobs {
			statsChart["musicJobs"][bucketKey(item.CreatedAt)] += 1
		}
	}

	stats.Chart = statsChart

	resp.SUCCESS(c, stats)
}
