package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"strings"
	"sync"
)

type RewardHandler struct {
	BaseHandler
	db   *gorm.DB
	lock sync.Mutex
}

func NewRewardHandler(server *core.AppServer, db *gorm.DB) *RewardHandler {
	h := RewardHandler{db: db, lock: sync.Mutex{}}
	h.App = server
	return &h
}

// Verify 打赏码核销
func (h *RewardHandler) Verify(c *gin.Context) {
	var data struct {
		TxId string `json:"tx_id"`
		Type string `json:"type"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.HACKER(c)
		return
	}

	// 移除转账单号中间的空格，防止有人复制的时候多复制了空格
	data.TxId = strings.ReplaceAll(data.TxId, " ", "")

	h.lock.Lock()
	defer h.lock.Unlock()

	var item model.Reward
	res := h.db.Where("tx_id = ?", data.TxId).First(&item)
	if res.Error != nil {
		resp.ERROR(c, "无效的众筹交易流水号！")
		return
	}

	if item.Status {
		resp.ERROR(c, "当前众筹交易流水号已经被核销，请不要重复核销！")
		return
	}

	tx := h.db.Begin()
	exchange := vo.RewardExchange{}
	if data.Type == "chat" {
		calls := math.Ceil(item.Amount / h.App.SysConfig.ChatCallPrice)
		exchange.Calls = int(calls)
		res = h.db.Model(&user).UpdateColumn("calls", gorm.Expr("calls + ?", calls))
	} else if data.Type == "img" {
		calls := math.Ceil(item.Amount / h.App.SysConfig.ImgCallPrice)
		exchange.ImgCalls = int(calls)
		res = h.db.Model(&user).UpdateColumn("img_calls", gorm.Expr("img_calls + ?", calls))
	}
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	// 更新核销状态
	item.Status = true
	item.UserId = user.Id
	item.Exchange = utils.JsonEncode(exchange)
	res = h.db.Updates(&item)
	if res.Error != nil {
		tx.Rollback()
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	tx.Commit()
	resp.SUCCESS(c)

}
