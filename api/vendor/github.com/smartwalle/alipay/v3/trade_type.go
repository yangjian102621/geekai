package alipay

import "encoding/json"

type Trade struct {
	AuxParam
	NotifyURL    string `json:"-"`
	ReturnURL    string `json:"-"`
	AppAuthToken string `json:"-"` // 可选

	// biz content，这四个参数是必须的
	Subject     string `json:"subject"`      // 订单标题
	OutTradeNo  string `json:"out_trade_no"` // 商户订单号，64个字符以内、可包含字母、数字、下划线；需保证在商户端不重复
	TotalAmount string `json:"total_amount"` // 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	ProductCode string `json:"product_code"` // 销售产品码，与支付宝签约的产品码名称。 参考官方文档, App 支付时默认值为 QUICK_MSECURITY_PAY

	Body               string          `json:"body,omitempty"`                 // 订单描述
	GoodsDetail        []*GoodsDetail  `json:"goods_detail,omitempty"`         // 可选 订单包含的商品列表信息，Json格式，详见商品明细说明
	BusinessParams     json.RawMessage `json:"business_params,omitempty"`      // 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	DisablePayChannels string          `json:"disable_pay_channels,omitempty"` // 禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔 注，与enable_pay_channels互斥
	EnablePayChannels  string          `json:"enable_pay_channels,omitempty"`  // 可用渠道，用户只能在指定渠道范围内支付  当有多个渠道时用“,”分隔 注，与disable_pay_channels互斥
	SpecifiedChannel   string          `json:"specified_channel,omitempty"`    // 指定渠道，目前仅支持传入pcredit  若由于用户原因渠道不可用，用户可选择是否用其他渠道支付。  注：该参数不可与花呗分期参数同时传入
	//ExtUserInfo        string `json:"ext_user_info,omitempty"`        // 外部指定买家
	ExtendParams        *ExtendParams `json:"extend_params,omitempty"`         // 可选 业务扩展参数，详见下面的“业务扩展参数说明”
	AgreementSignParams *SignParams   `json:"agreement_sign_params,omitempty"` // 签约参数。如果希望在sdk中支付并签约，需要在这里传入签约信息。 周期扣款场景 product_code 为 CYCLE_PAY_AUTH 时必填。
	GoodsType           string        `json:"goods_type,omitempty"`            // 商品主类型：0—虚拟类商品，1—实物类商品 注：虚拟类商品不支持使用花呗渠道
	InvoiceInfo         string        `json:"invoice_info,omitempty"`          // 开票信息
	PassbackParams      string        `json:"passback_params,omitempty"`       // 公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	PromoParams         string        `json:"promo_params,omitempty"`          // 优惠参数 注：仅与支付宝协商后可用
	RoyaltyInfo         string        `json:"royalty_info,omitempty"`          // 描述分账信息，json格式，详见分账参数说明
	SellerId            string        `json:"seller_id,omitempty"`             // 收款支付宝用户ID。 如果该值为空，则默认为商户签约账号对应的支付宝用户ID
	SettleInfo          string        `json:"settle_info,omitempty"`           // 描述结算信息，json格式，详见结算参数说明
	StoreId             string        `json:"store_id,omitempty"`              // 商户门店编号。该参数用于请求参数中以区分各门店，非必传项。
	SubMerchant         string        `json:"sub_merchant,omitempty"`          // 间连受理商户信息体，当前只对特殊银行机构特定场景下使用此字段
	TimeoutExpress      string        `json:"timeout_express,omitempty"`       // 该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	TimeExpire          string        `json:"time_expire,omitempty"`           // 该笔订单绝对超时时间，格式为yyyy-MM-dd HH:mm:ss
	MerchantOrderNo     string        `json:"merchant_order_no,omitempty"`     // 可选 商户的原始订单号
	ExtUserInfo         *ExtUserInfo  `json:"ext_user_info,omitempty"`         // 可选 外部指定买家
	QueryOptions        []string      `json:"query_options,omitempty"`         // 可选 通知参数选项。 商户通过传递该参数来定制需要异步通知的额外字段，数组格式。包括但不限于：["hyb_amount","enterprise_pay_info"]
}

type SignParams struct {
	PersonalProductCode string             `json:"personal_product_code"`        // 必选 个人签约产品码，商户和支付宝签约时确定。
	SignScene           string             `json:"sign_scene"`                   // 必选 协议签约场景，商户和支付宝签约时确定，商户可咨询技术支持。
	ExternalAgreementNo string             `json:"external_agreement_no"`        // 可选 商户签约号，代扣协议中标示用户的唯一签约号（确保在商户系统中唯一）。 格式规则：支持大写小写字母和数字，最长32位。 商户系统按需传入，如果同一用户在同一产品码、同一签约场景下，签订了多份代扣协议，那么需要指定并传入该值。
	ExternalLogonId     string             `json:"external_logon_id"`            // 可选 用户在商户网站的登录账号，用于在签约页面展示，如果为空，则不展示
	AccessParams        *AccessParams      `json:"access_params"`                // 必选 请按当前接入的方式进行填充，且输入值必须为文档中的参数取值范围。
	SubMerchant         *SubMerchantParams `json:"sub_merchant,omitempty"`       // 可选 此参数用于传递子商户信息，无特殊需求时不用关注。目前商户代扣、海外代扣、淘旅行信用住产品支持传入该参数（在销售方案中“是否允许自定义子商户信息”需要选是）。
	PeriodRuleParams    *PeriodRuleParams  `json:"period_rule_params,omitempty"` // 可选 周期管控规则参数period_rule_params，在签约周期扣款产品（如CYCLE_PAY_AUTH_P）时必传，在签约其他产品时无需传入。 周期扣款产品，会按照这里传入的参数提示用户，并对发起扣款的时间、金额、次数等做相应限制。
	SignNotifyURL       string             `json:"sign_notify_url"`              // 可选 签约成功后商户用于接收异步通知的地址。如果不传入，签约与支付的异步通知都会发到外层notify_url参数传入的地址；如果外层也未传入，签约与支付的异步通知都会发到商户appid配置的网关地址。
}

type ExtUserInfo struct {
	Name          string `json:"name"`            //  可选 指定买家姓名。 注： need_check_info=T时该参数才有效
	Mobile        string `json:"mobile"`          //  可选 指定买家手机号。 注：该参数暂不校验
	CertType      string `json:"cert_type"`       //  可选 指定买家证件类型。 枚举值：IDENTITY_CARD：身份证；PASSPORT：护照；OFFICER_CARD：军官证；SOLDIER_CARD：士兵证；HOKOU：户口本。如有其它类型需要支持，请与蚂蚁金服工作人员联系。注： need_check_info=T时该参数才有效，支付宝会比较买家在支付宝留存的证件类型与该参数传入的值是否匹配。
	CertNo        string `json:"cert_no"`         //  可选 买家证件号。 注：need_check_info=T时该参数才有效，支付宝会比较买家在支付宝留存的证件号码与该参数传入的值是否匹配。
	MinAge        string `json:"min_age"`         //  可选 允许的最小买家年龄。 买家年龄必须大于等于所传数值注：1. need_check_info=T时该参数才有效  2. min_age为整数，必须大于等于0
	NeedCheckInfo string `json:"need_check_info"` //  可选 是否强制校验买家信息； 需要强制校验传：T;不需要强制校验传：F或者不传；当传T时，支付宝会校验支付买家的信息与接口上传递的cert_type、cert_no、name或age是否匹配，只有接口传递了信息才会进行对应项的校验；只要有任何一项信息校验不匹配交易都会失败。如果传递了need_check_info，但是没有传任何校验项，则不进行任何校验。默认为不校验。
	IdentityHash  string `json:"identity_hash"`   //  可选 买家加密身份信息。当指定了此参数且指定need_check_info=T时，支付宝会对买家身份进行校验，校验逻辑为买家姓名、买家证件号拼接后的字符串，以sha256算法utf-8编码计算hash，若与传入的值不匹配则会拦截本次支付。注意：如果同时指定了用户明文身份信息（name，cert_type，cert_no中任意一个），则忽略identity_hash以明文参数校验。
}

// TradePagePay 统一收单下单并支付页面接口请求参数 https://opendocs.alipay.com/apis/api_1/alipay.trade.page.pay
type TradePagePay struct {
	Trade
	AuthToken   string `json:"auth_token,omitempty"`   // 针对用户授权接口，获取用户相关数据时，用于标识用户授权关系
	QRPayMode   string `json:"qr_pay_mode,omitempty"`  // PC扫码支付的方式，支持前置模式和跳转模式。
	QRCodeWidth string `json:"qrcode_width,omitempty"` // 商户自定义二维码宽度 注：qr_pay_mode=4时该参数生效
}

type GoodsDetail struct {
	GoodsId        string  `json:"goods_id"`
	AliPayGoodsId  string  `json:"alipay_goods_id,omitempty"`
	GoodsName      string  `json:"goods_name"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
	GoodsCategory  string  `json:"goods_category,omitempty"`
	CategoriesTree string  `json:"categories_tree,omitempty"`
	Body           string  `json:"body,omitempty"`
	ShowURL        string  `json:"show_url,omitempty"`
}

func (this TradePagePay) APIName() string {
	return "alipay.trade.page.pay"
}

func (this TradePagePay) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	m["return_url"] = this.ReturnURL
	return m
}

type TradeStatus string

const (
	TradeStatusWaitBuyerPay TradeStatus = "WAIT_BUYER_PAY" //（交易创建，等待买家付款）
	TradeStatusClosed       TradeStatus = "TRADE_CLOSED"   //（未付款交易超时关闭，或支付完成后全额退款）
	TradeStatusSuccess      TradeStatus = "TRADE_SUCCESS"  //（交易支付成功）
	TradeStatusFinished     TradeStatus = "TRADE_FINISHED" //（交易结束，不可退款）
)

// TradeQuery 统一收单线下交易查询接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.query/
type TradeQuery struct {
	AuxParam
	AppAuthToken string   `json:"-"`                       // 可选
	OutTradeNo   string   `json:"out_trade_no,omitempty"`  // 订单支付时传入的商户订单号, 与 TradeNo 二选一
	TradeNo      string   `json:"trade_no,omitempty"`      // 支付宝交易号
	OrgPid       string   `json:"org_pid,omitempty"`       // 可选 银行间联模式下有用，其它场景请不要使用； 双联通过该参数指定需要查询的交易所属收单机构的pid;
	QueryOptions []string `json:"query_options,omitempty"` // 可选 查询选项，商户传入该参数可定制本接口同步响应额外返回的信息字段，数组格式。支持枚举如下：trade_settle_info：返回的交易结算信息，包含分账、补差等信息； fund_bill_list：交易支付使用的资金渠道；voucher_detail_list：交易支付时使用的所有优惠券信息；discount_goods_detail：交易支付所使用的单品券优惠的商品优惠信息；mdiscount_amount：商家优惠金额；
}

func (this TradeQuery) APIName() string {
	return "alipay.trade.query"
}

func (this TradeQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// TradeQueryRsp 统一收单线下交易查询接口响应参数
type TradeQueryRsp struct {
	Error
	TradeNo               string           `json:"trade_no"`                      // 支付宝交易号
	OutTradeNo            string           `json:"out_trade_no"`                  // 商家订单号
	BuyerLogonId          string           `json:"buyer_logon_id"`                // 买家支付宝账号
	TradeStatus           TradeStatus      `json:"trade_status"`                  // 交易状态
	TotalAmount           string           `json:"total_amount"`                  // 交易的订单金额
	TransCurrency         string           `json:"trans_currency"`                // 标价币种
	SettleCurrency        string           `json:"settle_currency"`               // 订单结算币种
	SettleAmount          string           `json:"settle_amount"`                 // 结算币种订单金额
	PayCurrency           string           `json:"pay_currency"`                  // 订单支付币种
	PayAmount             string           `json:"pay_amount"`                    // 支付币种订单金额
	SettleTransRate       string           `json:"settle_trans_rate"`             // 结算币种兑换标价币种汇率
	TransPayRate          string           `json:"trans_pay_rate"`                // 标价币种兑换支付币种汇率
	BuyerPayAmount        string           `json:"buyer_pay_amount"`              // 买家实付金额，单位为元，两位小数。
	PointAmount           string           `json:"point_amount"`                  // 积分支付的金额，单位为元，两位小数。
	InvoiceAmount         string           `json:"invoice_amount"`                // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
	SendPayDate           string           `json:"send_pay_date"`                 // 本次交易打款给卖家的时间
	ReceiptAmount         string           `json:"receipt_amount"`                // 实收金额，单位为元，两位小数
	StoreId               string           `json:"store_id"`                      // 商户门店编号
	TerminalId            string           `json:"terminal_id"`                   // 商户机具终端编号
	FundBillList          []*FundBill      `json:"fund_bill_list,omitempty"`      // 交易支付使用的资金渠道
	StoreName             string           `json:"store_name"`                    // 请求交易支付中的商户店铺的名称
	BuyerUserId           string           `json:"buyer_user_id"`                 // 买家在支付宝的用户id
	BuyerUserName         string           `json:"buyer_user_name"`               // 买家名称；
	IndustrySepcDetailGov string           `json:"industry_sepc_detail_gov"`      // 行业特殊信息-统筹相关
	IndustrySepcDetailAcc string           `json:"industry_sepc_detail_acc"`      // 行业特殊信息-个账相关
	ChargeAmount          string           `json:"charge_amount"`                 // 该笔交易针对收款方的收费金额；
	ChargeFlags           string           `json:"charge_flags"`                  // 费率活动标识，当交易享受活动优惠费率时，返回该活动的标识；
	SettlementId          string           `json:"settlement_id"`                 // 支付清算编号，用于清算对账使用；
	TradeSettleInfo       *TradeSettleInfo `json:"trade_settle_info,omitempty"`   // 返回的交易结算信息，包含分账、补差等信息
	AuthTradePayMode      string           `json:"auth_trade_pay_mode"`           // 预授权支付模式，该参数仅在信用预授权支付场景下返回。信用预授权支付：CREDIT_PREAUTH_PAY
	BuyerUserType         string           `json:"buyer_user_type"`               // 买家用户类型。CORPORATE:企业用户；PRIVATE:个人用户。
	MdiscountAmount       string           `json:"mdiscount_amount"`              // 商家优惠金额
	DiscountAmount        string           `json:"discount_amount"`               // 平台优惠金额
	Subject               string           `json:"subject"`                       // 订单标题；
	Body                  string           `json:"body"`                          // 订单描述;
	AlipaySubMerchantId   string           `json:"alipay_sub_merchant_id"`        // 间连商户在支付宝端的商户编号；
	ExtInfos              string           `json:"ext_infos"`                     // 交易额外信息，特殊场景下与支付宝约定返回。
	PassbackParams        string           `json:"passback_params"`               // 公用回传参数。返回支付时传入的passback_params参数信息
	HBFQPayInfo           *HBFQPayInfo     `json:"hb_fq_pay_info"`                // 若用户使用花呗分期支付，且商家开通返回此通知参数，则会返回花呗分期信息。json格式其它说明详见花呗分期信息说明。 注意：商家需与支付宝约定后才返回本参数。
	CreditPayMode         string           `json:"credit_pay_mode"`               // 信用支付模式。表示订单是采用信用支付方式（支付时买家没有出资，需要后续履约）。"creditAdvanceV2"表示芝麻先用后付模式，用户后续需要履约扣款。 此字段只有信用支付场景才有值，商户需要根据字段值单独处理。此字段以后可能扩展其他值，建议商户使用白名单方式识别，对于未识别的值做失败处理，并联系支付宝技术支持人员。
	CreditBizOrderId      string           `json:"credit_biz_order_id"`           // 信用业务单号。信用支付场景才有值，先用后付产品里是芝麻订单号。
	HYBAmount             string           `json:"hyb_amount"`                    // 惠营宝回票金额
	BKAgentRespInfo       *BKAgentRespInfo `json:"bk_agent_resp_info"`            // 间联交易下，返回给机构的信 息
	ChargeInfoList        []*ChargeInfo    `json:"charge_info_list"`              // 计费信息列表
	DiscountGoodsDetail   string           `json:"discount_goods_detail"`         // 本次交易支付所使用的单品券优惠的商品优惠信息
	VoucherDetailList     []*VoucherDetail `json:"voucher_detail_list,omitempty"` // 本交易支付时使用的所有优惠券信息
}

type HBFQPayInfo struct {
	UserInstallNum string `json:"user_install_num"` // 用户使用花呗分期支付的分期数
}

type BKAgentRespInfo struct {
	BindtrxId        string `json:"bindtrx_id"`
	BindclrissrId    string `json:"bindclrissr_id"`
	BindpyeracctbkId string `json:"bindpyeracctbk_id"`
	BkpyeruserCode   string `json:"bkpyeruser_code"`
	EstterLocation   string `json:"estter_location"`
}

type ChargeInfo struct {
	ChargeFee               string          `json:"charge_fee"`
	OriginalChargeFee       string          `json:"original_charge_fee"`
	SwitchFeeRate           string          `json:"switch_fee_rate"`
	IsRatingOnTradeReceiver string          `json:"is_rating_on_trade_receiver"`
	IsRatingOnSwitch        string          `json:"is_rating_on_switch"`
	ChargeType              string          `json:"charge_type"`
	SubFeeDetailList        []*SubFeeDetail `json:"sub_fee_detail_list"`
}

type SubFeeDetail struct {
	ChargeFee         string `json:"charge_fee"`
	OriginalChargeFee string `json:"original_charge_fee"`
	SwitchFeeRate     string `json:"switch_fee_rate"`
}

type FundBill struct {
	FundChannel string  `json:"fund_channel"`       // 交易使用的资金渠道，详见 支付渠道列表
	BankCode    string  `json:"bank_code"`          // 银行卡支付时的银行代码
	Amount      string  `json:"amount"`             // 该支付工具类型所使用的金额
	RealAmount  float64 `json:"real_amount,string"` // 渠道实际付款金额
}

type VoucherDetail struct {
	Id                 string `json:"id"`                  // 券id
	Name               string `json:"name"`                // 券名称
	Type               string `json:"type"`                // 当前有三种类型： ALIPAY_FIX_VOUCHER - 全场代金券, ALIPAY_DISCOUNT_VOUCHER - 折扣券, ALIPAY_ITEM_VOUCHER - 单品优惠
	Amount             string `json:"amount"`              // 优惠券面额，它应该会等于商家出资加上其他出资方出资
	MerchantContribute string `json:"merchant_contribute"` // 商家出资（特指发起交易的商家出资金额）
	OtherContribute    string `json:"other_contribute"`    // 其他出资方出资金额，可能是支付宝，可能是品牌商，或者其他方，也可能是他们的一起出资
	Memo               string `json:"memo"`                // 优惠券备注信息
}

type TradeSettleInfo struct {
	TradeSettleDetailList []*TradeSettleDetail `json:"trade_settle_detail_list"`
}

type TradeSettleDetail struct {
	OperationType     string `json:"operation_type"`
	OperationSerialNo string `json:"operation_serial_no"`
	OperationDate     string `json:"operation_dt"`
	TransOut          string `json:"trans_out"`
	TransIn           string `json:"trans_in"`
	Amount            string `json:"amount"`
}

// TradeClose 统一收单交易关闭接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.close/
type TradeClose struct {
	AuxParam
	AppAuthToken string `json:"-"`                      // 可选
	NotifyURL    string `json:"-"`                      // 可选
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 与 TradeNo 二选一
	TradeNo      string `json:"trade_no,omitempty"`     // 与 OutTradeNo 二选一
	OperatorId   string `json:"operator_id,omitempty"`  // 可选
}

func (this TradeClose) APIName() string {
	return "alipay.trade.close"
}

func (this TradeClose) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// TradeCloseRsp 统一收单交易关闭接口响应参数
type TradeCloseRsp struct {
	Error
	OutTradeNo string `json:"out_trade_no"`
	TradeNo    string `json:"trade_no"`
}

// TradeRefund 统一收单交易退款接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.refund/
type TradeRefund struct {
	AuxParam
	AppAuthToken            string                    `json:"-"`                                   // 可选
	OutTradeNo              string                    `json:"out_trade_no,omitempty"`              // 与 TradeNo 二选一
	TradeNo                 string                    `json:"trade_no,omitempty"`                  // 与 OutTradeNo 二选一
	RefundAmount            string                    `json:"refund_amount"`                       // 必须 需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundReason            string                    `json:"refund_reason"`                       // 可选 退款的原因说明
	OutRequestNo            string                    `json:"out_request_no"`                      // 必须 标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	RefundRoyaltyParameters []*RefundRoyaltyParameter `json:"refund_royalty_parameters,omitempty"` // 可选 退分账明细信息。 注： 1.当面付且非直付通模式无需传入退分账明细，系统自动按退款金额与订单金额的比率，从收款方和分账收入方退款，不支持指定退款金额与退款方。2.直付通模式，电脑网站支付，手机 APP 支付，手机网站支付产品，须在退款请求中明确是否退分账，从哪个分账收入方退，退多少分账金额；如不明确，默认从收款方退款，收款方余额不足退款失败。不支持系统按比率退款。
	QueryOptions            []string                  `json:"query_options,omitempty"`             // 可选 查询选项。 商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。支持：refund_detail_item_list：退款使用的资金渠道；deposit_back_info：触发银行卡冲退信息通知；
	// OperatorId   string `json:"operator_id"`            // 可选 商户的操作员编号
	// StoreId    string `json:"store_id"`    // 可选 商户的门店编号
	// TerminalId string `json:"terminal_id"` // 可选 商户的终端编号
}

func (this TradeRefund) APIName() string {
	return "alipay.trade.refund"
}

func (this TradeRefund) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

type RefundRoyaltyParameter struct {
	RoyaltyType  string `json:"royalty_type,omitempty"`   // 可选 分账类型. 普通分账为：transfer;补差为：replenish;为空默认为分账transfer;
	TransOut     string `json:"trans_out,omitempty"`      // 可选 支出方账户。如果支出方账户类型为userId，本参数为支出方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果支出方类型为loginName，本参数为支出方的支付宝登录号。 泛金融类商户分账时，该字段不要上送。
	TransOutType string `json:"trans_out_type,omitempty"` // 可选 支出方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;loginName表示是支付宝登录号； 泛金融类商户分账时，该字段不要上送。
	TransInType  string `json:"trans_in_type,omitempty"`  // 可选 收入方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;cardAliasNo表示是卡编号;loginName表示是支付宝登录号；
	TransIn      string `json:"trans_in"`                 // 必选 收入方账户。如果收入方账户类型为userId，本参数为收入方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果收入方类型为cardAliasNo，本参数为收入方在支付宝绑定的卡编号；如果收入方类型为loginName，本参数为收入方的支付宝登录号；
	Amount       string `json:"amount,omitempty"`         // 可选 分账的金额，单位为元
	Desc         string `json:"desc,omitempty"`           // 可选 分账描述
	RoyaltyScene string `json:"royalty_scene,omitempty"`  // 可选 可选值：达人佣金、平台服务费、技术服务费、其他
	TransInName  string `json:"trans_in_name,omitempty"`  // 可选 分账收款方姓名，上送则进行姓名与支付宝账号的一致性校验，校验不一致则分账失败。不上送则不进行姓名校验
}

// TradeRefundRsp 统一收单交易退款接口响应参数
type TradeRefundRsp struct {
	Error
	TradeNo              string              `json:"trade_no"`                          // 支付宝交易号
	OutTradeNo           string              `json:"out_trade_no"`                      // 商户订单号
	BuyerLogonId         string              `json:"buyer_logon_id"`                    // 用户的登录id
	BuyerUserId          string              `json:"buyer_user_id"`                     // 买家在支付宝的用户id
	FundChange           string              `json:"fund_change"`                       // 本次退款是否发生了资金变化
	RefundFee            string              `json:"refund_fee"`                        // 退款总金额
	StoreName            string              `json:"store_name"`                        // 交易在支付时候的门店名称
	RefundDetailItemList []*TradeFundBill    `json:"refund_detail_item_list,omitempty"` // 退款使用的资金渠道
	SendBackFee          string              `json:"send_back_fee"`                     // 本次商户实际退回金额。 说明：如需获取该值，需在入参query_options中传入 refund_detail_item_list。
	RefundHYBAmount      string              `json:"refund_hyb_amount"`                 // 本次请求退惠营宝金额
	RefundChargeInfoList []*RefundChargeInfo `json:"refund_charge_info_list,omitempty"` // 退费信息
}

type TradeFundBill struct {
	FundChannel string `json:"fund_channel"` // 交易使用的资金渠道，详见 支付渠道列表
	Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
	FundType    string `json:"fund_type"`    // 渠道所使用的资金类型
}

type RefundChargeInfo struct {
	RefundChargeFee        string                `json:"refund_charge_fee"`                    // 实退费用
	SwitchFeeRate          string                `json:"switch_fee_rate"`                      // 签约费率
	ChargeType             string                `json:"charge_type"`                          // 收单手续费trade，花呗分期手续hbfq，其他手续费charge
	RefundSubFeeDetailList []*RefundSubFeeDetail `json:"refund_sub_fee_detail_list,omitempty"` // 组合支付退费明细
}

type RefundSubFeeDetail struct {
	RefundChargeFee string `json:"refund_charge_fee"` // 实退费用
	SwitchFeeRate   string `json:"switch_fee_rate"`   // 签约费率
}

// TradeFastPayRefundQuery 统一收单交易退款查询接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query
type TradeFastPayRefundQuery struct {
	AuxParam
	AppAuthToken string   `json:"-"`                       // 可选
	OutTradeNo   string   `json:"out_trade_no,omitempty"`  // 与 TradeNo 二选一
	TradeNo      string   `json:"trade_no,omitempty"`      // 与 OutTradeNo 二选一
	OutRequestNo string   `json:"out_request_no"`          // 必须 请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号
	QueryOptions []string `json:"query_options,omitempty"` // 可选 查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。 refund_detail_item_list
}

func (this TradeFastPayRefundQuery) APIName() string {
	return "alipay.trade.fastpay.refund.query"
}

func (this TradeFastPayRefundQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// TradeFastPayRefundQueryRsp 统一收单交易退款查询接口响应参数
type TradeFastPayRefundQueryRsp struct {
	Error
	TradeNo              string              `json:"trade_no"`                          // 支付宝交易号
	OutTradeNo           string              `json:"out_trade_no"`                      // 创建交易传入的商户订单号
	OutRequestNo         string              `json:"out_request_no"`                    // 本笔退款对应的退款请求号
	TotalAmount          string              `json:"total_amount"`                      // 发该笔退款所对应的交易的订单金额
	RefundAmount         string              `json:"refund_amount"`                     // 本次退款请求，对应的退款金额
	RefundStatus         string              `json:"refund_status"`                     // 退款状态。枚举值： REFUND_SUCCESS 退款处理成功； 未返回该字段表示退款请求未收到或者退款失败；
	RefundRoyaltys       []*RefundRoyalty    `json:"refund_royaltys"`                   // 退分账明细信息
	GMTRefundPay         string              `json:"gmt_refund_pay"`                    // 退款时间。
	RefundDetailItemList []*TradeFundBill    `json:"refund_detail_item_list"`           // 本次退款使用的资金渠道； 默认不返回该信息，需要在入参的query_options中指定"refund_detail_item_list"值时才返回该字段信息。
	SendBackFee          string              `json:"send_back_fee"`                     // 本次商户实际退回金额；
	DepositBackInfo      []*DepositBackInfo  `json:"deposit_back_info"`                 // 银行卡冲退信息
	RefundHYBAmount      string              `json:"refund_hyb_amount"`                 // 本次请求退惠营宝金额
	RefundChargeInfoList []*RefundChargeInfo `json:"refund_charge_info_list,omitempty"` // 退费信息
}

type RefundRoyalty struct {
	RefundAmount  string `json:"refund_amount"`
	RoyaltyType   string `json:"royalty_type"`
	ResultCode    string `json:"result_code"`
	TransOut      string `json:"trans_out"`
	TransOutEmail string `json:"trans_out_email"`
	TransIn       string `json:"trans_in"`
	TransInEmail  string `json:"trans_in_email"`
}

type DepositBackInfo struct {
	HasDepositBack     string `json:"has_deposit_back"`
	DBackStatus        string `json:"dback_status"`
	DBackAmount        string `json:"dback_amount"`
	BankAckTime        string `json:"bank_ack_time"`
	ESTBankReceiptTime string `json:"est_bank_receipt_time"`
}

// TradeOrderSettle 统一收单交易结算接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.order.settle
type TradeOrderSettle struct {
	AuxParam
	AppAuthToken      string              `json:"-"`                  // 可选
	OutRequestNo      string              `json:"out_request_no"`     // 必须 结算请求流水号 开发者自行生成并保证唯一性
	TradeNo           string              `json:"trade_no"`           // 必须 支付宝订单号
	RoyaltyParameters []*RoyaltyParameter `json:"royalty_parameters"` // 必须 分账明细信息
	OperatorId        string              `json:"operator_id"`        //可选 操作员id
}

func (this TradeOrderSettle) APIName() string {
	return "alipay.trade.order.settle"
}

func (this TradeOrderSettle) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

type RoyaltyParameter struct {
	TransOut         string  `json:"trans_out"`                   // 可选 分账支出方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	TransIn          string  `json:"trans_in"`                    // 可选 分账收入方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	Amount           float64 `json:"amount"`                      // 可选 分账的金额，单位为元
	AmountPercentage float64 `json:"amount_percentage,omitempty"` // 可选 分账信息中分账百分比。取值范围为大于0，少于或等于100的整数。
	Desc             string  `json:"desc"`                        // 可选 分账描述
}

// TradeOrderSettleRsp 统一收单交易结算接口响应参数
type TradeOrderSettleRsp struct {
	Error
	TradeNo string `json:"trade_no"`
}

// TradeCreate 统一收单交易创建接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.create/
type TradeCreate struct {
	Trade
	DiscountableAmount string             `json:"discountable_amount"` // 可打折金额. 参与优惠计算的金额，单位为元，精确到小数点后两位
	BuyerId            string             `json:"buyer_id"`
	BuyerOpenId        string             `json:"buyer_open_id"` // 新版接口无法获取user_id, 这里只能传递openid值
	GoodsDetail        []*GoodsDetailItem `json:"goods_detail,omitempty"`
	OperatorId         string             `json:"operator_id"`
	TerminalId         string             `json:"terminal_id"`
}

func (this TradeCreate) APIName() string {
	return "alipay.trade.create"
}

func (this TradeCreate) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// TradeCreateRsp 统一收单交易创建接口响应参数
type TradeCreateRsp struct {
	Error
	TradeNo    string `json:"trade_no"` // 支付宝交易号
	OutTradeNo string `json:"out_trade_no"`
}

type RoyaltyInfo struct {
	RoyaltyType       string                   `json:"royalty_type"`
	RoyaltyDetailInfo []*RoyaltyDetailInfoItem `json:"royalty_detail_infos,omitempty"`
}

type RoyaltyDetailInfoItem struct {
	SerialNo         string `json:"serial_no"`
	TransInType      string `json:"trans_in_type"`
	BatchNo          string `json:"batch_no"`
	OutRelationId    string `json:"out_relation_id"`
	TransOutType     string `json:"trans_out_type"`
	TransOut         string `json:"trans_out"`
	TransIn          string `json:"trans_in"`
	Amount           string `json:"amount"`
	Desc             string `json:"desc"`
	AmountPercentage string `json:"amount_percentage"`
	AliPayStoreId    string `json:"alipay_store_id"`
}

type SubMerchantItem struct {
	MerchantId string `json:"merchant_id"`
}

type GoodsDetailItem struct {
	GoodsId       string `json:"goods_id"`
	AliPayGoodsId string `json:"alipay_goods_id"`
	GoodsName     string `json:"goods_name"`
	Quantity      string `json:"quantity"`
	Price         string `json:"price"`
	GoodsCategory string `json:"goods_category"`
	Body          string `json:"body"`
	ShowUrl       string `json:"show_url"`
}

type AgreementParams struct {
	AgreementNo   string `json:"agreement_no,omitempty"`
	AuthConfirmNo string `json:"auth_confirm_no,omitempty"`
	ApplyToken    string `json:"apply_token,omitempty"`
}

// TradePay 统一收单交易支付接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.pay/
type TradePay struct {
	Trade
	Scene    string `json:"scene"`               // 必须 支付场景 条码支付，取值：bar_code 声波支付，取值：wave_code, bar_code, wave_code
	AuthCode string `json:"auth_code,omitempty"` // 必须 支付授权码
	AuthNo   string `json:"auth_no,omitempty"`   // 可选 预授权冻结交易号

	BuyerId            string             `json:"buyer_id"` // 可选 家的支付宝用户id，如果为空，会从传入了码值信息中获取买家ID
	TransCurrency      string             `json:"trans_currency,omitempty"`
	SettleCurrency     string             `json:"settle_currency,omitempty"`
	DiscountableAmount string             `json:"discountable_amount,omitempty"` // 可选 参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]。 如果该值未传入，但传入了【订单总金额】和【不可打折金额】，则该值默认为【订单总金额】-【不可打折金额】
	GoodsDetail        []*GoodsDetailItem `json:"goods_detail,omitempty"`        // 可选 订单包含的商品列表信息，Json格式，其它说明详见商品明细说明
	OperatorId         string             `json:"operator_id,omitempty"`         // 可选 商户操作员编号
	TerminalId         string             `json:"terminal_id,omitempty"`         // 可选 商户机具终端编号
	AuthConfirmMode    string             `json:"auth_confirm_mode,omitempty"`
	TerminalParams     string             `json:"terminal_params,omitempty"`
	AgreementParams    *AgreementParams   `json:"agreement_params,omitempty"`
}

func (this TradePay) APIName() string {
	return "alipay.trade.pay"
}

func (this TradePay) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// TradePayRsp 统一收单交易支付接口响应参数
type TradePayRsp struct {
	Error
	BuyerLogonId        string           `json:"buyer_logon_id"`           // 买家支付宝账号
	BuyerPayAmount      string           `json:"buyer_pay_amount"`         // 买家实付金额，单位为元，两位小数。
	BuyerUserId         string           `json:"buyer_user_id"`            // 买家在支付宝的用户id
	CardBalance         string           `json:"card_balance"`             // 支付宝卡余额
	DiscountGoodsDetail string           `json:"discount_goods_detail"`    // 本次交易支付所使用的单品券优惠的商品优惠信息
	FundBillList        []*FundBill      `json:"fund_bill_list,omitempty"` // 交易支付使用的资金渠道
	GmtPayment          string           `json:"gmt_payment"`
	InvoiceAmount       string           `json:"invoice_amount"`                // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
	OutTradeNo          string           `json:"out_trade_no"`                  // 创建交易传入的商户订单号
	TradeNo             string           `json:"trade_no"`                      // 支付宝交易号
	PointAmount         string           `json:"point_amount"`                  // 积分支付的金额，单位为元，两位小数。
	ReceiptAmount       string           `json:"receipt_amount"`                // 实收金额，单位为元，两位小数
	StoreName           string           `json:"store_name"`                    // 发生支付交易的商户门店名称
	TotalAmount         string           `json:"total_amount"`                  // 发该笔退款所对应的交易的订单金额
	VoucherDetailList   []*VoucherDetail `json:"voucher_detail_list,omitempty"` // 本交易支付时使用的所有优惠券信息
}

// TradeAppPay App支付接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.app.pay/
type TradeAppPay struct {
	Trade
}

func (this TradeAppPay) APIName() string {
	return "alipay.trade.app.pay"
}

func (this TradeAppPay) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

func (this TradeAppPay) NeedEncrypt() bool {
	return false
}

// TradePreCreate 统一收单线下交易预创建接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.precreate/
type TradePreCreate struct {
	Trade
	DiscountableAmount string             `json:"discountable_amount"`    // 可选 可打折金额. 参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000] 如果该值未传入，但传入了【订单总金额】，【不可打折金额】则该值默认为【订单总金额】-【不可打折金额】
	GoodsDetail        []*GoodsDetailItem `json:"goods_detail,omitempty"` // 可选 订单包含的商品列表信息.Json格式. 其它说明详见：“商品明细说明”
	OperatorId         string             `json:"operator_id"`            // 可选 商户操作员编号
	TerminalId         string             `json:"terminal_id"`            // 可选 商户机具终端编号
}

func (this TradePreCreate) APIName() string {
	return "alipay.trade.precreate"
}

func (this TradePreCreate) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// TradePreCreateRsp 统一收单线下交易预创建接口响应参数
type TradePreCreateRsp struct {
	Error
	OutTradeNo string `json:"out_trade_no"` // 创建交易传入的商户订单号
	QRCode     string `json:"qr_code"`      // 当前预下单请求生成的二维码码串，可以用二维码生成工具根据该码串值生成对应的二维码
}

// TradeCancel 统一收单交易撤销接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.cancel/
type TradeCancel struct {
	AuxParam
	AppAuthToken string `json:"-"` // 可选
	NotifyURL    string `json:"-"` // 可选

	OutTradeNo string `json:"out_trade_no"` // 原支付请求的商户订单号,和支付宝交易号不能同时为空
	TradeNo    string `json:"trade_no"`     // 支付宝交易号，和商户订单号不能同时为空
}

func (this TradeCancel) APIName() string {
	return "alipay.trade.cancel"
}

func (this TradeCancel) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// TradeCancelRsp 统一收单交易撤销接口响应参数
type TradeCancelRsp struct {
	Error
	TradeNo    string `json:"trade_no"`     // 支付宝交易号
	OutTradeNo string `json:"out_trade_no"` // 创建交易传入的商户订单号
	RetryFlag  string `json:"retry_flag"`   // 是否需要重试
	Action     string `json:"action"`       // 本次撤销触发的交易动作 close：关闭交易，无退款 refund：产生了退款
}

// TradeOrderInfoSync 支付宝订单信息同步接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync/
type TradeOrderInfoSync struct {
	AuxParam
	AppAuthToken string `json:"-"`              // 可选
	OutRequestNo string `json:"out_request_no"` // 必选 标识一笔交易多次请求，同一笔交易多次信息同步时需要保证唯一
	BizType      string `json:"biz_type"`       // 必选 交易信息同步对应的业务类型，具体值与支付宝约定；信用授权场景下传CREDIT_AUTH
	TradeNo      string `json:"trade_no"`       // 可选 支付宝交易号，和商户订单号不能同时为空
	OrderBizInfo string `json:"order_biz_info"` // 可选 商户传入同步信息，具体值要和支付宝约定；用于芝麻信用租车、单次授权等信息同步场景，格式为json格式
}

func (this TradeOrderInfoSync) APIName() string {
	return "alipay.trade.orderinfo.sync"
}

func (this TradeOrderInfoSync) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// TradeOrderInfoSyncRsp 支付宝订单信息同步接口响应参数
type TradeOrderInfoSyncRsp struct {
	Error
	TradeNo     string `json:"trade_no"`
	OutTradeNo  string `json:"out_trade_no"`
	BuyerUserId string `json:"buyer_user_id"`
}

// TradeMergePreCreate 统一收单合并支付预创建接口请求参数 https://opendocs.alipay.com/open/028xr9
type TradeMergePreCreate struct {
	AuxParam
	NotifyURL    string `json:"-"` // 可选
	AppAuthToken string `json:"-"` // 可选

	OutMergeNo     string         `json:"out_merge_no"`    // 可选 如果已经和支付宝约定要求子订单明细必须同时支付成功或者同时支付失败则必须传入此参数，且该参数必须在商户端唯一，否则可以不需要填。
	TimeExpire     string         `json:"time_expire"`     // 可选 订单绝对超时时间。格式为yyyy-MM-dd HH:mm:ss。
	TimeoutExpress string         `json:"timeout_express"` // 可选 合并支付订单相对超时时间。从商户合并预下单请求时间开始计算。 请求合并的所有订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。默认值为15d。
	OrderDetails   []*OrderDetail `json:"order_details"`   // 必选 子订单详情
}

func (this TradeMergePreCreate) APIName() string {
	return "alipay.trade.merge.precreate"
}

func (this TradeMergePreCreate) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

type OrderDetail struct {
	AppId          string         `json:"app_id"`                  // 必选 订单明细的应用唯一标识（16位纯数字），指商家的app_id。
	OutTradeNo     string         `json:"out_trade_no"`            // 必选 商户订单号,64个字符以内、只能包含字母、数字、下划线；需保证在商户端不重复
	SellerId       string         `json:"seller_id"`               // 可选 卖家支付宝用户ID。 如果该值与seller_logon_id同时为空，则卖家默认为app_id对应的支付宝用户ID
	SellerLogonId  string         `json:"seller_logon_id"`         // 可选 卖家支付宝logon_id。 支持手机和Email格式,如果该值与seller_id同时传入,将以seller_id为准
	ProductCode    string         `json:"product_code"`            // 必选 销售产品码，与支付宝签约的产品码名称
	TotalAmount    string         `json:"total_amount"`            // 必选 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	Subject        string         `json:"subject"`                 // 必选 订单标题
	Body           string         `json:"body"`                    // 可选 对交易或商品的描述
	ShowURL        string         `json:"show_url"`                // 可选 商品的展示地址
	GoodsDetail    []*GoodsDetail `json:"goods_detail,omitempty"`  // 可选 订单包含的商品列表信息，Json格式，详见商品明细说明
	ExtendParams   *ExtendParams  `json:"extend_params,omitempty"` // 可选 业务扩展参数
	SubMerchant    *Merchant      `json:"sub_merchant,omitempty"`  // 可选 二级商户信息，当前只对直付通特定场景下使用此字段
	SettleInfo     *SettleInfo    `json:"settle_info,omitempty"`   // 可选 描述结算信息，json格式，详见结算参数说明; 直付通场景下必传
	PassbackParams string         `json:"passback_params"`         // 可选 公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝只会在同步返回（包括跳转回商户网站）和异步通知时将该参数原样返回。
}

type ExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id"` // 可选 系统商编号 该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID
	HBFQNum              string `json:"hb_fq_num"`               // 可选 使用花呗分期要进行的分期数
	HBFQSellerPercent    string `json:"hb_fq_seller_percent"`    // 可选 使用花呗分期需要卖家承担的手续费比例的百分值，传入100代表100%
	IndustryRefluxInfo   string `json:"industry_reflux_info"`    // 可选 行业数据回流信息, 详见：地铁支付接口参数补充说明
	CardType             string `json:"card_type"`               // 可选 卡类型
	SpecifiedSellerName  string `json:"specified_seller_name"`   // 可选 特殊场景下，允许商户指定交易展示的卖家名称
}

type Merchant struct {
	MerchantId   string `json:"merchant_id"`   // 必选 支付宝二级商户编号。 间连受理商户的支付宝商户编号，通过间连商户入驻接口后由支付宝生成。直付通和机构间连业务场景下必传。
	MerchantType string `json:"merchant_type"` // 可选 二级商户编号类型。 枚举值：alipay:支付宝分配的间联商户编号；目前仅支持alipay，默认可以不传。
}

type SettleInfo struct {
	SettleDetailInfos []*SettleDetailInfo `json:"settle_detail_infos"` // 必选 结算详细信息，json数组，目前只支持一条。
	SettlePeriodTime  string              `json:"settle_period_time"`  // 可选 该笔订单的超期自动确认结算时间，到达期限后，将自动确认结算。此字段只在签约账期结算模式时有效。取值范围：1d～365d。d-天。 该参数数值不接受小数点。
}

type SettleDetailInfo struct {
	TransInType      string `json:"trans_in_type"`      // 必选 结算收款方的账户类型。 cardAliasNo：结算收款方的银行卡编号;userId：表示是支付宝账号对应的支付宝唯一用户号;loginName：表示是支付宝登录号；defaultSettle：表示结算到商户进件时设置的默认结算账号，结算主体为门店时不支持传defaultSettle；
	TransIn          string `json:"trans_in"`           // 必选 结算收款方。当结算收款方类型是cardAliasNo时，本参数为用户在支付宝绑定的卡编号；结算收款方类型是userId时，本参数为用户的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；当结算收款方类型是loginName时，本参数为用户的支付宝登录号；当结算收款方类型是defaultSettle时，本参数不能传值，保持为空。
	SummaryDimension string `json:"summary_dimension"`  // 可选 结算汇总维度，按照这个维度汇总成批次结算，由商户指定。 目前需要和结算收款方账户类型为cardAliasNo配合使用
	SettleEntityId   string `json:"settle_entity_id"`   // 可选 结算主体标识。当结算主体类型为SecondMerchant时，为二级商户的SecondMerchantID；当结算主体类型为Store时，为门店的外标。
	SettleEntityType string `json:"settle_entity_type"` // 可选 结算主体类型。 二级商户:SecondMerchant;商户或者直连商户门店:Store
	Amount           string `json:"amount"`             // 可选 结算的金额，单位为元。在创建订单和支付接口时必须和交易金额相同。在结算确认接口时必须等于交易金额减去已退款金额。
}

// TradeMergePreCreateRsp 统一收单合并支付预创建接口响应参数
type TradeMergePreCreateRsp struct {
	Error
	OutMergeNo         string            `json:"out_merge_no"`
	PreOrderNo         string            `json:"pre_order_no"`
	OrderDetailResults []*PreOrderResult `json:"order_detail_results"`
}

type PreOrderResult struct {
	AppId      string `json:"app_id"`
	OutTradeNo string `json:"out_trade_no"`
	Success    bool   `json:"success"`
	ResultCode string `json:"result_code"`
}

// TradeAppMergePay App合并支付接口 https://opendocs.alipay.com/open/028py8
type TradeAppMergePay struct {
	AuxParam
	AppAuthToken string `json:"-"`            // 可选
	PreOrderNo   string `json:"pre_order_no"` // 必选 预下单号。通过 alipay.trade.merge.precreate(统一收单合并支付预创建接口)返回。
}

func (this TradeAppMergePay) APIName() string {
	return "alipay.trade.app.merge.pay"
}

func (this TradeAppMergePay) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}
