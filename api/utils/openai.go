package utils

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"github.com/imroc/req/v3"
	"github.com/pkoukk/tiktoken-go"
	"gorm.io/gorm"
	"io"
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

func OpenAIRequest(db *gorm.DB, prompt string, modelName string) (string, error) {
	var apiKey model.ApiKey
	res := db.Where("type", "chat").Where("enabled", true).First(&apiKey)
	if res.Error != nil {
		return "", fmt.Errorf("error with fetch OpenAI API KEY：%v", res.Error)
	}

	messages := make([]interface{}, 1)
	messages[0] = types.Message{
		Role:    "user",
		Content: prompt,
	}

	var response apiRes
	client := req.C()
	if len(apiKey.ProxyURL) > 5 {
		client.SetProxyURL(apiKey.ApiURL)
	}
	apiURL := fmt.Sprintf("%s/v1/chat/completions", apiKey.ApiURL)
	logger.Debugf("Sending %s request, API KEY:%s, PROXY: %s, Model: %s", apiKey.ApiURL, apiURL, apiKey.ProxyURL, modelName)
	r, err := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(types.ApiRequest{
			Model:       modelName,
			Temperature: 0.9,
			MaxTokens:   1024,
			Stream:      false,
			Messages:    messages,
		}).Post(apiURL)
	if err != nil {
		return "", fmt.Errorf("请求 OpenAI API失败：%v", err)
	}

	if r.IsErrorState() {
		return "", fmt.Errorf("请求 OpenAI API失败：%v", r.Status)
	}

	body, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	// 更新 API KEY 的最后使用时间
	db.Model(&apiKey).UpdateColumn("last_used_at", time.Now().Unix())

	return response.Choices[0].Message.Content, nil
}
