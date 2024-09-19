package vo

import (
	"geekai/core/types"
)

type Order struct {
	BaseVo
	UserId    uint              `json:"user_id"`
	ProductId uint              `json:"product_id"`
	Username  string            `json:"username"`
	OrderNo   string            `json:"order_no"`
	TradeNo   string            `json:"trade_no"`
	Subject   string            `json:"subject"`
	Amount    float64           `json:"amount"`
	Status    types.OrderStatus `json:"status"`
	PayTime   int64             `json:"pay_time"`
	PayWay    string            `json:"pay_way"`
	PayType   string            `json:"pay_type"`
	PayMethod string            `json:"pay_method"`
	PayName   string            `json:"pay_name"`
	Remark    types.OrderRemark `json:"remark"`
}
