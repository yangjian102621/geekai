package types

type SMSConfig struct {
	Active string
	ALI    AliYunSmsConfig
	SMSBAO SmsBaoSmsConfig
}

type AliYunSmsConfig struct {
	AccessKey    string
	AccessSecret string
	Product      string
	Domain       string
	Sign         string // 短信签名
	CodeTempId   string // 验证码短信模板 ID
}

type SmsBaoSmsConfig struct {
	Account      string //短信包平台注册的用户名
	ApiKey       string //apiKey
	Domain       string //域名
	Sign         string // 短信签名
	CodeTemplate string // 验证码短信模板 匹配
	Num          string // 实效性
}
