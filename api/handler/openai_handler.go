package handler

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

type Usage struct {
	Prompt           string `json:"prompt,omitempty"`
	Content          string `json:"content,omitempty"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
}

type OpenAIResVo struct {
	Id                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

// OPenAI 消息发送实现
func (h *ChatHandler) sendOpenAiMessage(
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

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("请求 OpenAI API 失败：%d, %v", response.StatusCode, string(body))
	}

	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		replyCreatedAt := time.Now() // 记录回复时间
		// 循环读取 Chunk 消息
		var message = types.Message{Role: "assistant"}
		var contents = make([]string, 0)
		var function model.Function
		var toolCall = false
		var arguments = make([]string, 0)
		scanner := bufio.NewScanner(response.Body)
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
				utils.SendChunkMsg(ws, "抱歉😔😔😔，AI助手由于未知原因已经停止输出内容。")
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
					utils.SendChunkMsg(ws, callMsg)
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
				utils.SendChunkMsg(ws, responseBody.Choices[0].Delta.Content)
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
			r, err := req2.C().R().SetHeader("Body-Type", "application/json").
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
				errMsg = "调用函数工具出错：" + apiRes.Message + errMsg
				contents = append(contents, errMsg)
			} else {
				errMsg = utils.InterfaceToString(apiRes.Data)
				contents = append(contents, errMsg)
			}
			utils.SendChunkMsg(ws, errMsg)
		}

		// 消息发送成功
		if len(contents) > 0 {
			usage := Usage{
				Prompt:           prompt,
				Content:          strings.Join(contents, ""),
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			}
			message.Content = usage.Content
			h.saveChatHistory(req, usage, message, session, role, userVo, promptCreatedAt, replyCreatedAt)
		}
	} else { // 非流式输出
		var respVo OpenAIResVo
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("读取响应失败：%v", body)
		}
		err = json.Unmarshal(body, &respVo)
		if err != nil {
			return fmt.Errorf("解析响应失败：%v", body)
		}
		content := respVo.Choices[0].Message.Content
		if strings.HasPrefix(req.Model, "o1-") {
			content = fmt.Sprintf("AI思考结束，耗时：%d 秒。\n%s", time.Now().Unix()-session.Start, respVo.Choices[0].Message.Content)
		}
		utils.SendChunkMsg(ws, content)
		respVo.Usage.Prompt = prompt
		respVo.Usage.Content = content
		h.saveChatHistory(req, respVo.Usage, respVo.Choices[0].Message, session, role, userVo, promptCreatedAt, time.Now())
	}

	return nil
}
