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
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"strings"
	"time"
)

// 清华大学 ChatGML 消息发送实现

func (h *ChatHandler) sendChatGLMMessage(
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
		var event, content string
		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) < 5 || strings.HasPrefix(line, "id:") {
				continue
			}
			if strings.HasPrefix(line, "event:") {
				event = line[6:]
				continue
			}

			if strings.HasPrefix(line, "data:") {
				content = line[5:]
			}
			// 处理代码换行
			if len(content) == 0 {
				content = "\n"
			}
			switch event {
			case "add":
				if len(contents) == 0 {
					utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
				}
				utils.ReplyChunkMessage(ws, types.WsMessage{
					Type:    types.WsMiddle,
					Content: utils.InterfaceToString(content),
				})
				contents = append(contents, content)
			case "finish":
				break
			case "error":
				utils.ReplyMessage(ws, fmt.Sprintf("**调用 ChatGLM API 出错：%s**", content))
				break
			case "interrupted":
				utils.ReplyMessage(ws, "**调用 ChatGLM API 出错，当前输出被中断！**")
			}

		} // end for

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

func (h *ChatHandler) getChatGLMToken(apiKey string) (string, error) {
	ctx := context.Background()
	tokenString, err := h.redis.Get(ctx, apiKey).Result()
	if err == nil {
		return tokenString, nil
	}

	expr := time.Hour * 2
	key := strings.Split(apiKey, ".")
	if len(key) != 2 {
		return "", fmt.Errorf("invalid api key: %s", apiKey)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"api_key":   key[0],
		"timestamp": time.Now().Unix(),
		"exp":       time.Now().Add(expr).Add(time.Second * 10).Unix(),
	})
	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"
	delete(token.Header, "typ")
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString([]byte(key[1]))
	h.redis.Set(ctx, apiKey, tokenString, expr)
	return tokenString, err
}
