package payment

// 支付渠道定义
const PayChannelAL = "alipay" // 支付宝
const PayChannelWX = "wxpay"  // 微信支付
const PayChannelEpay = "epay" // 易支付

// 支付方式
const PayWayAL = "alipay"
const PayWayWX = "wxpay"

const (
	Success = 0
	Failure = 1
	Closed  = 2
)

type PayRequest struct {
	OutTradeNo string // 商户订单号
	Subject    string // 商品名称
	TotalFee   string // 商品金额
	ReturnURL  string // 回调地址
	NotifyURL  string // 回调地址

	// 易支付专有参数
	Method   string // 接口类型
	Device   string // 设备类型
	PayWay   string // 支付方式
	ClientIP string //用户IP地址
	OpenID   string // 用户openid

}

type OrderInfo struct {
	Mchid      string // 商户号
	OutTradeNo string // 商户订单号
	TradeId    string // 交易号
	Amount     string // 金额
	Status     int    // 状态 0: 未支付 1: 已支付 2: 已关闭
	PayTime    string // 完成支付时间
}

func (o OrderInfo) Closed() bool {
	return o.Status == Closed
}

func (o OrderInfo) Success() bool {
	return o.Status == Success
}

type PayService interface {
	Pay(params PayRequest) (string, error)      // 生成支付链接
	Query(outTradeNo string) (OrderInfo, error) // 查询订单
}
