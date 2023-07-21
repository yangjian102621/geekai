package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RewardHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewRewardHandler(server *core.AppServer, db *gorm.DB) *RewardHandler {
	h := RewardHandler{db: db}
	h.App = server
	return &h
}

// Verify 打赏码核销
func (h *RewardHandler) Verify(c *gin.Context) {
	var data struct {
		TxId string
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var item model.Reward
	res := h.db.Where("tx_id = ?", data.TxId).First(&item)
	if res.Error != nil {
		resp.ERROR(c, "无效的打赏交易流水号！")
		return
	}

	if item.Status {
		resp.ERROR(c, "当前打赏交易流水号已经被核销，请不要重复核销！")
		return
	}

	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.HACKER(c)
		return
	}

	tx := h.db.Begin()
	calls := (item.Amount + 0.01) * 10
	res = h.db.Model(&user).UpdateColumn("calls", gorm.Expr("calls + ?", calls))
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	// 更新核销状态
	item.Status = true
	res = h.db.Updates(&item)
	if res.Error != nil {
		tx.Rollback()
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	tx.Commit()
	resp.SUCCESS(c)

}
