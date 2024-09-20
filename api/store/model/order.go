package model

import (
	"geekai/core/types"
)

// Order 充值订单
type Order struct {
	BaseModel
	UserId    uint
	ProductId uint
	Username  string
	OrderNo   string
	TradeNo   string
	Subject   string
	Amount    float64
	Status    types.OrderStatus
	Remark    string
	PayTime   int64
	PayWay    string // 支付渠道
	PayType   string // 支付类型
}
