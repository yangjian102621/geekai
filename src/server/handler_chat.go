package server

import (
	"bufio"
	"bytes"
	"chatplus/types"
	"chatplus/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"net/url"
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
	chatId := c.Query("chatId")
	session, ok := s.ChatSession[sessionId]
	session.ChatId = chatId
	if !ok { // 用户未登录
		c.Abort()
		return
	}

	logger.Infof("New websocket connected, IP: %s, Username: %s", c.Request.RemoteAddr, session.Username)
	client := NewWsClient(ws)
	var roles = GetChatRoles()
	var chatRole = roles[roleKey]
	if !chatRole.Enable { // 角色未启用
		c.Abort()
		return
	}
	// 保存会话连接
	s.ChatClients[sessionId] = client

	// 加载历史消息，如果历史消息为空则发送打招呼消息
	_, err = GetChatHistory(session.Username, roleKey)
	if err != nil {
		replyMessage(client, chatRole.HelloMsg, true)
		// 发送项目地址
		replyMessage(client, "本项目已开放全部源代码：https://github.com/yangjian102621/chatgpt-plus，一分钟搭建自己的 ChatGPT 应用。", false)
	}
	go func() {
		for {
			_, message, err := client.Receive()
			if err != nil {
				logger.Error(err)
				client.Close()
				delete(s.ChatClients, sessionId)
				return
			}
			logger.Info("Receive a message: ", string(message))
			//replyMessage(client, "当前 TOKEN 无效，请使用合法的 TOKEN 登录！", false)
			//replyMessage(client, "![](images/wx.png)", false)
			ctx, cancel := context.WithCancel(context.Background())
			s.ReqCancelFunc[sessionId] = cancel
			// 回复消息
			err = s.sendMessage(ctx, session, chatRole, string(message), client)
			if err != nil {
				logger.Error(err)
			} else {
				replyChunkMessage(client, types.WsMessage{Type: types.WsEnd, IsHelloMsg: false})
				logger.Info("回答完毕: " + string(message))
			}

		}
	}()
}

// 将消息发送给 ChatGPT 并获取结果，通过 WebSocket 推送到客户端
func (s *Server) sendMessage(ctx context.Context, session types.ChatSession, role types.ChatRole, prompt string, ws Client) error {
	cancel := s.ReqCancelFunc[session.SessionId]
	defer func() {
		cancel()
		delete(s.ReqCancelFunc, session.SessionId)
	}()

	user, err := GetUser(session.Username)
	if err != nil {
		replyMessage(ws, "当前口令无效，请使用合法的口令登录！", false)
		return err
	}

	if user.Status == false {
		replyMessage(ws, "当前口令已经被禁用，如果疑问，请联系管理员！", false)
		replyMessage(ws, "![](images/wx.png)", false)
		return errors.New("当前口令" + user.Name + "已经被禁用")
	}

	if time.Now().Unix() > user.ExpiredTime {
		exTime := time.Unix(user.ExpiredTime, 0).Format("2006-01-02 15:04:05")
		replyMessage(ws, "当前口令已过期，过期时间为："+exTime+"，如果疑问，请联系管理员！", false)
		replyMessage(ws, "![](images/wx.png)", false)
		return errors.New("当前口令" + user.Name + "已过期")
	}

	if user.MaxCalls > 0 && user.RemainingCalls <= 0 {
		replyMessage(ws, "当前口令点数已经用尽，请联系管理员领取新的免费口令！", false)
		replyMessage(ws, "![](images/wx.png)", false)
		return nil
	}
	var req = types.ApiRequest{
		Model:       s.Config.Chat.Model,
		Temperature: s.Config.Chat.Temperature,
		MaxTokens:   s.Config.Chat.MaxTokens,
		Stream:      true,
	}
	var chatCtx []types.Message
	var ctxKey = fmt.Sprintf("%s-%s-%s", session.SessionId, role.Key, session.ChatId)
	if v, ok := s.ChatContexts[ctxKey]; ok && s.Config.Chat.EnableContext {
		chatCtx = v.Messages
	} else {
		chatCtx = role.Context
	}

	if s.DebugMode {
		logger.Infof("会话上下文：%+v", chatCtx)
	}

	req.Messages = append(chatCtx, types.Message{
		Role:    "user",
		Content: prompt,
	})

	// 创建 HttpClient 请求对象
	var client *http.Client
	var retryCount = 5 // 重试次数
	var response *http.Response
	var apiKey string
	var failedKey = ""
	var failedProxyURL = ""
	for retryCount > 0 {
		requestBody, err := json.Marshal(req)
		if err != nil {
			return err
		}

		request, err := http.NewRequest(http.MethodPost, s.Config.Chat.ApiURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}

		request = request.WithContext(ctx)
		request.Header.Add("Content-Type", "application/json")

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
		apiKey = s.getApiKey(failedKey)
		if apiKey == "" {
			logger.Info("Too many requests, all Api Key is not available")
			time.Sleep(time.Second)
			continue
		}
		logger.Infof("Use API KEY: %s, PROXY: %s", apiKey, proxyURL)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
		response, err = client.Do(request)
		if err == nil {
			break
		} else if strings.Contains(err.Error(), "context canceled") {
			return errors.New("用户取消了请求：" + prompt)
		} else {
			logger.Error("HTTP API 请求失败：" + err.Error())
			failedKey = apiKey
			failedProxyURL = proxyURL
		}
		retryCount--
	}

	if response != nil {
		defer response.Body.Close()
	}

	// 如果三次请求都失败的话，则返回对应的错误信息
	if err != nil {
		replyMessage(ws, ErrorMsg, false)
		replyMessage(ws, "![](images/wx.png)", false)
		return err
	}

	// 循环读取 Chunk 消息
	var message = types.Message{}
	var contents = make([]string, 0)
	var responseBody = types.ApiResponse{}
	reader := bufio.NewReader(response.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			logger.Error(err)
			break
		}
		if len(line) < 20 {
			continue
		} else if strings.Contains(line, "This key is associated with a deactivated account") || // 账号被禁用
			strings.Contains(line, "You exceeded your current quota") { // 当前 KEY 余额被用尽

			logger.Infof("API Key %s is deactivated", apiKey)
			// 移除当前 API key
			for i, v := range s.Config.Chat.ApiKeys {
				if v.Value == apiKey {
					s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys[:i], s.Config.Chat.ApiKeys[i+1:]...)
				}
			}
			// 更新配置文档
			_ = utils.SaveConfig(s.Config, s.ConfigPath)

			// 重发当前消息
			return s.sendMessage(ctx, session, role, prompt, ws)

			// 上下文超出长度了
		} else if strings.Contains(line, "This model's maximum context length is 4097 tokens") {
			logger.Infof("会话上下文长度超出限制, Username: %s", user.Name)
			replyMessage(ws, "温馨提示：当前会话上下文长度超出限制，已为您重置会话上下文！", false)
			// 重置上下文
			delete(s.ChatContexts, ctxKey)
			break
		} else if !strings.Contains(line, "data:") {
			continue
		}

		err = json.Unmarshal([]byte(line[6:]), &responseBody)
		if err != nil { // 数据解析出错
			logger.Error(err, line)
			replyMessage(ws, ErrorMsg, false)
			replyMessage(ws, "![](images/wx.png)", false)
			break
		}

		// 初始化 role
		if responseBody.Choices[0].Delta.Role != "" && message.Role == "" {
			message.Role = responseBody.Choices[0].Delta.Role
			replyChunkMessage(ws, types.WsMessage{Type: types.WsStart, IsHelloMsg: false})
			continue
		} else if responseBody.Choices[0].FinishReason != "" {
			break // 输出完成或者输出中断了
		} else {
			content := responseBody.Choices[0].Delta.Content
			contents = append(contents, content)
			replyChunkMessage(ws, types.WsMessage{
				Type:       types.WsMiddle,
				Content:    responseBody.Choices[0].Delta.Content,
				IsHelloMsg: false,
			})
		}

		// 监控取消信号
		select {
		case <-ctx.Done():
			_ = response.Body.Close() // 关闭响应流
			return errors.New("用户取消了请求：" + prompt)
		default:
			continue
		}
	} // end for

	// 消息发送成功
	if len(contents) > 0 {
		// 当前 Username 调用次数减 1
		if user.MaxCalls > 0 {
			user.RemainingCalls -= 1
			_ = PutUser(*user)
		}

		if message.Role == "" {
			message.Role = "assistant"
		}
		message.Content = strings.Join(contents, "")
		useMsg := types.Message{Role: "user", Content: prompt}

		// 更新上下文消息
		if s.Config.Chat.EnableContext {
			chatCtx = append(chatCtx, useMsg)  // 提问消息
			chatCtx = append(chatCtx, message) // 回复消息
			s.ChatContexts[ctxKey] = types.ChatContext{
				Messages:       chatCtx,
				LastAccessTime: time.Now().Unix(),
			}
		}

		// 追加历史消息
		if user.EnableHistory {
			err = AppendChatHistory(user.Name, role.Key, useMsg) // 提问消息
			if err != nil {
				return err
			}
			err = AppendChatHistory(user.Name, role.Key, message) // 回复消息
		}
	}

	return nil
}

// 随机获取一个 API Key，如果请求失败，则更换 API Key 重试
func (s *Server) getApiKey(failedKey string) string {
	var keys = make([]types.APIKey, 0)
	for _, key := range s.Config.Chat.ApiKeys {
		// 过滤掉刚刚失败的 Key
		if key.Value == failedKey {
			continue
		}

		// 保持每分钟访问不超过 15 次，控制调用频率
		if key.LastUsed > 0 && time.Now().Unix()-key.LastUsed <= 4 {
			continue
		}

		keys = append(keys, key)
	}
	// 从可用的 Key 中随机选一个
	rand.NewSource(time.Now().UnixNano())
	if len(keys) > 0 {
		key := keys[rand.Intn(len(keys))]
		// 更新选中 Key 的最后使用时间
		for i, item := range s.Config.Chat.ApiKeys {
			if item.Value == key.Value {
				s.Config.Chat.ApiKeys[i].LastUsed = time.Now().Unix()
			}
		}
		return key.Value
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
	user, err := GetUser(session.Username)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	if v, ok := user.ChatRoles[data.Role]; !ok || v != 1 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "No permission to access the history of role " + data.Role})
		return
	}

	history, err := GetChatHistory(session.Username, data.Role)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: nil, Message: "No history message"})
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

// StopGenerateHandle 停止生成
func (s *Server) StopGenerateHandle(c *gin.Context) {
	sessionId := c.GetHeader(types.TokenName)
	cancel := s.ReqCancelFunc[sessionId]
	cancel()
	delete(s.ReqCancelFunc, sessionId)
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success})
}

// GetHelloMsgHandle 获取角色的打招呼信息
func (s *Server) GetHelloMsgHandle(c *gin.Context) {
	role := strings.TrimSpace(c.Query("role"))
	if role == "" {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid args"})
		return
	}
	chatRole, err := GetChatRole(role)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Role not found"})
		return
	}
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: chatRole.HelloMsg})
}

// SetImgURLHandle SetImgURL 设置图片地址集合
func (s *Server) SetImgURLHandle(c *gin.Context) {
	var data struct {
		WechatCard  string `json:"wechat_card"`  // 个人微信二维码
		WechatGroup string `json:"wechat_group"` // 微信群聊二维码
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if data.WechatCard != "" {
		s.Config.ImgURL.WechatCard = data.WechatCard
	}
	if data.WechatGroup != "" {
		s.Config.ImgURL.WechatGroup = data.WechatGroup
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.ImgURL})

}

// GetImgURLHandle 获取图片地址集合
func (s *Server) GetImgURLHandle(c *gin.Context) {
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.ImgURL})
}
