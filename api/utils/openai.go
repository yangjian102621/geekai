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
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/pkoukk/tiktoken-go"
	"gorm.io/gorm"
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

// OpenAIResponse 非流式 chat/completions 响应；content 使用 RawMessage 以兼容 string 与多段式数组。
type OpenAIResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// assistantContentPart 兼容 OpenAI 多模态 / 多段文本：{"type":"text","text":"..."} 等。
type assistantContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
	// 少数网关使用 content 字段承载文本
	Inner string `json:"content"`
}

// NormalizeAssistantContent 将 message.content 规范为纯文本。
// 兼容：JSON string、null、OpenAI 数组段、以及单段对象。
func NormalizeAssistantContent(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return ""
	}
	var plain string
	if err := json.Unmarshal(raw, &plain); err == nil {
		return plain
	}
	var parts []assistantContentPart
	if err := json.Unmarshal(raw, &parts); err == nil && len(parts) > 0 {
		var b strings.Builder
		for _, p := range parts {
			switch {
			case p.Text != "":
				b.WriteString(p.Text)
			case p.Inner != "" && p.Type != "image_url" && p.Type != "image":
				b.WriteString(p.Inner)
			}
		}
		return b.String()
	}
	var one assistantContentPart
	if err := json.Unmarshal(raw, &one); err == nil {
		if one.Text != "" {
			return one.Text
		}
		if one.Inner != "" {
			return one.Inner
		}
	}
	return ""
}

func OpenAIRequest(db *gorm.DB, prompt string, modelId int) (string, error) {
	messages := make([]any, 1)
	messages[0] = types.Message{
		Role:    "user",
		Content: prompt,
	}
	return SendOpenAIMessage(db, messages, modelId)
}

func SendOpenAIMessage(db *gorm.DB, messages []any, modelId int) (string, error) {
	var chatModel model.ChatModel
	db.Where("id", modelId).First(&chatModel)
	if chatModel.Value == "" {
		chatModel.Value = "gpt-4o" // 默认使用 gpt-4o
	}
	var apiKey model.ApiKey
	session := db.Session(&gorm.Session{}).Where("type", "chat").Where("enabled", true)
	if chatModel.KeyId > 0 {
		session = session.Where("id", chatModel.KeyId)
	}
	err := session.First(&apiKey).Error
	if err != nil {
		return "", fmt.Errorf("error with fetch OpenAI API KEY：%v", err)
	}

	var response OpenAIResponse
	client := req.C()
	if len(apiKey.ProxyURL) > 5 {
		client.SetProxyURL(apiKey.ProxyURL)
	}
	var apiURL string
	p, _ := url.Parse(apiKey.ApiURL)
	// 如果设置的是 BASE_URL 没有路径，则添加 /v1/chat/completions
	if p.Path == "" {
		apiURL = fmt.Sprintf("%s/v1/chat/completions", apiKey.ApiURL)
	} else {
		apiURL = apiKey.ApiURL
	}
	logger.Infof("Sending %s request, API KEY:%s, PROXY: %s, Model: %s", apiURL, apiKey.ApiURL, apiKey.ProxyURL, chatModel.Name)
	r, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(types.ApiRequest{
			Model:       chatModel.Value,
			Temperature: 0.9,
			// gpt-5/o 系模型在复杂提示下可能先消耗大量 token 于推理过程，
			// 1024 容易导致 finish_reason=length 且 content 为空。
			MaxTokens:           4096,
			MaxCompletionTokens: 4096,
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
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("模型返回 choices 为空：%s", string(body))
	}

	out := strings.TrimSpace(NormalizeAssistantContent(response.Choices[0].Message.Content))
	if out == "" {
		return "", fmt.Errorf("模型返回内容为空，请重试或更换模型（若使用推理模型，请确认网关已返回 content 或 reasoning_content）")
	}

	// 更新 API KEY 的最后使用时间
	db.Model(&apiKey).UpdateColumn("last_used_at", time.Now().Unix())

	return out, nil
}
