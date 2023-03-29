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
	"openai/utils"
	"strings"
	"time"
)

const ErrorMsg = "抱歉，AI 助手开小差了，请马上联系管理员去盘它。"

// ChatHandle 处理聊天 WebSocket 请求
func (s *Server) ChatHandle(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Fatal(err)
		return
	}
	sessionId := c.Query("sessionId")
	roleKey := c.Query("role")
	session := s.ChatSession[sessionId]
	logger.Infof("New websocket connected, IP: %s, Username: %s", c.Request.RemoteAddr, session.Username)
	client := NewWsClient(ws)
	var roles = GetChatRoles()
	var chatRole = roles[roleKey]
	if !chatRole.Enable { // 角色未启用
		c.Abort()
		return
	}
	// 加载历史消息，如果历史消息为空则发送打招呼消息
	_, err = GetChatHistory(session.Username, roleKey)
	if err != nil {
		replyMessage(client, chatRole.HelloMsg, true)
	}
	go func() {
		for {
			_, message, err := client.Receive()
			if err != nil {
				logger.Error(err)
				client.Close()
				return
			}

			logger.Info("Receive a message: ", string(message))
			// TODO: 当前只保持当前会话的上下文，部保存用户的所有的聊天历史记录，后期要考虑保存所有的历史记录
			err = s.sendMessage(session, chatRole, string(message), client)
			if err != nil {
				logger.Error(err)
			}
		}
	}()
}

// 将消息发送给 ChatGPT 并获取结果，通过 WebSocket 推送到客户端
func (s *Server) sendMessage(session types.ChatSession, role types.ChatRole, prompt string, ws Client) error {
	user, err := GetUser(session.Username)
	if err != nil {
		replyMessage(ws, "当前 user 无效，请使用合法的 user 登录！", false)
		return err
	}

	if user.MaxCalls > 0 && user.RemainingCalls <= 0 {
		replyMessage(ws, "当前 TOKEN 点数已经用尽，加入我们的知识星球可以免费领取点卡！", false)
		replyMessage(ws, "![](images/start.png)", true)
		return nil
	}
	var r = types.ApiRequest{
		Model:       s.Config.Chat.Model,
		Temperature: s.Config.Chat.Temperature,
		MaxTokens:   s.Config.Chat.MaxTokens,
		Stream:      true,
	}
	var context []types.Message
	var ctxKey = fmt.Sprintf("%s-%s", session.SessionId, role.Key)
	if v, ok := s.ChatContexts[ctxKey]; ok && s.Config.Chat.EnableContext {
		context = v.Messages
	} else {
		context = role.Context
	}

	if s.DebugMode {
		logger.Infof("会话上下文：%+v", context)
	}

	r.Messages = append(context, types.Message{
		Role:    "user",
		Content: prompt,
	})

	requestBody, err := json.Marshal(r)
	if err != nil {
		return err
	}

	// 创建 HttpClient 请求对象
	var client *http.Client
	request, err := http.NewRequest(http.MethodPost, s.Config.Chat.ApiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	var retryCount = 5
	var response *http.Response
	var failedKey = ""
	var failedProxyURL = ""
	for retryCount > 0 {
		proxyURL := s.getProxyURL(failedProxyURL)
		if proxyURL == "" {
			client = &http.Client{}
		} else { // 使用代理
			uri := url.URL{}
			proxy, _ := uri.Parse(proxyURL)
			client = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxy),
				},
			}
		}
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
			failedProxyURL = proxyURL
		}
		retryCount--
	}

	// 如果三次请求都失败的话，则返回对应的错误信息
	if err != nil {
		replyMessage(ws, ErrorMsg, false)
		replyMessage(ws, "![](images/wx.png)", true)
		return err
	}

	// 循环读取 Chunk 消息
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
			replyChunkMessage(ws, types.WsMessage{Type: types.WsEnd})
			break
		} else if len(line) < 20 {
			continue
		}

		err = json.Unmarshal([]byte(line[6:]), &responseBody)
		if err != nil { // 数据解析出错
			logger.Error(err, line)
			replyChunkMessage(ws, types.WsMessage{Type: types.WsEnd, IsHelloMsg: false})
			replyMessage(ws, ErrorMsg, false)
			replyMessage(ws, "![](images/wx.png)", true)
			break
		}
		// 初始化 role
		if responseBody.Choices[0].Delta.Role != "" && message.Role == "" {
			message.Role = responseBody.Choices[0].Delta.Role
			replyChunkMessage(ws, types.WsMessage{Type: types.WsStart, IsHelloMsg: false})
			continue
		} else if responseBody.Choices[0].FinishReason != "" { // 输出完成或者输出中断了
			replyChunkMessage(ws, types.WsMessage{Type: types.WsEnd, IsHelloMsg: false})
			break
		} else {
			content := responseBody.Choices[0].Delta.Content
			contents = append(contents, content)
			replyChunkMessage(ws, types.WsMessage{
				Type:       types.WsMiddle,
				Content:    responseBody.Choices[0].Delta.Content,
				IsHelloMsg: false,
			})
		}
	}
	_ = response.Body.Close() // 关闭资源

	// 当前 Username 调用次数减 1
	if user.MaxCalls > 0 {
		user.RemainingCalls -= 1
		_ = PutUser(*user)
	}
	// 追加上下文消息
	useMsg := types.Message{Role: "user", Content: prompt}
	context = append(context, useMsg)
	message.Content = strings.Join(contents, "")
	context = append(context, message)
	// 更新上下文消息
	s.ChatContexts[ctxKey] = types.ChatContext{
		Messages:       context,
		LastAccessTime: time.Now().Unix(),
	}

	// 追加历史消息
	if user.EnableHistory {
		err = AppendChatHistory(user.Name, role.Key, useMsg)
		if err != nil {
			return err
		}
		err = AppendChatHistory(user.Name, role.Key, message)
	}
	return err
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

// 获取一个可用的代理
func (s *Server) getProxyURL(failedProxyURL string) string {
	if len(s.Config.ProxyURL) == 0 {
		return ""
	}

	if len(s.Config.ProxyURL) == 1 || failedProxyURL == "" {
		return s.Config.ProxyURL[0]
	}

	for i, v := range s.Config.ProxyURL {
		if failedProxyURL == v {
			if i == len(s.Config.ProxyURL)-1 {
				return s.Config.ProxyURL[0]
			} else {
				return s.Config.ProxyURL[i+1]
			}
		}
	}
	return ""
}

// 回复客户片段端消息
func replyChunkMessage(client Client, message types.WsMessage) {
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

// 回复客户端一条完整的消息
func replyMessage(ws Client, message string, isHelloMsg bool) {
	replyChunkMessage(ws, types.WsMessage{Type: types.WsStart, IsHelloMsg: isHelloMsg})
	replyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: message, IsHelloMsg: isHelloMsg})
	replyChunkMessage(ws, types.WsMessage{Type: types.WsEnd, IsHelloMsg: isHelloMsg})
}

func (s *Server) GetChatHistoryHandle(c *gin.Context) {
	sessionId := c.GetHeader(types.TokenName)
	var data struct {
		Role string `json:"role"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	session := s.ChatSession[sessionId]
	history, err := GetChatHistory(session.Username, data.Role)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "No history message"})
		return
	}

	var messages = make([]types.HistoryMessage, 0)
	role, err := GetChatRole(data.Role)
	if err == nil {
		// 先将打招呼的消息追加上去
		messages = append(messages, types.HistoryMessage{
			Type:    "reply",
			Id:      utils.RandString(32),
			Icon:    role.Icon,
			Content: role.HelloMsg,
		})

		for _, v := range history {
			if v.Role == "user" {
				messages = append(messages, types.HistoryMessage{
					Type:    "prompt",
					Id:      utils.RandString(32),
					Icon:    "images/avatar/user.png",
					Content: v.Content,
				})
			} else if v.Role == "assistant" {
				messages = append(messages, types.HistoryMessage{
					Type:    "reply",
					Id:      utils.RandString(32),
					Icon:    role.Icon,
					Content: v.Content,
				})
			}
		}
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: messages})
}

// ClearHistoryHandle 清空聊天记录
func (s *Server) ClearHistoryHandle(c *gin.Context) {
	sessionId := c.GetHeader(types.TokenName)
	var data struct {
		Role string `json:"role"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	session := s.ChatSession[sessionId]
	err = ClearChatHistory(session.Username, data.Role)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to remove data from DB"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success})
}
