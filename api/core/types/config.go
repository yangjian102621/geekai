package types

import (
	"fmt"
)

type AppConfig struct {
	Path      string `toml:"-"`
	Listen    string
	Session   Session
	ProxyURL  string
	MysqlDns  string                  // mysql 连接地址
	Manager   Manager                 // 后台管理员账户信息
	StaticDir string                  // 静态资源目录
	StaticUrl string                  // 静态资源 URL
	Redis     RedisConfig             // redis 连接信息
	ApiConfig ChatPlusApiConfig       // ChatPlus API authorization configs
	SmsConfig AliYunSmsConfig         // AliYun send message service config
	OSS       OSSConfig               // OSS config
	MjConfigs []MidJourneyConfig      // mj AI draw service pool
	WeChatBot bool                    // 是否启用微信机器人
	SdConfigs []StableDiffusionConfig // sd AI draw service pool

	XXLConfig     XXLConfig
	AlipayConfig  AlipayConfig
	HuPiPayConfig HuPiPayConfig
}

type ChatPlusApiConfig struct {
	ApiURL string
	AppId  string
	Token  string
}

type MidJourneyConfig struct {
	Enabled        bool
	UserToken      string
	BotToken       string
	GuildId        string // Server ID
	ChanelId       string // Chanel ID
	UseCDN         bool
	DiscordAPI     string
	DiscordCDN     string
	DiscordGateway string
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

type AlipayConfig struct {
	Enabled         bool   // 是否启用该支付通道
	SandBox         bool   // 是否沙盒环境
	AppId           string // 应用 ID
	UserId          string // 支付宝用户 ID
	PrivateKey      string // 用户私钥文件路径
	PublicKey       string // 用户公钥文件路径
	AlipayPublicKey string // 支付宝公钥文件路径
	RootCert        string // Root 秘钥路径
	NotifyURL       string // 异步通知回调
}

type HuPiPayConfig struct { //虎皮椒第四方支付配置
	Enabled   bool   // 是否启用该支付通道
	Name      string // 支付名称，如：wechat/alipay
	AppId     string // App ID
	AppSecret string // app 密钥
	NotifyURL string // 异步通知回调
	PayURL    string // 支付网关
}

type XXLConfig struct { // XXL 任务调度配置
	Enabled      bool
	ServerAddr   string
	ExecutorIp   string
	ExecutorPort string
	AccessToken  string
	RegistryKey  string
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

	EnableContext bool   `json:"enable_context"` // 是否开启聊天上下文
	EnableHistory bool   `json:"enable_history"` // 是否允许保存聊天记录
	ContextDeep   int    `json:"context_deep"`   // 上下文深度
	DallApiURL    string `json:"dall_api_url"`   // dall-e3 绘图 API 地址
	DallImgNum    int    `json:"dall_img_num"`   // dall-e3 出图数量
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

type InviteReward struct {
	ChatCalls int `json:"chat_calls"`
	ImgCalls  int `json:"img_calls"`
}

type ModelAPIConfig struct {
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	ApiKey      string  `json:"api_key"`
}

type SystemConfig struct {
	Title            string   `json:"title"`
	AdminTitle       string   `json:"admin_title"`
	Models           []string `json:"models"`
	InitChatCalls    int      `json:"init_chat_calls"`     // 新用户注册赠送对话次数
	InitImgCalls     int      `json:"init_img_calls"`      // 新用户注册赠送绘图次数
	VipMonthCalls    int      `json:"vip_month_calls"`     // VIP 会员每月赠送的对话次数
	VipMonthImgCalls int      `json:"vip_month_img_calls"` // VIP 会员每月赠送绘图次数
	EnabledRegister  bool     `json:"enabled_register"`    // 是否启用注册功能，关闭注册功能之后将无法注册
	EnabledMsg       bool     `json:"enabled_msg"`         // 是否启用短信验证码服务
	RewardImg        string   `json:"reward_img"`          // 众筹收款二维码地址
	EnabledReward    bool     `json:"enabled_reward"`      // 启用众筹功能
	ChatCallPrice    float64  `json:"chat_call_price"`     // 对话单次调用费用
	ImgCallPrice     float64  `json:"img_call_price"`      // 绘图单次调用费用
	EnabledAlipay    bool     `json:"enabled_alipay"`      // 是否启用支付宝支付通道
	OrderPayTimeout  int      `json:"order_pay_timeout"`   //订单支付超时时间
	DefaultModels    []string `json:"default_models"`      // 默认开通的 AI 模型
	OrderPayInfoText string   `json:"order_pay_info_text"` // 订单支付页面说明文字
	InviteChatCalls  int      `json:"invite_chat_calls"`   // 邀请用户注册奖励对话次数
	InviteImgCalls   int      `json:"invite_img_calls"`    // 邀请用户注册奖励绘图次数
	ForceInvite      bool     `json:"force_invite"`        // 是否强制必须使用邀请码才能注册

	ShowDemoNotice bool `json:"show_demo_notice"` // 显示演示站公告
}
