package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const maxModelRequestBytes = 16 << 20

func (h *ModelAPIHandler) ChatCompletions(c *gin.Context) {
	user := c.MustGet("model_api_user").(model.User)
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, maxModelRequestBytes))
	if err != nil {
		modelAPIError(c, http.StatusRequestEntityTooLarge, "invalid_request_error", "请求体超过限制或无法读取")
		return
	}
	var input struct {
		Model    string `json:"model"`
		Messages []struct {
			Role string `json:"role"`
		} `json:"messages"`
		Stream bool `json:"stream"`
		N      *int `json:"n"`
	}
	if json.Unmarshal(body, &input) != nil || input.Model == "" || len(input.Messages) == 0 || (input.N != nil && *input.N != 1) {
		modelAPIError(c, http.StatusBadRequest, "invalid_request_error", "请提供 model 和 messages，每次请求仅支持一条回复（n=1）")
		return
	}
	for _, message := range input.Messages {
		if message.Role == "" {
			modelAPIError(c, http.StatusBadRequest, "invalid_request_error", "每条消息必须包含 role")
			return
		}
	}
	var chatModel model.ChatModel
	err = h.availableModels(user).Where("value = ?", input.Model).Order("sort_num ASC, id ASC").First(&chatModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			modelAPIError(c, http.StatusNotFound, "model_not_found", "模型不存在或无权使用")
		} else {
			modelAPIError(c, http.StatusInternalServerError, "server_error", "查询模型失败")
		}
		return
	}
	// 使用 RawMessage 保留工具、多模态和厂商扩展字段。
	var payload map[string]json.RawMessage
	_ = json.Unmarshal(body, &payload)
	if err = limitChatTokens(payload, chatModel); err != nil {
		modelAPIError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}
	if _, ok := payload["temperature"]; !ok {
		payload["temperature"], _ = json.Marshal(chatModel.Temperature)
	}
	body, err = json.Marshal(payload)
	if err != nil {
		modelAPIError(c, http.StatusBadRequest, "invalid_request_error", "无效的模型请求")
		return
	}
	if chatModel.Power < 0 {
		modelAPIError(c, http.StatusInternalServerError, "server_error", "模型积分配置无效")
		return
	}
	if chatModel.Power > 0 {
		err = h.userService.DecreasePower(user.Id, chatModel.Power, model.PowerLog{
			Type: types.PowerConsume, Model: chatModel.Value, Remark: "ThinkBuddy 模型 API 调用",
		})
		if err != nil {
			if errors.Is(err, service.ErrInsufficientPower) {
				modelAPIError(c, http.StatusPaymentRequired, "insufficient_quota", "积分不足，请充值后重试")
			} else {
				modelAPIError(c, http.StatusInternalServerError, "server_error", "扣减积分失败")
			}
			return
		}
	}
	charged := false
	defer func() {
		// 已产生回复的流式中断仍按一次调用收费；没有有效回复时退回。
		if chatModel.Power > 0 && !charged {
			if err := h.userService.IncreasePower(user.Id, chatModel.Power, model.PowerLog{
				Type: types.PowerRefund, Model: chatModel.Value, Remark: "ThinkBuddy 模型 API 未产生回复，退回积分",
			}); err != nil {
				logger.Errorf("模型 API 退费失败: user=%d model=%s power=%d error=%v", user.Id, chatModel.Value, chatModel.Power, err)
			}
		}
	}()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()
	var key model.ApiKey
	response, err := service.RequestChatCompletion(ctx, h.DB, chatModel.KeyId, body, &key)
	if err != nil {
		modelAPIError(c, http.StatusBadGateway, "upstream_error", "模型服务暂时不可用，请稍后重试")
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		status := http.StatusBadGateway
		if response.StatusCode == http.StatusTooManyRequests {
			status = http.StatusTooManyRequests
		}
		modelAPIError(c, status, "upstream_error", "模型服务请求失败，请稍后重试")
		return
	}
	if input.Stream {
		if !strings.HasPrefix(response.Header.Get("Content-Type"), "text/event-stream") {
			modelAPIError(c, http.StatusBadGateway, "upstream_error", "模型服务未返回流式响应")
			return
		}
		charged = relayChatStream(c, response.Body)
		return
	}
	result, err := io.ReadAll(io.LimitReader(response.Body, (32<<20)+1))
	var reply struct {
		Choices []struct {
			Message json.RawMessage `json:"message"`
		} `json:"choices"`
		Error json.RawMessage `json:"error"`
	}
	if err != nil || len(result) > 32<<20 || json.Unmarshal(result, &reply) != nil || len(reply.Choices) == 0 || !hasJSONValue(reply.Choices[0].Message) || hasJSONValue(reply.Error) {
		modelAPIError(c, http.StatusBadGateway, "upstream_error", "模型服务返回了无效响应")
		return
	}
	charged = true
	c.Data(http.StatusOK, "application/json; charset=utf-8", result)
}

func limitChatTokens(payload map[string]json.RawMessage, config model.ChatModel) error {
	limit := config.MaxTokens
	if limit <= 0 {
		limit = 1024
	}
	_, hasMax := payload["max_tokens"]
	_, hasCompletionMax := payload["max_completion_tokens"]
	if hasMax && hasCompletionMax {
		return errors.New("max_tokens 和 max_completion_tokens 只能填写一个")
	}
	field := "max_tokens"
	if hasCompletionMax {
		field = "max_completion_tokens"
	}
	if raw, ok := payload[field]; ok {
		var value int
		if json.Unmarshal(raw, &value) != nil || value <= 0 || value > limit {
			return fmt.Errorf("%s 必须为 1 到 %d 的整数", field, limit)
		}
	} else {
		if strings.HasPrefix(config.Value, "o1") || strings.HasPrefix(config.Value, "o3") || strings.HasPrefix(config.Value, "gpt-5") {
			field = "max_completion_tokens"
		}
		payload[field], _ = json.Marshal(limit)
	}
	return nil
}

func hasJSONValue(value json.RawMessage) bool {
	switch strings.TrimSpace(string(value)) {
	case "", "null", `""`, "[]", "{}":
		return false
	}
	return true
}

func relayChatStream(c *gin.Context, body io.Reader) bool {
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache, no-store")
	c.Header("X-Accel-Buffering", "no")
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	var lines, data []string
	charged, hadChoice := false, false
	eventSize := 0
	streamError := func() {
		payload, _ := json.Marshal(gin.H{"error": gin.H{"message": "模型响应中断，请重试", "type": "upstream_error", "code": "upstream_error"}})
		fmt.Fprintf(c.Writer, "data: %s\n\n", payload)
		c.Writer.Flush()
	}
	forward := func() (stop bool) {
		text := strings.Join(data, "\n")
		billable := false
		if len(data) > 0 {
			if text == "[DONE]" {
				if !hadChoice {
					streamError()
					return true
				}
			} else {
				var chunk struct {
					Choices []struct {
						Delta        map[string]json.RawMessage `json:"delta"`
						FinishReason *string                    `json:"finish_reason"`
					} `json:"choices"`
					Error json.RawMessage `json:"error"`
					Usage json.RawMessage `json:"usage"`
				}
				if json.Unmarshal([]byte(text), &chunk) != nil || hasJSONValue(chunk.Error) || (chunk.Choices == nil && !hasJSONValue(chunk.Usage)) {
					streamError()
					return true
				}
				for _, choice := range chunk.Choices {
					hadChoice = true
					billable = billable || choice.FinishReason != nil
					for _, field := range []string{"content", "reasoning_content", "tool_calls", "function_call"} {
						billable = billable || hasJSONValue(choice.Delta[field])
					}
				}
			}
		}
		n, err := io.WriteString(c.Writer, strings.Join(lines, "\n")+"\n\n")
		if n > 0 && billable {
			charged = true
		}
		c.Writer.Flush()
		return err != nil || text == "[DONE]"
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(lines) > 0 && forward() {
				return charged
			}
			lines, data, eventSize = nil, nil, 0
			continue
		}
		eventSize += len(line)
		if eventSize > 8<<20 {
			streamError()
			return charged
		}
		lines = append(lines, line)
		if value, ok := strings.CutPrefix(line, "data:"); ok {
			data = append(data, strings.TrimPrefix(value, " "))
		}
	}
	if scanner.Err() == nil && len(lines) > 0 && forward() {
		return charged
	}
	streamError()
	return charged
}
