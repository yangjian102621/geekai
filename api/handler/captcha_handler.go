package handler

import (
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
)

// 今日头条函数实现

type CaptchaHandler struct {
	service *service.CaptchaService
}

func NewCaptchaHandler(s *service.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{service: s}
}

func (h *CaptchaHandler) Get(c *gin.Context) {
	data, err := h.service.Get()
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// Check verify the captcha data
func (h *CaptchaHandler) Check(c *gin.Context) {
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
