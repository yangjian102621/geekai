package utils

import (
	"chatplus/core/types"
	"chatplus/store/model"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/pkoukk/tiktoken-go"
	"gorm.io/gorm"
	"time"
)

func CalcTokens(text string, model string) (int, error) {
	encoding, ok := tiktoken.MODEL_TO_ENCODING[model]
	if !ok {
		encoding = "cl100k_base"
	}
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, fmt.Errorf("getEncoding: %v", err)
	}

	token := tke.Encode(text, nil, nil)
	return len(token), nil
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

func OpenAIRequest(db *gorm.DB, prompt string) (string, error) {
	var apiKey model.ApiKey
	res := db.Where("platform = ?", types.OpenAI).Where("type = ?", "chat").Where("enabled = ?", true).First(&apiKey)
	if res.Error != nil {
		return "", fmt.Errorf("error with fetch OpenAI API KEY：%v", res.Error)
	}

	messages := make([]interface{}, 1)
	messages[0] = types.Message{
		Role:    "user",
		Content: prompt,
	}

	var response apiRes
	var errRes apiErrRes
	client := req.C()
	if len(apiKey.ProxyURL) > 5 {
		client.SetProxyURL(apiKey.ApiURL)
	}
	r, err := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(types.ApiRequest{
			Model:       "gpt-3.5-turbo-0125",
			Temperature: 0.9,
			MaxTokens:   1024,
			Stream:      false,
			Messages:    messages,
		}).
		SetErrorResult(&errRes).
		SetSuccessResult(&response).Post(apiKey.ApiURL)
	if err != nil || r.IsErrorState() {
		return "", fmt.Errorf("error with http request: %v%v%s", err, r.Err, errRes.Error.Message)
	}

	// 更新 API KEY 的最后使用时间
	db.Model(&apiKey).UpdateColumn("last_used_at", time.Now().Unix())

	return response.Choices[0].Message.Content, nil
}
