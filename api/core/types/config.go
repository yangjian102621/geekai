package types

import (
	"fmt"
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
	SmsConfig     AliYunSmsConfig       // AliYun send message service config
	OSS           OSSConfig             // OSS config
	MjConfig      MidJourneyConfig      // mj 绘画配置
	WeChatBot     bool                  // 是否启用微信机器人
	SdConfig      StableDiffusionConfig // sd 绘画配置
}

type ChatPlusApiConfig struct {
	ApiURL string
	AppId  string
	Token  string
}

type MidJourneyConfig struct {
	Enabled   bool
	UserToken string
	BotToken  string
	GuildId   string // Server ID
	ChanelId  string // Chanel ID
}

type WeChatConfig struct {
	Enabled bool
}

type StableDiffusionConfig struct {
	Enabled         bool
	ApiURL          string
	ApiKey          string
	Txt2ImgJsonPath string
}

type AliYunSmsConfig struct {
	AccessKey    string
	AccessSecret string
	Product      string
	Domain       string
	Sign         string // 短信签名
	CodeTempId   string // 验证码短信模板 ID
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func (c RedisConfig) Url() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Manager 管理员
type Manager struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChatConfig 系统默认的聊天配置
type ChatConfig struct {
	OpenAI  ModelAPIConfig `json:"open_ai"`
	Azure   ModelAPIConfig `json:"azure"`
	ChatGML ModelAPIConfig `json:"chat_gml"`
	Baidu   ModelAPIConfig `json:"baidu"`
	XunFei  ModelAPIConfig `json:"xun_fei"`

	EnableContext bool `json:"enable_context"` // 是否开启聊天上下文
	EnableHistory bool `json:"enable_history"` // 是否允许保存聊天记录
	ContextDeep   int  `json:"context_deep"`   // 上下文深度
}

type Platform string

const OpenAI = Platform("OpenAI")
const Azure = Platform("Azure")
const ChatGLM = Platform("ChatGLM")
const Baidu = Platform("Baidu")
const XunFei = Platform("XunFei")

// UserChatConfig 用户的聊天配置
type UserChatConfig struct {
	ApiKeys map[Platform]string `json:"api_keys"`
}

type ModelAPIConfig struct {
	ApiURL      string  `json:"api_url,omitempty"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	ApiKey      string  `json:"api_key"`
}

type SystemConfig struct {
	Title           string   `json:"title"`
	AdminTitle      string   `json:"admin_title"`
	Models          []string `json:"models"`
	UserInitCalls   int      `json:"user_init_calls"` // 新用户注册默认总送多少次调用
	InitImgCalls    int      `json:"init_img_calls"`
	VipMonthCalls   int      `json:"vip_month_calls"` // 会员每个赠送的调用次数
	EnabledRegister bool     `json:"enabled_register"`
	EnabledMsg      bool     `json:"enabled_msg"`      // 启用短信验证码服务
	EnabledDraw     bool     `json:"enabled_draw"`     // 启动 AI 绘画功能
	RewardImg       string   `json:"reward_img"`       // 众筹收款二维码地址
	EnabledFunction bool     `json:"enabled_function"` // 启用 API 函数功能
	EnabledReward   bool     `json:"enabled_reward"`   // 启用众筹功能
	DefaultModels   []string `json:"default_models"`   // 默认开通的 AI 模型
}
