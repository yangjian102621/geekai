package chatimpl

import (
	"bytes"
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	logger2 "chatplus/logger"
	"chatplus/service/oss"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const ErrorMsg = "抱歉，AI 助手开小差了，请稍后再试。"

var ErrImg = "![](/images/wx.png)"

var logger = logger2.GetLogger()

type ChatHandler struct {
	handler.BaseHandler
	redis         *redis.Client
	uploadManager *oss.UploaderManager
}

func NewChatHandler(app *core.AppServer, db *gorm.DB, redis *redis.Client, manager *oss.UploaderManager) *ChatHandler {
	return &ChatHandler{
		BaseHandler:   handler.BaseHandler{App: app, DB: db},
		redis:         redis,
		uploadManager: manager,
	}
}

func (h *ChatHandler) Init() {
	// 如果后台有上传微信客服微信二维码，则覆盖
	if h.App.SysConfig.WechatCardURL != "" {
		ErrImg = fmt.Sprintf("![](%s)", h.App.SysConfig.WechatCardURL)
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
	// get model info
	var chatModel model.ChatModel
	res := h.DB.First(&chatModel, modelId)
	if res.Error != nil || chatModel.Enabled == false {
		utils.ReplyMessage(client, "当前AI模型暂未启用，连接已关闭！！！")
		c.Abort()
		return
	}

	session := h.App.ChatSession.Get(sessionId)
	if session == nil {
		user, err := h.GetLoginUser(c)
		if err != nil {
			logger.Info("用户未登录")
			c.Abort()
			return
		}
		session = &types.ChatSession{
			SessionId: sessionId,
			ClientIP:  c.ClientIP(),
			Username:  user.Username,
			UserId:    user.Id,
		}
		h.App.ChatSession.Put(sessionId, session)
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
		Platform:    types.Platform(chatModel.Platform)}
	logger.Infof("New websocket connected, IP: %s, Username: %s", c.ClientIP(), session.Username)
	var chatRole model.ChatRole
	res = h.DB.First(&chatRole, roleId)
	if res.Error != nil || !chatRole.Enable {
		utils.ReplyMessage(client, "当前聊天角色不存在或者未启用，连接已关闭！！！")
		c.Abort()
		return
	}

	h.Init()

	// 保存会话连接
	h.App.ChatClients.Put(sessionId, client)
	go func() {
		for {
			_, msg, err := client.Receive()
			if err != nil {
				client.Close()
				h.App.ChatClients.Delete(sessionId)
				cancelFunc := h.App.ReqCancelFunc.Get(sessionId)
				if cancelFunc != nil {
					cancelFunc()
					h.App.ReqCancelFunc.Delete(sessionId)
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
			h.App.ReqCancelFunc.Put(sessionId, cancel)
			// 回复消息
			err = h.sendMessage(ctx, session, chatRole, utils.InterfaceToString(message.Content), client)
			if err != nil {
				logger.Error(err)
				utils.ReplyChunkMessage(client, types.WsMessage{Type: types.WsEnd})
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
		utils.ReplyMessage(ws, "未授权用户，您正在进行非法操作！")
		return res.Error
	}
	var userVo vo.User
	err := utils.CopyObject(user, &userVo)
	userVo.Id = user.Id
	if err != nil {
		return errors.New("User 对象转换失败，" + err.Error())
	}

	if userVo.Status == false {
		utils.ReplyMessage(ws, "您的账号已经被禁用，如果疑问，请联系管理员！")
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}

	if userVo.Power < session.Model.Power {
		utils.ReplyMessage(ws, fmt.Sprintf("您当前剩余算力（%d）已不足以支付当前模型的单次对话需要消耗的算力（%d）！", userVo.Power, session.Model.Power))
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}

	if userVo.ExpiredTime > 0 && userVo.ExpiredTime <= time.Now().Unix() {
		utils.ReplyMessage(ws, "您的账号已经过期，请联系管理员！")
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}

	// 检查 prompt 长度是否超过了当前模型允许的最大上下文长度
	promptTokens, err := utils.CalcTokens(prompt, session.Model.Value)
	if promptTokens > session.Model.MaxContext {
		utils.ReplyMessage(ws, "对话内容超出了当前模型允许的最大上下文长度！")
		return nil
	}

	var req = types.ApiRequest{
		Model:  session.Model.Value,
		Stream: true,
	}
	switch session.Model.Platform {
	case types.Azure, types.ChatGLM, types.Baidu, types.XunFei:
		req.Temperature = session.Model.Temperature
		req.MaxTokens = session.Model.MaxTokens
		break
	case types.OpenAI:
		req.Temperature = session.Model.Temperature
		req.MaxTokens = session.Model.MaxTokens
		// OpenAI 支持函数功能
		var items []model.Function
		res := h.DB.Where("enabled", true).Find(&items)
		if res.Error != nil {
			break
		}

		var tools = make([]interface{}, 0)
		for _, v := range items {
			var parameters map[string]interface{}
			err = utils.JsonDecode(v.Parameters, &parameters)
			if err != nil {
				continue
			}
			required := parameters["required"]
			delete(parameters, "required")
			tools = append(tools, gin.H{
				"type": "function",
				"function": gin.H{
					"name":        v.Name,
					"description": v.Description,
					"parameters":  parameters,
					"required":    required,
				},
			})
		}

		if len(tools) > 0 {
			req.Tools = tools
			req.ToolChoice = "auto"
		}
	case types.QWen:
		req.Parameters = map[string]interface{}{
			"max_tokens":  session.Model.MaxTokens,
			"temperature": session.Model.Temperature,
		}
		break

	default:
		utils.ReplyMessage(ws, "不支持的平台："+session.Model.Platform+"，请联系管理员！")
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}

	// 加载聊天上下文
	chatCtx := make([]types.Message, 0)
	messages := make([]types.Message, 0)
	if h.App.SysConfig.EnableContext {
		if h.App.ChatContexts.Has(session.ChatId) {
			messages = h.App.ChatContexts.Get(session.ChatId)
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

	if session.Model.Platform == types.QWen {
		req.Input = make(map[string]interface{})
		reqMgs = append(reqMgs, types.Message{
			Role:    "user",
			Content: prompt,
		})
		req.Input["messages"] = reqMgs
	} else {
		req.Messages = append(reqMgs, map[string]interface{}{
			"role":    "user",
			"content": prompt,
		})
	}

	switch session.Model.Platform {
	case types.Azure:
		return h.sendAzureMessage(chatCtx, req, userVo, ctx, session, role, prompt, ws)
	case types.OpenAI:
		return h.sendOpenAiMessage(chatCtx, req, userVo, ctx, session, role, prompt, ws)
	case types.ChatGLM:
		return h.sendChatGLMMessage(chatCtx, req, userVo, ctx, session, role, prompt, ws)
	case types.Baidu:
		return h.sendBaiduMessage(chatCtx, req, userVo, ctx, session, role, prompt, ws)
	case types.XunFei:
		return h.sendXunFeiMessage(chatCtx, req, userVo, ctx, session, role, prompt, ws)
	case types.QWen:
		return h.sendQWenMessage(chatCtx, req, userVo, ctx, session, role, prompt, ws)
	}
	utils.ReplyChunkMessage(ws, types.WsMessage{
		Type:    types.WsMiddle,
		Content: fmt.Sprintf("Not supported platform: %s", session.Model.Platform),
	})
	return nil
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
	if h.App.ReqCancelFunc.Has(sessionId) {
		h.App.ReqCancelFunc.Get(sessionId)()
		h.App.ReqCancelFunc.Delete(sessionId)
	}
	resp.SUCCESS(c, types.OkMsg)
}

// 发送请求到 OpenAI 服务器
// useOwnApiKey: 是否使用了用户自己的 API KEY
func (h *ChatHandler) doRequest(ctx context.Context, req types.ApiRequest, platform types.Platform, apiKey *model.ApiKey) (*http.Response, error) {
	res := h.DB.Where("platform = ?", platform).Where("type = ?", "chat").Where("enabled = ?", true).Order("last_used_at ASC").First(apiKey)
	if res.Error != nil {
		return nil, errors.New("no available key, please import key")
	}
	var apiURL string
	switch platform {
	case types.Azure:
		md := strings.Replace(req.Model, ".", "", 1)
		apiURL = strings.Replace(apiKey.ApiURL, "{model}", md, 1)
		break
	case types.ChatGLM:
		apiURL = strings.Replace(apiKey.ApiURL, "{model}", req.Model, 1)
		req.Prompt = req.Messages // 使用 prompt 字段替代 message 字段
		req.Messages = nil
		break
	case types.Baidu:
		apiURL = strings.Replace(apiKey.ApiURL, "{model}", req.Model, 1)
		break
	case types.QWen:
		apiURL = apiKey.ApiURL
		req.Messages = nil
		break
	default:
		apiURL = apiKey.ApiURL
	}
	// 更新 API KEY 的最后使用时间
	h.DB.Model(apiKey).UpdateColumn("last_used_at", time.Now().Unix())
	// 百度文心，需要串接 access_token
	if platform == types.Baidu {
		token, err := h.getBaiduToken(apiKey.Value)
		if err != nil {
			return nil, err
		}
		logger.Info("百度文心 Access_Token：", token)
		apiURL = fmt.Sprintf("%s?access_token=%s", apiURL, token)
	}

	logger.Debugf(utils.JsonEncode(req))

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
	var proxyURL string
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
	logger.Debugf("Sending %s request, ApiURL:%s, API KEY:%s, PROXY: %s, Model: %s", platform, apiURL, apiKey.Value, proxyURL, req.Model)
	switch platform {
	case types.Azure:
		request.Header.Set("api-key", apiKey.Value)
		break
	case types.ChatGLM:
		token, err := h.getChatGLMToken(apiKey.Value)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		break
	case types.Baidu:
		request.RequestURI = ""
	case types.OpenAI:
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey.Value))
		break
	case types.QWen:
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey.Value))
		request.Header.Set("X-DashScope-SSE", "enable")
		break
	}
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

		newImgURL, err := h.uploadManager.GetUploadHandler().PutImg(imageURL, false)
		if err != nil {
			logger.Error("error with download image: ", err)
			continue
		}

		text = strings.ReplaceAll(text, imageURL, newImgURL)
	}
	return text
}
