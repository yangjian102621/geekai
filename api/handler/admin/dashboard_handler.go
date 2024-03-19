package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type DashboardHandler struct {
	handler.BaseHandler
}

func NewDashboardHandler(app *core.AppServer, db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

type statsVo struct {
	Users  int64                         `json:"users"`
	Chats  int64                         `json:"chats"`
	Tokens int                           `json:"tokens"`
	Income float64                       `json:"income"`
	Chart  map[string]map[string]float64 `json:"chart"`
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	stats := statsVo{}
	// new users statistic
	var userCount int64
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	res := h.DB.Model(&model.User{}).Where("created_at > ?", zeroTime).Count(&userCount)
	if res.Error == nil {
		stats.Users = userCount
	}

	// new chats statistic
	var chatCount int64
	res = h.DB.Model(&model.ChatItem{}).Where("created_at > ?", zeroTime).Count(&chatCount)
	if res.Error == nil {
		stats.Chats = chatCount
	}

	// tokens took stats
	var historyMessages []model.ChatMessage
	res = h.DB.Where("created_at > ?", zeroTime).Find(&historyMessages)
	for _, item := range historyMessages {
		stats.Tokens += item.Tokens
	}

	// 众筹收入
	var rewards []model.Reward
	res = h.DB.Where("created_at > ?", zeroTime).Find(&rewards)
	for _, item := range rewards {
		stats.Income += item.Amount
	}

	// 订单收入
	var orders []model.Order
	res = h.DB.Where("status = ?", types.OrderPaidSuccess).Where("created_at > ?", zeroTime).Find(&orders)
	for _, item := range orders {
		stats.Income += item.Amount
	}

	// 统计7天的订单的图表
	startDate := now.Add(-7 * 24 * time.Hour).Format("2006-01-02")
	var statsChart = make(map[string]map[string]float64)
	//// 初始化
	var userStatistic, historyMessagesStatistic, incomeStatistic = make(map[string]float64), make(map[string]float64), make(map[string]float64)
	for i := 0; i < 7; i++ {
		var initTime = time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
		userStatistic[initTime] = float64(0)
		historyMessagesStatistic[initTime] = float64(0)
		incomeStatistic[initTime] = float64(0)
	}

	// 统计用户7天增加的曲线
	var users []model.User
	res = h.DB.Model(&model.User{}).Where("created_at > ?", startDate).Find(&users)
	if res.Error == nil {
		for _, item := range users {
			userStatistic[item.CreatedAt.Format("2006-01-02")] += 1
		}
	}

	// 统计7天Token 消耗
	res = h.DB.Where("created_at > ?", startDate).Find(&historyMessages)
	for _, item := range historyMessages {
		historyMessagesStatistic[item.CreatedAt.Format("2006-01-02")] += float64(item.Tokens)
	}

	// 浮点数相加？
	// 统计最近7天的众筹
	res = h.DB.Where("created_at > ?", startDate).Find(&rewards)
	for _, item := range rewards {
		incomeStatistic[item.CreatedAt.Format("2006-01-02")], _ = decimal.NewFromFloat(incomeStatistic[item.CreatedAt.Format("2006-01-02")]).Add(decimal.NewFromFloat(item.Amount)).Float64()
	}

	// 统计最近7天的订单
	res = h.DB.Where("status = ?", types.OrderPaidSuccess).Where("created_at > ?", startDate).Find(&orders)
	for _, item := range orders {
		incomeStatistic[item.CreatedAt.Format("2006-01-02")], _ = decimal.NewFromFloat(incomeStatistic[item.CreatedAt.Format("2006-01-02")]).Add(decimal.NewFromFloat(item.Amount)).Float64()
	}

	statsChart["users"] = userStatistic
	statsChart["historyMessage"] = historyMessagesStatistic
	statsChart["orders"] = incomeStatistic

	stats.Chart = statsChart

	resp.SUCCESS(c, stats)
}
