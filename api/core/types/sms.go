package types

type SMSConfig struct {
	Active string
	Ali    SmsConfigAli
	Bao    SmsConfigBao
}

// SmsConfigAli 阿里云短信平台配置
type SmsConfigAli struct {
	AccessKey    string
	AccessSecret string
	Product      string
	Domain       string
	Sign         string // 短信签名
	CodeTempId   string // 验证码短信模板 ID
}

// SmsConfigBao 短信宝平台配置
type SmsConfigBao struct {
	Username     string //短信宝平台注册的用户名
	Password     string //短信宝平台注册的密码
	Domain       string //域名
	Sign         string // 短信签名
	CodeTemplate string // 验证码短信模板 匹配
}
