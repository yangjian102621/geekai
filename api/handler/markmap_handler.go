package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MarkMapHandler 生成思维导图
type MarkMapHandler struct {
	BaseHandler
	clients *types.LMap[int, *types.WsClient]
}

func NewMarkMapHandler(app *core.AppServer, db *gorm.DB) *MarkMapHandler {
	return &MarkMapHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
		clients:     types.NewLMap[int, *types.WsClient](),
	}
}

func (h *MarkMapHandler) Client(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	modelId := h.GetInt(c, "model_id", 0)
	userId := h.GetInt(c, "user_id", 0)

	client := types.NewWsClient(ws)
	h.clients.Put(userId, client)
	go func() {
		for {
			_, msg, err := client.Receive()
			if err != nil {
				client.Close()
				h.clients.Delete(userId)
				return
			}

			var message types.WsMessage
			err = utils.JsonDecode(string(msg), &message)
			if err != nil {
				continue
			}

			// 心跳消息
			if message.Type == "heartbeat" {
				logger.Debug("收到 MarkMap 心跳消息：", message.Content)
				continue
			}
			// change model
			if message.Type == "model_id" {
				modelId = utils.IntValue(utils.InterfaceToString(message.Content), 0)
				continue
			}

			logger.Info("Receive a message: ", message.Content)
			err = h.sendMessage(client, utils.InterfaceToString(message.Content), modelId, userId)
			if err != nil {
				logger.Error(err)
				utils.ReplyChunkMessage(client, types.WsMessage{Type: types.WsErr, Content: err.Error()})
			}

		}
	}()
}

func (h *MarkMapHandler) sendMessage(client *types.WsClient, prompt string, modelId int, userId int) error {
	var user model.User
	res := h.DB.Model(&model.User{}).First(&user, userId)
	if res.Error != nil {
		return fmt.Errorf("error with query user info: %v", res.Error)
	}
	var chatModel model.ChatModel
	res = h.DB.Where("id", modelId).First(&chatModel)
	if res.Error != nil {
		return fmt.Errorf("error with query chat model: %v", res.Error)
	}

	if user.Power < chatModel.Power {
		return fmt.Errorf("您当前剩余算力（%d）已不足以支付当前模型算力（%d）！", user.Power, chatModel.Power)
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

另外，除此之外不要任何解释性语句。
`})
	messages = append(messages, types.Message{Role: "user", Content: fmt.Sprintf("请生成一份有关【%s】一份思维导图，要求结构清晰，有条理", prompt)})
	var req = types.ApiRequest{
		Model:    chatModel.Value,
		Stream:   true,
		Messages: messages,
	}

	var apiKey model.ApiKey
	response, err := h.doRequest(req, chatModel, &apiKey)
	if err != nil {
		return fmt.Errorf("请求 OpenAI API 失败: %s", err)
	}

	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		// 循环读取 Chunk 消息
		scanner := bufio.NewScanner(response.Body)
		var isNew = true
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(line, "data:") || len(line) < 30 {
				continue
			}

			var responseBody = types.ApiResponse{}
			err = json.Unmarshal([]byte(line[6:]), &responseBody)
			if err != nil { // 数据解析出错
				return fmt.Errorf("error with decode data: %v", line)
			}

			if len(responseBody.Choices) == 0 { // Fixed: 兼容 Azure API 第一个输出空行
				continue
			}

			if responseBody.Choices[0].FinishReason == "stop" {
				break
			}

			if isNew {
				utils.ReplyChunkMessage(client, types.WsMessage{Type: types.WsStart})
				isNew = false
			}
			utils.ReplyChunkMessage(client, types.WsMessage{
				Type:    types.WsMiddle,
				Content: utils.InterfaceToString(responseBody.Choices[0].Delta.Content),
			})
		} // end for

		utils.ReplyChunkMessage(client, types.WsMessage{Type: types.WsEnd})

	} else {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("请求 OpenAI API 失败：%s", string(body))
	}

	// 扣减算力
	if chatModel.Power > 0 {
		res = h.DB.Model(&model.User{}).Where("id", userId).UpdateColumn("power", gorm.Expr("power - ?", chatModel.Power))
		if res.Error == nil {
			// 记录算力消费日志
			var u model.User
			h.DB.Where("id", userId).First(&u)
			h.DB.Create(&model.PowerLog{
				UserId:    u.Id,
				Username:  u.Username,
				Type:      types.PowerConsume,
				Amount:    chatModel.Power,
				Mark:      types.PowerSub,
				Balance:   u.Power,
				Model:     chatModel.Value,
				Remark:    fmt.Sprintf("AI绘制思维导图，模型名称：%s, ", chatModel.Value),
				CreatedAt: time.Now(),
			})
		}
	}

	return nil
}

func (h *MarkMapHandler) doRequest(req types.ApiRequest, chatModel model.ChatModel, apiKey *model.ApiKey) (*http.Response, error) {

	session := h.DB.Session(&gorm.Session{})
	// if the chat model bind a KEY, use it directly
	if chatModel.KeyId > 0 {
		session = session.Where("id", chatModel.KeyId)
	} else { // use the last unused key
		session = session.Where("type", "chat").
			Where("enabled", true).Order("last_used_at ASC")
	}

	res := session.First(apiKey)
	if res.Error != nil {
		return nil, errors.New("no available key, please import key")
	}
	apiURL := fmt.Sprintf("%s/v1/chat/completions", apiKey.ApiURL)
	// 更新 API KEY 的最后使用时间
	h.DB.Model(apiKey).UpdateColumn("last_used_at", time.Now().Unix())

	// 创建 HttpClient 请求对象
	var client *http.Client
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	if len(apiKey.ProxyURL) > 5 { // 使用代理
		proxy, _ := url.Parse(apiKey.ProxyURL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	} else {
		client = http.DefaultClient
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey.Value))
	logger.Debugf("Sending %s request, API KEY:%s, PROXY: %s, Model: %s", apiKey.ApiURL, apiURL, apiKey.ProxyURL, req.Model)
	return client.Do(request)
}
