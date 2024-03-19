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

	req2 "github.com/imroc/req/v3"
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
		if response.Body != nil {
			all, _ := io.ReadAll(response.Body)
			logger.Error(string(all))
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
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.Contains(line, "data:") || len(line) < 30 {
				continue
			}

			var responseBody = types.ApiResponse{}
			err = json.Unmarshal([]byte(line[6:]), &responseBody)
			if err != nil || len(responseBody.Choices) == 0 { // 数据解析出错
				logger.Error(err, line)
				utils.ReplyMessage(ws, ErrorMsg)
				utils.ReplyMessage(ws, ErrImg)
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
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: fmt.Sprintf("正在调用工具 `%s` 作答 ...\n\n", function.Label)})
				}
				continue
			}

			if responseBody.Choices[0].FinishReason == "tool_calls" ||
				responseBody.Choices[0].FinishReason == "function_call" { // 函数调用完毕
				break
			}

			// 初始化 role
			if responseBody.Choices[0].Delta.Role != "" && message.Role == "" {
				message.Role = responseBody.Choices[0].Delta.Role
				utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
				continue
			} else if responseBody.Choices[0].FinishReason != "" {
				break // 输出完成或者输出中断了
			} else {
				content := responseBody.Choices[0].Delta.Content
				contents = append(contents, utils.InterfaceToString(content))
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
			var params map[string]interface{}
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
			if message.Role == "" {
				message.Role = "assistant"
			}
			message.Content = strings.Join(contents, "")
			useMsg := types.Message{Role: "user", Content: prompt}

			// 更新上下文消息，如果是调用函数则不需要更新上下文
			if h.App.SysConfig.EnableContext && toolCall == false {
				chatCtx = append(chatCtx, useMsg)  // 提问消息
				chatCtx = append(chatCtx, message) // 回复消息
				h.App.ChatContexts.Put(session.ChatId, chatCtx)
			}

			// 追加聊天记录
			useContext := true
			if toolCall {
				useContext = false
			}

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
				UseContext: useContext,
				Model:      req.Model,
			}
			historyUserMsg.CreatedAt = promptCreatedAt
			historyUserMsg.UpdatedAt = promptCreatedAt
			res := h.DB.Save(&historyUserMsg)
			if res.Error != nil {
				logger.Error("failed to save prompt history message: ", res.Error)
			}

			// 计算本次对话消耗的总 token 数量
			var replyTokens = 0
			if toolCall { // prompt + 函数名 + 参数 token
				tokens, _ := utils.CalcTokens(function.Name, req.Model)
				replyTokens += tokens
				tokens, _ = utils.CalcTokens(utils.InterfaceToString(arguments), req.Model)
				replyTokens += tokens
			} else {
				replyTokens, _ = utils.CalcTokens(message.Content, req.Model)
			}
			replyTokens += getTotalTokens(req)

			historyReplyMsg := model.ChatMessage{
				UserId:     userVo.Id,
				ChatId:     session.ChatId,
				RoleId:     role.Id,
				Type:       types.ReplyMsg,
				Icon:       role.Icon,
				Content:    h.extractImgUrl(message.Content),
				Tokens:     replyTokens,
				UseContext: useContext,
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
			utils.ReplyMessage(ws, "请求 OpenAI API 失败："+err.Error())
			return fmt.Errorf("error with reading response: %v", err)
		}
		var res types.ApiError
		err = json.Unmarshal(body, &res)
		if err != nil {
			utils.ReplyMessage(ws, "请求 OpenAI API 失败：\n"+"```\n"+string(body)+"```")
			return fmt.Errorf("error with decode response: %v", err)
		}

		// OpenAI API 调用异常处理
		if strings.Contains(res.Error.Message, "This key is associated with a deactivated account") {
			utils.ReplyMessage(ws, "请求 OpenAI API 失败：API KEY 所关联的账户被禁用。")
			// 移除当前 API key
			h.DB.Where("value = ?", apiKey).Delete(&model.ApiKey{})
		} else if strings.Contains(res.Error.Message, "You exceeded your current quota") {
			utils.ReplyMessage(ws, "请求 OpenAI API 失败：API KEY 触发并发限制，请稍后再试。")
		} else if strings.Contains(res.Error.Message, "This model's maximum context length") {
			logger.Error(res.Error.Message)
			utils.ReplyMessage(ws, "当前会话上下文长度超出限制，已为您清空会话上下文！")
			h.App.ChatContexts.Delete(session.ChatId)
			return h.sendMessage(ctx, session, role, prompt, ws)
		} else {
			utils.ReplyMessage(ws, "请求 OpenAI API 失败："+res.Error.Message)
		}
	}

	return nil
}
