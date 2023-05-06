package types

import (
	"net/http"
)

type Config struct {
	Title        string
	ConsoleTitle string
	Listen       string
	Session      Session
	ProxyURL     []string
	ImgURL       ImgURL  // 各种图片资源链接地址，比如微信二维码，群二维码
	AccessKey    string  // 管理员访问 AccessKey, 通过传入这个参数可以访问系统管理 API
	Manager      Manager // 后台管理员账户信息
	Chat         Chat
}

type User struct {
	Name           string         `json:"name"`
	MaxCalls       int            `json:"max_calls"`         // 最多调用次数，如果为 0 则表示不限制
	RemainingCalls int            `json:"remaining_calls"`   // 剩余调用次数
	EnableHistory  bool           `json:"enable_history"`    // 是否启用聊天记录
	Status         bool           `json:"status"`            // 当前状态
	Term           int            `json:"term" default:"30"` // 会员有效期，单位：天
	ActiveTime     int64          `json:"active_time"`       // 激活时间
	ExpiredTime    int64          `json:"expired_time"`      // 到期时间
	ApiKey         string         `json:"api_key"`           // OpenAI  API KEY
	ChatRoles      map[string]int `json:"chat_roles"`        // 当前用户已订阅的聊天角色 map[role_key] => 0/1
}

// Manager 管理员
type Manager struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Chat configs struct
type Chat struct {
	ApiURL                string
	Model                 string
	Temperature           float32
	MaxTokens             int
	EnableContext         bool // 是否保持聊天上下文
	ChatContextExpireTime int  // 聊天上下文过期时间，单位：秒
	ApiKeys               []APIKey
}

type APIKey struct {
	Value    string `json:"value"`     // Key value
	LastUsed int64  `json:"last_used"` // 最后使用时间
}

// Session configs struct
type Session struct {
	SecretKey string // session encryption key
	Name      string
	Path      string
	Domain    string
	MaxAge    int
	Secure    bool
	HttpOnly  bool
	SameSite  http.SameSite
}

type ImgURL struct {
	WechatCard  string `json:"wechat_card"`  // 个人微信二维码
	WechatGroup string `json:"wechat_group"` // 微信群二维码
}
