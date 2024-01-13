package alipay

// FundTransToAccountTransfer 单笔转账到支付宝账户接口请求参数 https://opendocs.alipay.com/apis/api_33/alipay.fund.trans.toaccount.transfer
type FundTransToAccountTransfer struct {
	AuxParam
	AppAuthToken  string `json:"-"`               // 可选
	OutBizNo      string `json:"out_biz_no"`      // 必选 商户转账唯一订单号
	PayeeType     string `json:"payee_type"`      // 必选 收款方账户类型,"ALIPAY_LOGONID":支付宝帐号
	PayeeAccount  string `json:"payee_account"`   // 必选 收款方账户。与payee_type配合使用
	Amount        string `json:"amount"`          // 必选 转账金额,元
	PayerShowName string `json:"payer_show_name"` // 可选 付款方显示姓名
	PayeeRealName string `json:"payee_real_name"` // 可选 收款方真实姓名,如果本参数不为空，则会校验该账户在支付宝登记的实名是否与收款方真实姓名一致。
	Remark        string `json:"remark"`          // 可选 转账备注,金额大于50000时必填
}

func (this FundTransToAccountTransfer) APIName() string {
	return "alipay.fund.trans.toaccount.transfer"
}

func (this FundTransToAccountTransfer) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// FundTransToAccountTransferRsp 单笔转账到支付宝账户接口响应参数
type FundTransToAccountTransferRsp struct {
	Error
	OutBizNo string `json:"out_biz_no"` // 商户转账唯一订单号：发起转账来源方定义的转账单据号。请求时对应的参数，原样返回
	OrderId  string `json:"order_id"`   // 支付宝转账单据号，成功一定返回，失败可能不返回也可能返回
	PayDate  string `json:"pay_date"`   // 支付时间：格式为yyyy-MM-dd HH:mm:ss，仅转账成功返回
}

// FundTransOrderQuery 查询转账订单接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.trans.order.query/
type FundTransOrderQuery struct {
	AuxParam
	AppAuthToken string `json:"-"`                    // 可选
	OutBizNo     string `json:"out_biz_no,omitempty"` // 与 OrderId 二选一
	OrderId      string `json:"order_id,omitempty"`   // 与 OutBizNo 二选一
}

func (this FundTransOrderQuery) APIName() string {
	return "alipay.fund.trans.order.query"
}

func (this FundTransOrderQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// FundTransOrderQueryRsp 查询转账订单接口响应参数
type FundTransOrderQueryRsp struct {
	Error
	OutBizNo       string `json:"out_biz_no"`       // 发起转账来源方定义的转账单据号。 该参数的赋值均以查询结果中 的 out_biz_no 为准。 如果查询失败，不返回该参数
	OrderId        string `json:"order_id"`         // 支付宝转账单据号，查询失败不返回。
	Status         string `json:"status"`           // 转账单据状态
	PayDate        string `json:"pay_date"`         // 支付时间
	ArrivalTimeEnd string `json:"arrival_time_end"` // 预计到账时间，转账到银行卡专用
	OrderFree      string `json:"order_fee"`        // 预计收费金额（元），转账到银行卡专用
	FailReason     string `json:"fail_reason"`      // 查询到的订单状态为FAIL失败或REFUND退票时，返回具体的原因。
	ErrorCode      string `json:"error_code"`       // 查询失败时，本参数为错误代 码。 查询成功不返回。 对于退票订单，不返回该参数。
}

// FundAuthOrderVoucherCreate 资金授权发码接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.voucher.create/
type FundAuthOrderVoucherCreate struct {
	AuxParam
	NotifyURL         string `json:"-"`
	AppAuthToken      string `json:"-"`                             // 可选
	OutOrderNo        string `json:"out_order_no"`                  // 必选, 商户授权资金订单号，创建后不能修改，需要保证在商户端不重复。
	OutRequestNo      string `json:"out_request_no"`                // 必选, 商户本次资金操作的请求流水号，用于标示请求流水的唯一性，需要保证在商户端不重复。
	ProductCode       string `json:"product_code,omitempty"`        // 必选, 销售产品码，后续新接入预授权当面付的业务，本字段取值固定为PRE_AUTH。
	OrderTitle        string `json:"order_title"`                   // 必选, 业务订单的简单描述，如商品名称等 长度不超过100个字母或50个汉字
	Amount            string `json:"amount"`                        // 必选, 需要冻结的金额，单位为：元（人民币），精确到小数点后两位 取值范围：[0.01,100000000.00]
	PayeeUserId       string `json:"payee_user_id,omitempty"`       // 可选, 收款方的支付宝唯一用户号,以2088开头的16位纯数字组成，如果非空则会在支付时校验交易的的收款方与此是否一致，如果商户有勾选花呗渠道，收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)不能同时为空。
	PayeeLogonId      string `json:"payee_logon_id,omitempty"`      // 可选, 收款方支付宝账号（Email或手机号），如果收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)同时传递，则以用户号(payee_user_id)为准，如果商户有勾选花呗渠道，收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)不能同时为空。
	PayTimeout        string `json:"pay_timeout,omitempty"`         // 可选, 该笔订单允许的最晚付款时间，逾期将关闭该笔订单 取值范围：1m～15d。m-分钟，h-小时，d-天。 该参数数值不接受小数点， 如 1.5h，可转换为90m 如果为空，默认15m
	ExtraParam        string `json:"extra_param,omitempty"`         // 可选, 业务扩展参数，用于商户的特定业务信息的传递，json格式。 1.授权业务对应的类目，key为category，value由支付宝分配，比如充电桩业务传 "CHARGE_PILE_CAR"； 2. 外部商户的门店编号，key为outStoreCode，可选； 3. 外部商户的门店简称，key为outStoreAlias，可选。
	TransCurrency     string `json:"trans_currency,omitempty"`      // 可选, 标价币种, amount 对应的币种单位。支持澳元：AUD, 新西兰元：NZD, 台币：TWD, 美元：USD, 欧元：EUR, 英镑：GBP
	SettleCurrency    string `json:"settle_currency,omitempty"`     // 可选, 商户指定的结算币种。支持澳元：AUD, 新西兰元：NZD, 台币：TWD, 美元：USD, 欧元：EUR, 英镑：GBP
	EnablePayChannels string `json:"enable_pay_channels,omitempty"` // 可选, 商户可用该参数指定用户可使用的支付渠道，本期支持商户可支持三种支付渠道，余额宝（MONEY_FUND）、花呗（PCREDIT_PAY）以及芝麻信用（CREDITZHIMA）。商户可设置一种支付渠道，也可设置多种支付渠道。
}

func (this FundAuthOrderVoucherCreate) APIName() string {
	return "alipay.fund.auth.order.voucher.create"
}

func (this FundAuthOrderVoucherCreate) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// FundAuthOrderVoucherCreateRsp 资金授权发码接口响应参数
type FundAuthOrderVoucherCreateRsp struct {
	Error
	OutOrderNo   string `json:"out_order_no"`
	OutRequestNo string `json:"out_request_no"`
	CodeType     string `json:"code_type"`
	CodeValue    string `json:"code_value"`
	CodeURL      string `json:"code_url"`
}

// FundAuthOrderFreeze 资金授权冻结接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.freeze/
type FundAuthOrderFreeze struct {
	AuxParam
	NotifyURL    string `json:"-"`
	AppAuthToken string `json:"-"`                        // 可选
	AuthCode     string `json:"auth_code"`                // 必选, 支付授权码，25~30开头的长度为16~24位的数字，实际字符串长度以开发者获取的付款码长度为准
	AuthCodeType string `json:"auth_code_type"`           // 必选, 授权码类型 目前仅支持"bar_code"
	OutOrderNo   string `json:"out_order_no"`             // 必选, 商户授权资金订单号 ,不能包含除中文、英文、数字以外的字符，创建后不能修改，需要保证在商户端不重复。
	OutRequestNo string `json:"out_request_no"`           // 必选, 商户本次资金操作的请求流水号，用于标示请求流水的唯一性，不能包含除中文、英文、数字以外的字符，需要保证在商户端不重复。
	OrderTitle   string `json:"order_title"`              // 必选, 业务订单的简单描述，如商品名称等 长度不超过100个字母或50个汉字
	Amount       string `json:"amount"`                   // 必选, 需要冻结的金额，单位为：元（人民币），精确到小数点后两位 取值范围：[0.01,100000000.00]
	PayeeLogonId string `json:"payee_logon_id,omitempty"` // 可选, 收款方支付宝账号（Email或手机号），如果收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)同时传递，则以用户号(payee_user_id)为准，如果商户有勾选花呗渠道，收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)不能同时为空。
	PayeeUserId  string `json:"payee_user_id,omitempty"`  // 可选, 收款方的支付宝唯一用户号,以2088开头的16位纯数字组成，如果非空则会在支付时校验交易的的收款方与此是否一致，如果商户有勾选花呗渠道，收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)不能同时为空。
	PayTimeout   string `json:"pay_timeout,omitempty"`    // 可选, 该笔订单允许的最晚付款时间，逾期将关闭该笔订单 取值范围：1m～15d。m-分钟，h-小时，d-天。 该参数数值不接受小数点， 如 1.5h，可转换为90m 如果为空，默认15m
	ExtraParam   string `json:"extra_param,omitempty"`    // 可选, 业务扩展参数，用于商户的特定业务信息的传递，json格式。 1.授权业务对应的类目，key为category，value由支付宝分配，比如充电桩业务传 "CHARGE_PILE_CAR"； 2. 外部商户的门店编号，key为outStoreCode，可选； 3. 外部商户的门店简称，key为outStoreAlias，可选。
	ProductCode  string `json:"product_code,omitempty"`   // 可选, 销售产品码，后续新接入预授权当面付的业务，本字段取值固定为PRE_AUTH。
}

func (this FundAuthOrderFreeze) APIName() string {
	return "alipay.fund.auth.order.app.freeze"
}

func (this FundAuthOrderFreeze) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// FundAuthOrderFreezeRsp 资金授权冻结接口响应参数
type FundAuthOrderFreezeRsp struct {
	Error
	AuthNo       string `json:"auth_no"`
	OutOrderNo   string `json:"out_order_no"`
	OperationId  string `json:"operation_id"`
	OutRequestNo string `json:"out_request_no"`
	Amount       string `json:"amount"`
	Status       string `json:"status"`
	PayerUserId  string `json:"payer_user_id"`
	GMTTrans     string `json:"gmt_trans"`
}

// FundAuthOrderUnfreeze 资金授权解冻接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.unfreeze/
type FundAuthOrderUnfreeze struct {
	AuxParam
	NotifyURL    string `json:"-"`
	AuthNo       string `json:"auth_no"`               // 必选,支付宝资金授权订单号,支付宝冻结时返回的交易号，数字格式 2016101210002001810258115912
	AppAuthToken string `json:"-"`                     // 可选
	OutRequestNo string `json:"out_request_no"`        // 必选, 商户本次资金操作的请求流水号，用于标示请求流水的唯一性，不能包含除中文、英文、数字以外的字符，需要保证在商户端不重复。
	Amount       string `json:"amount"`                // 必选, 本次操作解冻的金额，单位为：元（人民币），精确到小数点后两位，取值范围：[0.01,100000000.00]
	Remark       string `json:"remark"`                // 必选, 商户对本次解冻操作的附言描述，长度不超过100个字母或50个汉字
	ExtraParam   string `json:"extra_param,omitempty"` // 可选, 解冻扩展信息，json格式；unfreezeBizInfo 目前为芝麻消费字段，支持Key值如下： "bizComplete":"true" -- 选填：标识本次解冻用户是否履约，如果true信用单会完结为COMPLETE
}

func (this FundAuthOrderUnfreeze) APIName() string {
	return "alipay.fund.auth.order.unfreeze"
}

func (this FundAuthOrderUnfreeze) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// FundAuthOrderUnfreezeRsp 资金授权解冻接口响应参数
type FundAuthOrderUnfreezeRsp struct {
	Error
	AuthNo       string `json:"auth_no"`
	OutOrderNo   string `json:"out_order_no"`
	OperationId  string `json:"operation_id"`
	OutRequestNo string `json:"out_request_no"`
	Amount       string `json:"amount"`
	Status       string `json:"status"`
	GMTTrans     string `json:"gmt_trans"`
	CreditAmount string `json:"credit_amount"`
	FundAmount   string `json:"fund_amount"`
}

// FundAuthOperationCancel 资金授权撤销接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.auth.operation.cancel/
type FundAuthOperationCancel struct {
	AuxParam
	NotifyURL    string `json:"-"`
	AppAuthToken string `json:"-"`                        // 可选
	AuthNo       string `json:"auth_no,omitempty"`        // 特殊可选, 支付宝授权资金订单号，与商户的授权资金订单号不能同时为空，二者都存在时，以支付宝资金授权订单号为准，该参数与支付宝授权资金操作流水号配对使用。
	OutOrderNo   string `json:"out_order_no,omitempty"`   // 特殊可选,  商户的授权资金订单号，与支付宝的授权资金订单号不能同时为空，二者都存在时，以支付宝的授权资金订单号为准，该参数与商户的授权资金操作流水号配对使用。
	OperationId  string `json:"operation_id,omitempty"`   // 特殊可选, 支付宝的授权资金操作流水号，与商户的授权资金操作流水号不能同时为空，二者都存在时，以支付宝的授权资金操作流水号为准，该参数与支付宝授权资金订单号配对使用。
	OutRequestNo string `json:"out_request_no,omitempty"` // 特殊可选, 商户的授权资金操作流水号，与支付宝的授权资金操作流水号不能同时为空，二者都存在时，以支付宝的授权资金操作流水号为准，该参数与商户的授权资金订单号配对使用。
	Remark       string `json:"remark"`                   // 必选, 商户对本次撤销操作的附言描述，长度不超过100个字母或50个汉字
}

func (this FundAuthOperationCancel) APIName() string {
	return "alipay.fund.auth.operation.cancel"
}

func (this FundAuthOperationCancel) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// FundAuthOperationCancelRsp 资金授权撤销接口响应参数
type FundAuthOperationCancelRsp struct {
	Error
	AuthNo       string `json:"auth_no"`
	OutOrderNo   string `json:"out_order_no"`
	OperationId  string `json:"operation_id"`
	OutRequestNo string `json:"out_request_no"`
	Action       string `json:"action"`
}

// FundAuthOperationDetailQuery 资金授权操作查询接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.auth.operation.detail.query/
type FundAuthOperationDetailQuery struct {
	AuxParam
	AppAuthToken string `json:"-"`              // 可选
	AuthNo       string `json:"auth_no"`        // 特殊可选, 支付宝授权资金订单号，与商户的授权资金订单号不能同时为空，二者都存在时，以支付宝资金授权订单号为准，该参数与支付宝授权资金操作流水号配对使用。
	OutOrderNo   string `json:"out_order_no"`   // 特殊可选, 商户的授权资金订单号，与支付宝的授权资金订单号不能同时为空，二者都存在时，以支付宝的授权资金订单号为准，该参数与商户的授权资金操作流水号配对使用。
	OperationId  string `json:"operation_id"`   // 特殊可选, 支付宝的授权资金操作流水号，与商户的授权资金操作流水号不能同时为空，二者都存在时，以支付宝的授权资金操作流水号为准，该参数与支付宝授权资金订单号配对使用。
	OutRequestNo string `json:"out_request_no"` // 特殊可选, 商户的授权资金操作流水号，与支付宝的授权资金操作流水号不能同时为空，二者都存在时，以支付宝的授权资金操作流水号为准，该参数与商户的授权资金订单号配对使用。
}

func (this FundAuthOperationDetailQuery) APIName() string {
	return "alipay.fund.auth.operation.detail.query"
}

func (this FundAuthOperationDetailQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

type OrderStatus string

const (
	OrderStatusInit       OrderStatus = "INIT"       //（初始状态：已创建未授权）
	OrderStatusAuthorized OrderStatus = "AUTHORIZED" //（已授权状态：授权成功，可以进行转支付或解冻操作）
	OrderStatusFinish     OrderStatus = "FINISH"     //（完成状态：转支付完成且无剩余冻结资金）
	OrderStatusClosed     OrderStatus = "CLOSED"     //（关闭状态：授权未完成超时关闭或冻结资金全额解冻）
)

// FundAuthOperationDetailQueryRsp 资金授权操作查询接口响应参数
type FundAuthOperationDetailQueryRsp struct {
	Error
	AuthNo                  string      `json:"auth_no"`
	OutOrderNo              string      `json:"out_order_no"`
	OrderStatus             OrderStatus `json:"order_status"`
	TotalFreezeAmount       string      `json:"total_freeze_amount"`
	RestAmount              string      `json:"rest_amount"`
	TotalPayAmount          string      `json:"total_pay_amount"`
	OrderTitle              string      `json:"order_title"`
	PayerLogonId            string      `json:"payer_logon_id"`
	PayerUserId             string      `json:"payer_user_id"`
	ExtraParam              string      `json:"extra_param"`
	OperationId             string      `json:"operation_id"`
	OutRequestNo            string      `json:"out_request_no"`
	Amount                  string      `json:"amount"`
	OperationType           string      `json:"operation_type"`
	Status                  string      `json:"status"`
	Remark                  string      `json:"remark"`
	GMTCreate               string      `json:"gmt_create"`
	GMTTrans                string      `json:"gmt_trans"`
	PreAuthType             string      `json:"pre_auth_type"`
	TransCurrency           string      `json:"trans_currency"`
	TotalFreezeCreditAmount string      `json:"total_freeze_credit_amount"`
	TotalFreezeFundAmount   string      `json:"total_freeze_fund_amount"`
	TotalPayCreditAmount    string      `json:"total_pay_credit_amount"`
	TotalPayFundAmount      string      `json:"total_pay_fund_amount"`
	RestCreditAmount        string      `json:"rest_credit_amount"`
	RestFundAmount          string      `json:"rest_fund_amount"`
	CreditAmount            string      `json:"credit_amount"`
	FundAmount              string      `json:"fund_amount"`
}

// FundAuthOrderAppFreeze 线上资金授权冻结接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.auth.order.app.freeze
type FundAuthOrderAppFreeze struct {
	AuxParam
	NotifyURL         string `json:"-"`
	AppAuthToken      string `json:"-"`                             // 可选
	OutOrderNo        string `json:"out_order_no"`                  // 必选, 商户授权资金订单号 ,不能包含除中文、英文、数字以外的字符，创建后不能修改，需要保证在商户端不重复。
	OutRequestNo      string `json:"out_request_no"`                // 必选, 商户本次资金操作的请求流水号，用于标示请求流水的唯一性，不能包含除中文、英文、数字以外的字符，需要保证在商户端不重复。
	OrderTitle        string `json:"order_title"`                   // 必选, 业务订单的简单描述，如商品名称等 长度不超过100个字母或50个汉字
	Amount            string `json:"amount"`                        // 必选, 需要冻结的金额，单位为：元（人民币），精确到小数点后两位 取值范围：[0.01,100000000.00]
	ProductCode       string `json:"product_code"`                  // 必选, 销售产品码，新接入线上预授权的业务，本字段取值固定为PRE_AUTH_ONLINE 。
	PayeeLogonId      string `json:"payee_logon_id,omitempty"`      // 收款方支付宝账号（Email或手机号），如果收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)同时传递，则以用户号(payee_user_id)为准，如果商户有勾选花呗渠道，收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)不能同时为空。
	PayeeUserId       string `json:"payee_user_id,omitempty"`       // 收款方的支付宝唯一用户号,以2088开头的16位纯数字组成，如果非空则会在支付时校验交易的的收款方与此是否一致，如果商户有勾选花呗渠道，收款方支付宝登录号(payee_logon_id)和用户号(payee_user_id)不能同时为空。
	PayTimeout        string `json:"pay_timeout,omitempty"`         // 该笔订单允许的最晚付款时间，逾期将关闭该笔订单 取值范围：1m～15d。m-分钟，h-小时，d-天。 该参数数值不接受小数点， 如 1.5h，可转换为90m 如果为空，默认15m
	ExtraParam        string `json:"extra_param,omitempty"`         // 业务扩展参数，用于商户的特定业务信息的传递，json格式。 1.授权业务对应的类目，key为category，value由支付宝分配，比如充电桩业务传 "CHARGE_PILE_CAR"； 2. 外部商户的门店编号，key为outStoreCode，可选； 3. 外部商户的门店简称，key为outStoreAlias，可选。
	EnablePayChannels string `json:"enable_pay_channels,omitempty"` // 商户可用该参数指定用户可使用的支付渠道，本期支持商户可支持三种支付渠道，余额宝（MONEY_FUND）、花呗（PCREDIT_PAY）以及芝麻信用（CREDITZHIMA）。商户可设置一种支付渠道，也可设置多种支付渠道。
}

func (this FundAuthOrderAppFreeze) APIName() string {
	return "alipay.fund.auth.order.app.freeze"
}

func (this FundAuthOrderAppFreeze) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// FundAuthOrderAppFreezeRsp 线上资金授权冻结接口响应参数
type FundAuthOrderAppFreezeRsp struct {
	Error
	AuthNo       string `json:"auth_no"`
	OutOrderNo   string `json:"out_order_no"`
	OperationId  string `json:"operation_id"`
	OutRequestNo string `json:"out_request_no"`
	Amount       string `json:"amount"`
	Status       string `json:"status"`
	PayerUserId  string `json:"payer_user_id"`
	GMTTrans     string `json:"gmt_trans"`
	PreAuthType  string `json:"pre_auth_type"`
	CreditAmount string `json:"credit_amount"`
	FundAmount   string `json:"fund_amount"`
}

// FundTransUniTransfer 单笔转账接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.trans.uni.transfer/
type FundTransUniTransfer struct {
	AuxParam
	AppAuthToken    string     `json:"-"`                 // 可选
	OutBizNo        string     `json:"out_biz_no"`        // 必选 商户端的唯一订单号，对于同一笔转账请求，商户需保证该订单号唯一。
	TransAmount     string     `json:"trans_amount"`      // 必选 订单总金额，单位为元，精确到小数点后两位，STD_RED_PACKET 产品取值范围[0.01,100000000]； TRANS_ACCOUNT_NO_PWD产品取值范围[0.1,100000000]
	ProductCode     string     `json:"product_code"`      // 必选 业务产品码， 收发现金红包固定为：STD_RED_PACKET； 单笔无密转账到支付宝账户固定为：TRANS_ACCOUNT_NO_PWD； 单笔无密转账到银行卡固定为：TRANS_BANKCARD_NO_PWD
	BizScene        string     `json:"biz_scene"`         // 可选 描述特定的业务场景，可传的参数如下： PERSONAL_COLLECTION：C2C现金红包-领红包； DIRECT_TRANSFER：B2C现金红包、单笔无密转账到支付宝/银行卡
	OrderTitle      string     `json:"order_title"`       // 可选 转账业务的标题，用于在支付宝用户的账单里显示
	OriginalOrderId string     `json:"original_order_id"` // 可选 原支付宝业务单号。C2C现金红包-红包领取时，传红包支付时返回的支付宝单号；B2C现金红包、单笔无密转账到支付宝/银行卡不需要该参数。
	PayeeInfo       *PayeeInfo `json:"payee_info"`        // 必选 收款方信息
	Remark          string     `json:"remark"`            // 可选 业务备注
	BusinessParams  string     `json:"business_params"`   // 可选 转账业务请求的扩展参数，支持传入的扩展参数如下： 1、sub_biz_scene 子业务场景，红包业务必传，取值REDPACKET，C2C现金红包、B2C现金红包均需传入； 2、withdraw_timeliness为转账到银行卡的预期到账时间，可选（不传入则默认为T1），取值T0表示预期T+0到账，取值T1表示预期T+1到账，因到账时效受银行机构处理影响，支付宝无法保证一定是T0或者T1到账；
}

func (this FundTransUniTransfer) APIName() string {
	return "alipay.fund.trans.uni.transfer"
}

func (this FundTransUniTransfer) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

type PayeeInfo struct {
	Identity     string `json:"identity"`      // 必填 参与方的唯一标识
	IdentityType string `json:"identity_type"` // 必填 参与方的标识类型，目前支持如下类型： 1、ALIPAY_USER_ID 支付宝的会员ID 2、ALIPAY_LOGON_ID：支付宝登录号，支持邮箱和手机号格式
	Name         string `json:"name"`          // 可选 参与方真实姓名，如果非空，将校验收款支付宝账号姓名一致性。当identity_type=ALIPAY_LOGON_ID时，本字段必填。
}

// FundTransUniTransferRsp 单笔转账接口响应参数
type FundTransUniTransferRsp struct {
	Error
	OutBizNo       string `json:"out_biz_no"`        // 用户订单号
	OrderId        string `json:"order_id"`          // 支付宝转账订单号
	PayFundOrderId string `json:"pay_fund_order_id"` // 支付宝支付资金流水号
	Status         string `json:"status"`            // 转账单据状态。 SUCCESS：成功（对转账到银行卡的单据, 该状态可能变为退票[REFUND]状态）； FAIL：失败（具体失败原因请参见error_code以及fail_reason返回值）； DEALING：处理中； REFUND：退票；
	TransDate      string `json:"trans_date"`        // 订单支付时间，格式为yyyy-MM-dd HH:mm:ss
}

// FundTransCommonQuery 转账业务单据查询接口请求参数 https://docs.open.alipay.com/api_28/alipay.fund.trans.common.query/
type FundTransCommonQuery struct {
	AuxParam
	AppAuthToken   string `json:"-"`                 // 可选
	ProductCode    string `json:"product_code"`      // 必选 业务产品码， 收发现金红包固定为：STD_RED_PACKET； 单笔无密转账到支付宝账户固定为：TRANS_ACCOUNT_NO_PWD； 单笔无密转账到银行卡固定为：TRANS_BANKCARD_NO_PWD
	BizScene       string `json:"biz_scene"`         // 必选 描述特定的业务场景，可传的参数如下： PERSONAL_COLLECTION：C2C现金红包-领红包； DIRECT_TRANSFER：B2C现金红包、单笔无密转账到支付宝/银行卡
	OutBizNo       string `json:"out_biz_no"`        // 可选 商户端的唯一订单号，对于同一笔转账请求，商户需保证该订单号唯一。
	OrderId        string `json:"order_id"`          // 可选 支付宝转账单据号
	PayFundOrderId string `json:"pay_fund_order_id"` // 可选 支付宝支付资金流水号
}

func (this FundTransCommonQuery) APIName() string {
	return "alipay.fund.trans.common.query"
}

func (this FundTransCommonQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// FundTransCommonQueryRsp 转账业务单据查询接口响应参数
type FundTransCommonQueryRsp struct {
	Error
	OrderId          string `json:"order_id"`           // 支付宝转账订单号
	PayFundOrderId   string `json:"pay_fund_order_id"`  // 支付宝支付资金流水号
	OutBizNo         string `json:"out_biz_no"`         // 用户订单号
	TransAmount      string `json:"trans_amount"`       // 付款金额
	Status           string `json:"status"`             // 转账单据状态。 SUCCESS：成功（对转账到银行卡的单据, 该状态可能变为退票[REFUND]状态）； FAIL：失败（具体失败原因请参见error_code以及fail_reason返回值）； DEALING：处理中； REFUND：退票；
	PayDate          string `json:"pay_date"`           // 支付时间
	ArrivalTimeEnd   string `json:"arrival_time_end"`   // 预计到账时间
	OrderFee         string `json:"order_fee"`          // 预计收费金额
	ErrorCode        string `json:"error_code"`         // 查询到的订单状态为FAIL失败或REFUND退票时，返回错误代码
	FailReason       string `json:"fail_reason"`        // 查询到的订单状态为FAIL失败或REFUND退票时，返回具体的原因。
	DeductBillInfo   string `json:"deduct_bill_info"`   // 商户查询代扣订单信息时返回其在代扣请求中传入的账单属性
	TransferBillInfo string `json:"transfer_bill_info"` // 商户在查询代发订单信息时返回其在代发请求中传入的账单属性。
}

// FundAccountQuery 支付宝资金账户资产查询接口请求参数  https://docs.open.alipay.com/api_28/alipay.fund.account.query
type FundAccountQuery struct {
	AuxParam
	AppAuthToken string `json:"-"`              // 可选
	AliPayUserId string `json:"alipay_user_id"` // 必选 蚂蚁统一会员ID
	AccountType  string `json:"account_type"`   // 特殊可选 查询的账号类型，如查询托管账户值为TRUSTEESHIP_ACCOUNT，查询余额账户值为ACCTRANS_ACCOUNT。查询余额账户时必填。
}

func (this FundAccountQuery) APIName() string {
	return "alipay.fund.account.query"
}

func (this FundAccountQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// FundAccountQueryRsp 支付宝资金账户资产查询接口响应参数
type FundAccountQueryRsp struct {
	Error
	AvailableAmount string `json:"available_amount"`
	FreezeAmount    string `json:"freeze_amount"`
}

// FundTransAppPay https://opendocs.alipay.com/open/03rbyf https://opendocs.alipay.com/open/03rbyf
type FundTransAppPay struct {
	AuxParam
	AppAuthToken     string `json:"-"`                  // 可选
	OutBizNo         string `json:"out_biz_no"`         // 必选 商户端的唯一订单号，对于同一笔转账请求，商户需保证该订单号唯一。
	TransAmount      string `json:"trans_amount"`       // 必选 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,9999999999999.99]
	OrderId          string `json:"order_id"`           // 可选 支付宝订单号
	ProductCode      string `json:"product_code"`       // 必选 销售产品码，商家和支付宝签约的产品码。 STD_RED_PACKET：现金红包
	BizScene         string `json:"biz_scene"`          // 必选 描述特定的业务场景，可传值如下： PERSONAL_PAY: 发红包
	Remark           string `json:"remark"`             // 可选 支付备注
	OrderTitle       string `json:"order_title"`        // 可选 支付订单的标题，用于在收银台和消费记录展示
	TimeExpire       string `json:"time_expire"`        // 可选 绝对超时时间，格式为yyyy-MM-dd HH:mm
	RefundTimeExpire string `json:"refund_time_expire"` // 可选 退款超时时间，格式yyyy-MM-dd HH:mm。到指定时间后，系统会自动触发退款，并原路退回到付款账户。如果指定了退款时间，必须早于销售方案里设置的最晚退款时间。
	BusinessParams   string `json:"business_params"`    // 可选 JSON格式，传递业务扩展参数. 业务扩展字段，JSON格式。支持如下属性： sub_biz_scene 子场景，必填，传REDPACKET payer_binded_alipay_uid 创建红包的商户会员绑定的支付宝userId，必填
}

func (this FundTransAppPay) APIName() string {
	return "alipay.fund.trans.app.pay"
}

func (this FundTransAppPay) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}
