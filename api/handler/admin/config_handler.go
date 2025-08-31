package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service"
	"geekai/service/moderation"
	"geekai/service/oss"
	"geekai/service/payment"
	"geekai/service/sms"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	handler.BaseHandler
	licenseService    *service.LicenseService
	sysConfig         *types.SystemConfig
	alipayService     *payment.AlipayService
	wxpayService      *payment.WxPayService
	epayService       *payment.EPayService
	smsManager        *sms.SmsManager
	uploaderManager   *oss.UploaderManager
	smtpService       *service.SmtpService
	captchaService    *service.CaptchaService
	wxLoginService    *service.WxLoginService
	moderationManager *moderation.ServiceManager
}

func NewConfigHandler(
	app *core.AppServer,
	db *gorm.DB,
	licenseService *service.LicenseService,
	sysConfig *types.SystemConfig,
	alipayService *payment.AlipayService,
	wxpayService *payment.WxPayService,
	epayService *payment.EPayService,
	smsManager *sms.SmsManager,
	uploaderManager *oss.UploaderManager,
	smtpService *service.SmtpService,
	captchaService *service.CaptchaService,
	wxLoginService *service.WxLoginService,
	moderationManager *moderation.ServiceManager,
) *ConfigHandler {
	return &ConfigHandler{
		BaseHandler:       handler.BaseHandler{App: app, DB: db},
		licenseService:    licenseService,
		sysConfig:         sysConfig,
		alipayService:     alipayService,
		wxpayService:      wxpayService,
		epayService:       epayService,
		smsManager:        smsManager,
		uploaderManager:   uploaderManager,
		moderationManager: moderationManager,
		smtpService:       smtpService,
		captchaService:    captchaService,
		wxLoginService:    wxLoginService,
	}
}

// RegisterRoutes 注册路由
func (h *ConfigHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/config")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.POST("update/base", h.UpdateBase)
		rg.POST("update/power", h.UpdatePower)
		rg.POST("update/notice", h.UpdateNotice)
		rg.POST("update/agreement", h.UpdateAgreement)
		rg.POST("update/privacy", h.UpdatePrivacy)
		rg.POST("update/mark_map", h.UpdateMarkMap)
		rg.POST("update/captcha", h.UpdateCaptcha)
		rg.POST("update/wx_login", h.UpdateWxLogin)
		rg.POST("update/payment", h.UpdatePayment)
		rg.POST("update/sms", h.UpdateSms)
		rg.POST("update/oss", h.UpdateOss)
		rg.POST("update/smtp", h.UpdateStmp)
		rg.POST("update/moderation", h.UpdateModeration)
		rg.POST("moderation/test", h.TestModeration)
		rg.GET("get", h.Get)
		rg.POST("license/active", h.Active)
		rg.GET("license/get", h.GetLicense)
	}
}

// UpdateBase 更新基础配置
func (h *ConfigHandler) UpdateBase(c *gin.Context) {
	var data types.BaseConfig

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 未授权的话不允许修改版权
	license := h.licenseService.GetLicense()
	if !license.IsActive && data.Copyright != h.sysConfig.Base.Copyright {
		resp.ERROR(c, "未授权系统不允许修改版权信息")
		return
	}

	// 未授权的话不允许修改 Logo
	if !license.IsActive && data.Logo != h.sysConfig.Base.Logo {
		resp.ERROR(c, "未授权系统不允许修改 Logo")
		return
	}

	err := h.Update(types.ConfigKeySystem, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	h.sysConfig.Base = data

	resp.SUCCESS(c, data)
}

// UpdatePower 更新系统配置
func (h *ConfigHandler) UpdatePower(c *gin.Context) {
	var data struct {
		InitPower     int            `json:"init_power,omitempty"`      // 新用户注册赠送算力值
		DailyPower    int            `json:"daily_power,omitempty"`     // 每日签到赠送算力
		InvitePower   int            `json:"invite_power,omitempty"`    // 邀请新用户赠送算力值
		MjPower       int            `json:"mj_power,omitempty"`        // MJ 绘画消耗算力
		MjActionPower int            `json:"mj_action_power,omitempty"` // MJ 操作（放大，变换）消耗算力
		SdPower       int            `json:"sd_power,omitempty"`        // SD 绘画消耗算力
		SunoPower     int            `json:"suno_power,omitempty"`      // Suno 生成歌曲消耗算力
		LumaPower     int            `json:"luma_power,omitempty"`      // Luma 生成视频消耗算力
		KeLingPowers  map[string]int `json:"keling_powers,omitempty"`   // 可灵生成视频消耗算力
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	h.sysConfig.Base.InitPower = data.InitPower
	h.sysConfig.Base.DailyPower = data.DailyPower
	h.sysConfig.Base.InvitePower = data.InvitePower
	h.sysConfig.Base.MjPower = data.MjPower
	h.sysConfig.Base.MjActionPower = data.MjActionPower
	h.sysConfig.Base.SdPower = data.SdPower
	h.sysConfig.Base.SunoPower = data.SunoPower
	h.sysConfig.Base.LumaPower = data.LumaPower
	h.sysConfig.Base.KeLingPowers = data.KeLingPowers

	err := h.Update(types.ConfigKeySystem, h.sysConfig.Base)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, h.sysConfig.Base)
}

// UpdateNotice 更新公告配置
func (h *ConfigHandler) UpdateNotice(c *gin.Context) {
	var data struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyNotice, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// UpdateAgreement 更新用户协议配置
func (h *ConfigHandler) UpdateAgreement(c *gin.Context) {
	var data struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyAgreement, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// UpdatePrivacy 更新隐私政策配置
func (h *ConfigHandler) UpdatePrivacy(c *gin.Context) {
	var data struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyPrivacy, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// UpdateMarkMap 更新思维导图配置
func (h *ConfigHandler) UpdateMarkMap(c *gin.Context) {
	var data struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyMarkMap, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// UpdateCaptcha 更新行为验证码配置
func (h *ConfigHandler) UpdateCaptcha(c *gin.Context) {
	var data types.CaptchaConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyCaptcha, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	h.captchaService.UpdateConfig(data)
	resp.SUCCESS(c, data)

}

// UpdatePayment 更新支付配置
func (h *ConfigHandler) UpdatePayment(c *gin.Context) {
	var data types.PaymentConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyPayment, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 如果启用状态发生改变，则需要更新支付服务配置
	if data.WxPay.Enabled {
		err = h.wxpayService.UpdateConfig(&data.WxPay)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	if data.Epay.Enabled {
		h.epayService.UpdateConfig(&data.Epay)
	}
	if data.Alipay.Enabled {
		err = h.alipayService.UpdateConfig(&data.Alipay)
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	h.sysConfig.Payment = data
	resp.SUCCESS(c, data)
}

// UpdateSms 更新短信配置
func (h *ConfigHandler) UpdateSms(c *gin.Context) {
	var data types.SMSConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeySms, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 更新服务配置
	h.smsManager.UpdateConfig(data)

	resp.SUCCESS(c, data)
}

// UpdateOss 更新 Oss 配置
func (h *ConfigHandler) UpdateOss(c *gin.Context) {
	var data types.OSSConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyOss, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 更新服务配置
	h.uploaderManager.UpdateConfig(data)
	h.sysConfig.OSS = data

	resp.SUCCESS(c, data)
}

// UpdateStmp 更新 Stmp 配置
func (h *ConfigHandler) UpdateStmp(c *gin.Context) {
	var data types.SmtpConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeySmtp, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 更新服务配置
	h.smtpService.UpdateConfig(&data)
	h.sysConfig.SMTP = data
	resp.SUCCESS(c, data)
}

// UpdateWxLogin 更新微信登录配置
func (h *ConfigHandler) UpdateWxLogin(c *gin.Context) {
	var data types.WxLoginConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	err := h.Update(types.ConfigKeyWxLogin, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if data.Enabled {
		h.wxLoginService.UpdateConfig(data)
	}

	h.sysConfig.WxLogin = data
	resp.SUCCESS(c, data)
}

// Update 更新系统配置
func (h *ConfigHandler) Update(name string, value any) error {
	var config model.Config
	err := h.DB.Where("name", name).First(&config).Error
	if err != nil { // 不存在则创建
		config.Name = name
		config.Value = utils.JsonEncode(value)
		return h.DB.Create(&config).Error
	} else { // 存在则更新
		config.Value = utils.JsonEncode(value)
		return h.DB.Updates(&config).Error
	}

}

// Get 获取指定名称的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	name := c.Query("key")
	var config model.Config
	res := h.DB.Where("name", name).First(&config)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var value map[string]any
	err := utils.JsonDecode(config.Value, &value)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, value)
}

// Active 激活系统
func (h *ConfigHandler) Active(c *gin.Context) {
	var data struct {
		License string `json:"license"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.licenseService.ActiveLicense(data.License)
	license := h.licenseService.GetLicense()
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	if err := h.Update(types.ConfigKeyLicense, license); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 更新系统配置
	h.sysConfig.License = *license

	resp.SUCCESS(c, license.MachineId)

}

// GetLicense 获取 License 信息
func (h *ConfigHandler) GetLicense(c *gin.Context) {
	license := h.licenseService.GetLicense()
	resp.SUCCESS(c, license)
}

// UpdateModeration 更新文本审查配置
func (h *ConfigHandler) UpdateModeration(c *gin.Context) {
	var data types.ModerationConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.Update(types.ConfigKeyModeration, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	h.moderationManager.UpdateConfig(data)
	h.sysConfig.Moderation = data

	resp.SUCCESS(c, data)
}

// 测试结果类型，用于前端显示
type ModerationTestResult struct {
	IsAbnormal bool                   `json:"isAbnormal"`
	Details    []ModerationTestDetail `json:"details"`
}

type ModerationTestDetail struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Confidence  string `json:"confidence"`
	IsCategory  bool   `json:"isCategory"`
}

// TestModeration 测试文本审查服务
func (h *ConfigHandler) TestModeration(c *gin.Context) {
	var data struct {
		Text    string `json:"text"`
		Service string `json:"service"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Text == "" {
		resp.ERROR(c, "测试文本不能为空")
		return
	}

	// 检查是否启用了文本审查
	if !h.sysConfig.Moderation.Enable {
		resp.ERROR(c, "文本审查服务未启用")
		return
	}

	// 获取当前激活的审核服务
	service := h.moderationManager.GetService()
	// 执行文本审核
	result, err := service.Moderate(data.Text)
	if err != nil {
		resp.ERROR(c, "审核服务调用失败: "+err.Error())
		return
	}

	// 转换为前端需要的格式
	testResult := ModerationTestResult{
		IsAbnormal: result.Flagged,
		Details:    make([]ModerationTestDetail, 0),
	}

	// 构建详细信息
	for category, description := range types.ModerationCategories {
		score := result.CategoryScores[category]
		isCategory := result.Categories[category]

		testResult.Details = append(testResult.Details, ModerationTestDetail{
			Category:    category,
			Description: description,
			Confidence:  fmt.Sprintf("%.2f", score),
			IsCategory:  isCategory,
		})
	}

	resp.SUCCESS(c, testResult)
}
