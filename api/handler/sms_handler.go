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
	redis   *redis.Client
	sms     *sms.ServiceManager
	smtp    *service.SmtpService
	captcha *service.CaptchaService
}

func NewSmsHandler(
	app *core.AppServer,
	client *redis.Client,
	sms *sms.ServiceManager,
	smtp *service.SmtpService,
	captcha *service.CaptchaService) *SmsHandler {
	return &SmsHandler{
		redis:       client,
		sms:         sms,
		captcha:     captcha,
		smtp:        smtp,
		BaseHandler: BaseHandler{App: app}}
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
	var check bool
	if data.X != 0 {
		check = h.captcha.SlideCheck(data)
	} else {
		check = h.captcha.Check(data)
	}
	if !check {
		resp.ERROR(c, "验证码错误，请先完人机验证")
		return
	}

	code := utils.RandomNumber(6)
	var err error
	if strings.Contains(data.Receiver, "@") { // email
		if !utils.Contains(h.App.SysConfig.RegisterWays, "email") {
			resp.ERROR(c, "系统已禁用邮箱注册！")
			return
		}
		err = h.smtp.SendVerifyCode(data.Receiver, code)
	} else {
		if !utils.Contains(h.App.SysConfig.RegisterWays, "mobile") {
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

	if h.App.Debug {
		resp.SUCCESS(c, code)
	} else {
		resp.SUCCESS(c)
	}
}
