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
	MysqlDns        string       // mysql 连接地址
	StaticDir       string       // 静态资源目录
	StaticUrl       string       // 静态资源 URL
	Redis           RedisConfig  // redis 连接信息
	SMS             SMSConfig    // send mobile message config
	OSS             OSSConfig    // OSS config
	SmtpConfig      SmtpConfig   // 邮件发送配置
	AlipayConfig    AlipayConfig // 支付宝支付渠道配置
	GeekPayConfig   EpayConfig   // GEEK 支付配置
	WechatPayConfig WxPayConfig  // 微信支付渠道配置
	TikaHost        string       // TiKa 服务器地址
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

type BaseConfig struct {
	Title      string `json:"title,omitempty"`       // 网站标题
	Slogan     string `json:"slogan,omitempty"`      // 网站 slogan
	AdminTitle string `json:"admin_title,omitempty"` // 管理后台标题
	Logo       string `json:"logo,omitempty"`        // 圆形 Logo
	BarLogo    string `json:"bar_logo,omitempty"`    // 条形 Logo

	RegisterWays    []string `json:"register_ways,omitempty"`    // 注册方式：支持手机（mobile），邮箱注册（email），账号密码注册
	EnabledRegister bool     `json:"enabled_register,omitempty"` // 是否开放注册

	OrderPayTimeout int `json:"order_pay_timeout,omitempty"` //订单支付超时时间，单位：分钟

	InitPower         int            `json:"init_power,omitempty"`          // 新用户注册赠送算力值
	DailyPower        int            `json:"daily_power,omitempty"`         // 每日签到赠送算力
	InvitePower       int            `json:"invite_power,omitempty"`        // 邀请新用户赠送算力值
	MjPower           int            `json:"mj_power,omitempty"`            // MJ 绘画消耗算力
	MjActionPower     int            `json:"mj_action_power,omitempty"`     // MJ 操作（放大，变换）消耗算力，未配置分项时回退用
	MjUpscalePower    int            `json:"mj_upscale_power,omitempty"`    // MJ 放大/变换消耗算力
	MjBlendPower      int            `json:"mj_blend_power,omitempty"`      // MJ 融图消耗算力
	MjSwapFacePower   int            `json:"mj_swap_face_power,omitempty"`  // MJ 换脸消耗算力
	MjModalPower      int            `json:"mj_modal_power,omitempty"`      // MJ 局部重绘消耗算力
	SunoPower         int            `json:"suno_power,omitempty"`          // Suno 生成歌曲消耗算力
	LumaPower         int            `json:"luma_power,omitempty"`          // Luma 生成视频消耗算力
	KeLingPowers      map[string]int `json:"keling_powers,omitempty"`       // 可灵生成视频消耗算力
	AdvanceVoicePower int            `json:"advance_voice_power,omitempty"` // 高级语音对话消耗算力

	WechatCardURL string `json:"wechat_card_url,omitempty"` // 微信客服地址

	EnableContext bool `json:"enable_context,omitempty"`
	ContextDeep   int  `json:"context_deep,omitempty"`

	MjMode string `json:"mj_mode"` // midjourney 默认的API模式，relax, fast, turbo

	IndexNavs []int  `json:"index_navs"` // 首页显示的导航菜单
	IndexPage string `json:"index_page"` // 首页显示的页面
	Copyright string `json:"copyright"`  // 版权信息
	ICP       string `json:"icp"`        // ICP 备案号
	GaBeian   string `json:"ga_beian"`   // 公安备案号

	EmailWhiteList   []string `json:"email_white_list"`   // 邮箱白名单列表
	AssistantModelId int      `json:"assistant_model_id"` // 用来做提示词,翻译的AI模型 id
	MaxFileSize      int      `json:"max_file_size"`      // 最大文件大小,单位：MB

	EnableMobileSite bool `json:"enable_mobile_site"` // 是否开启手机站点
}

type SystemConfig struct {
	Base       BaseConfig
	Payment    PaymentConfig
	OSS        OSSConfig
	SMS        SMSConfig
	SMTP       SmtpConfig
	Captcha    CaptchaConfig
	WxLogin    WxLoginConfig
	Jimeng     JimengConfig
	Moderation ModerationConfig
	WxGzh      WxGzhConfig
}

// 配置键名常量
const (
	ConfigKeySystem     = "system"     // 系统配置
	ConfigKeyNotice     = "notice"     // 公告配置
	ConfigKeyAgreement  = "agreement"  // 用户协议配置
	ConfigKeyPrivacy    = "privacy"    // 隐私政策配置
	ConfigKeyMarkMap    = "mark_map"   // 水印配置
	ConfigKeyCaptcha    = "captcha"    // 验证码配置
	ConfigKeyWxLogin    = "wx_login"   // 微信扫码登录配置
	ConfigKeyWxGzh      = "wx_gzh"     // 微信公众号配置
	ConfigKeySms        = "sms"        // 短信配置
	ConfigKeySmtp       = "smtp"       // SMTP 配置
	ConfigKeyOss        = "oss"        // OSS 配置
	ConfigKeyPayment    = "payment"    // 支付配置
	ConfigKeyModeration = "moderation" // 文本审查配置
	ConfigKeyAI3D       = "ai3d"       // AI3D 配置
	ConfigKeyJimeng     = "jimeng"     // 即梦AI配置
	ConfigKeyVideo      = "video"      // 视频生成配置
	ConfigKeyPPT        = "ppt"        // PPT 生成配置
)
