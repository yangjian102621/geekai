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
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	req2 "github.com/imroc/req/v3"
	"io"
	"strings"
	"time"
)

// OPenAI 消息发送实现
func (h *ChatHandler) sendOpenAiMessage(
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
	var apiKey = model.ApiKey{}
	response, err := h.doRequest(ctx, req, session, &apiKey)
	logger.Info("HTTP请求完成，耗时：", time.Now().Sub(start))
	if err != nil {
		if strings.Contains(err.Error(), "context canceled") {
			return fmt.Errorf("用户取消了请求：%s", prompt)
		} else if strings.Contains(err.Error(), "no available key") {
			return errors.New("抱歉😔😔😔，系统已经没有可用的 API KEY，请联系管理员！")
		}
		return err
	} else {
		defer response.Body.Close()
	}

	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		replyCreatedAt := time.Now() // 记录回复时间
		// 循环读取 Chunk 消息
		var message = types.Message{}
		var contents = make([]string, 0)
		var function model.Function
		var toolCall = false
		var arguments = make([]string, 0)
		scanner := bufio.NewScanner(response.Body)
		var isNew = true
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(line, "data:") || len(line) < 30 {
				continue
			}
			var responseBody = types.ApiResponse{}
			err = json.Unmarshal([]byte(line[6:]), &responseBody)
			if err != nil { // 数据解析出错
				return errors.New(line)
			}
			if len(responseBody.Choices) == 0 { // Fixed: 兼容 Azure API 第一个输出空行
				continue
			}
			if responseBody.Choices[0].Delta.Content == nil && responseBody.Choices[0].Delta.ToolCalls == nil {
				continue
			}

			if responseBody.Choices[0].FinishReason == "stop" && len(contents) == 0 {
				utils.ReplyMessage(ws, "抱歉😔😔😔，AI助手由于未知原因已经停止输出内容。")
				break
			}

			var tool types.ToolCall
			if len(responseBody.Choices[0].Delta.ToolCalls) > 0 {
				tool = responseBody.Choices[0].Delta.ToolCalls[0]
				if toolCall && tool.Function.Name == "" {
					arguments = append(arguments, tool.Function.Arguments)
					continue
				}
			}

			// 兼容 Function Call
			fun := responseBody.Choices[0].Delta.FunctionCall
			if fun.Name != "" {
				tool = *new(types.ToolCall)
				tool.Function.Name = fun.Name
			} else if toolCall {
				arguments = append(arguments, fun.Arguments)
				continue
			}

			if !utils.IsEmptyValue(tool) {
				res := h.DB.Where("name = ?", tool.Function.Name).First(&function)
				if res.Error == nil {
					toolCall = true
					callMsg := fmt.Sprintf("正在调用工具 `%s` 作答 ...\n\n", function.Label)
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: callMsg})
					contents = append(contents, callMsg)
				}
				continue
			}

			if responseBody.Choices[0].FinishReason == "tool_calls" ||
				responseBody.Choices[0].FinishReason == "function_call" { // 函数调用完毕
				break
			}

			// output stopped
			if responseBody.Choices[0].FinishReason != "" {
				break // 输出完成或者输出中断了
			} else {
				content := responseBody.Choices[0].Delta.Content
				contents = append(contents, utils.InterfaceToString(content))
				if isNew {
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
					isNew = false
				}
				utils.ReplyChunkMessage(ws, types.WsMessage{
					Type:    types.WsMiddle,
					Content: utils.InterfaceToString(responseBody.Choices[0].Delta.Content),
				})
			}
		} // end for

		if err := scanner.Err(); err != nil {
			if strings.Contains(err.Error(), "context canceled") {
				logger.Info("用户取消了请求：", prompt)
			} else {
				logger.Error("信息读取出错：", err)
			}
		}

		if toolCall { // 调用函数完成任务
			params := make(map[string]interface{})
			_ = utils.JsonDecode(strings.Join(arguments, ""), &params)
			logger.Debugf("函数名称: %s, 函数参数：%s", function.Name, params)
			params["user_id"] = userVo.Id
			var apiRes types.BizVo
			r, err := req2.C().R().SetHeader("Content-Type", "application/json").
				SetHeader("Authorization", function.Token).
				SetBody(params).
				SetSuccessResult(&apiRes).Post(function.Action)
			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			} else if r.IsErrorState() {
				errMsg = r.Status
			}
			if errMsg != "" || apiRes.Code != types.Success {
				msg := "调用函数工具出错：" + apiRes.Message + errMsg
				utils.ReplyChunkMessage(ws, types.WsMessage{
					Type:    types.WsMiddle,
					Content: msg,
				})
				contents = append(contents, msg)
			} else {
				utils.ReplyChunkMessage(ws, types.WsMessage{
					Type:    types.WsMiddle,
					Content: apiRes.Data,
				})
				contents = append(contents, utils.InterfaceToString(apiRes.Data))
			}
		}

		// 消息发送成功
		if len(contents) > 0 {
			h.saveChatHistory(req, prompt, contents, message, chatCtx, session, role, userVo, promptCreatedAt, replyCreatedAt)
		}
	} else {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("请求 OpenAI API 失败：%s", body)
	}

	return nil
}
