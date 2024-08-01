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
	"io"
	"net/http"
	"strings"
	"time"
)

type baiduResp struct {
	Id               string `json:"id"`
	Object           string `json:"object"`
	Created          int    `json:"created"`
	SentenceId       int    `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	IsTruncated      bool   `json:"is_truncated"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// 百度文心一言消息发送实现

func (h *ChatHandler) sendBaiduMessage(
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
		logger.Error(err)
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
		var content string
		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) < 5 || strings.HasPrefix(line, "id:") {
				continue
			}

			if strings.HasPrefix(line, "data:") {
				content = line[5:]
			}

			// 处理代码换行
			if len(content) == 0 {
				content = "\n"
			}

			var resp baiduResp
			err := utils.JsonDecode(content, &resp)
			if err != nil {
				logger.Error("error with parse data line: ", err)
				utils.ReplyMessage(ws, fmt.Sprintf("**解析数据行失败：%s**", err))
				break
			}

			if len(contents) == 0 {
				utils.ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
			}
			utils.ReplyChunkMessage(ws, types.WsMessage{
				Type:    types.WsMiddle,
				Content: utils.InterfaceToString(resp.Result),
			})
			contents = append(contents, resp.Result)

			if resp.IsTruncated {
				utils.ReplyMessage(ws, "AI 输出异常中断")
				break
			}

			if resp.IsEnd {
				break
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

func (h *ChatHandler) getBaiduToken(apiKey string) (string, error) {
	ctx := context.Background()
	tokenString, err := h.redis.Get(ctx, apiKey).Result()
	if err == nil {
		return tokenString, nil
	}

	expr := time.Hour * 24 * 20 // access_token 有效期
	key := strings.Split(apiKey, "|")
	if len(key) != 2 {
		return "", fmt.Errorf("invalid api key: %s", apiKey)
	}
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", key[0], key[1])
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error with send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error with read response: %w", err)
	}
	var r map[string]interface{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("error with parse response: %w", err)
	}

	if r["error"] != nil {
		return "", fmt.Errorf("error with api response: %s", r["error_description"])
	}

	tokenString = fmt.Sprintf("%s", r["access_token"])
	h.redis.Set(ctx, apiKey, tokenString, expr)
	return tokenString, nil
}
