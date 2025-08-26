package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
)

// 今日头条函数实现

type CaptchaHandler struct {
	App     *core.AppServer
	service *service.CaptchaService
	config  types.CaptchaConfig
}

func NewCaptchaHandler(app *core.AppServer, s *service.CaptchaService, sysConfig *types.SystemConfig) *CaptchaHandler {
	return &CaptchaHandler{App: app, service: s, config: sysConfig.GeekAPI.Captcha}
}

// RegisterRoutes 注册路由
func (h *CaptchaHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/captcha/")

	// 需要用户授权的接口
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.GET("get", h.Get)
		group.POST("check", h.Check)
		group.GET("slide/get", h.SlideGet)
		group.POST("slide/check", h.SlideCheck)
	}
}

func (h *CaptchaHandler) Get(c *gin.Context) {
	if !h.config.Enabled {
		resp.ERROR(c, "验证码服务未启用")
		return
	}

	data, err := h.service.Get()
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// Check verify the captcha data
func (h *CaptchaHandler) Check(c *gin.Context) {
	if !h.config.Enabled {
		resp.ERROR(c, "验证码服务未启用")
		return
	}

	var data struct {
		Key  string `json:"key"`
		Dots string `json:"dots"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if h.service.Check(data) {
		resp.SUCCESS(c)
	} else {
		resp.ERROR(c)
	}

}

// SlideGet 获取滑动验证图片
func (h *CaptchaHandler) SlideGet(c *gin.Context) {
	if !h.config.Enabled {
		resp.ERROR(c, "验证码服务未启用")
		return
	}

	data, err := h.service.SlideGet()
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// SlideCheck 滑动验证结果校验
func (h *CaptchaHandler) SlideCheck(c *gin.Context) {
	if !h.config.Enabled {
		resp.ERROR(c, "验证码服务未启用")
		return
	}

	var data struct {
		Key string `json:"key"`
		X   int    `json:"x"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if h.service.SlideCheck(data) {
		resp.SUCCESS(c)
	} else {
		resp.ERROR(c)
	}

}
