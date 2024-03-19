package chatimpl

import (
	"bufio"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"
	"unicode/utf8"
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
	response, err := h.doRequest(ctx, req, session.Model.Platform, &apiKey)
	logger.Info("HTTP请求完成，耗时：", time.Now().Sub(start))
	if err != nil {
		if strings.Contains(err.Error(), "context canceled") {
			logger.Info("用户取消了请求：", prompt)
			return nil
		} else if strings.Contains(err.Error(), "no available key") {
			utils.ReplyMessage(ws, "抱歉😔😔😔，系统已经没有可用的 API KEY，请联系管理员！")
			return nil
		} else {
			logger.Error(err)
		}

		utils.ReplyMessage(ws, ErrorMsg)
		utils.ReplyMessage(ws, ErrImg)
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

			if strings.HasPrefix(line, "data:") {
				content = line[5:]
			}

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
			if message.Role == "" {
				message.Role = "assistant"
			}
			message.Content = strings.Join(contents, "")
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
	} else {
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
		utils.ReplyMessage(ws, "请求通义千问大模型 API 失败："+res.Msg)
	}

	return nil
}
