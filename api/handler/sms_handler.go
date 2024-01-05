package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strings"
)

const CodeStorePrefix = "/verify/codes/"

type SmsHandler struct {
	BaseHandler
	redis   *redis.Client
	sms     *service.AliYunSmsService
	smtp    *service.SmtpService
	captcha *service.CaptchaService
}

func NewSmsHandler(
	app *core.AppServer,
	client *redis.Client,
	sms *service.AliYunSmsService,
	smtp *service.SmtpService,
	captcha *service.CaptchaService) *SmsHandler {
	handler := &SmsHandler{redis: client, sms: sms, captcha: captcha, smtp: smtp}
	handler.App = app
	return handler
}

// SendCode 发送验证码
func (h *SmsHandler) SendCode(c *gin.Context) {
	var data struct {
		Receiver string `json:"receiver"` // 接收者
		Key      string `json:"key"`
		Dots     string `json:"dots"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if !h.captcha.Check(data) {
		resp.ERROR(c, "验证码错误，请先完人机验证")
		return
	}

	code := utils.RandomNumber(6)
	var err error
	if strings.Contains(data.Receiver, "@") { // email
		err = h.smtp.SendVerifyCode(data.Receiver, code)
	} else {
		err = h.sms.SendVerifyCode(data.Receiver, code)
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
