package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
)

type AppConfig struct {
	Path            string `toml:"-"`
	Listen          string
	Session         Session
	AdminSession    Session
	ProxyURL        string
	MysqlDns        string      // mysql 连接地址
	StaticDir       string      // 静态资源目录
	StaticUrl       string      // 静态资源 URL
	Redis           RedisConfig // redis 连接信息
	ApiConfig       ApiConfig   // ChatPlus API authorization configs
	SMS             SMSConfig   // send mobile message config
	OSS             OSSConfig   // OSS config
	SmtpConfig      SmtpConfig  // 邮件发送配置
	XXLConfig       XXLConfig
	AlipayConfig    AlipayConfig    // 支付宝支付渠道配置
	HuPiPayConfig   HuPiPayConfig   // 虎皮椒支付配置
	GeekPayConfig   GeekPayConfig   // GEEK 支付配置
	WechatPayConfig WechatPayConfig // 微信支付渠道配置
	TikaHost        string          // TiKa 服务器地址
}

type SmtpConfig struct {
	UseTls   bool // 是否使用 TLS 发送
	Host     string
	Port     int
	AppName  string // 应用名称
	From     string // 发件人邮箱地址
	Password string // 发件人邮箱密码
}

type ApiConfig struct {
	ApiURL string
	AppId  string
	Token  string
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
	NotifyURL       string // 异步通知地址
	ReturnURL       string // 同步通知地址
}

type WechatPayConfig struct {
	Enabled    bool   // 是否启用该支付通道
	AppId      string // 公众号的APPID,如：wxd678efh567hg6787
	MchId      string // 直连商户的商户号，由微信支付生成并下发
	SerialNo   string // 商户证书的证书序列号
	PrivateKey string // 用户私钥文件路径
	ApiV3Key   string // API V3 秘钥
	NotifyURL  string // 异步通知地址
}

type HuPiPayConfig struct { //虎皮椒第四方支付配置
	Enabled   bool   // 是否启用该支付通道
	AppId     string // App ID
	AppSecret string // app 密钥
	ApiURL    string // 支付网关
	NotifyURL string // 异步通知地址
	ReturnURL string // 同步通知地址
}

// GeekPayConfig GEEK支付配置
type GeekPayConfig struct {
	Enabled    bool
	AppId      string   // 商户 ID
	PrivateKey string   // 私钥
	ApiURL     string   // API 网关
	NotifyURL  string   // 异步通知地址
	ReturnURL  string   // 同步通知地址
	Methods    []string // 支付方式
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

// LicenseKey 存储许可证书的 KEY
const LicenseKey = "Geek-AI-License"

type License struct {
	Key       string        `json:"key"`        // 许可证书密钥
	MachineId string        `json:"machine_id"` // 机器码
	ExpiredAt int64         `json:"expired_at"` // 过期时间
	IsActive  bool          `json:"is_active"`  // 是否激活
	Configs   LicenseConfig `json:"configs"`
}

type LicenseConfig struct {
	UserNum int  `json:"user_num"` // 用户数量
	DeCopy  bool `json:"de_copy"`  // 去版权
}

func (c RedisConfig) Url() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type SystemConfig struct {
	Title         string `json:"title,omitempty"`           // 网站标题
	Slogan        string `json:"slogan,omitempty"`          // 网站 slogan
	AdminTitle    string `json:"admin_title,omitempty"`     // 管理后台标题
	Logo          string `json:"logo,omitempty"`            // 方形 Logo
	InitPower     int    `json:"init_power,omitempty"`      // 新用户注册赠送算力值
	DailyPower    int    `json:"daily_power,omitempty"`     // 每日赠送算力
	InvitePower   int    `json:"invite_power,omitempty"`    // 邀请新用户赠送算力值
	VipMonthPower int    `json:"vip_month_power,omitempty"` // VIP 会员每月赠送的算力值

	RegisterWays    []string `json:"register_ways,omitempty"`    // 注册方式：支持手机（mobile），邮箱注册（email），账号密码注册
	EnabledRegister bool     `json:"enabled_register,omitempty"` // 是否开放注册

	OrderPayTimeout int    `json:"order_pay_timeout,omitempty"` //订单支付超时时间
	VipInfoText     string `json:"vip_info_text,omitempty"`     // 会员页面充值说明

	MjPower       int `json:"mj_power,omitempty"`        // MJ 绘画消耗算力
	MjActionPower int `json:"mj_action_power,omitempty"` // MJ 操作（放大，变换）消耗算力
	SdPower       int `json:"sd_power,omitempty"`        // SD 绘画消耗算力
	DallPower     int `json:"dall_power,omitempty"`      // DALL-E-3 绘图消耗算力
	SunoPower     int `json:"suno_power,omitempty"`      // Suno 生成歌曲消耗算力
	LumaPower     int `json:"luma_power,omitempty"`      // Luma 生成视频消耗算力

	WechatCardURL string `json:"wechat_card_url,omitempty"` // 微信客服地址

	EnableContext bool `json:"enable_context,omitempty"`
	ContextDeep   int  `json:"context_deep,omitempty"`

	SdNegPrompt string `json:"sd_neg_prompt"` // SD 默认反向提示词
	MjMode      string `json:"mj_mode"`       // midjourney 默认的API模式，relax, fast, turbo

	IndexBgURL  string `json:"index_bg_url"`  // 前端首页背景图片
	IndexNavs   []int  `json:"index_navs"`    // 首页显示的导航菜单
	Copyright   string `json:"copyright"`     // 版权信息
	MarkMapText string `json:"mark_map_text"` // 思维导入的默认文本

	EnabledVerify    bool     `json:"enabled_verify"`     // 是否启用验证码
	EmailWhiteList   []string `json:"email_white_list"`   // 邮箱白名单列表
	TranslateModelId int      `json:"translate_model_id"` // 用来做提示词翻译的大模型 id
}
