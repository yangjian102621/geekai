package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/store"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
)

const CodeStorePrefix = "/verify/codes/"

type SmsHandler struct {
	BaseHandler
	leveldb *store.LevelDB
	sms     *service.AliYunSmsService
	captcha *service.CaptchaService
}

func NewSmsHandler(app *core.AppServer, db *store.LevelDB, sms *service.AliYunSmsService, captcha *service.CaptchaService) *SmsHandler {
	handler := &SmsHandler{leveldb: db, sms: sms, captcha: captcha}
	handler.App = app
	return handler
}

// SendCode 发送验证码短信
func (h *SmsHandler) SendCode(c *gin.Context) {
	var data struct {
		Mobile string `json:"mobile"`
		Key    string `json:"key"`
		Dots   string `json:"dots"`
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
	err := h.sms.SendVerifyCode(data.Mobile, code)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 存储验证码，等待后面注册验证
	err = h.leveldb.Put(CodeStorePrefix+data.Mobile, code)
	if err != nil {
		resp.ERROR(c, "验证码保存失败")
		return
	}

	resp.SUCCESS(c)
}

type statusVo struct {
	EnabledMsgService bool `json:"enabled_msg_service"`
	EnabledRegister   bool `json:"enabled_register"`
}

// Status check if the message service is enabled
func (h *SmsHandler) Status(c *gin.Context) {
	resp.SUCCESS(c, statusVo{EnabledMsgService: h.App.SysConfig.EnabledMsg, EnabledRegister: h.App.SysConfig.EnabledRegister})
}
