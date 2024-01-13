package alipay

const (
	NotifyTypeTradeStatusSync = "trade_status_sync"
)

// TradeNotification
// Deprecated: use Notification instead.
type TradeNotification Notification

// Notification 通知响应参数 https://doc.open.alipay.com/docs/doc.htm?spm=a219a.7629140.0.0.8AmJwg&treeId=203&articleId=105286&docType=1
type Notification struct {
	AuthAppId           string      `json:"auth_app_id"`           // App Id
	NotifyTime          string      `json:"notify_time"`           // 通知时间
	NotifyType          string      `json:"notify_type"`           // 通知类型
	NotifyId            string      `json:"notify_id"`             // 通知校验ID
	AppId               string      `json:"app_id"`                // 开发者的app_id
	Charset             string      `json:"charset"`               // 编码格式
	Version             string      `json:"version"`               // 接口版本
	SignType            string      `json:"sign_type"`             // 签名类型
	Sign                string      `json:"sign"`                  // 签名
	TradeNo             string      `json:"trade_no"`              // 支付宝交易号
	OutTradeNo          string      `json:"out_trade_no"`          // 商户订单号
	OutBizNo            string      `json:"out_biz_no"`            // 商户业务号
	BuyerId             string      `json:"buyer_id"`              // 买家支付宝用户号
	BuyerLogonId        string      `json:"buyer_logon_id"`        // 买家支付宝账号
	SellerId            string      `json:"seller_id"`             // 卖家支付宝用户号
	SellerEmail         string      `json:"seller_email"`          // 卖家支付宝账号
	TradeStatus         TradeStatus `json:"trade_status"`          // 交易状态
	TotalAmount         string      `json:"total_amount"`          // 订单金额
	ReceiptAmount       string      `json:"receipt_amount"`        // 实收金额
	InvoiceAmount       string      `json:"invoice_amount"`        // 开票金额
	BuyerPayAmount      string      `json:"buyer_pay_amount"`      // 付款金额
	PointAmount         string      `json:"point_amount"`          // 集分宝金额
	RefundFee           string      `json:"refund_fee"`            // 总退款金额
	Subject             string      `json:"subject"`               // 商品的标题/交易标题/订单标题/订单关键字等，是请求时对应的参数，原样通知回来。
	Body                string      `json:"body"`                  // 商品描述
	GmtCreate           string      `json:"gmt_create"`            // 交易创建时间
	GmtPayment          string      `json:"gmt_payment"`           // 交易付款时间
	GmtRefund           string      `json:"gmt_refund"`            // 交易退款时间
	GmtClose            string      `json:"gmt_close"`             // 交易结束时间
	FundBillList        string      `json:"fund_bill_list"`        // 支付金额信息
	PassbackParams      string      `json:"passback_params"`       // 回传参数
	VoucherDetailList   string      `json:"voucher_detail_list"`   // 优惠券信息
	AgreementNo         string      `json:"agreement_no"`          // 支付宝签约号
	ExternalAgreementNo string      `json:"external_agreement_no"` // 商户自定义签约号
}
