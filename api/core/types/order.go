package types

type OrderStatus int

const (
	OrderNotPaid     = OrderStatus(0)
	OrderScanned     = OrderStatus(1) // 已扫码
	OrderPaidSuccess = OrderStatus(2)
)

type OrderRemark struct {
	Days     int     `json:"days"`      // 有效期
	Calls    int     `json:"calls"`     // 增加对话次数
	ImgCalls int     `json:"img_calls"` // 增加绘图次数
	Name     string  `json:"name"`      // 产品名称
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
}
