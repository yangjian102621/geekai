package utils

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/store/model"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"time"
)

var logger = logger2.GetLogger()

// ReplyChunkMessage 回复客户片段端消息
func ReplyChunkMessage(client *types.WsClient, message interface{}) {
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error for decoding json data: %v", err.Error())
		return
	}
	err = client.Send(msg)
	if err != nil {
		logger.Errorf("Error for reply message: %v", err.Error())
	}
}

// ReplyMessage 回复客户端一条完整的消息
func ReplyMessage(ws *types.WsClient, message interface{}) {
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: message})
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsEnd})
}

func DownloadImage(imageURL string, proxy string) ([]byte, error) {
	var client *http.Client
	if proxy == "" {
		client = http.DefaultClient
	} else {
		proxyURL, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}
	request, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
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

func OpenAIRequest(db *gorm.DB, prompt string, proxy string) (string, error) {
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
	if apiKey.UseProxy && proxy != "" {
		client.SetProxyURL(proxy)
	}
	r, err := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(types.ApiRequest{
			Model:       "gpt-3.5-turbo",
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
