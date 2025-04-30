package model

import (
	"geekai/core/types"
	"time"
)

// Order 充值订单
type Order struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户ID" json:"user_id"`
	ProductId uint       `gorm:"column:product_id;type:int;not null;comment:产品ID" json:"product_id"`
	Username  string    `gorm:"column:username;type:varchar(30);not null;comment:用户名" json:"username"`
	OrderNo   string    `gorm:"column:order_no;type:varchar(30);uniqueIndex;not null;comment:订单ID" json:"order_no"`
	TradeNo   string    `gorm:"column:trade_no;type:varchar(60);comment:支付平台交易流水号" json:"trade_no"`
	Subject   string    `gorm:"column:subject;type:varchar(100);not null;comment:订单产品" json:"subject"`
	Amount    float64     `gorm:"column:amount;type:decimal(10,2);not null;default:0.00;comment:订单金额" json:"amount"`
	Status    types.OrderStatus         `gorm:"column:status;type:tinyint(1);not null;default:0;comment:订单状态（0：待支付，1：已扫码，2：支付成功）" json:"status"`
	Remark    string      `gorm:"column:remark;type:varchar(255);not null;comment:备注" json:"remark"`
	PayTime   int64       `gorm:"column:pay_time;type:int;comment:支付时间" json:"pay_time"`
	PayWay    string      `gorm:"column:pay_way;type:varchar(20);not null;comment:支付方式" json:"pay_way"`
	PayType   string      `gorm:"column:pay_type;type:varchar(30);not null;comment:支付类型" json:"pay_type"`
	CreatedAt time.Time   `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time   `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *Order) TableName() string {
	return "chatgpt_orders"
}
