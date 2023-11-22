package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type DashboardHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewDashboardHandler(app *core.AppServer, db *gorm.DB) *DashboardHandler {
	h := DashboardHandler{db: db}
	h.App = app
	return &h
}

type statsVo struct {
	Users  int64   `json:"users"`
	Chats  int64   `json:"chats"`
	Tokens int     `json:"tokens"`
	Income float64 `json:"income"`
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	stats := statsVo{}
	// new users statistic
	var userCount int64
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	res := h.db.Model(&model.User{}).Where("created_at > ?", zeroTime).Count(&userCount)
	if res.Error == nil {
		stats.Users = userCount
	}

	// new chats statistic
	var chatCount int64
	res = h.db.Model(&model.ChatItem{}).Where("created_at > ?", zeroTime).Count(&chatCount)
	if res.Error == nil {
		stats.Chats = chatCount
	}

	// tokens took stats
	var historyMessages []model.HistoryMessage
	res = h.db.Where("created_at > ?", zeroTime).Find(&historyMessages)
	for _, item := range historyMessages {
		stats.Tokens += item.Tokens
	}

	// 众筹收入
	var rewards []model.Reward
	res = h.db.Where("created_at > ?", zeroTime).Find(&rewards)
	for _, item := range rewards {
		stats.Income += item.Amount
	}

	// 订单收入
	var orders []model.Order
	res = h.db.Where("status = ?", types.OrderPaidSuccess).Where("created_at > ?", zeroTime).Find(&orders)
	for _, item := range orders {
		stats.Income += item.Amount
	}
	resp.SUCCESS(c, stats)
}
