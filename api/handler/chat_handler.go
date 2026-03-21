package handler

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
	"geekai/service"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

const (
	ChatEventStart        = "start"
	ChatEventEnd          = "end"
	ChatEventError        = "error"
	ChatEventMessageDelta = "message_delta"
	ChatEventTitle        = "title"
)

type ChatInput struct {
	UserId    uint            `json:"user_id"`
	RoleId    uint            `json:"role_id"`
	ModelId   uint            `json:"model_id"`
	ChatId    string          `json:"chat_id"`
	Prompt    string          `json:"prompt"`
	Tools     []uint          `json:"tools"`
	Stream    bool            `json:"stream"`
	Files     []vo.File       `json:"files"`
	ChatModel model.ChatModel `json:"chat_model,omitempty"`
	ChatRole  model.ChatRole  `json:"chat_role,omitempty"`
	LastMsgId uint            `json:"last_msg_id,omitempty"` // 最后的消息ID，用于重新生成答案的时候过滤上下文
}

type ChatHandler struct {
	BaseHandler
	redis          *redis.Client
	uploadManager  *oss.UploaderManager
	licenseService *service.LicenseService
	ReqCancelFunc  *types.LMap[string, context.CancelFunc] // HttpClient 请求取消 handle function
	userService    *service.UserService
}

func NewChatHandler(app *core.AppServer, db *gorm.DB, redis *redis.Client, manager *oss.UploaderManager, licenseService *service.LicenseService, userService *service.UserService) *ChatHandler {
	return &ChatHandler{
		BaseHandler:    BaseHandler{App: app, DB: db},
		redis:          redis,
		uploadManager:  manager,
		licenseService: licenseService,
		ReqCancelFunc:  types.NewLMap[string, context.CancelFunc](),
		userService:    userService,
	}
}

// Chat 处理聊天请求
func (h *ChatHandler) Chat(c *gin.Context) {
	var input ChatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 设置SSE响应头
	c.Header("Prompt-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// 使用旧的聊天数据覆盖模型和角色ID
	var chat model.ChatItem
	h.DB.Where("chat_id", input.ChatId).First(&chat)
	if chat.Id > 0 {
		input.ModelId = chat.ModelId
		input.RoleId = chat.RoleId
	}

	// 验证聊天角色
	var chatRole model.ChatRole
	err := h.DB.First(&chatRole, input.RoleId).Error
	if err != nil || !chatRole.Enable {
		pushMessage(c, ChatEventError, "当前聊天角色不存在或者未启用，请更换角色之后再发起对话！")
		return
	}
	input.ChatRole = chatRole

	// 获取模型信息
	var chatModel model.ChatModel
	err = h.DB.Where("id", input.ModelId).First(&chatModel).Error
	if err != nil || !chatModel.Enabled {
		pushMessage(c, ChatEventError, "当前AI模型暂未启用，请更换模型后再发起对话！")
		return
	}
	input.ChatModel = chatModel

	// 发送消息
	err = h.sendMessage(ctx, input, c)
	if err != nil {
		pushMessage(c, ChatEventError, err.Error())
		return
	}

	pushMessage(c, ChatEventEnd, "对话完成")
}

func pushMessage(c *gin.Context, msgType string, content interface{}) {
	c.SSEvent("message", map[string]interface{}{
		"type": msgType,
		"body": content,
	})
	c.Writer.Flush()
}

func (h *ChatHandler) sendMessage(ctx context.Context, input ChatInput, c *gin.Context) error {
	var user model.User
	res := h.DB.Model(&model.User{}).First(&user, input.UserId)
	if res.Error != nil {
		return errors.New("未授权用户，您正在进行非法操作！")
	}
	var userVo vo.User
	err := utils.CopyObject(user, &userVo)
	userVo.Id = user.Id
	if err != nil {
		return errors.New("User 对象转换失败，" + err.Error())
	}

	if !userVo.Status {
		return errors.New("您的账号已经被禁用，如果疑问，请联系管理员！")
	}

	if userVo.Power < input.ChatModel.Power {
		return fmt.Errorf("您当前剩余算力 %d 已不足以支付当前模型的单次对话需要消耗的算力 %d，[立即购买](/member)。", userVo.Power, input.ChatModel.Power)
	}

	if userVo.ExpiredTime > 0 && userVo.ExpiredTime <= time.Now().Unix() {
		return errors.New("您的账号已经过期，请联系管理员！")
	}

	// 检查 prompt 长度是否超过了当前模型允许的最大上下文长度
	promptTokens, _ := utils.CalcTokens(input.Prompt, input.ChatModel.Value)
	if promptTokens > input.ChatModel.MaxContext {

		return errors.New("对话内容超出了当前模型允许的最大上下文长度！")
	}

	var req = types.ApiRequest{
		Model:       input.ChatModel.Value,
		Stream:      input.Stream,
		Temperature: input.ChatModel.Temperature,
	}
	// 兼容 MiniMax 模型，温度范围 0-1
	if strings.HasPrefix(strings.ToLower(input.ChatModel.Value), "minimax") {
		if req.Temperature > 1.0 {
			req.Temperature = 1.0
		}
		if req.Temperature < 0 {
			req.Temperature = 0
		}
	}
	// 兼容 OpenAI 模型
	if strings.HasPrefix(input.ChatModel.Value, "o1-") ||
		strings.HasPrefix(input.ChatModel.Value, "o3-") ||
		strings.HasPrefix(input.ChatModel.Value, "gpt") {
		req.MaxCompletionTokens = input.ChatModel.MaxTokens
	} else {
		req.MaxTokens = input.ChatModel.MaxTokens
	}

	if len(input.Tools) > 0 && !strings.HasPrefix(input.ChatModel.Value, "o1-") {
		var items []model.Function
		res = h.DB.Where("enabled", true).Where("id IN ?", input.Tools).Find(&items)
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
	}

	// 加载聊天上下文
	chatCtx := make([]any, 0)
	messages := make([]any, 0)
	if h.App.SysConfig.EnableContext {
		_ = utils.JsonDecode(input.ChatRole.Context, &messages)
		if h.App.SysConfig.ContextDeep > 0 {
			var historyMessages []model.ChatMessage
			dbSession := h.DB.Session(&gorm.Session{}).Where("chat_id", input.ChatId)
			if input.LastMsgId > 0 { // 重新生成逻辑
				dbSession = dbSession.Where("id < ?", input.LastMsgId)
				// 删除对应的聊天记录
				h.DB.Debug().Where("chat_id", input.ChatId).Where("id >= ?", input.LastMsgId).Delete(&model.ChatMessage{})
			}
			err = dbSession.Limit(h.App.SysConfig.ContextDeep).Order("id DESC").Find(&historyMessages).Error
			if err == nil {
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

		// 计算当前请求的 token 总长度，确保不会超出最大上下文长度
		// MaxContextLength = Response + Tool + Prompt + Context
		tokens := req.MaxTokens // 最大响应长度
		tks, _ := utils.CalcTokens(utils.JsonEncode(req.Tools), req.Model)
		tokens += tks + promptTokens

		for i := len(messages) - 1; i >= 0; i-- {
			v := messages[i]
			tks, _ = utils.CalcTokens(utils.JsonEncode(v), req.Model)
			// 上下文 token 超出了模型的最大上下文长度
			if tokens+tks >= input.ChatModel.MaxContext {
				break
			}

			// 上下文的深度超出了模型的最大上下文深度
			if len(chatCtx) >= h.App.SysConfig.ContextDeep {
				break
			}

			tokens += tks
			chatCtx = append(chatCtx, v)
		}
	}
	reqMgs := make([]any, 0)

	for i := len(chatCtx) - 1; i >= 0; i-- {
		reqMgs = append(reqMgs, chatCtx[i])
	}

	fileContents := make([]string, 0) // 文件内容
	var finalPrompt = input.Prompt
	imgList := make([]any, 0)
	for _, file := range input.Files {
		logger.Debugf("detected file: %+v", file.URL)
		// 处理图片
		if isImageURL(file.URL) {
			imgList = append(imgList, gin.H{
				"type": "image_url",
				"image_url": gin.H{
					"url": file.URL,
				},
			})
		} else {
			// 如果不是逆向模型，则提取文件内容
			modelValue := input.ChatModel.Value
			if !(strings.Contains(modelValue, "-all") || strings.HasPrefix(modelValue, "gpt-4-gizmo") || strings.HasPrefix(modelValue, "claude")) {
				content, err := utils.ReadFileContent(file.URL, h.App.Config.TikaHost)
				if err != nil {
					logger.Error("error with read file: ", err)
					continue
				} else {
					fileContents = append(fileContents, fmt.Sprintf("%s 文件内容：%s", file.Name, content))
				}
			}
		}
	}

	if len(fileContents) > 0 {
		finalPrompt = fmt.Sprintf("请根据提供的文件内容信息回答问题(其中Excel 已转成 HTML)：\n\n %s\n\n 问题：%s", strings.Join(fileContents, "\n"), input.Prompt)
		tokens, _ := utils.CalcTokens(finalPrompt, req.Model)
		if tokens > input.ChatModel.MaxContext {
			return fmt.Errorf("文件的长度超出模型允许的最大上下文长度，请减少文件内容数量或文件大小。")
		}
	} else {
		finalPrompt = input.Prompt
	}

	if len(imgList) > 0 {
		imgList = append(imgList, map[string]interface{}{
			"type": "text",
			"text": input.Prompt,
		})
		req.Messages = append(reqMgs, map[string]interface{}{
			"role":    "user",
			"content": imgList,
		})
	} else {
		req.Messages = append(reqMgs, map[string]interface{}{
			"role":    "user",
			"content": finalPrompt,
		})
	}

	return h.sendOpenAiMessage(req, userVo, ctx, input, c)
}

// 判断一个 URL 是否图片链接
func isImageURL(url string) bool {
	// 检查是否是有效的URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	// 检查文件扩展名
	ext := strings.ToLower(path.Ext(url))
	validImageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".svg":  true,
		".ico":  true,
	}

	if !validImageExts[ext] {
		return false
	}

	// 发送HEAD请求检查Content-Type
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "image/")
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
func (h *ChatHandler) doRequest(ctx context.Context, req types.ApiRequest, input ChatInput, apiKey *model.ApiKey) (*http.Response, error) {
	// if the chat model bind a KEY, use it directly
	if input.ChatModel.KeyId > 0 {
		h.DB.Where("id", input.ChatModel.KeyId).Find(apiKey)
	} else { // use the last unused key
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
	logger.Debugf("对话请求消息体：%+v", req)
	var apiURL string
	p, _ := url.Parse(apiKey.ApiURL)
	// 如果设置的是 BASE_URL 没有路径，则添加 /v1/chat/completions
	if p.Path == "" {
		apiURL = fmt.Sprintf("%s/v1/chat/completions", apiKey.ApiURL)
	} else {
		apiURL = apiKey.ApiURL
	}
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
	logger.Infof("Sending %s request, API KEY:%s, PROXY: %s, Model: %s", apiKey.ApiURL, apiURL, apiKey.ProxyURL, req.Model)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey.Value))
	// 更新API KEY 最后使用时间
	h.DB.Model(&model.ApiKey{}).Where("id", apiKey.Id).UpdateColumn("last_used_at", time.Now().Unix())
	return client.Do(request)
}

// 扣减用户算力
func (h *ChatHandler) subUserPower(userVo vo.User, input ChatInput, promptTokens int, replyTokens int) {
	power := 1
	if input.ChatModel.Power > 0 {
		power = input.ChatModel.Power
	}

	err := h.userService.DecreasePower(userVo.Id, power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  input.ChatModel.Value,
		Remark: fmt.Sprintf("模型名称：%s, 提问长度：%d，回复长度：%d", input.ChatModel.Name, promptTokens, replyTokens),
	})
	if err != nil {
		logger.Error(err)
	}
}

func (h *ChatHandler) saveChatHistory(
	req types.ApiRequest,
	usage Usage,
	message types.Message,
	input ChatInput,
	userVo vo.User,
	promptCreatedAt time.Time,
	replyCreatedAt time.Time) {

	// 追加聊天记录
	// for prompt
	var promptTokens, replyTokens, totalTokens int
	if usage.PromptTokens > 0 {
		promptTokens = usage.PromptTokens
	} else {
		promptTokens, _ = utils.CalcTokens(usage.Content, req.Model)
	}

	historyUserMsg := model.ChatMessage{
		UserId: userVo.Id,
		ChatId: input.ChatId,
		RoleId: input.RoleId,
		Type:   types.PromptMsg,
		Icon:   userVo.Avatar,
		Content: utils.JsonEncode(vo.MsgContent{
			Text:  usage.Prompt,
			Files: input.Files,
		}),
		Tokens:      promptTokens,
		TotalTokens: promptTokens,
		UseContext:  true,
		Model:       req.Model,
	}
	historyUserMsg.CreatedAt = promptCreatedAt
	historyUserMsg.UpdatedAt = promptCreatedAt
	err := h.DB.Save(&historyUserMsg).Error
	if err != nil {
		logger.Error("failed to save prompt history message: ", err)
	}

	// for reply
	// 计算本次对话消耗的总 token 数量
	if usage.CompletionTokens > 0 {
		replyTokens = usage.CompletionTokens
		totalTokens = usage.TotalTokens
	} else {
		replyTokens, _ = utils.CalcTokens(message.Content, req.Model)
		totalTokens = replyTokens + getTotalTokens(req)
	}
	historyReplyMsg := model.ChatMessage{
		UserId: userVo.Id,
		ChatId: input.ChatId,
		RoleId: input.RoleId,
		Type:   types.ReplyMsg,
		Icon:   input.ChatRole.Icon,
		Content: utils.JsonEncode(vo.MsgContent{
			Text:  message.Content,
			Files: input.Files,
		}),
		Tokens:      replyTokens,
		TotalTokens: totalTokens,
		UseContext:  true,
		Model:       req.Model,
	}
	historyReplyMsg.CreatedAt = replyCreatedAt
	historyReplyMsg.UpdatedAt = replyCreatedAt
	err = h.DB.Create(&historyReplyMsg).Error
	if err != nil {
		logger.Error("failed to save reply history message: ", err)
	}

	// 更新用户算力
	if input.ChatModel.Power > 0 {
		h.subUserPower(userVo, input, promptTokens, replyTokens)
	}
	// 保存当前会话
	var chatItem model.ChatItem
	err = h.DB.Where("chat_id = ?", input.ChatId).First(&chatItem).Error
	if err != nil {
		chatItem.ChatId = input.ChatId
		chatItem.UserId = userVo.Id
		chatItem.RoleId = input.RoleId
		chatItem.ModelId = input.ModelId
		if utf8.RuneCountInString(usage.Prompt) > 30 {
			chatItem.Title = string([]rune(usage.Prompt)[:30]) + "..."
		} else {
			chatItem.Title = usage.Prompt
		}
		chatItem.Model = req.Model
		err = h.DB.Create(&chatItem).Error
		if err != nil {
			logger.Error("failed to save chat item: ", err)
		}
	}
}

// TextToSpeech 文本生成语音
func (h *ChatHandler) TextToSpeech(c *gin.Context) {
	var data struct {
		ModelId int    `json:"model_id"`
		Text    string `json:"text"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	textHash := utils.Sha256(fmt.Sprintf("%d/%s", data.ModelId, data.Text))
	audioFile := fmt.Sprintf("%s/audio", h.App.Config.StaticDir)
	if _, err := os.Stat(audioFile); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if err := os.MkdirAll(audioFile, 0755); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	audioFile = fmt.Sprintf("%s/%s.mp3", audioFile, textHash)
	if _, err := os.Stat(audioFile); err == nil {
		// 设置响应头
		c.Header("Prompt-Type", "audio/mpeg")
		c.Header("Prompt-Disposition", "attachment; filename=speech.mp3")
		c.File(audioFile)
		return
	}

	// 查询模型
	var chatModel model.ChatModel
	err := h.DB.Where("id", data.ModelId).First(&chatModel).Error
	if err != nil {
		resp.ERROR(c, "找不到语音模型")
		return
	}

	// 调用 DeepSeek 的 API 接口
	var apiKey model.ApiKey
	if chatModel.KeyId > 0 {
		h.DB.Where("id", chatModel.KeyId).First(&apiKey)
	}
	if apiKey.Id == 0 {
		h.DB.Where("type", "tts").Where("enabled", true).First(&apiKey)
	}
	if apiKey.Id == 0 {
		resp.ERROR(c, "no TTS API key, please import key")
		return
	}

	logger.Debugf("chatModel: %+v, apiKey: %+v", chatModel, apiKey)

	// 调用 openai tts api
	config := openai.DefaultConfig(apiKey.Value)
	config.BaseURL = apiKey.ApiURL + "/v1"
	client := openai.NewClientWithConfig(config)
	voice := openai.VoiceAlloy
	var options map[string]string
	err = utils.JsonDecode(chatModel.Options, &options)
	if err == nil {
		voice = openai.SpeechVoice(options["voice"])
	}
	req := openai.CreateSpeechRequest{
		Model: openai.SpeechModel(chatModel.Value),
		Input: data.Text,
		Voice: voice,
	}

	audioData, err := client.CreateSpeech(context.Background(), req)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 先将音频数据读取到内存
	audioBytes, err := io.ReadAll(audioData)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 保存到音频文件
	err = os.WriteFile(audioFile, audioBytes, 0644)
	if err != nil {
		logger.Error("failed to save audio file: ", err)
	}

	// 设置响应头
	c.Header("Prompt-Type", "audio/mpeg")
	c.Header("Prompt-Disposition", "attachment; filename=speech.mp3")

	// 直接写入完整的音频数据到响应
	_, err = c.Writer.Write(audioBytes)
	if err != nil {
		logger.Error("写入音频数据到响应失败:", err)
	}
}

// // OPenAI 消息发送实现
// func (h *ChatHandler) sendOpenAiMessage(
// 	req types.ApiRequest,
// 	userVo vo.User,
// 	ctx context.Context,
// 	session *types.ChatSession,
// 	role model.ChatRole,
// 	prompt string,
// 	c *gin.Context) error {
// 	promptCreatedAt := time.Now() // 记录提问时间
// 	start := time.Now()
// 	var apiKey = model.ApiKey{}
// 	response, err := h.doRequest(ctx, req, session, &apiKey)
// 	logger.Info("HTTP请求完成，耗时：", time.Since(start))
// 	if err != nil {
// 		if strings.Contains(err.Error(), "context canceled") {
// 			return fmt.Errorf("用户取消了请求：%s", prompt)
// 		} else if strings.Contains(err.Error(), "no available key") {
// 			return errors.New("抱歉😔😔😔，系统已经没有可用的 API KEY，请联系管理员！")
// 		}
// 		return err
// 	} else {
// 		defer response.Body.Close()
// 	}

// 	if response.StatusCode != 200 {
// 		body, _ := io.ReadAll(response.Body)
// 		return fmt.Errorf("请求 OpenAI API 失败：%d, %v", response.StatusCode, string(body))
// 	}

// 	contentType := response.Header.Get("Prompt-Type")
// 	if strings.Contains(contentType, "text/event-stream") {
// 		replyCreatedAt := time.Now() // 记录回复时间
// 		// 循环读取 Chunk 消息
// 		var message = types.Message{Role: "assistant"}
// 		var contents = make([]string, 0)
// 		var function model.Function
// 		var toolCall = false
// 		var arguments = make([]string, 0)
// 		var reasoning = false

// 		pushMessage(c, ChatEventStart, "开始响应")
// 		scanner := bufio.NewScanner(response.Body)
// 		for scanner.Scan() {
// 			line := scanner.Text()
// 			if !strings.Contains(line, "data:") || len(line) < 30 {
// 				continue
// 			}
// 			var responseBody = types.ApiResponse{}
// 			err = json.Unmarshal([]byte(line[6:]), &responseBody)
// 			if err != nil { // 数据解析出错
// 				return errors.New(line)
// 			}
// 			if len(responseBody.Choices) == 0 { // Fixed: 兼容 Azure API 第一个输出空行
// 				continue
// 			}
// 			if responseBody.Choices[0].Delta.Prompt == nil &&
// 				responseBody.Choices[0].Delta.ToolCalls == nil &&
// 				responseBody.Choices[0].Delta.ReasoningContent == "" {
// 				continue
// 			}

// 			if responseBody.Choices[0].FinishReason == "stop" && len(contents) == 0 {
// 				pushMessage(c, ChatEventError, "抱歉😔😔😔，AI助手由于未知原因已经停止输出内容。")
// 				break
// 			}

// 			var tool types.ToolCall
// 			if len(responseBody.Choices[0].Delta.ToolCalls) > 0 {
// 				tool = responseBody.Choices[0].Delta.ToolCalls[0]
// 				if toolCall && tool.Function.Name == "" {
// 					arguments = append(arguments, tool.Function.Arguments)
// 					continue
// 				}
// 			}

// 			// 兼容 Function Call
// 			fun := responseBody.Choices[0].Delta.FunctionCall
// 			if fun.Name != "" {
// 				tool = *new(types.ToolCall)
// 				tool.Function.Name = fun.Name
// 			} else if toolCall {
// 				arguments = append(arguments, fun.Arguments)
// 				continue
// 			}

// 			if !utils.IsEmptyValue(tool) {
// 				res := h.DB.Where("name = ?", tool.Function.Name).First(&function)
// 				if res.Error == nil {
// 					toolCall = true
// 					callMsg := fmt.Sprintf("正在调用工具 `%s` 作答 ...\n\n", function.Label)
// 					pushMessage(c, ChatEventMessageDelta, map[string]interface{}{
// 						"type":    "text",
// 						"content": callMsg,
// 					})
// 					contents = append(contents, callMsg)
// 				}
// 				continue
// 			}

// 			if responseBody.Choices[0].FinishReason == "tool_calls" ||
// 				responseBody.Choices[0].FinishReason == "function_call" { // 函数调用完毕
// 				break
// 			}

// 			// output stopped
// 			if responseBody.Choices[0].FinishReason != "" {
// 				break // 输出完成或者输出中断了
// 			} else { // 正常输出结果
// 				// 兼容思考过程
// 				if responseBody.Choices[0].Delta.ReasoningContent != "" {
// 					reasoningContent := responseBody.Choices[0].Delta.ReasoningContent
// 					if !reasoning {
// 						reasoningContent = fmt.Sprintf("<think>%s", reasoningContent)
// 						reasoning = true
// 					}

// 					pushMessage(c, ChatEventMessageDelta, map[string]interface{}{
// 						"type":    "text",
// 						"content": reasoningContent,
// 					})
// 					contents = append(contents, reasoningContent)
// 				} else if responseBody.Choices[0].Delta.Prompt != "" {
// 					finalContent := responseBody.Choices[0].Delta.Prompt
// 					if reasoning {
// 						finalContent = fmt.Sprintf("</think>%s", responseBody.Choices[0].Delta.Prompt)
// 						reasoning = false
// 					}
// 					contents = append(contents, utils.InterfaceToString(finalContent))
// 					pushMessage(c, ChatEventMessageDelta, map[string]interface{}{
// 						"type":    "text",
// 						"content": finalContent,
// 					})
// 				}
// 			}
// 		} // end for

// 		if err := scanner.Err(); err != nil {
// 			if strings.Contains(err.Error(), "context canceled") {
// 				logger.Info("用户取消了请求：", prompt)
// 			} else {
// 				logger.Error("信息读取出错：", err)
// 			}
// 		}

// 		if toolCall { // 调用函数完成任务
// 			params := make(map[string]any)
// 			_ = utils.JsonDecode(strings.Join(arguments, ""), &params)
// 			logger.Debugf("函数名称: %s, 函数参数：%s", function.Name, params)
// 			params["user_id"] = userVo.Id
// 			var apiRes types.BizVo
// 			r, err := req2.C().R().SetHeader("Body-Type", "application/json").
// 				SetHeader("Authorization", function.Token).
// 				SetBody(params).Post(function.Action)
// 			errMsg := ""
// 			if err != nil {
// 				errMsg = err.Error()
// 			} else {
// 				all, _ := io.ReadAll(r.Body)
// 				err = json.Unmarshal(all, &apiRes)
// 				if err != nil {
// 					errMsg = err.Error()
// 				} else if apiRes.Code != types.Success {
// 					errMsg = apiRes.Message
// 				}
// 			}

// 			if errMsg != "" {
// 				errMsg = "调用函数工具出错：" + errMsg
// 				contents = append(contents, errMsg)
// 			} else {
// 				errMsg = utils.InterfaceToString(apiRes.Data)
// 				contents = append(contents, errMsg)
// 			}
// 			pushMessage(c, ChatEventMessageDelta, map[string]interface{}{
// 				"type":    "text",
// 				"content": errMsg,
// 			})
// 		}

// 		// 消息发送成功
// 		if len(contents) > 0 {
// 			usage := Usage{
// 				Prompt:           prompt,
// 				Prompt:          strings.Join(contents, ""),
// 				PromptTokens:     0,
// 				CompletionTokens: 0,
// 				TotalTokens:      0,
// 			}
// 			message.Prompt = usage.Prompt
// 			h.saveChatHistory(req, usage, message, session, role, userVo, promptCreatedAt, replyCreatedAt)
// 		}
// 	} else {
// 		var respVo OpenAIResVo
// 		body, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			return fmt.Errorf("读取响应失败：%v", body)
// 		}
// 		err = json.Unmarshal(body, &respVo)
// 		if err != nil {
// 			return fmt.Errorf("解析响应失败：%v", body)
// 		}
// 		content := respVo.Choices[0].Message.Prompt
// 		if strings.HasPrefix(req.Model, "o1-") {
// 			content = fmt.Sprintf("AI思考结束，耗时：%d 秒。\n%s", time.Now().Unix()-session.Start, respVo.Choices[0].Message.Prompt)
// 		}
// 		pushMessage(c, ChatEventMessageDelta, map[string]interface{}{
// 			"type":    "text",
// 			"content": content,
// 		})
// 		respVo.Usage.Prompt = prompt
// 		respVo.Usage.Prompt = content
// 		h.saveChatHistory(req, respVo.Usage, respVo.Choices[0].Message, session, role, userVo, promptCreatedAt, time.Now())
// 	}

// 	return nil
// }
