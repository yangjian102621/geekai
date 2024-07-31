package chatimpl

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

type ChatHandler struct {
	handler.BaseHandler
	redis          *redis.Client
	uploadManager  *oss.UploaderManager
	licenseService *service.LicenseService
	ReqCancelFunc  *types.LMap[string, context.CancelFunc] // HttpClient 请求取消 handle function
	ChatContexts   *types.LMap[string, []types.Message]    // 聊天上下文 Map [chatId] => []Message
}

func NewChatHandler(app *core.AppServer, db *gorm.DB, redis *redis.Client, manager *oss.UploaderManager, licenseService *service.LicenseService) *ChatHandler {
	return &ChatHandler{
		BaseHandler:    handler.BaseHandler{App: app, DB: db},
		redis:          redis,
		uploadManager:  manager,
		licenseService: licenseService,
		ReqCancelFunc:  types.NewLMap[string, context.CancelFunc](),
		ChatContexts:   types.NewLMap[string, []types.Message](),
	}
}

// ChatHandle 处理聊天 WebSocket 请求
func (h *ChatHandler) ChatHandle(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	sessionId := c.Query("session_id")
	roleId := h.GetInt(c, "role_id", 0)
	chatId := c.Query("chat_id")
	modelId := h.GetInt(c, "model_id", 0)

	client := types.NewWsClient(ws)
	var chatRole model.ChatRole
	res := h.DB.First(&chatRole, roleId)
	if res.Error != nil || !chatRole.Enable {
		utils.ReplyMessage(client, "当前聊天角色不存在或者未启用，连接已关闭！！！")
		c.Abort()
		return
	}
	// if the role bind a model_id, use role's bind model_id
	if chatRole.ModelId > 0 {
		modelId = chatRole.ModelId
	}
	// get model info
	var chatModel model.ChatModel
	res = h.DB.First(&chatModel, modelId)
	if res.Error != nil || chatModel.Enabled == false {
		utils.ReplyMessage(client, "当前AI模型暂未启用，连接已关闭！！！")
		c.Abort()
		return
	}

	session := &types.ChatSession{
		SessionId: sessionId,
		ClientIP:  c.ClientIP(),
		UserId:    h.GetLoginUserId(c),
	}

	// use old chat data override the chat model and role ID
	var chat model.ChatItem
	res = h.DB.Where("chat_id = ?", chatId).First(&chat)
	if res.Error == nil {
		chatModel.Id = chat.ModelId
		roleId = int(chat.RoleId)
	}

	session.ChatId = chatId
	session.Model = types.ChatModel{
		Id:          chatModel.Id,
		Name:        chatModel.Name,
		Value:       chatModel.Value,
		Power:       chatModel.Power,
		MaxTokens:   chatModel.MaxTokens,
		MaxContext:  chatModel.MaxContext,
		Temperature: chatModel.Temperature,
		KeyId:       chatModel.KeyId}
	logger.Infof("New websocket connected, IP: %s", c.ClientIP())

	go func() {
		for {
			_, msg, err := client.Receive()
			if err != nil {
				logger.Debugf("close connection: %s", client.Conn.RemoteAddr())
				client.Close()
				cancelFunc := h.ReqCancelFunc.Get(sessionId)
				if cancelFunc != nil {
					cancelFunc()
					h.ReqCancelFunc.Delete(sessionId)
				}
				return
			}

			var message types.WsMessage
			err = utils.JsonDecode(string(msg), &message)
			if err != nil {
				continue
			}

			// 心跳消息
			if message.Type == "heartbeat" {
				logger.Debug("收到 Chat 心跳消息：", message.Content)
				continue
			}

			logger.Info("Receive a message: ", message.Content)

			ctx, cancel := context.WithCancel(context.Background())
			h.ReqCancelFunc.Put(sessionId, cancel)
			// 回复消息
			err = h.sendMessage(ctx, session, chatRole, utils.InterfaceToString(message.Content), client)
			if err != nil {
				logger.Error(err)
				utils.ReplyMessage(client, err.Error())
			} else {
				utils.ReplyChunkMessage(client, types.WsMessage{Type: types.WsEnd})
				logger.Infof("回答完毕: %v", message.Content)
			}

		}
	}()
}

func (h *ChatHandler) sendMessage(ctx context.Context, session *types.ChatSession, role model.ChatRole, prompt string, ws *types.WsClient) error {
	if !h.App.Debug {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recover message from error: ", r)
			}
		}()
	}

	var user model.User
	res := h.DB.Model(&model.User{}).First(&user, session.UserId)
	if res.Error != nil {
		return errors.New("未授权用户，您正在进行非法操作！")
	}
	var userVo vo.User
	err := utils.CopyObject(user, &userVo)
	userVo.Id = user.Id
	if err != nil {
		return errors.New("User 对象转换失败，" + err.Error())
	}

	if userVo.Status == false {
		return errors.New("您的账号已经被禁用，如果疑问，请联系管理员！")
	}

	if userVo.Power < session.Model.Power {
		return fmt.Errorf("您当前剩余算力 %d 已不足以支付当前模型的单次对话需要消耗的算力 %d，[立即购买](/member)。", userVo.Power, session.Model.Power)
	}

	if userVo.ExpiredTime > 0 && userVo.ExpiredTime <= time.Now().Unix() {
		return errors.New("您的账号已经过期，请联系管理员！")
	}

	// 检查 prompt 长度是否超过了当前模型允许的最大上下文长度
	promptTokens, err := utils.CalcTokens(prompt, session.Model.Value)
	if promptTokens > session.Model.MaxContext {

		return errors.New("对话内容超出了当前模型允许的最大上下文长度！")
	}

	var req = types.ApiRequest{
		Model:  session.Model.Value,
		Stream: true,
	}
	req.Temperature = session.Model.Temperature
	req.MaxTokens = session.Model.MaxTokens
	// OpenAI 支持函数功能
	var items []model.Function
	res = h.DB.Where("enabled", true).Find(&items)
	if res.Error == nil {
		var tools = make([]types.Tool, 0)
		for _, v := range items {
			var parameters map[string]interface{}
			err = utils.JsonDecode(v.Parameters, &parameters)
			if err != nil {
				continue
			}
			tool := types.Tool{
				Type: "function",
				Function: types.Function{
					Name:        v.Name,
					Description: v.Description,
					Parameters:  parameters,
				},
			}
			if v, ok := parameters["required"]; v == nil || !ok {
				tool.Function.Parameters["required"] = []string{}
			}
			tools = append(tools, tool)
		}

		if len(tools) > 0 {
			req.Tools = tools
			req.ToolChoice = "auto"
		}
	}

	// 加载聊天上下文
	chatCtx := make([]types.Message, 0)
	messages := make([]types.Message, 0)
	if h.App.SysConfig.EnableContext {
		if h.ChatContexts.Has(session.ChatId) {
			messages = h.ChatContexts.Get(session.ChatId)
		} else {
			_ = utils.JsonDecode(role.Context, &messages)
			if h.App.SysConfig.ContextDeep > 0 {
				var historyMessages []model.ChatMessage
				res := h.DB.Where("chat_id = ? and use_context = 1", session.ChatId).Limit(h.App.SysConfig.ContextDeep).Order("id DESC").Find(&historyMessages)
				if res.Error == nil {
					for i := len(historyMessages) - 1; i >= 0; i-- {
						msg := historyMessages[i]
						ms := types.Message{Role: "user", Content: msg.Content}
						if msg.Type == types.ReplyMsg {
							ms.Role = "assistant"
						}
						chatCtx = append(chatCtx, ms)
					}
				}
			}
		}

		// 计算当前请求的 token 总长度，确保不会超出最大上下文长度
		// MaxContextLength = Response + Tool + Prompt + Context
		tokens := req.MaxTokens // 最大响应长度
		tks, _ := utils.CalcTokens(utils.JsonEncode(req.Tools), req.Model)
		tokens += tks + promptTokens

		for _, v := range messages {
			tks, _ := utils.CalcTokens(v.Content, req.Model)
			// 上下文 token 超出了模型的最大上下文长度
			if tokens+tks >= session.Model.MaxContext {
				break
			}

			// 上下文的深度超出了模型的最大上下文深度
			if len(chatCtx) >= h.App.SysConfig.ContextDeep {
				break
			}

			tokens += tks
			chatCtx = append(chatCtx, v)
		}

		logger.Debugf("聊天上下文：%+v", chatCtx)
	}
	reqMgs := make([]interface{}, 0)
	for _, m := range chatCtx {
		reqMgs = append(reqMgs, m)
	}

	fullPrompt := prompt
	text := prompt
	// extract files in prompt
	files := utils.ExtractFileURLs(prompt)
	logger.Debugf("detected FILES: %+v", files)
	// 如果不是逆向模型，则提取文件内容
	if len(files) > 0 && !(session.Model.Value == "gpt-4-all" ||
		strings.HasPrefix(session.Model.Value, "gpt-4-gizmo") ||
		strings.HasSuffix(session.Model.Value, "claude-3")) {
		contents := make([]string, 0)
		var file model.File
		for _, v := range files {
			h.DB.Where("url = ?", v).First(&file)
			content, err := utils.ReadFileContent(v, h.App.Config.TikaHost)
			if err != nil {
				logger.Error("error with read file: ", err)
			} else {
				contents = append(contents, fmt.Sprintf("%s 文件内容：%s", file.Name, content))
			}
			text = strings.Replace(text, v, "", 1)
		}
		if len(contents) > 0 {
			fullPrompt = fmt.Sprintf("请根据提供的文件内容信息回答问题(其中Excel 已转成 HTML)：\n\n %s\n\n 问题：%s", strings.Join(contents, "\n"), text)
		}

		tokens, _ := utils.CalcTokens(fullPrompt, req.Model)
		if tokens > session.Model.MaxContext {
			return fmt.Errorf("文件的长度超出模型允许的最大上下文长度，请减少文件内容数量或文件大小。")
		}
	}
	logger.Debug("最终Prompt：", fullPrompt)

	// extract images from prompt
	imgURLs := utils.ExtractImgURLs(prompt)
	logger.Debugf("detected IMG: %+v", imgURLs)
	var content interface{}
	if len(imgURLs) > 0 {
		data := make([]interface{}, 0)
		for _, v := range imgURLs {
			text = strings.Replace(text, v, "", 1)
			data = append(data, gin.H{
				"type": "image_url",
				"image_url": gin.H{
					"url": v,
				},
			})
		}
		data = append(data, gin.H{
			"type": "text",
			"text": strings.TrimSpace(text),
		})
		content = data
	} else {
		content = fullPrompt
	}
	req.Messages = append(reqMgs, map[string]interface{}{
		"role":    "user",
		"content": content,
	})

	logger.Debugf("%+v", req.Messages)

	return h.sendOpenAiMessage(chatCtx, req, userVo, ctx, session, role, prompt, ws)
}

// Tokens 统计 token 数量
func (h *ChatHandler) Tokens(c *gin.Context) {
	var data struct {
		Text   string `json:"text"`
		Model  string `json:"model"`
		ChatId string `json:"chat_id"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 如果没有传入 text 字段，则说明是获取当前 reply 总的 token 消耗（带上下文）
	if data.Text == "" && data.ChatId != "" {
		var item model.ChatMessage
		userId, _ := c.Get(types.LoginUserID)
		res := h.DB.Where("user_id = ?", userId).Where("chat_id = ?", data.ChatId).Last(&item)
		if res.Error != nil {
			resp.ERROR(c, res.Error.Error())
			return
		}
		resp.SUCCESS(c, item.Tokens)
		return
	}

	tokens, err := utils.CalcTokens(data.Text, data.Model)
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
	if h.ReqCancelFunc.Has(sessionId) {
		h.ReqCancelFunc.Get(sessionId)()
		h.ReqCancelFunc.Delete(sessionId)
	}
	resp.SUCCESS(c, types.OkMsg)
}

// 发送请求到 OpenAI 服务器
// useOwnApiKey: 是否使用了用户自己的 API KEY
func (h *ChatHandler) doRequest(ctx context.Context, req types.ApiRequest, session *types.ChatSession, apiKey *model.ApiKey) (*http.Response, error) {
	// if the chat model bind a KEY, use it directly
	if session.Model.KeyId > 0 {
		h.DB.Where("id", session.Model.KeyId).Find(apiKey)
	}
	// use the last unused key
	if apiKey.Id == 0 {
		h.DB.Where("type", "chat").Where("enabled", true).Order("last_used_at ASC").First(apiKey)
	}
	if apiKey.Id == 0 {
		return nil, errors.New("no available key, please import key")
	}

	// ONLY allow apiURL in blank list
	err := h.licenseService.IsValidApiURL(apiKey.ApiURL)
	if err != nil {
		return nil, err
	}
	logger.Debugf(utils.JsonEncode(req))

	apiURL := fmt.Sprintf("%s/v1/chat/completions", apiKey.ApiURL)
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

	request = request.WithContext(ctx)
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
	logger.Debugf("Sending %s request, API KEY:%s, PROXY: %s, Model: %s", apiKey.ApiURL, apiURL, apiKey.ProxyURL, req.Model)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey.Value))
	// 更新API KEY 最后使用时间
	h.DB.Model(&model.ApiKey{}).Where("id", apiKey.Id).UpdateColumn("last_used_at", time.Now().Unix())
	return client.Do(request)
}

// 扣减用户算力
func (h *ChatHandler) subUserPower(userVo vo.User, session *types.ChatSession, promptTokens int, replyTokens int) {
	power := 1
	if session.Model.Power > 0 {
		power = session.Model.Power
	}
	res := h.DB.Model(&model.User{}).Where("id = ?", userVo.Id).UpdateColumn("power", gorm.Expr("power - ?", power))
	if res.Error == nil {
		// 记录算力消费日志
		var u model.User
		h.DB.Where("id", userVo.Id).First(&u)
		h.DB.Create(&model.PowerLog{
			UserId:    userVo.Id,
			Username:  userVo.Username,
			Type:      types.PowerConsume,
			Amount:    power,
			Mark:      types.PowerSub,
			Balance:   u.Power,
			Model:     session.Model.Value,
			Remark:    fmt.Sprintf("模型名称：%s, 提问长度：%d，回复长度：%d", session.Model.Name, promptTokens, replyTokens),
			CreatedAt: time.Now(),
		})
	}

}

func (h *ChatHandler) saveChatHistory(
	req types.ApiRequest,
	prompt string,
	contents []string,
	message types.Message,
	chatCtx []types.Message,
	session *types.ChatSession,
	role model.ChatRole,
	userVo vo.User,
	promptCreatedAt time.Time,
	replyCreatedAt time.Time) {
	if message.Role == "" {
		message.Role = "assistant"
	}
	message.Content = strings.Join(contents, "")
	useMsg := types.Message{Role: "user", Content: prompt}

	// 更新上下文消息，如果是调用函数则不需要更新上下文
	if h.App.SysConfig.EnableContext {
		chatCtx = append(chatCtx, useMsg)  // 提问消息
		chatCtx = append(chatCtx, message) // 回复消息
		h.ChatContexts.Put(session.ChatId, chatCtx)
	}

	// 追加聊天记录
	// for prompt
	promptToken, err := utils.CalcTokens(prompt, req.Model)
	if err != nil {
		logger.Error(err)
	}
	historyUserMsg := model.ChatMessage{
		UserId:     userVo.Id,
		ChatId:     session.ChatId,
		RoleId:     role.Id,
		Type:       types.PromptMsg,
		Icon:       userVo.Avatar,
		Content:    template.HTMLEscapeString(prompt),
		Tokens:     promptToken,
		UseContext: true,
		Model:      req.Model,
	}
	historyUserMsg.CreatedAt = promptCreatedAt
	historyUserMsg.UpdatedAt = promptCreatedAt
	res := h.DB.Save(&historyUserMsg)
	if res.Error != nil {
		logger.Error("failed to save prompt history message: ", res.Error)
	}

	// for reply
	// 计算本次对话消耗的总 token 数量
	replyTokens, _ := utils.CalcTokens(message.Content, req.Model)
	totalTokens := replyTokens + getTotalTokens(req)
	historyReplyMsg := model.ChatMessage{
		UserId:     userVo.Id,
		ChatId:     session.ChatId,
		RoleId:     role.Id,
		Type:       types.ReplyMsg,
		Icon:       role.Icon,
		Content:    message.Content,
		Tokens:     totalTokens,
		UseContext: true,
		Model:      req.Model,
	}
	historyReplyMsg.CreatedAt = replyCreatedAt
	historyReplyMsg.UpdatedAt = replyCreatedAt
	res = h.DB.Create(&historyReplyMsg)
	if res.Error != nil {
		logger.Error("failed to save reply history message: ", res.Error)
	}

	// 更新用户算力
	if session.Model.Power > 0 {
		h.subUserPower(userVo, session, promptToken, replyTokens)
	}
	// 保存当前会话
	var chatItem model.ChatItem
	res = h.DB.Where("chat_id = ?", session.ChatId).First(&chatItem)
	if res.Error != nil {
		chatItem.ChatId = session.ChatId
		chatItem.UserId = userVo.Id
		chatItem.RoleId = role.Id
		chatItem.ModelId = session.Model.Id
		if utf8.RuneCountInString(prompt) > 30 {
			chatItem.Title = string([]rune(prompt)[:30]) + "..."
		} else {
			chatItem.Title = prompt
		}
		chatItem.Model = req.Model
		h.DB.Create(&chatItem)
	}
}

// 将AI回复消息中生成的图片链接下载到本地
func (h *ChatHandler) extractImgUrl(text string) string {
	pattern := `!\[([^\]]*)]\(([^)]+)\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	// 下载图片并替换链接地址
	for _, match := range matches {
		imageURL := match[2]
		logger.Debug(imageURL)
		// 对于相同地址的图片，已经被替换了，就不再重复下载了
		if !strings.Contains(text, imageURL) {
			continue
		}

		newImgURL, err := h.uploadManager.GetUploadHandler().PutUrlFile(imageURL, false)
		if err != nil {
			logger.Error("error with download image: ", err)
			continue
		}

		text = strings.ReplaceAll(text, imageURL, newImgURL)
	}
	return text
}
