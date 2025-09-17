package types

type PaymentConfig struct {
	Alipay AlipayConfig `json:"alipay,omitempty"` // 支付宝支付渠道配置
	Epay   EpayConfig   `json:"epay,omitempty"`   // 易支付配置
	WxPay  WxPayConfig  `json:"wxpay,omitempty"`  // 微信支付渠道配置
}

// AlipayConfig 支付宝支付配置
type AlipayConfig struct {
	Enabled         bool   `json:"enabled,omitempty"`           // 是否启用该支付通道
	SandBox         bool   `json:"sandbox,omitempty"`           // 是否沙盒环境
	AppId           string `json:"app_id,omitempty"`            // 应用 ID
	PrivateKey      string `json:"private_key,omitempty"`       // 应用私钥
	AlipayPublicKey string `json:"alipay_public_key,omitempty"` // 支付宝公钥
	Domain          string `json:"domain,omitempty"`            // 支付回调域名
}

func (c *AlipayConfig) Equal(other *AlipayConfig) bool {
	return c.AppId == other.AppId &&
		c.PrivateKey == other.PrivateKey &&
		c.AlipayPublicKey == other.AlipayPublicKey &&
		c.Domain == other.Domain
}

// WxPayConfig 微信支付配置
type WxPayConfig struct {
	Enabled    bool   `json:"enabled,omitempty"`     // 是否启用该支付通道
	AppId      string `json:"app_id,omitempty"`      // 公众号的APPID,如：wxd678efh567hg6787
	MchId      string `json:"mch_id,omitempty"`      // 直连商户的商户号，由微信支付生成并下发
	SerialNo   string `json:"serial_no,omitempty"`   // 商户证书的证书序列号
	PrivateKey string `json:"private_key,omitempty"` // 商户证书私钥
	ApiV3Key   string `json:"api_v3_key,omitempty"`  // API V3 秘钥
	Domain     string `json:"domain,omitempty"`      // 支付回调域名
}

func (c *WxPayConfig) Equal(other *WxPayConfig) bool {
	return c.AppId == other.AppId &&
		c.MchId == other.MchId &&
		c.SerialNo == other.SerialNo &&
		c.PrivateKey == other.PrivateKey &&
		c.ApiV3Key == other.ApiV3Key &&
		c.Domain == other.Domain
}

// EpayConfig 易支付配置
type EpayConfig struct {
	Enabled    bool   `json:"enabled,omitempty"`     // 是否启用该支付通道
	AppId      string `json:"app_id,omitempty"`      // 商户 ID
	PrivateKey string `json:"private_key,omitempty"` // 私钥
	ApiURL     string `json:"api_url,omitempty"`     // z支付 API 网关
	Domain     string `json:"domain,omitempty"`      // 支付回调域名
}

func (c *EpayConfig) Equal(other *EpayConfig) bool {
	return c.AppId == other.AppId &&
		c.PrivateKey == other.PrivateKey &&
		c.ApiURL == other.ApiURL &&
		c.Domain == other.Domain
}
