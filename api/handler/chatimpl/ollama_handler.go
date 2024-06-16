package chatimpl

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/levigross/grequests"

	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
)

// ChatResponse is the response returned by [Client.Chat]. Its fields are
// similar to [GenerateResponse].
type ChatResponse struct {
	Model      string        `json:"model"`
	CreatedAt  time.Time     `json:"created_at"`
	Message    types.Message `json:"message"`
	DoneReason string        `json:"done_reason,omitempty"`

	Done bool `json:"done"`

	Metrics
}

type Metrics struct {
	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
}

type Message struct {
	types.Message
	Images []ImageData `json:"images,omitempty"`
}

type ImageData []byte

// 通义千问消息发送实现
func (h *ChatHandler) sendOllamaMessage(
	chatCtx []types.Message,
	req types.ApiRequest,
	userVo vo.User,
	ctx context.Context,
	session *types.ChatSession,
	role model.ChatRole,
	prompt string,
	ws *types.WsClient) error {

	promptCreatedAt := time.Now() // 记录提问时间
	start := time.Now()

	//var apiKey = model.ApiKey{}
	//response, err := h.doRequest(ctx, req, session, &apiKey)
	response, err := h.sendOllamaRequest(chatCtx, session, prompt)

	logger.Info("HTTP请求完成，耗时：", time.Now().Sub(start))

	if err != nil {
		h.processError(err, prompt, ws)
		return nil
	} else {
		defer response.Body.Close()
	}

	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/x-ndjson") {

		h.processOllamaStreamResponse(chatCtx, req, userVo, response, ws, prompt, session, role, promptCreatedAt)

	} else {
		if err = h.processOllamaJsonResponse(response, ws); err != nil {
			return err
		}
	}

	return nil
}

func (h *ChatHandler) sendOllamaRequest(chatCtx []types.Message, session *types.ChatSession, prompt string) (*http.Response, error) {
	apiKey, err := h.queryApiKey(session)
	if err != nil {
		return nil, err
	}

	chatCtx = append(chatCtx, types.Message{
		Role:    "user",
		Content: prompt,
	})
	messages := make([]Message, 0)
	for _, ctx := range chatCtx {
		if ctx.Role == "" {
			continue
		}

		m := Message{
			Message: ctx,
		}

		url := h.parseURL(ctx.Content)
		if url != "" {
			encode, err := h.downImgAndBase64Encode(url)
			if err != nil {
				logger.Infof("img url convert to binary err：%s, will not send image to ollama", err)
				continue
			}
			m.Content = strings.Replace(ctx.Content, url, "", 1)
			m.Images = []ImageData{encode}
		}

		messages = append(messages, m)
	}

	postData := map[string]interface{}{
		"model":    session.Model.Value,
		"stream":   true,
		"messages": messages,
		"options": map[string]interface{}{
			"temperature": session.Model.Temperature,
		},
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	// 兼容ollama原生11343端口，与ollama webui api-key的方式
	if strings.HasPrefix(apiKey.Value, "sk-") {
		headers["Authorization"] = "Bearer " + apiKey.Value
	}

	ro := &grequests.RequestOptions{
		JSON:    postData,
		Headers: headers,
	}
	requestBody, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	logger.Debugf("ollama request body: %s", string(requestBody))

	resp, err := grequests.Post(apiKey.ApiURL, ro)
	if err != nil {
		return nil, err
	}

	if !resp.Ok {
		return nil, resp.Error
	}

	return resp.RawResponse, nil
}

func (h *ChatHandler) queryApiKey(session *types.ChatSession) (*model.ApiKey, error) {
	apiKey := &model.ApiKey{}

	// if the chat model bind a KEY, use it directly
	if session.Model.KeyId > 0 {
		h.DB.Debug().Where("id", session.Model.KeyId).Where("enabled", true).Find(apiKey)
	}
	// use the last unused key
	if apiKey.Id == 0 {
		h.DB.Debug().Where("platform", session.Model.Platform).Where("type", "chat").Where("enabled", true).Order("last_used_at ASC").First(apiKey)
	}
	if apiKey.Id == 0 {
		return nil, errors.New("no available key, please import key")
	}

	h.DB.Model(apiKey).UpdateColumn("last_used_at", time.Now().Unix())

	return apiKey, nil
}

func (h *ChatHandler) processOllamaStreamResponse(
	chatCtx []types.Message, req types.ApiRequest, userVo vo.User,
	response *http.Response, ws *types.WsClient, prompt string,
	session *types.ChatSession, role model.ChatRole, promptCreatedAt time.Time) {

	// 记录回复时间
	replyCreatedAt := time.Now()
	scanner := bufio.NewScanner(response.Body)

	var contents = make([]string, 0)
	var content string
	var outPutStart = true

	// 循环读取 返回 消息
	for scanner.Scan() {
		var resp ChatResponse
		line := scanner.Text()

		err := utils.JsonDecode(line, &resp)
		if err != nil {
			logger.Error("error with parse data line: ", line)
			utils.ReplyMessage(ws, fmt.Sprintf("**解析数据行失败：%s**", err))
			break
		}

		if resp.Done == true && resp.DoneReason == "stop" {
			utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsEnd})

			// 消息发送成功后做记录工作
			h.recordInfoAfterSendMessage(chatCtx, req, userVo, prompt, session, role, promptCreatedAt, replyCreatedAt, strings.Join(contents, ""))

			break
		} else if resp.Done == true && resp.DoneReason != "stop" {
			utils.ReplyMessage(ws, fmt.Sprintf("**API 返回错误：%s**", resp.DoneReason))
			break
		}

		if len(contents) == 0 && outPutStart {
			logger.Infof("开始输出消息：%s", resp.Message.Content)
			utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
			outPutStart = false
		}

		if len(resp.Message.Content) > 0 {
			utils.ReplyChunkMessage(ws, types.WsMessage{
				Type:    types.WsMiddle,
				Content: utils.InterfaceToString(resp.Message.Content),
			})

			content += resp.Message.Content
			contents = append(contents, resp.Message.Content)
		}

	}

	if err := scanner.Err(); err != nil {
		if strings.Contains(err.Error(), "context canceled") {
			logger.Info("用户取消了请求：", prompt)
		} else {
			logger.Error("信息读取出错：", err)
		}
	}
}

func (h *ChatHandler) processOllamaJsonResponse(response *http.Response, ws *types.WsClient) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error with reading response: %v", err)
	}

	var res struct {
		Code int    `json:"error_code"`
		Msg  string `json:"error_msg"`
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("error with decode response: %v", err)
	}
	utils.ReplyMessage(ws, "请求Ollama大模型 API 失败："+res.Msg)
	return nil
}

func (h *ChatHandler) recordInfoAfterSendMessage(chatCtx []types.Message, req types.ApiRequest, userVo vo.User,
	prompt string, session *types.ChatSession, role model.ChatRole,
	promptCreatedAt time.Time, replyCreatedAt time.Time, content string) {

	message := types.Message{Role: "assistant", Content: content}
	useMsg := types.Message{Role: "user", Content: prompt}

	// 更新上下文消息，如果是调用函数则不需要更新上下文
	if h.App.SysConfig.EnableContext {
		chatCtx = append(chatCtx, useMsg)  // 提问消息
		chatCtx = append(chatCtx, message) // 回复消息
		h.App.ChatContexts.Put(session.ChatId, chatCtx)
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
		Content:    content,
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
	h.subUserPower(userVo, session, promptToken, replyTokens)

	// 保存当前会话
	var chatItem model.ChatItem
	res = h.DB.Where("chat_id = ?", session.ChatId).First(&chatItem)
	if res.Error != nil {
		chatItem.ChatId = session.ChatId
		chatItem.UserId = session.UserId
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

func (h *ChatHandler) processError(err error, prompt string, ws *types.WsClient) {
	if strings.Contains(err.Error(), "context canceled") {
		logger.Info("用户取消了请求：", prompt)
		return
	} else if strings.Contains(err.Error(), "no available key") {
		utils.ReplyMessage(ws, "抱歉😔😔😔，系统已经没有可用的 API KEY，请联系管理员！")
		return
	} else {
		logger.Error(err)
	}

	utils.ReplyMessage(ws, ErrorMsg)
	utils.ReplyMessage(ws, ErrImg)
	return
}

func (h *ChatHandler) downImgAndBase64Encode(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("download img failed")
	}

	return ioutil.ReadAll(resp.Body)
}

func (h *ChatHandler) parseURL(input string) string {
	// 正则表达式模式匹配包含 HTTP URL 的字符串
	regexStr := `(?i)\b((https?://|www\.)[-A-Za-z0-9+&@#/%?=~_|!:,.;]*[-A-Za-z0-9+&@#/%=~_|]\.(jpg|jpeg|png|gif|bmp|webp))`

	// 创建正则表达式对象，并验证输入字符串是否以 URL 开始
	re := regexp.MustCompile(regexStr)

	matches := re.FindStringSubmatch(input)
	if len(matches) > 0 {
		return matches[0] // 返回第一个匹配的URL
	}

	return ""
}
