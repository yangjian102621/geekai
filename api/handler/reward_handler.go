package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
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

func (h *RewardHandler) Notify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != h.App.Config.ExtConfig.Token {
		resp.NotAuth(c)
		return
	}

	var data struct {
		TransId string  `json:"trans_id"` // 微信转账交易 ID
		Amount  float64 `json:"amount"`   // 微信转账交易金额
		Remark  string  `json:"remark"`   // 转账备注
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Amount <= 0 {
		resp.ERROR(c, "Amount should not be 0")
		return
	}

	logger.Infof("收到众筹收款信息: %+v", data)
	var item model.Reward
	res := h.db.Where("tx_id = ?", data.TransId).First(&item)
	if res.Error == nil {
		resp.ERROR(c, "当前交易 ID 己经存在！")
		return
	}

	res = h.db.Create(&model.Reward{
		TxId:   data.TransId,
		Amount: data.Amount,
		Remark: data.Remark,
		Status: false,
	})
	if res.Error != nil {
		logger.Errorf("交易保存失败: %v", res.Error)
		resp.ERROR(c, "交易保存失败")
		return
	}
	resp.SUCCESS(c)
}

// Verify 打赏码核销
func (h *RewardHandler) Verify(c *gin.Context) {
	var data struct {
		TxId string `json:"tx_id"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 移除转账单号中间的空格，防止有人复制的时候多复制了空格
	data.TxId = strings.ReplaceAll(data.TxId, " ", "")

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

	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.HACKER(c)
		return
	}

	tx := h.db.Begin()
	calls := (item.Amount + 0.1) * 10
	res = h.db.Model(&user).UpdateColumn("calls", gorm.Expr("calls + ?", calls))
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	// 更新核销状态
	item.Status = true
	item.UserId = user.Id
	res = h.db.Updates(&item)
	if res.Error != nil {
		tx.Rollback()
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	tx.Commit()
	resp.SUCCESS(c)

}
