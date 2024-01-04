package chatimpl

import (
	"bytes"
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	logger2 "chatplus/logger"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ErrorMsg = "抱歉，AI 助手开小差了，请稍后再试。"
const ErrImg = "![](/images/wx.png)"

var logger = logger2.GetLogger()

type ChatHandler struct {
	handler.BaseHandler
	db    *gorm.DB
	redis *redis.Client
}

func NewChatHandler(app *core.AppServer, db *gorm.DB, redis *redis.Client) *ChatHandler {
	h := ChatHandler{
		db:    db,
		redis: redis,
	}
	h.App = app
	return &h
}

var chatConfig types.ChatConfig

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
	res := h.db.First(&chatModel, modelId)
	if res.Error != nil || chatModel.Enabled == false {
		utils.ReplyMessage(client, "当前AI模型暂未启用，连接已关闭！！！")
		c.Abort()
		return
	}

	session := h.App.ChatSession.Get(sessionId)
	if session == nil {
		user, err := utils.GetLoginUser(c, h.db)
		if err != nil {
			logger.Info("用户未登录")
			c.Abort()
			return
		}
		session = &types.ChatSession{
			SessionId: sessionId,
			ClientIP:  c.ClientIP(),
			Username:  user.Mobile,
			UserId:    user.Id,
		}
		h.App.ChatSession.Put(sessionId, session)
	}

	// use old chat data override the chat model and role ID
	var chat model.ChatItem
	res = h.db.Where("chat_id=?", chatId).First(&chat)
	if res.Error == nil {
		chatModel.Id = chat.ModelId
		roleId = int(chat.RoleId)
	}

	session.ChatId = chatId
	session.Model = types.ChatModel{
		Id:       chatModel.Id,
		Value:    chatModel.Value,
		Weight:   chatModel.Weight,
		Platform: types.Platform(chatModel.Platform)}
	logger.Infof("New websocket connected, IP: %s, Username: %s", c.ClientIP(), session.Username)
	var chatRole model.ChatRole
	res = h.db.First(&chatRole, roleId)
	if res.Error != nil || !chatRole.Enable {
		utils.ReplyMessage(client, "当前聊天角色不存在或者未启用，连接已关闭！！！")
		c.Abort()
		return
	}

	// 初始化聊天配置
	var config model.Config
	h.db.Where("marker", "chat").First(&config)
	err = utils.JsonDecode(config.Config, &chatConfig)
	if err != nil {
		utils.ReplyMessage(client, "加载系统配置失败，连接已关闭！！！")
		c.Abort()
		return
	}

	// 保存会话连接
	h.App.ChatClients.Put(sessionId, client)
	go func() {
		for {
			_, msg, err := client.Receive()
			if err != nil {
				logger.Error(err)
				client.Close()
				h.App.ChatClients.Delete(sessionId)
				cancelFunc := h.App.ReqCancelFunc.Get(sessionId)
				if cancelFunc != nil {
					cancelFunc()
					h.App.ReqCancelFunc.Delete(sessionId)
				}
				return
			}

			message := string(msg)
			logger.Info("Receive a message: ", message)
			//utils.ReplyMessage(client, "这是一条测试消息！")
			ctx, cancel := context.WithCancel(context.Background())
			h.App.ReqCancelFunc.Put(sessionId, cancel)
			// 回复消息
			err = h.sendMessage(ctx, session, chatRole, message, client)
			if err != nil {
				logger.Error(err)
				utils.ReplyChunkMessage(client, types.WsMessage{Type: types.WsEnd})
			} else {
				utils.ReplyChunkMessage(client, types.WsMessage{Type: types.WsEnd})
				logger.Info("回答完毕: " + string(message))
			}

		}
	}()
}

func (h *ChatHandler) sendMessage(ctx context.Context, session *types.ChatSession, role model.ChatRole, prompt string, ws *types.WsClient) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recover message from error: ", r)
		}
	}()

	var user model.User
	res := h.db.Model(&model.User{}).First(&user, session.UserId)
	if res.Error != nil {
		utils.ReplyMessage(ws, "非法用户，请联系管理员！")
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

	if userVo.Calls < session.Model.Weight {
		utils.ReplyMessage(ws, fmt.Sprintf("您当前剩余对话次数（%d）已不足以支付当前模型的单次对话需要消耗的对话额度（%d）！", userVo.Calls, session.Model.Weight))
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}

	if userVo.Calls <= 0 && userVo.ChatConfig.ApiKeys[session.Model.Platform] == "" {
		utils.ReplyMessage(ws, "您的对话次数已经用尽，请联系管理员或者充值点卡继续对话！")
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}

	if userVo.ExpiredTime > 0 && userVo.ExpiredTime <= time.Now().Unix() {
		utils.ReplyMessage(ws, "您的账号已经过期，请联系管理员！")
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}
	var req = types.ApiRequest{
		Model:  session.Model.Value,
		Stream: true,
	}
	switch session.Model.Platform {
	case types.Azure:
		req.Temperature = h.App.ChatConfig.Azure.Temperature
		req.MaxTokens = h.App.ChatConfig.Azure.MaxTokens
		break
	case types.ChatGLM:
		req.Temperature = h.App.ChatConfig.ChatGML.Temperature
		req.MaxTokens = h.App.ChatConfig.ChatGML.MaxTokens
		break
	case types.Baidu:
		req.Temperature = h.App.ChatConfig.OpenAI.Temperature
		// TODO： 目前只支持 ERNIE-Bot-turbo 模型，如果是 ERNIE-Bot 模型则需要增加函数支持
		break
	case types.OpenAI:
		req.Temperature = h.App.ChatConfig.OpenAI.Temperature
		req.MaxTokens = h.App.ChatConfig.OpenAI.MaxTokens
		// OpenAI 支持函数功能
		var items []model.Function
		res := h.db.Where("enabled", true).Find(&items)
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
	case types.XunFei:
		req.Temperature = h.App.ChatConfig.XunFei.Temperature
		req.MaxTokens = h.App.ChatConfig.XunFei.MaxTokens
		break
	default:
		utils.ReplyMessage(ws, "不支持的平台："+session.Model.Platform+"，请联系管理员！")
		utils.ReplyMessage(ws, ErrImg)
		return nil
	}

	// 加载聊天上下文
	var chatCtx []interface{}
	if h.App.ChatConfig.EnableContext {
		if h.App.ChatContexts.Has(session.ChatId) {
			chatCtx = h.App.ChatContexts.Get(session.ChatId)
		} else {
			// calculate the tokens of current request, to prevent to exceeding the max tokens num
			tokens := req.MaxTokens
			tks, _ := utils.CalcTokens(utils.JsonEncode(req.Tools), req.Model)
			tokens += tks
			// loading the role context
			var messages []types.Message
			err := utils.JsonDecode(role.Context, &messages)
			if err == nil {
				for _, v := range messages {
					tks, _ := utils.CalcTokens(v.Content, req.Model)
					if tokens+tks >= types.GetModelMaxToken(req.Model) {
						break
					}
					tokens += tks
					chatCtx = append(chatCtx, v)
				}
			}

			// loading recent chat history as chat context
			if chatConfig.ContextDeep > 0 {
				var historyMessages []model.HistoryMessage
				res := h.db.Debug().Where("chat_id = ? and use_context = 1", session.ChatId).Limit(chatConfig.ContextDeep).Order("id desc").Find(&historyMessages)
				if res.Error == nil {
					for i := len(historyMessages) - 1; i >= 0; i-- {
						msg := historyMessages[i]
						if tokens+msg.Tokens >= types.GetModelMaxToken(session.Model.Value) {
							break
						}
						tokens += msg.Tokens
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
		var item model.HistoryMessage
		userId, _ := c.Get(types.LoginUserID)
		res := h.db.Where("user_id = ?", userId).Where("chat_id = ?", data.ChatId).Last(&item)
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
	res := h.db.Where("platform = ?", platform).Where("type = ?", "chat").Where("enabled = ?", true).Order("last_used_at ASC").First(apiKey)
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
	default:
		apiURL = apiKey.ApiURL
	}
	// 更新 API KEY 的最后使用时间
	h.db.Model(apiKey).UpdateColumn("last_used_at", time.Now().Unix())
	// 百度文心，需要串接 access_token
	if platform == types.Baidu {
		token, err := h.getBaiduToken(apiKey.Value)
		if err != nil {
			return nil, err
		}
		logger.Info("百度文心 Access_Token：", token)
		apiURL = fmt.Sprintf("%s?access_token=%s", apiURL, token)
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
	proxyURL := h.App.Config.ProxyURL
	if proxyURL != "" && platform == types.OpenAI { // 使用代理
		proxy, _ := url.Parse(proxyURL)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	} else {
		client = http.DefaultClient
	}
	logger.Debugf("Sending %s request, ApiURL:%s, ApiKey:%s, PROXY: %s, Model: %s", platform, apiURL, apiKey.Value, proxyURL, req.Model)
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
	}
	return client.Do(request)
}

// 扣减用户的对话次数
func (h *ChatHandler) subUserCalls(userVo vo.User, session *types.ChatSession) {
	// 仅当用户没有导入自己的 API KEY 时才进行扣减
	if userVo.ChatConfig.ApiKeys[session.Model.Platform] == "" {
		num := 1
		if session.Model.Weight > 0 {
			num = session.Model.Weight
		}
		h.db.Model(&model.User{}).Where("id = ?", userVo.Id).UpdateColumn("calls", gorm.Expr("calls - ?", num))
	}
}

func (h *ChatHandler) incUserTokenFee(userId uint, tokens int) {
	h.db.Model(&model.User{}).Where("id = ?", userId).
		UpdateColumn("total_tokens", gorm.Expr("total_tokens + ?", tokens))
	h.db.Model(&model.User{}).Where("id = ?", userId).
		UpdateColumn("tokens", gorm.Expr("tokens + ?", tokens))
}
