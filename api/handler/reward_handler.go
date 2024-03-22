package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"strings"
	"sync"
	"time"
)

type RewardHandler struct {
	BaseHandler
	lock sync.Mutex
}

func NewRewardHandler(app *core.AppServer, db *gorm.DB) *RewardHandler {
	return &RewardHandler{BaseHandler: BaseHandler{App: app, DB: db}}
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

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.HACKER(c)
		return
	}

	// 移除转账单号中间的空格，防止有人复制的时候多复制了空格
	data.TxId = strings.ReplaceAll(data.TxId, " ", "")

	h.lock.Lock()
	defer h.lock.Unlock()

	var item model.Reward
	res := h.DB.Where("tx_id = ?", data.TxId).First(&item)
	if res.Error != nil {
		resp.ERROR(c, "无效的众筹交易流水号！")
		return
	}

	if item.Status {
		resp.ERROR(c, "当前众筹交易流水号已经被核销，请不要重复核销！")
		return
	}

	tx := h.DB.Begin()
	exchange := vo.RewardExchange{}
	power := math.Ceil(item.Amount / h.App.SysConfig.PowerPrice)
	exchange.Power = int(power)
	res = tx.Model(&user).UpdateColumn("power", gorm.Expr("power + ?", exchange.Power))
	if res.Error != nil {
		tx.Rollback()
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	// 更新核销状态
	item.Status = true
	item.UserId = user.Id
	item.Exchange = utils.JsonEncode(exchange)
	res = tx.Updates(&item)
	if res.Error != nil {
		tx.Rollback()
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	// 记录算力充值日志
	h.DB.Create(&model.PowerLog{
		UserId:    user.Id,
		Username:  user.Username,
		Type:      types.PowerReward,
		Amount:    exchange.Power,
		Balance:   user.Power + exchange.Power,
		Mark:      types.PowerAdd,
		Model:     "众筹支付",
		Remark:    fmt.Sprintf("众筹充值算力，金额：%f，价格：%f", item.Amount, h.App.SysConfig.PowerPrice),
		CreatedAt: time.Now(),
	})
	tx.Commit()
	resp.SUCCESS(c)

}
