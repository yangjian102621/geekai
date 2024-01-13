package alipay

type ZmAuthParams struct {
	BuckleAppId      string `json:"buckle_app_id,omitempty"`      // 商户在芝麻端申请的appId
	BuckleMerchantId string `json:"buckle_merchant_id,omitempty"` // 商户在芝麻端申请的 merchantId,omitempty
}

type ProdParams struct {
	AuthBizParams string `json:"auth_biz_params,omitempty"` // 预授权业务信息
}

type AccessParams struct {
	Channel string `json:"channel,omitempty"` // 目前支持以下值：1.ALIPAYAPP（钱包h5页面签约）2.QRCODE(扫码签约)3.QRCODEORSMS(扫码签约或者短信签约)
}

type SubMerchantParams struct {
	SubMerchantId                 string `json:"sub_merchant_id,omitempty"`
	SubMerchantName               string `json:"sub_merchant_name,omitempty"`
	SubMerchantServiceName        string `json:"sub_merchant_service_name,omitempty"`
	SubMerchantServiceDescription string `json:"sub_merchant_service_description,omitempty"`
}

type DeviceParams struct {
	DeviceId   string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	DeviceType string `json:"device_type,omitempty"` // 设备类型，目前有四种值：VR一体机：VR_MACHINE、电视：TV、身份证：ID_CARD、工牌：WORK_CARD
}

type IdentityParams struct {
	UserName     string `json:"user_name,omitempty"`
	CertNo       string `json:"cert_no,omitempty"`
	IdentityHash string `json:"identity_hash,omitempty"`
	SignUserId   string `json:"sign_user_id,omitempty"`
}

type PeriodRuleParams struct {
	PeriodType    string `json:"period_type,omitempty"`
	Period        string `json:"period,omitempty"`
	ExecuteTime   string `json:"execute_time,omitempty"`
	SingleAmount  string `json:"single_amount,omitempty"`
	TotalAmount   string `json:"total_amount,omitempty"`
	TotalPayments int    `json:"total_payments,omitempty"`
}

// AgreementPageSign 支付宝个人协议页面签约接口请求参数 https://docs.open.alipay.com/api_2/alipay.user.agreement.page.sign
type AgreementPageSign struct {
	AuxParam
	AppAuthToken        string             `json:"-"`
	ReturnURL           string             `json:"-"`
	NotifyURL           string             `json:"-"`
	SignValidityPeriod  string             `json:"sign_validity_period,omitempty"`  // 当前用户签约请求的协议有效周期。整形数字加上时间单位的协议有效期，从发起签约请求的时间开始算起。目前支持的时间单位：1.d：天2.m：月如果未传入，默认为长期有效
	ProductCode         string             `json:"product_code,omitempty"`          // 销售产品码，商户签约的支付宝合同所对应的产品码
	ExternalLogonId     string             `json:"external_logon_id,omitempty"`     // 用户在商户网站的登录账号，用于在签约页面展示，如果为空，则不展示
	PersonalProductCode string             `json:"personal_product_code,omitempty"` // 个人签约产品码，商户和支付宝签约时确定，商户可咨询技 术支持
	SignScene           string             `json:"sign_scene,omitempty"`            // 协议签约场景，商户和支付宝 签约时确定，商户可咨询技术支持。当传入商户签约号external_agreement_no时，场景不能为默认值DEFAULT|DEFAULT
	ExternalAgreementNo string             `json:"external_agreement_no,omitempty"` // 商户签约号，代扣协议中标示 用户的唯一签约号（确保在商 户系统中唯一）。
	ThirdPartyType      string             `json:"third_party_type,omitempty"`      // 签约第三方主体类型。对于三方协议，表示当前用户和哪一类的第三方主体进行签约。1.PARTNER（平台商户）;2.MERCHANT（集团商户）;默认为PARTNER
	ZmAuthParams        *ZmAuthParams      `json:"zm_auth_params,omitempty"`        // 芝麻授权信息，针对于信用代扣签约。json格式。
	ProdParams          *ProdParams        `json:"prod_params,omitempty"`           // 签约产品属性
	PromoParams         string             `json:"promo_params,omitempty"`          // 签约营销参数，此值为json格式；具体的key需与营销约定
	AccessParams        *AccessParams      `json:"access_params,omitempty"`
	SubMerchantParams   *SubMerchantParams `json:"sub_merchant_params,omitempty"`
	DeviceParams        *DeviceParams      `json:"device_params,omitempty"`
	MerchantProcessUrl  string             `json:"merchant_process_url,omitempty"`
	IdentityParams      *IdentityParams    `json:"identity_params,omitempty"`
	AgreementEffectType string             `json:"agreement_effect_type,omitempty"`
	UserAgeRange        string             `json:"user_age_range,omitempty"`
	PeriodRuleParams    *PeriodRuleParams  `json:"period_rule_params,omitempty"`
}

// AgreementPageSignRsp 支付宝个人协议页面签约接口请响应参数
type AgreementPageSignRsp struct {
	Error
	ExternalAgreementNo string `json:"external_agreement_no"` // 代扣协议中标示用户的唯一签约号(确保在商户系统中唯一)
	PersonalProductCode string `json:"personal_product_code"` // 协议产品码，商户和支付宝签约时确定，不同业务场景对应不同的签约产品码
	ValidTime           string `json:"valid_time"`            // 协议生效时间，格式为 yyyyMM-dd HH:mm:ss
	SignScene           string `json:"sign_scene"`            // 签约协议的场景
	AgreementNo         string `json:"agreement_no"`          // 用户签约成功后的协议号
	ZmOpenId            string `json:"zm_open_id"`            // 用户的芝麻信用 openId，供商户查询用户芝麻信用使用。
	InvalidTime         string `json:"invalid_time"`          // 协议失效时间，格式为 yyyyMM-dd HH:mm:ss
	SignTime            string `json:"sign_time"`             // 协议签约时间，格式为 yyyyMM-dd HH:mm:ss
	AlipayUserId        string `json:"alipay_user_id"`        // 用户的支付宝账号对应的支付宝唯一用户号，以2088开头的16位纯数字组成;本参数与alipay_logon_id不可同时为空，若都填写，则以本参数为准，优先级高于alipay_logon_id
	Status              string `json:"status"`                // 协议当前状态 1.TEMP：暂存，协议未生效过；2.NORMAL：正常；3.STOP：暂停
	ForexEligible       string `json:"forex_eligible"`        // 是否海外购汇身份。值：T/F（只有在签约成功时才会返回）
	ExternalLogonId     string `json:"external_logon_id"`     // 外部登录Id
	AlipayLogonId       string `json:"alipay_logon_id"`       // 返回脱敏的支付宝账号
}

func (this AgreementPageSign) APIName() string {
	return "alipay.user.agreement.page.sign"
}

func (this AgreementPageSign) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	m["return_url"] = this.ReturnURL
	return m
}

// AgreementQuery 支付宝个人代扣协议查询接口请求参数 https://opendocs.alipay.com/open/02fkao?scene=8837b4183390497f84bb53783b488ecc
type AgreementQuery struct {
	AuxParam
	AppAuthToken        string `json:"-"`
	PersonalProductCode string `json:"personal_product_code,omitempty"` // 协议产品码，商户和支付宝签约时确定，商户可咨询技术支持
	AlipayUserId        string `json:"alipay_user_id,omitempty"`        // 用户的支付宝账号对应的支付宝唯一用户号，以2088开头的16位纯数字组成;本参数与alipay_logon_id不可同时为空，若都填写，则以本参数为准，优先级高于alipay_logon_id
	AlipayLogonId       string `json:"alipay_logon_id,omitempty"`       // 用户的支付宝登录账号，支持邮箱或手机号码格式。本参数与alipay_user_id不可同时为空，若都填写，则以alipay_user_id为准
	SignScene           string `json:"sign_scene,omitempty"`            // 签约协议场景，商户和支付宝签约时确定，商户可咨询技术支持
	ExternalAgreementNo string `json:"external_agreement_no,omitempty"` // 代扣协议中标示用户的唯一签约号(确保在商户系统中 唯一)。
	ThirdPartyType      string `json:"third_party_type,omitempty"`      // 签约第三方主体类型。对于三方协议，表示当前用户和哪一类的第三方主体进行签约
	AgreementNo         string `json:"agreement_no,omitempty"`          // 支付宝系统中用以唯一标识用户签约记录的编号（用户签约成功后的协议号），如果传了该参数，其他参数会被忽略
}

// AgreementQueryRsp 支付宝个人代扣协议查询接口响应参数
type AgreementQueryRsp struct {
	Error
	PrincipalId         string `json:"principal_id"`          // 签约主体标识。当principal_type为CARD时，该字段为支付宝用户号;当principal_type为CUSTOMER时，该字段为支付宝用户标识。
	ValidTime           string `json:"valid_time"`            // 协议生效时间，格式为 yyyyMM-dd HH:mm:ss
	AlipayLogonId       string `json:"alipay_logon_id"`       // 返回脱敏的支付宝账号
	InvalidTime         string `json:"invalid_time"`          // 协议失效时间，格式为 yyyyMM-dd HH:mm:ss
	PricipalType        string `json:"pricipal_type"`         // 签约主体类型。 CARD:支付宝账号 CUSTOMER:支付宝用户
	DeviceId            string `json:"device_id"`             // 设备Id
	SignScene           string `json:"sign_scene"`            // 签约协议的场景
	AgreementNo         string `json:"agreement_no"`          // 用户签约成功后的协议号
	ThirdPartyType      string `json:"third_party_type"`      // 签约第三方主体类型。对于三方协议，表示当前用户和哪一类的第三方主体进行签约。 1.PARTNER（平台商户）;2.MERCHANT（集团商户），集团下子商户可共享用户签约内容;默认为PARTNER
	Status              string `json:"status"`                // 协议当前状态 1.TEMP：暂存，协议未生效过；2.NORMAL：正常；3.STOP：暂停
	SignTime            string `json:"sign_time"`             // 协议签约时间，格式为 yyyyMM-dd HH:mm:ss
	PersonalProductCode string `json:"personal_product_code"` // 协议产品码，商户和支付宝签约时确定，不同业务场景对应不同的签约产品码
	ExternalAgreementNo string `json:"external_agreement_no"` // 代扣协议中标示用户的唯一签约号(确保在商户系统中唯一)
	ZmOpenId            string `json:"zm_open_id"`            // 用户的芝麻信用 openId，供商户查询用户芝麻信用使用。
	ExternalLogonId     string `json:"external_logon_id"`     // 外部登录Id
	CreditAuthMode      string `json:"credit_auth_mode"`      // 授信模式，取值：DEDUCT_HUAZHI-花芝GO。目前只在花芝代扣（即花芝go）协议时才会返回
	SingleQuota         string `json:"single_quota"`          // 单笔代扣额度
	LastDeductTime      string `json:"last_deduct_time"`      // 周期扣协议，上次扣款成功时间
	NextDeductTime      string `json:"next_deduct_time"`      // 周期扣协议，预计下次扣款时间
}

func (this AgreementQuery) APIName() string {
	return "alipay.user.agreement.query"
}

func (this AgreementQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

// AgreementUnsign 支付宝个人代扣协议解约接口请求参数 https://docs.open.alipay.com/api_2/alipay.user.agreement.unsign
type AgreementUnsign struct {
	AuxParam
	AppAuthToken        string `json:"-"`
	NotifyURL           string `json:"-"`
	AlipayUserId        string `json:"alipay_user_id,omitempty"`        // 用户的支付宝账号对应的支付宝唯一用户号，以2088开头的16位纯数字组成;本参数与alipay_logon_id不可同时为空，若都填写，则以本参数为准，优先级高于alipay_logon_id
	AlipayLogonId       string `json:"alipay_logon_id,omitempty"`       // 返回脱敏的支付宝账号
	PersonalProductCode string `json:"personal_product_code,omitempty"` // 协议产品码，商户和支付宝签约时确定，不同业务场景对应不同的签约产品码
	SignScene           string `json:"sign_scene,omitempty"`            // 签约协议的场景
	ExternalAgreementNo string `json:"external_agreement_no,omitempty"` // 代扣协议中标示用户的唯一签约号(确保在商户系统中唯一)
	ThirdPartyType      string `json:"third_party_type,omitempty"`      // 签约第三方主体类型。对于三方协议，表示当前用户和哪一类的第三方主体进行签约
	AgreementNo         string `json:"agreement_no,omitempty"`          // 支付宝系统中用以唯一标识用户签约记录的编号（用户签约成功后的协议号），如果传了该参数，其他参数会被忽略
	ExtendParams        string `json:"extend_params,omitempty"`
	OperateType         string `json:"operate_type,omitempty"`
}

func (this AgreementUnsign) APIName() string {
	return "alipay.user.agreement.unsign"
}

func (this AgreementUnsign) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// AgreementUnsignRsp 支付宝个人代扣协议解约接口响应参数
type AgreementUnsignRsp struct {
	Error
}

// AgreementExecutionPlanModify 周期性扣款协议执行计划修改接口请求参数 https://docs.open.alipay.com/api_2/alipay.user.agreement.executionplan.modify
type AgreementExecutionPlanModify struct {
	AuxParam
	AppAuthToken string `json:"-"`
	NotifyURL    string `json:"-"`
	AgreementNo  string `json:"agreement_no,omitempty"` // 支付宝系统中用以唯一标识用户签约记录的编号（用户签约成功后的协议号）
	DeductTime   string `json:"deduct_time,omitempty"`  // 商户下一次扣款时间
	Memo         string `json:"memo,omitempty"`         // 具体修改原因
}

func (this AgreementExecutionPlanModify) APIName() string {
	return "alipay.user.agreement.executionplan.modify"
}

func (this AgreementExecutionPlanModify) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["notify_url"] = this.NotifyURL
	return m
}

// AgreementExecutionPlanModifyRsp 周期性扣款协议执行计划修改响应参数
type AgreementExecutionPlanModifyRsp struct {
	Error
	AgreementNo string `json:"agreement_no"`
	DeductTime  string `json:"deduct_time"`
}

// MobileNumber 小程序获取会员手机号  https://opendocs.alipay.com/mini/api/getphonenumber
type MobileNumber struct {
	Code   Code   `json:"code"`
	Msg    string `json:"msg"`
	Mobile string `json:"mobile"`
}
