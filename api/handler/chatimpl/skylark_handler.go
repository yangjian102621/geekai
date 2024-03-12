package chatimpl

import (
	"bufio"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	Region  = "cn-beijing"
	Service = "ml_maas"
)

type skylarkResp struct {
	Id     string `json:"req_id"`
	Choice struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choice"`
}

// Skylark 消息发送实现
func (h *ChatHandler) sendSkylarkMessage(
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
	var apiKey model.ApiKey
	response, err := h.doRequest(ctx, req, session.Model.Platform, &apiKey)
	logger.Info("HTTP请求完成，耗时：", time.Since(start))
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

			if content == "[DONE]" {
				break
			}

			// 处理代码换行
			if len(content) == 0 {
				content = "\n"
			}

			var resp skylarkResp
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
				Content: utils.InterfaceToString(resp.Choice.Message.Content),
			})
			contents = append(contents, resp.Choice.Message.Content)

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
			// 更新用户的对话次数
			h.subUserCalls(userVo, session)

			if message.Role == "" {
				message.Role = "assistant"
			}
			message.Content = strings.Join(contents, "")
			useMsg := types.Message{Role: "user", Content: prompt}

			// 更新上下文消息，如果是调用函数则不需要更新上下文
			if h.App.ChatConfig.EnableContext {
				chatCtx = append(chatCtx, useMsg)  // 提问消息
				chatCtx = append(chatCtx, message) // 回复消息
				h.App.ChatContexts.Put(session.ChatId, chatCtx)
			}

			// 追加聊天记录
			if h.App.ChatConfig.EnableHistory {
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
				res := h.db.Save(&historyUserMsg)
				if res.Error != nil {
					logger.Error("failed to save prompt history message: ", res.Error)
				}

				// for reply
				// 计算本次对话消耗的总 token 数量
				replyToken, _ := utils.CalcTokens(message.Content, req.Model)
				totalTokens := replyToken + getTotalTokens(req)
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
				res = h.db.Create(&historyReplyMsg)
				if res.Error != nil {
					logger.Error("failed to save reply history message: ", res.Error)
				}
				// 更新用户信息
				h.incUserTokenFee(userVo.Id, totalTokens)
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
				chatItem.Model = req.Model
				h.db.Create(&chatItem)
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
		utils.ReplyMessage(ws, "请求 Skylark API 失败："+res.Msg)
	}

	return nil
}

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func getSignedKey(secretKey, date, region, service string) []byte {
	kDate := hmacSHA256([]byte(secretKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "request")

	return kSigning
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		logger.Debugf("input hash err:%s", err.Error())
	}

	return hash.Sum(nil)
}

type SkylarkRequest struct {
	Method          string
	Url             string
	Body            []byte
	AccessKeyID     string
	SecretAccessKey string
}

func (h *ChatHandler) makeSkylarkRequestUrl(params SkylarkRequest) (http.Header, error) {
	ul, err := url.Parse(params.Url)
	if err != nil {
		return nil, err
	}
	requestAddr := fmt.Sprintf("%s%s?%s", ul.Host, ul.Path, ul.Query().Encode())
	logger.Debugf("request addr: %s\n", requestAddr)

	var reqHeaders = http.Header{}

	now := time.Now()
	date := now.UTC().Format("20060102T150405Z")
	authDate := date[:8]
	reqHeaders.Set("X-Date", date)
	payload := hex.EncodeToString(hashSHA256(params.Body))
	reqHeaders.Set("X-Content-Sha256", payload)
	reqHeaders.Set("Content-Type", "application/json")

	queryString := strings.Replace(ul.Query().Encode(), "+", "%20", -1)
	signedHeaders := []string{"content-type", "host", "x-content-sha256", "x-date"}
	var headerList []string
	for _, header := range signedHeaders {
		if header == "host" {
			headerList = append(headerList, header+":"+ul.Host)
		} else {
			v := reqHeaders.Get(header)
			headerList = append(headerList, header+":"+strings.TrimSpace(v))
		}
	}
	headerString := strings.Join(headerList, "\n")

	canonicalString := strings.Join([]string{
		params.Method,
		ul.Path,
		queryString,
		headerString + "\n",
		strings.Join(signedHeaders, ";"),
		payload,
	}, "\n")
	logger.Debugf("canonical string:\n%s\n", canonicalString)

	hashedCanonicalString := hex.EncodeToString(hashSHA256([]byte(canonicalString)))
	logger.Debugf("hashed canonical string: %s\n", hashedCanonicalString)

	credentialScope := authDate + "/" + Region + "/" + Service + "/request"
	signString := strings.Join([]string{
		"HMAC-SHA256",
		date,
		credentialScope,
		hashedCanonicalString,
	}, "\n")
	logger.Debugf("sign string:\n%s\n", signString)

	// 3. 构建认证请求头
	signedKey := getSignedKey(params.SecretAccessKey, authDate, Region, Service)
	signature := hex.EncodeToString(hmacSHA256(signedKey, signString))
	logger.Debugf("signature: %s\n", signature)

	authorization := "HMAC-SHA256" +
		" Credential=" + params.AccessKeyID + "/" + credentialScope +
		", SignedHeaders=" + strings.Join(signedHeaders, ";") +
		", Signature=" + signature
	reqHeaders.Set("Authorization", authorization)
	logger.Debugf("authorization: %s\n", authorization)

	return reqHeaders, nil
}
