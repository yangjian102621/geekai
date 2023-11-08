package model

import (
	"chatplus/core/types"
	"gorm.io/gorm"
)

// Order 充值订单
type Order struct {
	BaseModel
	UserId    uint
	ProductId uint
	Mobile    string
	OrderNo   string
	Subject   string
	Amount    float64
	Status    types.OrderStatus
	Remark    string
	PayTime   int64
	PayWay    string // 支付方式
	DeletedAt gorm.DeletedAt
}
