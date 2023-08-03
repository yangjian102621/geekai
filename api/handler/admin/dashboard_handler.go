package admin

import (
	"chatplus/core"
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
	Users   int64   `json:"users"`
	Chats   int64   `json:"chats"`
	Tokens  int64   `json:"tokens"`
	Rewards float64 `json:"rewards"`
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
	var tokenCount int64
	res = h.db.Model(&model.HistoryMessage{}).Select("sum(tokens) as total").Where("created_at > ?", zeroTime).Scan(&tokenCount)
	if res.Error == nil {
		stats.Tokens = tokenCount
	}

	// reward revenue
	var amount float64
	res = h.db.Model(&model.Reward{}).Select("sum(amount) as total").Where("created_at > ?", zeroTime).Scan(&amount)
	if res.Error == nil {
		stats.Rewards = amount
	}
	resp.SUCCESS(c, stats)
}
