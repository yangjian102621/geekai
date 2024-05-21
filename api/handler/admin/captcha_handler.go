package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/handler"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type CaptchaHandler struct {
	handler.BaseHandler
}

func NewCaptchaHandler(app *core.AppServer) *CaptchaHandler {
	return &CaptchaHandler{BaseHandler: handler.BaseHandler{App: app}}
}

type CaptchaVo struct {
	CaptchaId string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}

// GetCaptcha 获取验证码
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	var captchaVo CaptchaVo
	driver := base64Captcha.NewDriverDigit(48, 130, 4, 0.4, 10)
	cp := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	// b64s是图片的base64编码
	id, b64s, err := cp.Generate()
	if err != nil {
		resp.ERROR(c, "生成验证码错误!")
		return
	}
	captchaVo.CaptchaId = id
	captchaVo.PicPath = b64s

	resp.SUCCESS(c, captchaVo)
}
