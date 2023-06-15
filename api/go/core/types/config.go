package types

import (
	"net/http"
)

type AppConfig struct {
	Path     string `toml:"-"`
	Listen   string
	Session  Session
	ProxyURL string
	MysqlDns string  // mysql 连接地址
	Manager  Manager // 后台管理员账户信息
}

// Manager 管理员
type Manager struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

// ChatConfig 系统默认的聊天配置
type ChatConfig struct {
	ApiURL        string  `json:"api_url"`
	Model         string  `json:"model"` // 默认模型
	Temperature   float32 `json:"temperature"`
	MaxTokens     int     `json:"max_tokens"`
	EnableContext bool    `json:"enable_context"` // 是否开启聊天上下文
	EnableHistory bool    `json:"enable_history"` // 是否允许保存聊天记录
	ApiKey        string  `json:"api_key"`        // OpenAI  API key
}

type SystemConfig struct {
	Title      string   `json:"title"`
	AdminTitle string   `json:"admin_title"`
	Models     []string `json:"models"`
}

var GptModels = []string{"gpt-3.5-turbo", "gpt-3.5-turbo-16k", "gpt-3.5-turbo-0613", "gpt-3.5-turbo-16k-0613", "gpt-4", "gpt-4-0613", "gpt-4-32k", "gpt-4-32k-0613"}
