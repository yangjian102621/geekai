package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/sms"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const CodeStorePrefix = "/verify/codes/"

type SmsHandler struct {
	BaseHandler
	redis          *redis.Client
	sms            *sms.SmsManager
	smtp           *service.SmtpService
	captchaService *service.CaptchaService
}

func NewSmsHandler(
	app *core.AppServer,
	client *redis.Client,
	sms *sms.SmsManager,
	smtp *service.SmtpService,
	captcha *service.CaptchaService) *SmsHandler {
	return &SmsHandler{
		redis:          client,
		sms:            sms,
		captchaService: captcha,
		smtp:           smtp,
		BaseHandler:    BaseHandler{App: app}}
}

// RegisterRoutes 注册路由
func (h *SmsHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/sms/")
	// 无需授权的接口
	group.POST("code", h.SendCode)
}

// SendCode 发送验证码
func (h *SmsHandler) SendCode(c *gin.Context) {
	var data struct {
		Receiver string `json:"receiver"` // 接收者
		Key      string `json:"key"`
		Dots     string `json:"dots,omitempty"`
		X        int    `json:"x,omitempty"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if h.captchaService.GetConfig().Enabled {
		var check bool
		if data.X != 0 {
			check = h.captchaService.SlideCheck(data)
		} else {
			check = h.captchaService.Check(data)
		}
		if !check {
			resp.ERROR(c, "请先完人机验证")
			return
		}
	}

	code := utils.RandomNumber(6)
	var err error
	if strings.Contains(data.Receiver, "@") { // email
		if !utils.Contains(h.App.SysConfig.Base.RegisterWays, "email") {
			resp.ERROR(c, "系统已禁用邮箱注册！")
			return
		}
		// 检查邮箱后缀是否在白名单
		if len(h.App.SysConfig.Base.EmailWhiteList) > 0 {
			inWhiteList := false
			for _, suffix := range h.App.SysConfig.Base.EmailWhiteList {
				if strings.HasSuffix(data.Receiver, suffix) {
					inWhiteList = true
					break
				}
			}
			if !inWhiteList {
				resp.ERROR(c, "邮箱后缀不在白名单中")
				return
			}
		}
		err = h.smtp.SendVerifyCode(data.Receiver, code)
	} else {
		if !utils.Contains(h.App.SysConfig.Base.RegisterWays, "mobile") {
			resp.ERROR(c, "系统已禁用手机号注册！")
			return
		}
		err = h.sms.GetService().SendVerifyCode(data.Receiver, code)

	}
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 存储验证码，等待后面注册验证
	_, err = h.redis.Set(c, CodeStorePrefix+data.Receiver, code, 0).Result()
	if err != nil {
		resp.ERROR(c, "验证码保存失败")
		return
	}

	resp.SUCCESS(c)
}
