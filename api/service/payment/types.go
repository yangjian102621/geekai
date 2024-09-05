package payment

type NotifyVo struct {
	Status     int
	OutTradeNo string // 商户订单号
	TradeId    string // 交易ID
	Amount     string // 交易金额
	Message    string
	Subject    string
}

func (v NotifyVo) Success() bool {
	return v.Status == Success
}

const (
	Success = 0
	Failure = 1
)
