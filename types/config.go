package types

import (
	"net/http"
)

type Config struct {
	Listen     string
	Session    Session
	ProxyURL   []string
	EnableAuth bool   // 是否开启鉴权
	AccessKey  string // 管理员访问 AccessKey, 通过传入这个参数可以访问系统管理 API
	Chat       Chat
}

type User struct {
	Name           string `json:"name"`
	MaxCalls       int    `json:"max_calls"`       // 最多调用次数，如果为 0 则表示不限制
	RemainingCalls int    `json:"remaining_calls"` // 剩余调用次数
	EnableHistory  bool   `json:"enable_history"`  // 是否启用聊天记录
}

// Chat configs struct
type Chat struct {
	ApiURL                string
	ApiKeys               []string
	Model                 string
	Temperature           float32
	MaxTokens             int
	EnableContext         bool // 是否保持聊天上下文
	ChatContextExpireTime int  // 聊天上下文过期时间，单位：秒
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
