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
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io"
	"strings"
	"time"
)

type qWenResp struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
	} `json:"output,omitempty"`
	Usage struct {
		TotalTokens  int `json:"total_tokens"`
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage,omitempty"`
	RequestID string `json:"request_id"`

	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// 通义千问消息发送实现
func (h *ChatHandler) sendQWenMessage(
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
		scanner := bufio.NewScanner(response.Body)

		var content, lastText, newText string
		var outPutStart = false

		for scanner.Scan() {
			line := scanner.Text()
			if len(line) < 5 || strings.HasPrefix(line, "id:") ||
				strings.HasPrefix(line, "event:") || strings.HasPrefix(line, ":HTTP_STATUS/200") {
				continue
			}

			if !strings.HasPrefix(line, "data:") {
				continue
			}

			content = line[5:]
			var resp qWenResp
			if len(contents) == 0 { // 发送消息头
				if !outPutStart {
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
					outPutStart = true
					continue
				} else {
					// 处理代码换行
					content = "\n"
				}
			} else {
				err := utils.JsonDecode(content, &resp)
				if err != nil {
					logger.Error("error with parse data line: ", content)
					utils.ReplyMessage(ws, fmt.Sprintf("**解析数据行失败：%s**", err))
					break
				}
				if resp.Message != "" {
					utils.ReplyMessage(ws, fmt.Sprintf("**API 返回错误：%s**", resp.Message))
					break
				}
			}

			//通过比较 lastText（上一次的文本）和 currentText（当前的文本），
			//提取出新添加的文本部分。然后只将这部分新文本发送到客户端。
			//每次循环结束后，lastText 会更新为当前的完整文本，以便于下一次循环进行比较。
			currentText := resp.Output.Text
			if currentText != lastText {
				// 提取新增文本
				newText = strings.Replace(currentText, lastText, "", 1)
				utils.ReplyChunkMessage(ws, types.WsMessage{
					Type:    types.WsMiddle,
					Content: utils.InterfaceToString(newText),
				})
				lastText = currentText // 更新 lastText
			}
			contents = append(contents, newText)

			if resp.Output.FinishReason == "stop" {
				break
			}

		} //end for

		if err := scanner.Err(); err != nil {
			if strings.Contains(err.Error(), "context canceled") {
				logger.Info("用户取消了请求：", prompt)
			} else {
				logger.Error("信息读取出错：", err)
			}
		}

		// 消息发送成功
		if len(contents) > 0 {
			h.saveChatHistory(req, prompt, contents, message, chatCtx, session, role, userVo, promptCreatedAt, replyCreatedAt)
		}
	} else {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("请求大模型 API 失败：%s", body)
	}

	return nil
}
