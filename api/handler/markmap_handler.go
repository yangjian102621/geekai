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
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MarkMapHandler 生成思维导图
type MarkMapHandler struct {
	BaseHandler
	clients     *types.LMap[int, *types.WsClient]
	userService *service.UserService
}

func NewMarkMapHandler(app *core.AppServer, db *gorm.DB, userService *service.UserService) *MarkMapHandler {
	return &MarkMapHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
		clients:     types.NewLMap[int, *types.WsClient](),
		userService: userService,
	}
}

// Generate 生成思维导图
func (h *MarkMapHandler) Generate(c *gin.Context) {
	var data struct {
		Prompt  string `json:"prompt"`
		ModelId int    `json:"model_id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	userId := h.GetLoginUserId(c)
	var user model.User
	err := h.DB.Where("id", userId).First(&user, userId).Error
	if err != nil {
		resp.ERROR(c, "error with query user info")
		return
	}
	var chatModel model.ChatModel
	err = h.DB.Where("id", data.ModelId).First(&chatModel).Error
	if err != nil {
		resp.ERROR(c, "error with query chat model")
		return
	}

	if user.Power < chatModel.Power {
		resp.ERROR(c, fmt.Sprintf("您当前剩余算力（%d）已不足以支付当前模型算力（%d）！", user.Power, chatModel.Power))
		return
	}

	messages := make([]interface{}, 0)
	messages = append(messages, types.Message{Role: "system", Content: `
你是一位非常优秀的思维导图助手， 你能帮助用户整理思路，根据用户提供的主题或内容，快速生成结构清晰，有条理的思维导图，然后以 Markdown 格式输出。markdown 只需要输出一级标题，二级标题，三级标题，四级标题，最多输出四级，除此之外不要输出任何其他 markdown 标记。下面是一个合格的例子：
# Geek-AI 助手

## 完整的开源系统
### 前端开源
### 后端开源

## 支持各种大模型
### OpenAI 
### Azure 
### 文心一言
### 通义千问

## 集成多种收费方式
### 支付宝
### 微信

请直接生成结果，不要任何解释性语句。
`})
	messages = append(messages, types.Message{Role: "user", Content: fmt.Sprintf("请生成一份有关【%s】一份思维导图，要求结构清晰，有条理", data.Prompt)})
	content, err := utils.SendOpenAIMessage(h.DB, messages, data.ModelId)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("请求 OpenAI API 失败: %s", err))
		return
	}

	// 扣减算力
	if chatModel.Power > 0 {
		err = h.userService.DecreasePower(int(userId), chatModel.Power, model.PowerLog{
			Type:   types.PowerConsume,
			Model:  chatModel.Value,
			Remark: fmt.Sprintf("AI绘制思维导图，模型名称：%s, ", chatModel.Value),
		})
		if err != nil {
			resp.ERROR(c, "error with save power log, "+err.Error())
			return
		}
	}

	resp.SUCCESS(c, content)
}
