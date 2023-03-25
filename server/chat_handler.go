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
	"net/url"
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
	token := c.Query("token")
	role := c.Query("role")
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
			// TODO: 当前只保持当前会话的上下文，部保存用户的所有的聊天历史记录，后期要考虑保存所有的历史记录
			err = s.sendMessage(token, role, string(message), client)
			if err != nil {
				logger.Error(err)
			}
		}
	}()
}

// 将消息发送给 ChatGPT 并获取结果，通过 WebSocket 推送到客户端
func (s *Server) sendMessage(sessionId string, role string, text string, ws Client) error {
	var r = types.ApiRequest{
		Model:       s.Config.Chat.Model,
		Temperature: s.Config.Chat.Temperature,
		MaxTokens:   s.Config.Chat.MaxTokens,
		Stream:      true,
	}
	var context []types.Message
	var key = sessionId + role
	if v, ok := s.ChatContext[key]; ok && s.Config.Chat.EnableContext {
		context = v
	} else {
		context = s.Config.ChatRoles[role].Context
	}
	logger.Info(context)
	r.Messages = append(context, types.Message{
		Role:    "user",
		Content: text,
	})

	requestBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	// 创建 HttpClient 请求对象
	var client *http.Client
	if s.Config.ProxyURL == "" {
		client = &http.Client{}
	} else { // 使用代理
		uri := url.URL{}
		proxy, _ := uri.Parse(s.Config.ProxyURL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}
	request, err := http.NewRequest(http.MethodPost, s.Config.Chat.ApiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	var retryCount = 3
	var response *http.Response
	var failedKey = ""
	for retryCount > 0 {
		apiKey := s.getApiKey(failedKey)
		if apiKey == "" {
			logger.Info("Too many requests, all Api Key is not available")
			time.Sleep(time.Second)
			continue
		}
		logger.Infof("Use API KEY: %s", apiKey)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
		response, err = client.Do(request)
		if err == nil {
			break
		} else {
			logger.Error(err)
			failedKey = apiKey
		}
		retryCount--
	}

	// 如果三次请求都失败的话，则返回对应的错误信息
	if err != nil {
		replyError(ws)
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
			replyError(ws)
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
	context = append(context, types.Message{
		Role:    "user",
		Content: text,
	})
	message.Content = strings.Join(contents, "")
	context = append(context, message)
	// 保存上下文
	s.ChatContext[key] = context
	return nil
}

func replyError(ws Client) {
	replyMessage(types.WsMessage{Type: types.WsStart}, ws)
	replyMessage(types.WsMessage{Type: types.WsMiddle, Content: "抱歉，AI 助手开小差了，我马上找人去盘它。"}, ws)
	replyMessage(types.WsMessage{Type: types.WsEnd}, ws)
}

// 随机获取一个 API Key，如果请求失败，则更换 API Key 重试
func (s *Server) getApiKey(failedKey string) string {
	var keys = make([]string, 0)
	for _, v := range s.Config.Chat.ApiKeys {
		// 过滤掉刚刚失败的 Key
		if v == failedKey {
			continue
		}

		// 获取 API Key 的上次调用时间，控制调用频率
		var lastAccess int64
		if t, ok := s.ApiKeyAccessStat[v]; ok {
			lastAccess = t
		}
		// 保持每分钟访问不超过 15 次
		if time.Now().Unix()-lastAccess <= 4 {
			continue
		}

		keys = append(keys, v)
	}
	rand.Seed(time.Now().UnixNano())
	if len(keys) > 0 {
		key := keys[rand.Intn(len(keys))]
		s.ApiKeyAccessStat[key] = time.Now().Unix()
		return key
	}
	return ""
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
