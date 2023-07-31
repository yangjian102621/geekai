package handler

import (
	"bufio"
	"bytes"
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const ErrorMsg = "抱歉，AI 助手开小差了，请稍后再试。"

type ChatHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewChatHandler(app *core.AppServer, db *gorm.DB) *ChatHandler {
	handler := ChatHandler{db: db}
	handler.App = app
	return &handler
}

var chatConfig types.ChatConfig

// ChatHandle 处理聊天 WebSocket 请求
func (h *ChatHandler) ChatHandle(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	// 设置读写超时时间
	_ = ws.SetWriteDeadline(time.Now().Add(300 * time.Second))
	_ = ws.SetReadDeadline(time.Now().Add(300 * time.Second))

	sessionId := c.Query("session_id")
	roleId := h.GetInt(c, "role_id", 0)
	chatId := c.Query("chat_id")
	chatModel := c.Query("model")

	session := h.App.ChatSession.Get(sessionId)
	if session.SessionId == "" {
		user, err := utils.GetLoginUser(c, h.db)
		if err != nil {
			logger.Info("用户未登录")
			c.Abort()
			return
		}
		session = types.ChatSession{
			SessionId: sessionId,
			ClientIP:  c.ClientIP(),
			Username:  user.Username,
			UserId:    user.Id,
		}
		h.App.ChatSession.Put(sessionId, session)
	}

	// use old chat data override the chat model and role ID
	var chat model.ChatItem
	res := h.db.Where("chat_id=?", chatId).First(&chat)
	if res.Error == nil {
		chatModel = chat.Model
		roleId = int(chat.RoleId)
	}

	session.ChatId = chatId
	session.Model = chatModel
	logger.Infof("New websocket connected, IP: %s, Username: %s", c.Request.RemoteAddr, session.Username)
	client := types.NewWsClient(ws)
	var chatRole model.ChatRole
	res = h.db.First(&chatRole, roleId)
	if res.Error != nil || !chatRole.Enable {
		replyMessage(client, "当前聊天角色不存在或者未启用，连接已关闭！！！")
		c.Abort()
		return
	}

	// 初始化聊天配置
	var config model.Config
	h.db.Where("marker", "chat").First(&config)
	err = utils.JsonDecode(config.Config, &chatConfig)
	if err != nil {
		replyMessage(client, "加载系统配置失败，连接已关闭！！！")
		c.Abort()
		return
	}

	// 保存会话连接
	h.App.ChatClients.Put(sessionId, client)
	go func() {
		for {
			_, message, err := client.Receive()
			if err != nil {
				logger.Error(err)
				client.Close()
				h.App.ChatClients.Delete(sessionId)
				h.App.ReqCancelFunc.Delete(sessionId)
				return
			}
			logger.Info("Receive a message: ", string(message))
			//replyMessage(client, "这是一条测试消息！")
			ctx, cancel := context.WithCancel(context.Background())
			h.App.ReqCancelFunc.Put(sessionId, cancel)
			// 回复消息
			err = h.sendMessage(ctx, session, chatRole, string(message), client)
			if err != nil {
				logger.Error(err)
			} else {
				replyChunkMessage(client, types.WsMessage{Type: types.WsEnd})
				logger.Info("回答完毕: " + string(message))
			}

		}
	}()
}

// 将消息发送给 ChatGPT 并获取结果，通过 WebSocket 推送到客户端
func (h *ChatHandler) sendMessage(ctx context.Context, session types.ChatSession, role model.ChatRole, prompt string, ws types.Client) error {
	promptCreatedAt := time.Now() // 记录提问时间

	var user model.User
	res := h.db.Model(&model.User{}).First(&user, session.UserId)
	if res.Error != nil {
		replyMessage(ws, "非法用户，请联系管理员！")
		return res.Error
	}
	var userVo vo.User
	err := utils.CopyObject(user, &userVo)
	userVo.Id = user.Id
	if err != nil {
		return errors.New("User 对象转换失败，" + err.Error())
	}

	if userVo.Status == false {
		replyMessage(ws, "您的账号已经被禁用，如果疑问，请联系管理员！")
		replyMessage(ws, "![](/images/wx.png)")
		return nil
	}

	if userVo.Calls <= 0 && userVo.ChatConfig.ApiKey == "" {
		replyMessage(ws, "您的对话次数已经用尽，请联系管理员或者点击左下角菜单加入众筹获得100次对话！")
		replyMessage(ws, "![](/images/wx.png)")
		return nil
	}

	if userVo.ExpiredTime > 0 && userVo.ExpiredTime <= time.Now().Unix() {
		replyMessage(ws, "您的账号已经过期，请联系管理员！")
		replyMessage(ws, "![](/images/wx.png)")
		return nil
	}
	var req = types.ApiRequest{
		Model:       session.Model,
		Temperature: userVo.ChatConfig.Temperature,
		MaxTokens:   userVo.ChatConfig.MaxTokens,
		Stream:      true,
		Functions:   types.InnerFunctions,
	}

	// 加载聊天上下文
	var chatCtx []interface{}
	if userVo.ChatConfig.EnableContext {
		if h.App.ChatContexts.Has(session.ChatId) {
			chatCtx = h.App.ChatContexts.Get(session.ChatId)
		} else {
			// 加载角色信息
			var messages []types.Message
			err := utils.JsonDecode(role.Context, &messages)
			if err == nil {
				for _, v := range messages {
					chatCtx = append(chatCtx, v)
				}
			}

			// 加载最近的聊天记录作为聊天上下文
			if chatConfig.ContextDeep > 0 {
				var historyMessages []model.HistoryMessage
				res := h.db.Where("chat_id = ? and use_context = 1", session.ChatId).Limit(2).Order("created_at desc").Find(&historyMessages)
				if res.Error == nil {
					for _, msg := range historyMessages {
						ms := types.Message{Role: "user", Content: msg.Content}
						if msg.Type == types.ReplyMsg {
							ms.Role = "assistant"
						}
						chatCtx = append(chatCtx, ms)
					}
				}
			}
		}

		logger.Debugf("聊天上下文：%+v", chatCtx)
	}
	reqMgs := make([]interface{}, 0)
	for _, m := range chatCtx {
		reqMgs = append(reqMgs, m)
	}

	req.Messages = append(reqMgs, map[string]interface{}{
		"role":    "user",
		"content": prompt,
	})
	var apiKey string
	response, err := h.doRequest(ctx, userVo, &apiKey, req)
	if err != nil {
		if strings.Contains(err.Error(), "context canceled") {
			logger.Info("用户取消了请求：", prompt)
			return nil
		} else if strings.Contains(err.Error(), "no available key") {
			replyMessage(ws, "抱歉😔😔😔，系统已经没有可用的 API KEY🔑，您可以导入自己的 API KEY🔑 继续使用！🙏🙏🙏")
			return nil
		} else {
			logger.Error(err)
		}

		replyMessage(ws, ErrorMsg)
		replyMessage(ws, "![](/images/wx.png)")
		return err
	} else {
		defer response.Body.Close()
	}

	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		if true {
			replyCreatedAt := time.Now()
			// 循环读取 Chunk 消息
			var message = types.Message{}
			var contents = make([]string, 0)
			var functionCall = false
			var functionName string
			var arguments = make([]string, 0)
			reader := bufio.NewReader(response.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if strings.Contains(err.Error(), "context canceled") {
						logger.Info("用户取消了请求：", prompt)
					} else if err != io.EOF {
						logger.Error("信息读取出错：", err)
					}
					break
				}
				if !strings.Contains(line, "data:") || len(line) < 30 {
					continue
				}

				var responseBody = types.ApiResponse{}
				err = json.Unmarshal([]byte(line[6:]), &responseBody)
				if err != nil || len(responseBody.Choices) == 0 { // 数据解析出错
					logger.Error(err, line)
					replyMessage(ws, ErrorMsg)
					replyMessage(ws, "![](/images/wx.png)")
					break
				}

				fun := responseBody.Choices[0].Delta.FunctionCall
				if functionCall && fun.Name == "" {
					arguments = append(arguments, fun.Arguments)
					continue
				}

				if !utils.IsEmptyValue(fun) {
					functionCall = true
					functionName = fun.Name
					f := h.App.Functions[functionName]
					replyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
					replyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: fmt.Sprintf("正在调用函数 `%s` 作答 ...\n\n", f.Name())})
					continue
				}

				if responseBody.Choices[0].FinishReason == "function_call" { // 函数调用完毕
					break
				}

				// 初始化 role
				if responseBody.Choices[0].Delta.Role != "" && message.Role == "" {
					message.Role = responseBody.Choices[0].Delta.Role
					replyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
					continue
				} else if responseBody.Choices[0].FinishReason != "" {
					break // 输出完成或者输出中断了
				} else {
					content := responseBody.Choices[0].Delta.Content
					contents = append(contents, utils.InterfaceToString(content))
					replyChunkMessage(ws, types.WsMessage{
						Type:    types.WsMiddle,
						Content: utils.InterfaceToString(responseBody.Choices[0].Delta.Content),
					})
				}
			} // end for

			if functionCall { // 调用函数完成任务
				logger.Info(functionName)
				logger.Info(arguments)
				f := h.App.Functions[functionName]
				data, err := f.Invoke(arguments)
				if err != nil {
					msg := "调用函数出错：" + err.Error()
					replyChunkMessage(ws, types.WsMessage{
						Type:    types.WsMiddle,
						Content: msg,
					})
					contents = append(contents, msg)
				} else {
					replyChunkMessage(ws, types.WsMessage{
						Type:    types.WsMiddle,
						Content: data,
					})
					contents = append(contents, data)
				}
			}

			// 消息发送成功
			if len(contents) > 0 {
				// 更新用户的对话次数
				if userVo.ChatConfig.ApiKey == "" { // 如果用户使用的是自己绑定的 API KEY 则不扣减对话次数
					res := h.db.Model(&user).UpdateColumn("calls", gorm.Expr("calls - ?", 1))
					if res.Error != nil {
						return res.Error
					}
				}

				if message.Role == "" {
					message.Role = "assistant"
				}
				message.Content = strings.Join(contents, "")
				useMsg := types.Message{Role: "user", Content: prompt}

				// 更新上下文消息，如果是调用函数则不需要更新上下文
				if userVo.ChatConfig.EnableContext && functionCall == false {
					chatCtx = append(chatCtx, useMsg)  // 提问消息
					chatCtx = append(chatCtx, message) // 回复消息
					h.App.ChatContexts.Put(session.ChatId, chatCtx)
				}

				// 追加聊天记录
				if userVo.ChatConfig.EnableHistory {
					useContext := true
					if functionCall {
						useContext = false
					}

					// for prompt
					promptToken, err := utils.CalcTokens(prompt, req.Model)
					if err != nil {
						logger.Error(err)
					}
					historyUserMsg := model.HistoryMessage{
						UserId:     userVo.Id,
						ChatId:     session.ChatId,
						RoleId:     role.Id,
						Type:       types.PromptMsg,
						Icon:       user.Avatar,
						Content:    prompt,
						Tokens:     promptToken,
						UseContext: useContext,
					}
					historyUserMsg.CreatedAt = promptCreatedAt
					historyUserMsg.UpdatedAt = promptCreatedAt
					res := h.db.Save(&historyUserMsg)
					if res.Error != nil {
						logger.Error("failed to save prompt history message: ", res.Error)
					}

					// for reply
					// 计算本次对话消耗的总 token 数量
					var replyToken = 0
					if functionCall { // 函数名 + 参数 token
						tokens, _ := utils.CalcTokens(functionName, req.Model)
						replyToken += tokens
						tokens, _ = utils.CalcTokens(utils.InterfaceToString(arguments), req.Model)
						replyToken += tokens
					} else {
						replyToken, _ = utils.CalcTokens(message.Content, req.Model)
					}

					historyReplyMsg := model.HistoryMessage{
						UserId:     userVo.Id,
						ChatId:     session.ChatId,
						RoleId:     role.Id,
						Type:       types.ReplyMsg,
						Icon:       role.Icon,
						Content:    message.Content,
						Tokens:     replyToken,
						UseContext: useContext,
					}
					historyReplyMsg.CreatedAt = replyCreatedAt
					historyReplyMsg.UpdatedAt = replyCreatedAt
					res = h.db.Create(&historyReplyMsg)
					if res.Error != nil {
						logger.Error("failed to save reply history message: ", res.Error)
					}

					// 计算本次对话消耗的总 token 数量
					var totalTokens = 0
					if functionCall { // prompt + 函数名 + 参数 token
						totalTokens = promptToken + replyToken
					} else {
						totalTokens = replyToken + getTotalTokens(req)
					}
					//replyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: fmt.Sprintf("\n\n `本轮对话共消耗 Token 数量: %d`", totalTokens+11)})
					if userVo.ChatConfig.ApiKey != "" { // 调用自己的 API KEY 不计算 token 消耗
						h.db.Model(&user).UpdateColumn("tokens", gorm.Expr("tokens + ?",
							totalTokens))
					}
				}

				// 保存当前会话
				var chatItem model.ChatItem
				res = h.db.Where("chat_id = ?", session.ChatId).First(&chatItem)
				if res.Error != nil {
					chatItem.ChatId = session.ChatId
					chatItem.UserId = session.UserId
					chatItem.RoleId = role.Id
					chatItem.Model = session.Model
					if utf8.RuneCountInString(prompt) > 30 {
						chatItem.Title = string([]rune(prompt)[:30]) + "..."
					} else {
						chatItem.Title = prompt
					}
					h.db.Create(&chatItem)
				}
			}
		}
	} else {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error with reading response: %v", err)
		}
		var res types.ApiError
		err = json.Unmarshal(body, &res)
		if err != nil {
			return fmt.Errorf("error with decode response: %v", err)
		}

		// OpenAI API 调用异常处理
		// TODO: 是否考虑重发消息？
		if strings.Contains(res.Error.Message, "This key is associated with a deactivated account") {
			replyMessage(ws, "请求 OpenAI API 失败：API KEY 所关联的账户被禁用。")
			// 移除当前 API key
			h.db.Where("value = ?", apiKey).Delete(&model.ApiKey{})
		} else if strings.Contains(res.Error.Message, "You exceeded your current quota") {
			replyMessage(ws, "请求 OpenAI API 失败：API KEY 触发并发限制，请稍后再试。")
		} else if strings.Contains(res.Error.Message, "This model's maximum context length") {
			replyMessage(ws, "当前会话上下文长度超出限制，已为您删减会话上下文！")
			// 只保留最近的三条记录
			chatContext := h.App.ChatContexts.Get(session.ChatId)
			if len(chatContext) > 3 {
				chatContext = chatContext[len(chatContext)-3:]
			}
			h.App.ChatContexts.Put(session.ChatId, chatContext)
			return h.sendMessage(ctx, session, role, prompt, ws)
		} else {
			replyMessage(ws, "请求 OpenAI API 失败："+res.Error.Message)
		}
	}

	return nil
}

// 发送请求到 OpenAI 服务器
// useOwnApiKey: 是否使用了用户自己的 API KEY
func (h *ChatHandler) doRequest(ctx context.Context, user vo.User, apiKey *string, req types.ApiRequest) (*http.Response, error) {
	var client *http.Client
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	// 创建 HttpClient 请求对象
	request, err := http.NewRequest(http.MethodPost, h.App.ChatConfig.ApiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	request.Header.Add("Content-Type", "application/json")

	proxyURL := h.App.Config.ProxyURL
	if proxyURL == "" {
		client = &http.Client{}
	} else { // 使用代理
		proxy, _ := url.Parse(proxyURL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}
	// 查询当前用户是否导入了自己的 API KEY
	if user.ChatConfig.ApiKey != "" {
		logger.Info("使用用户自己的 API KEY: ", user.ChatConfig.ApiKey)
		*apiKey = user.ChatConfig.ApiKey
	} else { // 获取系统的 API KEY
		var key model.ApiKey
		res := h.db.Where("user_id = ?", 0).Order("last_used_at ASC").First(&key)
		if res.Error != nil {
			return nil, errors.New("no available key, please import key")
		}
		*apiKey = key.Value
		// 更新 API KEY 的最后使用时间
		h.db.Model(&key).UpdateColumn("last_used_at", time.Now().Unix())
	}

	logger.Infof("Sending OpenAI request, KEY: %s, PROXY: %s, Model: %s", *apiKey, proxyURL, req.Model)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *apiKey))
	return client.Do(request)
}

// 回复客户片段端消息
func replyChunkMessage(client types.Client, message types.WsMessage) {
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error for decoding json data: %v", err.Error())
		return
	}
	err = client.(*types.WsClient).Send(msg)
	if err != nil {
		logger.Errorf("Error for reply message: %v", err.Error())
	}
}

// 回复客户端一条完整的消息
func replyMessage(ws types.Client, message string) {
	replyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
	replyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: message})
	replyChunkMessage(ws, types.WsMessage{Type: types.WsEnd})
}

// Tokens 统计 token 数量
func (h *ChatHandler) Tokens(c *gin.Context) {
	text := c.Query("text")
	md := c.Query("model")
	tokens, err := utils.CalcTokens(text, md)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, tokens)
}

func getTotalTokens(req types.ApiRequest) int {
	encode := utils.JsonEncode(req.Messages)
	var items []map[string]interface{}
	err := utils.JsonDecode(encode, &items)
	if err != nil {
		return 0
	}
	tokens := 0
	for _, item := range items {
		content, ok := item["content"]
		if ok && !utils.IsEmptyValue(content) {
			t, err := utils.CalcTokens(utils.InterfaceToString(content), req.Model)
			if err == nil {
				tokens += t
			}
		}
	}
	return tokens
}

// StopGenerate 停止生成
func (h *ChatHandler) StopGenerate(c *gin.Context) {
	sessionId := c.Query("session_id")
	if h.App.ReqCancelFunc.Has(sessionId) {
		h.App.ReqCancelFunc.Get(sessionId)()
		h.App.ReqCancelFunc.Delete(sessionId)
	}
	resp.SUCCESS(c, types.OkMsg)
}
