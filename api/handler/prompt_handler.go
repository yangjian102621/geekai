package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/oss"
	"geekai/service/suno"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// 提示词生成 handler
// 使用 AI 生成绘画指令，歌词，视频生成指令等

type PromptHandler struct {
	BaseHandler
	sunoService *suno.Service
	uploader    *oss.UploaderManager
	userService *service.UserService
}

func NewPromptHandler(app *core.AppServer, db *gorm.DB, userService *service.UserService) *PromptHandler {
	return &PromptHandler{
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
		userService: userService,
	}
}

// Lyric 生成歌词
func (h *PromptHandler) Lyric(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	content, err := utils.OpenAIRequest(h.DB, fmt.Sprintf(service.LyricPromptTemplate, data.Prompt), h.App.SysConfig.TranslateModelId)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, content)
}

// Image 生成 AI 绘画提示词
func (h *PromptHandler) Image(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	content, err := utils.OpenAIRequest(h.DB, fmt.Sprintf(service.ImagePromptOptimizeTemplate, data.Prompt), h.App.SysConfig.TranslateModelId)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, strings.Trim(content, `"`))
}

// Video 生成视频提示词
func (h *PromptHandler) Video(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	content, err := utils.OpenAIRequest(h.DB, fmt.Sprintf(service.VideoPromptTemplate, data.Prompt), h.App.SysConfig.TranslateModelId)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, strings.Trim(content, `"`))
}

// MetaPrompt 生成元提示词
func (h *PromptHandler) MetaPrompt(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	messages := make([]interface{}, 0)
	messages = append(messages, types.Message{
		Role:    "system",
		Content: service.MetaPromptTemplate,
	})
	messages = append(messages, types.Message{
		Role:    "user",
		Content: "Task, Goal, or the Role to actor is:\n" + data.Prompt,
	})
	content, err := utils.SendOpenAIMessage(h.DB, messages, 0)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, strings.Trim(content, `"`))
}
