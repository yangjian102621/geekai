package types

import (
	"fmt"
)

type AppConfig struct {
	Path          string `toml:"-"`
	Listen        string
	Session       Session
	ProxyURL      string
	MysqlDns      string                  // mysql 连接地址
	Manager       Manager                 // 后台管理员账户信息
	StaticDir     string                  // 静态资源目录
	StaticUrl     string                  // 静态资源 URL
	Redis         RedisConfig             // redis 连接信息
	ApiConfig     ChatPlusApiConfig       // ChatPlus API authorization configs
	SMS           SMSConfig               // send mobile message config
	OSS           OSSConfig               // OSS config
	MjConfigs     []MidJourneyConfig      // mj AI draw service pool
	MjPlusConfigs []MidJourneyPlusConfig  // MJ plus config
	WeChatBot     bool                    // 是否启用微信机器人
	SdConfigs     []StableDiffusionConfig // sd AI draw service pool

	XXLConfig     XXLConfig
	AlipayConfig  AlipayConfig
	HuPiPayConfig HuPiPayConfig
	SmtpConfig    SmtpConfig // 邮件发送配置
	JPayConfig    JPayConfig // payjs 支付配置
}

type SmtpConfig struct {
	Host     string
	Port     int
	AppName  string // 应用名称
	From     string // 发件人邮箱地址
	Password string // 发件人邮箱密码
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
	ImgCdnURL      string // 图片反代加速地址
	DiscordAPI     string
	DiscordGateway string
}

type StableDiffusionConfig struct {
	Enabled         bool
	ApiURL          string
	ApiKey          string
	Txt2ImgJsonPath string
}

type MidJourneyPlusConfig struct {
	Enabled   bool   // 如果启用了 MidJourney Plus，将会自动禁用原生的MidJourney服务
	ApiURL    string // api 地址
	Mode      string // 绘画模式，可选值：fast/turbo/relax
	CdnURL    string // CDN 加速地址
	ApiKey    string
	NotifyURL string // 任务进度更新回调地址
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
	ApiURL    string // 支付网关
	NotifyURL string // 异步通知回调
}

// JPayConfig PayJs 支付配置
type JPayConfig struct {
	Enabled    bool
	Name       string // 支付名称，默认 wechat
	AppId      string // 商户 ID
	PrivateKey string // 私钥
	ApiURL     string // API 网关
	NotifyURL  string // 异步回调地址
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

	EnableContext bool `json:"enable_context"` // 是否开启聊天上下文
	EnableHistory bool `json:"enable_history"` // 是否允许保存聊天记录
	ContextDeep   int  `json:"context_deep"`   // 上下文深度
	DallImgNum    int  `json:"dall_img_num"`   // dall-e3 出图数量
}

type Platform string

const OpenAI = Platform("OpenAI")
const Azure = Platform("Azure")
const ChatGLM = Platform("ChatGLM")
const Baidu = Platform("Baidu")
const XunFei = Platform("XunFei")
const QWen = Platform("QWen")

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
}

type SystemConfig struct {
	Title            string `json:"title"`
	AdminTitle       string `json:"admin_title"`
	InitChatCalls    int    `json:"init_chat_calls"`     // 新用户注册赠送对话次数
	InitImgCalls     int    `json:"init_img_calls"`      // 新用户注册赠送绘图次数
	VipMonthCalls    int    `json:"vip_month_calls"`     // VIP 会员每月赠送的对话次数
	VipMonthImgCalls int    `json:"vip_month_img_calls"` // VIP 会员每月赠送绘图次数

	RegisterWays    []string `json:"register_ways"`    // 注册方式：支持手机，邮箱注册
	EnabledRegister bool     `json:"enabled_register"` // 是否开放注册

	RewardImg     string  `json:"reward_img"`      // 众筹收款二维码地址
	EnabledReward bool    `json:"enabled_reward"`  // 启用众筹功能
	ChatCallPrice float64 `json:"chat_call_price"` // 对话单次调用费用
	ImgCallPrice  float64 `json:"img_call_price"`  // 绘图单次调用费用

	OrderPayTimeout  int      `json:"order_pay_timeout"`   //订单支付超时时间
	DefaultModels    []string `json:"default_models"`      // 默认开通的 AI 模型
	OrderPayInfoText string   `json:"order_pay_info_text"` // 订单支付页面说明文字
	InviteChatCalls  int      `json:"invite_chat_calls"`   // 邀请用户注册奖励对话次数
	InviteImgCalls   int      `json:"invite_img_calls"`    // 邀请用户注册奖励绘图次数

	WechatCardURL string `json:"wechat_card_url"` // 微信客服地址
}
