package alipay

import (
	"encoding/json"
	"fmt"
)

const (
	kPastSandboxGateway    = "https://openapi.alipaydev.com/gateway.do"
	kNewSandboxGateway     = "https://openapi-sandbox.dl.alipaydev.com/gateway.do"
	kProductionGateway     = "https://openapi.alipay.com/gateway.do"
	kProductionMAPIGateway = "https://mapi.alipay.com/gateway.do"

	kFormat       = "JSON"
	kCharset      = "utf-8"
	kVersion      = "1.0"
	kSignTypeRSA2 = "RSA2"
	kContentType  = "application/x-www-form-urlencoded;charset=utf-8"
	kTimeFormat   = "2006-01-02 15:04:05"
)

const (
	kResponseSuffix        = "_response"
	kErrorResponse         = "error_response"
	kFieldAppId            = "app_id"
	kFieldMethod           = "method"
	kFieldFormat           = "format"
	kFieldCharset          = "charset"
	kFieldSign             = "sign"
	kFieldSignType         = "sign_type"
	kFieldTimestamp        = "timestamp"
	kFieldVersion          = "version"
	kFieldBizContent       = "biz_content"
	kFieldAppCertSN        = "app_cert_sn"
	kFieldEncryptType      = "encrypt_type"
	kFieldAliPayRootCertSN = "alipay_root_cert_sn"
	kFieldAlyPayCertSN     = "alipay_cert_sn"
	kCertificateBegin      = "-----BEGIN CERTIFICATE-----"
	kCertificateEnd        = "-----END CERTIFICATE-----"
)

// Code 支付宝接口响应错误码 https://doc.open.alipay.com/docs/doc.htm?treeId=291&articleId=105806&docType=1
type Code string

func (c Code) IsSuccess() bool {
	return c == CodeSuccess
}

func (c Code) IsFailure() bool {
	return c != CodeSuccess
}

// 公共错误码 https://opendocs.alipay.com/common/02km9f#API%20%E5%85%AC%E5%85%B1%E9%94%99%E8%AF%AF%E7%A0%81
const (
	CodeSuccess                Code = "10000" // 接口调用成功
	CodeUnknowError            Code = "20000" // 服务不可用
	CodeInvalidAuthToken       Code = "20001" // 授权权限不足
	CodeMissingParam           Code = "40001" // 缺少必选参数
	CodeInvalidParam           Code = "40002" // 非法的参数
	CodeInsufficientConditions Code = "40003" // 条件异常
	CodeBusinessFailed         Code = "40004" // 业务处理失败
	CodeCallLimited            Code = "40005" // 调用频次超限
	CodePermissionDenied       Code = "40006" // 权限不足
)

type Param interface {
	// APIName 用于提供访问的 method，即接口名称
	APIName() string

	// Params 公共请求参数
	Params() map[string]string

	// FileParams 文件参数
	FileParams() map[string]*FileItem

	// NeedEncrypt 该接口是否支持内容加密，有的接口不支持内容加密，比如文件上传接口：alipay.open.file.upload
	NeedEncrypt() bool

	// NeedVerify 是否对支付宝接口返回的数据进行签名验证， 为了安全建议都需要对签名进行验证，本方法存在是因为部分接口不支持签名验证。
	// 比如：应用支付宝公钥证书下载 https://opendocs.alipay.com/common/06ue2z
	NeedVerify() bool
}

type AuxParam struct {
}

func (this AuxParam) FileParams() map[string]*FileItem {
	return nil
}

func (this AuxParam) NeedEncrypt() bool {
	return true
}

func (this AuxParam) NeedVerify() bool {
	return true
}

type FileItem struct {
	Name     string
	Filename string
	Filepath string
}

type Payload struct {
	method  string                 // 接口名称
	Encrypt bool                   // 是否进行内容加密
	Verify  bool                   // 是否验证签名
	param   map[string]string      // 请求参数
	biz     map[string]interface{} // biz_content 请求参数
	files   map[string]*FileItem   // 文件参数
}

func NewPayload(method string) *Payload {
	var nPayload = &Payload{}
	nPayload.method = method
	nPayload.Encrypt = true
	nPayload.Verify = true
	nPayload.param = make(map[string]string)
	nPayload.biz = make(map[string]interface{})
	return nPayload
}

func (this *Payload) APIName() string {
	return this.method
}

func (this *Payload) Params() map[string]string {
	return this.param
}

func (this *Payload) FileParams() map[string]*FileItem {
	return this.files
}

func (this *Payload) NeedEncrypt() bool {
	return this.Encrypt
}

func (this *Payload) NeedVerify() bool {
	return this.Verify
}

// AddParam 添加公共请求参数。
//
// 这里添加的参数一般为支付宝接口文档中的【公共请求参数】，参考：https://opendocs.alipay.com/apis/api_1/alipay.trade.query/#%E5%85%AC%E5%85%B1%E8%AF%B7%E6%B1%82%E5%8F%82%E6%95%B0。
//
// 一般情况下，不需要调用本方法添加公共请求参数，因为公共参数基本都是必须且其值相对固定，都已处理。除了个别参数，如：app_auth_token。
func (this *Payload) AddParam(key, value string) *Payload {
	this.param[key] = value
	return this
}

// AddBizField 添加请求参数 biz_content 的字段，这里添加的信息会序列化成 JSON 字符串，然后通过 biz_content 参数传递。
//
// 这里添加的参数一般为支付宝接口文档中的【请求参数】，参考：https://opendocs.alipay.com/apis/api_1/alipay.trade.query/#%E8%AF%B7%E6%B1%82%E5%8F%82%E6%95%B0。
//
// 一般情况下，支付宝接口文档中的【请求参数】都是通过调用本方法添加。但是也有例外，如 https://opendocs.alipay.com/mini/05snwo#%E8%AF%B7%E6%B1%82%E5%8F%82%E6%95%B0 中的【请求参数】就需要调用 AddParam 进行添加，因为其【公共请求参数】中没有 biz_content 字段。
func (this *Payload) AddBizField(key string, value interface{}) *Payload {
	this.biz[key] = value
	return this
}

// Set 参考 AddBizField。
//
// Deprecated: use AddBizField instead.
func (this *Payload) Set(key string, value interface{}) *Payload {
	this.biz[key] = value
	return this
}

// AddFile 添加需要上传的文件。
//
// name: 参数名称。
//
// filename: 文件名称。
//
// filepath: 本地文件完整路径。
func (this *Payload) AddFile(name, filename, filepath string) {
	if this.files == nil {
		this.files = make(map[string]*FileItem)
	}

	if filename == "" {
		filename = name
	}
	this.files[name] = &FileItem{Name: name, Filename: filename, Filepath: filepath}
}

func (this *Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.biz)
}

type Error struct {
	Code    Code   `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (this Error) Error() string {
	return fmt.Sprintf("%s - %s", this.Code, this.SubMsg)
}

func (this Error) IsSuccess() bool {
	return this.Code.IsSuccess()
}

func (this Error) IsFailure() bool {
	return this.Code.IsFailure()
}

// CertDownload 应用支付宝公钥证书下载 https://opendocs.alipay.com/common/06ue2z
type CertDownload struct {
	AuxParam
	AppAuthToken string `json:"-"`              // 可选
	AliPayCertSN string `json:"alipay_cert_sn"` // 支付宝公钥证书序列号
}

func (this CertDownload) APIName() string {
	return "alipay.open.app.alipaycert.download"
}

func (this CertDownload) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

func (this CertDownload) NeedEncrypt() bool {
	return false
}

func (this CertDownload) NeedVerify() bool {
	return false
}

type CertDownloadRsp struct {
	Error
	AliPayCertContent string `json:"alipay_cert_content"`
}
