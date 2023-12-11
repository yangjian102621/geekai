package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/utils/resp"
	"fmt"

	"github.com/imroc/req/v3"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const translatePromptTemplate = "Please rewrite the following text into AI painting prompt words, and please try to add detailed description of the picture, painting style, scene, rendering effect, picture light and other elements. Please output directly in English without any explanation, within 150 words. The text to be rewritten is: [%s]"

type PromptHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewPromptHandler(app *core.AppServer, db *gorm.DB) *PromptHandler {
	h := &PromptHandler{db: db}
	h.App = app
	return h
}

type apiRes struct {
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type apiErrRes struct {
	Error struct {
		Code    interface{} `json:"code"`
		Message string      `json:"message"`
		Param   interface{} `json:"param"`
		Type    string      `json:"type"`
	} `json:"error"`
}

func (h *PromptHandler) Translate(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 获取 OpenAI 的 API KEY
	var apiKey model.ApiKey
	res := h.db.Where("platform = ?", types.OpenAI).First(&apiKey)
	if res.Error != nil {
		resp.ERROR(c, "找不到可用 OpenAI API KEY")
		return
	}

	messages := make([]interface{}, 1)
	messages[0] = types.Message{
		Role:    "user",
		Content: fmt.Sprintf(translatePromptTemplate, data.Prompt),
	}

	var response apiRes
	var errRes apiErrRes
	r, err := req.C().SetProxyURL(h.App.Config.ProxyURL).R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(types.ApiRequest{
			Model:       "gpt-3.5-turbo",
			Temperature: 0.9,
			MaxTokens:   1024,
			Stream:      false,
			Messages:    messages,
		}).
		SetErrorResult(&errRes).
		SetSuccessResult(&response).Post(h.App.ChatConfig.OpenAI.ApiURL)
	if err != nil || r.IsErrorState() {
		resp.ERROR(c, fmt.Sprintf("error with http request: %v%v%s", err, r.Err, errRes.Error.Message))
		return
	}

	resp.SUCCESS(c, response.Choices[0].Message.Content)
}
