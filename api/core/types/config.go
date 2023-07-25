package types

import (
	"fmt"
	"net/http"
)

type AppConfig struct {
	Path              string `toml:"-"`
	Listen            string
	Session           Session
	ProxyURL          string
	MysqlDns          string            // mysql 连接地址
	Manager           Manager           // 后台管理员账户信息
	StaticDir         string            // 静态资源目录
	StaticUrl         string            // 静态资源 URL
	Redis             RedisConfig       // redis 连接信息
	Func              FunctionApiConfig // function api configs
	AesEncryptKey     string
	SmsConfig         AliYunSmsConfig // 短信发送配置
	StartWechatBot    bool            // 是否启动微信机器人
	EnabledMsgService bool            // 是否启用短信服务
}

type AliYunSmsConfig struct {
	AccessKey    string
	AccessSecret string
	Product      string
	Domain       string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

func (c RedisConfig) Url() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Manager 管理员
type Manager struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionDriver string

const (
	SessionDriverMem    = SessionDriver("mem")
	SessionDriverRedis  = SessionDriver("redis")
	SessionDriverCookie = SessionDriver("cookie")
)

// Session configs struct
type Session struct {
	Driver    SessionDriver // session 存储驱动 mem|cookie|redis
	SecretKey string        // session encryption key
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
	ApiURL        string  `json:"api_url,omitempty"`
	Model         string  `json:"model"` // 默认模型
	Temperature   float32 `json:"temperature"`
	MaxTokens     int     `json:"max_tokens"`
	EnableContext bool    `json:"enable_context"` // 是否开启聊天上下文
	EnableHistory bool    `json:"enable_history"` // 是否允许保存聊天记录
	ApiKey        string  `json:"api_key"`
	ContextDeep   int     `json:"context_deep"` // 上下文深度

}

type SystemConfig struct {
	Title         string   `json:"title"`
	AdminTitle    string   `json:"admin_title"`
	Models        []string `json:"models"`
	UserInitCalls int      `json:"user_init_calls"` // 新用户注册默认总送多少次调用
}

type FunctionApiConfig struct {
	ApiURL string
	AppId  string
	Token  string
}

const UserInitCalls = 1000
