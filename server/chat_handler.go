package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"math/rand"
	"net/http"
	"openai/types"
	"strings"
	"time"
)

// ChatHandle 处理聊天 WebSocket 请求
func (s *Server) ChatHandle(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Fatal(err)
		return
	}
	logger.Infof("New websocket connected, IP: %s", c.Request.RemoteAddr)
	client := NewWsClient(ws)
	go func() {
		for {
			_, message, err := client.Receive()
			if err != nil {
				logger.Error(err)
				client.Close()
				return
			}

			logger.Info(string(message))
			// TODO: 根据会话请求，传入不同的用户 ID
			err = s.sendMessage("test", string(message), client)
			if err != nil {
				logger.Error(err)
			}
		}
	}()
}

// 将消息发送给 ChatGPT 并获取结果，通过 WebSocket 推送到客户端
func (s *Server) sendMessage(userId string, text string, ws Client) error {
	var r = types.ApiRequest{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.9,
		MaxTokens:   1024,
		Stream:      true,
	}
	var history []types.Message
	if v, ok := s.History[userId]; ok && s.Config.Chat.EnableContext {
		history = v
		//logger.Infof("上下文历史消息:%+v", history)
	} else {
		history = make([]types.Message, 0)
	}
	r.Messages = append(history, types.Message{
		Role:    "user",
		Content: text,
	})

	requestBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, s.Config.Chat.ApiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	// 随机获取一个 API Key，如果请求失败，则更换 API Key 重试
	// TODO: 需要将失败的 Key 移除列表
	rand.Seed(time.Now().UnixNano())
	var retryCount = 3
	var response *http.Response
	for retryCount > 0 {
		index := rand.Intn(len(s.Config.Chat.ApiKeys))
		apiKey := s.Config.Chat.ApiKeys[index]
		logger.Infof("Use API KEY: %s", apiKey)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
		response, err = s.Client.Do(request)
		if err == nil {
			break
		} else {
			logger.Error(err)
		}
		retryCount--
	}
	if err != nil {
		return err
	}

	var message = types.Message{}
	var contents = make([]string, 0)
	var responseBody = types.ApiResponse{}
	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			logger.Error(err)
			break
		}

		if line == "" {
			replyMessage(types.WsMessage{Type: types.WsEnd}, ws)
			break
		} else if len(line) < 20 {
			continue
		}

		err = json.Unmarshal([]byte(line[6:]), &responseBody)
		if err != nil {
			logger.Error(line)
			continue
		}
		// 初始化 role
		if responseBody.Choices[0].Delta.Role != "" && message.Role == "" {
			message.Role = responseBody.Choices[0].Delta.Role
			replyMessage(types.WsMessage{Type: types.WsStart}, ws)
			continue
		} else if responseBody.Choices[0].FinishReason != "" { // 输出完成或者输出中断了
			replyMessage(types.WsMessage{Type: types.WsEnd}, ws)
			break
		} else {
			content := responseBody.Choices[0].Delta.Content
			contents = append(contents, content)
			replyMessage(types.WsMessage{
				Type:    types.WsMiddle,
				Content: responseBody.Choices[0].Delta.Content,
			}, ws)
		}
	}

	// 追加历史消息
	history = append(history, types.Message{
		Role:    "user",
		Content: text,
	})
	message.Content = strings.Join(contents, "")
	history = append(history, message)
	s.History[userId] = history
	return nil
}

// 回复客户端消息
func replyMessage(message types.WsMessage, client Client) {
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error for decoding json data: %v", err.Error())
		return
	}
	err = client.(*WsClient).Send(msg)
	if err != nil {
		logger.Errorf("Error for reply message: %v", err.Error())
	}
}
