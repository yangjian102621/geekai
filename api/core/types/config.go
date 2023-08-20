package types

import (
	"fmt"
	"net/http"
)

type AppConfig struct {
	Path          string `toml:"-"`
	Listen        string
	Session       Session
	ProxyURL      string
	MysqlDns      string            // mysql 连接地址
	Manager       Manager           // 后台管理员账户信息
	StaticDir     string            // 静态资源目录
	StaticUrl     string            // 静态资源 URL
	Redis         RedisConfig       // redis 连接信息
	ApiConfig     ChatPlusApiConfig // ChatPlus API authorization configs
	AesEncryptKey string
	SmsConfig     AliYunSmsConfig   // AliYun send message service config
	ExtConfig     ChatPlusExtConfig // ChatPlus extensions callback api config
	MinioConfig   MinioConfig
}

type ChatPlusApiConfig struct {
	ApiURL string
	AppId  string
	Token  string
}

type ChatPlusExtConfig struct {
	ApiURL string
	Token  string
}

type AliYunSmsConfig struct {
	AccessKey    string
	AccessSecret string
	Product      string
	Domain       string
}

type MinioConfig struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	Bucket       string
	UseSSL       bool
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
	Title             string   `json:"title"`
	AdminTitle        string   `json:"admin_title"`
	Models            []string `json:"models"`
	UserInitCalls     int      `json:"user_init_calls"` // 新用户注册默认总送多少次调用
	InitImgCalls      int      `json:"init_img_calls"`
	EnabledRegister   bool     `json:"enabled_register"`
	EnabledMsgService bool     `json:"enabled_msg_service"`
}
