package alipay

// BillDownloadURLQuery 查询对账单下载地址接口请求参数 https://docs.open.alipay.com/api_15/alipay.data.dataservice.bill.downloadurl.query
type BillDownloadURLQuery struct {
	AuxParam
	AppAuthToken string `json:"-"`         // 可选
	BillType     string `json:"bill_type"` // 必选 账单类型，商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型：trade、signcustomer；trade指商户基于支付宝交易收单的业务账单；signcustomer是指基于商户支付宝余额收入及支出等资金变动的帐务账单。
	BillDate     string `json:"bill_date"` // 必选 账单时间：日账单格式为yyyy-MM-dd，最早可下载2016年1月1日开始的日账单；月账单格式为yyyy-MM，最早可下载2016年1月开始的月账单。
}

func (this BillDownloadURLQuery) APIName() string {
	return "alipay.data.dataservice.bill.downloadurl.query"
}

func (this BillDownloadURLQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// BillAccountLogQuery 查询账户账务明细接口请求参数 https://opendocs.alipay.com/apis/api_15/alipay.data.bill.accountlog.query
type BillAccountLogQuery struct {
	AuxParam
	AppAuthToken         string `json:"-"`                      // 可选
	StartTime            string `json:"start_time"`             // 账务流水创建时间的起始范围	2019-01-01 00:00:00
	EndTime              string `json:"end_time"`               // 账务流水创建时间的结束范围。与起始时间间隔不超过31天。查询结果为起始时间至结束时间的左闭右开区间	2019-01-02 00:00:00
	AliPayOrderNo        string `json:"alipay_order_no"`        // 可选 支付宝订单号。对账使用，不脱敏	20190101***
	MerchantOrderNo      string `json:"merchant_order_no"`      // 可选 商户订单号，创建支付宝交易时传入的信息。对账使用，不脱敏 TX***
	PageNo               string `json:"page_no"`                // 可选 分页号，从1开始	1
	PageSize             string `json:"page_size"`              // 可选 分页大小1000-2000，默认2000	2000
	TransCode            string `json:"trans_code"`             // 可选 账务的类型代码，特殊场景下使用	101101,301101
	AgreementNo          string `json:"agreement_no"`           // 可选 协议授权码，特殊场景下使用	20215606000635888888
	AgreementProductCode string `json:"agreement_product_code"` // 可选 协议产品码。特殊场景下使用	FUND_SIGN_WITHHOLDING
	BillUserId           string `json:"bill_user_id"`           // 可选 指定用户做账单查询	2088123456789012
}

func (this BillAccountLogQuery) APIName() string {
	return "alipay.data.bill.accountlog.query"
}

func (this BillAccountLogQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

type AccountLogItem struct {
	TransDt             string `json:"trans_dt"`               // 入账时间	2019-01-01 00:00:00
	AccountLogId        string `json:"account_log_id"`         // 支付宝账务流水号。对账使用，不脱敏	1***
	AliPayOrderNo       string `json:"alipay_order_no"`        // 支付宝订单号。对账使用，不脱敏	20190101***
	MerchantOrderNo     string `json:"merchant_order_no"`      // 商户订单号，创建支付宝交易时传入的信息。对账使用，不脱敏 TX***
	TransAmount         string `json:"trans_amount"`           // 金额 1000.00
	Balance             string `json:"balance"`                // 余额，仅供参考。由于架构原因，余额变动并不保证连续。也就是余额不一定等于上一笔余额减去当笔金额。但是会保证最终一致。	10000.00
	Type                string `json:"type"`                   // 账务记录的类型，仅供参考	交易
	OtherAccount        string `json:"other_account"`          // 对方账户	张*(a*@******.com)
	TransMemo           string `json:"trans_memo"`             // 账务备注。由上游业务决定，不可依赖此字段进行对账	备注1
	Direction           string `json:"direction"`              // 收入/支出。表示收支。amount是正数，返回“收入”。amount是负数，返回“支出”	收入
	BillSource          string `json:"bill_source"`            // 业务账单来源，资金收支对应的上游业务订单数据来源，确认业务订单出处。此字段供商户对账使用，不脱敏。	商家中心
	BizNos              string `json:"biz_nos"`                // 业务订单号，资金收支相关的业务场景订单号明细，字母大写；M：平台交易主单号，S：平台交易子单号，O：业务系统单据号（如退款订单号）。此字段供商户对账使用，不脱敏。	M{330***}|S{330***}|O{192***}
	BizOrigNo           string `json:"biz_orig_no"`            // 业务基础订单号，资金收支对应的原始业务订单唯一识别编号。此字段供商户对账使用，不脱敏。	330***
	BizDesc             string `json:"biz_desc"`               // 业务描述，资金收支对应的详细业务场景信息。此字段供商户对账使用，不脱敏。	002***|交易退款
	MerchantOutRefundNo string `json:"merchant_out_refund_no"` // 支付宝交易商户退款请求号。对应商户在调用收单退款接口openApi请求传入的outRequestNo参数值	20211119***
	ComplementInfo      string `json:"complement_info"`        // 账单的补全信息，用于特殊场景，普通商户不需要传值，对账时可忽略。	商家中心
	StoreName           string `json:"store_name"`
}

type BillAccountLogQueryRsp struct {
	Error
	PageNo     string            `json:"page_no"`
	PageSize   string            `json:"page_size"`
	TotalSize  string            `json:"total_size"`
	DetailList []*AccountLogItem `json:"detail_list"`
}

// BillDownloadURLQueryRsp 查询对账单下载地址接口响应参数
type BillDownloadURLQueryRsp struct {
	Error
	BillDownloadURL string `json:"bill_download_url"`
}

// BillBalanceQuery 支付宝商家账户当前余额查询 https://opendocs.alipay.com/apis/api_15/alipay.data.bill.balance.query
type BillBalanceQuery struct {
	AuxParam
	AppAuthToken string `json:"-"` // 可选
}

func (this BillBalanceQuery) APIName() string {
	return "alipay.data.bill.balance.query"
}

func (this BillBalanceQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// BillBalanceQueryRsp 支付宝商家账户当前余额查询响应参数
type BillBalanceQueryRsp struct {
	Error
	TotalAmount     string `json:"total_amount"`
	AvailableAmount string `json:"available_amount"`
	FreezeAmount    string `json:"freeze_amount"`
}
