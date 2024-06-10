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
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/levigross/grequests"

	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
)

type ollamaResp struct {
	Id    string `json:"id"`
	Model string `json:"model"`

	CreatedAt  string `json:"created_at"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason"`
	Context    []int  `json:"context"`

	TotalDuration      int64 `json:"total_duration"`       // 生成响应所花费的总时间
	LoadDuration       int64 `json:"load_duration"`        // 以纳秒为单位加载模型所花费的时间
	PromptEvalCount    int   `json:"prompt_eval_count"`    // 提示文本中的标记（tokens）数量
	PromptEvalDuration int64 `json:"prompt_eval_duration"` // 以纳秒为单位评估提示文本所花费的时间
	EvalCount          int64 `json:"eval_count"`           // 生成响应中的标记数量
	EvalDuration       int64 `json:"eval_duration"`        // 以纳秒为单位生成响应所花费的时间
}

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
	response, err := h.sendOllamaRequest(session, prompt)
	defer response.Body.Close()

	logger.Info("HTTP请求完成，耗时：", time.Now().Sub(start))

	if err != nil {
		h.processError(err, prompt, ws)
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

func (h *ChatHandler) sendOllamaRequest(session *types.ChatSession, prompt string) (*http.Response, error) {
	apiKey, err := h.queryApiKey(session)
	if err != nil {
		return nil, err
	}

	// todo add context to request body
	postData := map[string]interface{}{
		"model":  session.Model.Value,
		"stream": true,
		"prompt": prompt,
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey.Value,
	}

	ro := &grequests.RequestOptions{
		JSON:    postData,
		Headers: headers,
	}
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
	// 循环读取 Chunk 消息
	var message = types.Message{}
	scanner := bufio.NewScanner(response.Body)

	var content string
	var replyTokens int

	for scanner.Scan() {
		var resp ollamaResp
		line := scanner.Text()

		err := utils.JsonDecode(line, &resp)
		if err != nil {
			logger.Error("error with parse data line: ", content)
			utils.ReplyMessage(ws, fmt.Sprintf("**解析数据行失败：%s**", err))
			break
		}

		if resp.Done == true && resp.DoneReason == "stop" {
			utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsEnd})
			message.Content = utils.InterfaceToString(resp.Context)
			replyTokens = resp.PromptEvalCount

			// 消息发送成功后做记录工作
			h.recordInfoAfterSendMessage(chatCtx, req, userVo, message, prompt, session, role, promptCreatedAt, replyTokens, replyCreatedAt)

			break
		} else if resp.Done == true && resp.DoneReason != "stop" {
			utils.ReplyMessage(ws, fmt.Sprintf("**API 返回错误：%s**", resp.DoneReason))
			break
		}

		if len(resp.Id) > 0 {
			utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
		}

		if len(resp.Response) > 0 {
			utils.ReplyChunkMessage(ws, types.WsMessage{
				Type:    types.WsMiddle,
				Content: utils.InterfaceToString(resp.Response),
			})
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

func (h *ChatHandler) recordInfoAfterSendMessage(chatCtx []types.Message, req types.ApiRequest, userVo vo.User, message types.Message, prompt string, session *types.ChatSession, role model.ChatRole, promptCreatedAt time.Time, replyTokens int, replyCreatedAt time.Time) {
	if message.Role == "" {
		message.Role = "assistant"
	}

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
	//totalTokens := replyTokens + getTotalTokens(req)
	// todo rebuild the tokens
	historyReplyMsg := model.ChatMessage{
		UserId:     userVo.Id,
		ChatId:     session.ChatId,
		RoleId:     role.Id,
		Type:       types.ReplyMsg,
		Icon:       role.Icon,
		Content:    message.Content,
		Tokens:     replyTokens,
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
