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
	"gorm.io/gorm"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

// OPenAI 消息发送实现
func (h *ChatHandler) sendOpenAiMessage(
	chatCtx []interface{},
	req types.ApiRequest,
	userVo vo.User,
	ctx context.Context,
	session *types.ChatSession,
	role model.ChatRole,
	prompt string,
	ws *types.WsClient) error {
	promptCreatedAt := time.Now() // 记录提问时间
	start := time.Now()
	var apiKey = userVo.ChatConfig.ApiKeys[session.Model.Platform]
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
		var functionCall = false
		var functionName string
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

			fun := responseBody.Choices[0].Delta.FunctionCall
			if functionCall && fun.Name == "" {
				arguments = append(arguments, fun.Arguments)
				continue
			}

			if !utils.IsEmptyValue(fun) {
				functionName = fun.Name
				f := h.App.Functions[functionName]
				if f != nil {
					functionCall = true
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: fmt.Sprintf("正在调用函数 `%s` 作答 ...\n\n", f.Name())})
				}
				continue
			}

			if responseBody.Choices[0].FinishReason == "function_call" { // 函数调用完毕
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

		if functionCall { // 调用函数完成任务
			var params map[string]interface{}
			_ = utils.JsonDecode(strings.Join(arguments, ""), &params)
			logger.Debugf("函数名称: %s, 函数参数：%s", functionName, params)

			// for creating image, check if the user's img_calls > 0
			if functionName == types.FuncMidJourney && userVo.ImgCalls <= 0 {
				utils.ReplyMessage(ws, "**当前用户剩余绘图次数已用尽，请扫描下面二维码联系管理员！**")
				utils.ReplyMessage(ws, ErrImg)
			} else {
				f := h.App.Functions[functionName]
				if functionName == types.FuncMidJourney {
					params["user_id"] = userVo.Id
					params["role_id"] = role.Id
					params["chat_id"] = session.ChatId
					params["icon"] = "/images/avatar/mid_journey.png"
					params["session_id"] = session.SessionId
				}
				data, err := f.Invoke(params)
				if err != nil {
					msg := "调用函数出错：" + err.Error()
					utils.ReplyChunkMessage(ws, types.WsMessage{
						Type:    types.WsMiddle,
						Content: msg,
					})
					contents = append(contents, msg)
				} else {
					content := data
					if functionName == types.FuncMidJourney {
						content = fmt.Sprintf("绘画提示词：%s 已推送任务到 MidJourney 机器人，请耐心等待任务执行...", data)
						h.mjService.ChatClients.Put(session.SessionId, ws)
						// update user's img_calls
						h.db.Model(&model.User{}).Where("id = ?", userVo.Id).UpdateColumn("img_calls", gorm.Expr("img_calls - ?", 1))
					}

					utils.ReplyChunkMessage(ws, types.WsMessage{
						Type:    types.WsMiddle,
						Content: content,
					})
					contents = append(contents, content)
				}
			}
		}

		// 消息发送成功
		if len(contents) > 0 {
			// 更新用户的对话次数
			h.subUserCalls(userVo, session)

			if message.Role == "" {
				message.Role = "assistant"
			}
			message.Content = strings.Join(contents, "")
			useMsg := types.Message{Role: "user", Content: prompt}

			// 更新上下文消息，如果是调用函数则不需要更新上下文
			if h.App.ChatConfig.EnableContext && functionCall == false {
				chatCtx = append(chatCtx, useMsg)  // 提问消息
				chatCtx = append(chatCtx, message) // 回复消息
				h.App.ChatContexts.Put(session.ChatId, chatCtx)
			}

			// 追加聊天记录
			if h.App.ChatConfig.EnableHistory {
				useContext := true
				if functionCall {
					useContext = false
				}

				// for prompt
				promptToken, err := utils.CalcTokens(prompt, req.Model)
				if err != nil {
					logger.Error(err)
				}
				historyUserMsg := model.HistoryMessage{
					UserId:     userVo.Id,
					ChatId:     session.ChatId,
					RoleId:     role.Id,
					Type:       types.PromptMsg,
					Icon:       userVo.Avatar,
					Content:    prompt,
					Tokens:     promptToken,
					UseContext: useContext,
				}
				historyUserMsg.CreatedAt = promptCreatedAt
				historyUserMsg.UpdatedAt = promptCreatedAt
				res := h.db.Save(&historyUserMsg)
				if res.Error != nil {
					logger.Error("failed to save prompt history message: ", res.Error)
				}

				// 计算本次对话消耗的总 token 数量
				var totalTokens = 0
				if functionCall { // prompt + 函数名 + 参数 token
					tokens, _ := utils.CalcTokens(functionName, req.Model)
					totalTokens += tokens
					tokens, _ = utils.CalcTokens(utils.InterfaceToString(arguments), req.Model)
					totalTokens += tokens
				} else {
					totalTokens, _ = utils.CalcTokens(message.Content, req.Model)
				}
				totalTokens += getTotalTokens(req)

				historyReplyMsg := model.HistoryMessage{
					UserId:     userVo.Id,
					ChatId:     session.ChatId,
					RoleId:     role.Id,
					Type:       types.ReplyMsg,
					Icon:       role.Icon,
					Content:    message.Content,
					Tokens:     totalTokens,
					UseContext: useContext,
				}
				historyReplyMsg.CreatedAt = replyCreatedAt
				historyReplyMsg.UpdatedAt = replyCreatedAt
				res = h.db.Create(&historyReplyMsg)
				if res.Error != nil {
					logger.Error("failed to save reply history message: ", res.Error)
				}

				// 更新用户信息
				h.db.Model(&model.User{}).Where("id = ?", userVo.Id).
					UpdateColumn("total_tokens", gorm.Expr("total_tokens + ?", totalTokens))
			}

			// 保存当前会话
			var chatItem model.ChatItem
			res := h.db.Where("chat_id = ?", session.ChatId).First(&chatItem)
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
				h.db.Create(&chatItem)
			}
		}
	} else {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("error with reading response: %v", err)
		}
		var res types.ApiError
		err = json.Unmarshal(body, &res)
		if err != nil {
			return fmt.Errorf("error with decode response: %v", err)
		}

		// OpenAI API 调用异常处理
		if strings.Contains(res.Error.Message, "This key is associated with a deactivated account") {
			utils.ReplyMessage(ws, "请求 OpenAI API 失败：API KEY 所关联的账户被禁用。")
			// 移除当前 API key
			h.db.Where("value = ?", apiKey).Delete(&model.ApiKey{})
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
