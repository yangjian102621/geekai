package adapters

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// KelingAdapter 可灵视频生成适配器
type KelingAdapter struct {
	db *gorm.DB
}

// NewKelingAdapter 创建可灵适配器
func NewKelingAdapter(db *gorm.DB) *KelingAdapter {
	return &KelingAdapter{
		db: db,
	}
}

// GetProvider 获取服务提供商名称
func (a *KelingAdapter) GetProvider() string {
	return "keling"
}

// KelingCreateRequest 可灵创建任务请求
type KelingCreateRequest struct {
	ModelName      string  `json:"model_name"`
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt,omitempty"`
	CfgScale       float64 `json:"cfg_scale,omitempty"`
	Mode           string  `json:"mode,omitempty"`
	AspectRatio    string  `json:"aspect_ratio,omitempty"`
	Duration       string  `json:"duration,omitempty"`
	Sound          bool    `json:"sound,omitempty"`
	Image          string  `json:"image,omitempty"`
	ImageTail      string  `json:"image_tail,omitempty"`
}

// KelingCreateResponse 可灵创建任务响应
type KelingCreateResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Data      struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
		CreatedAt  int64  `json:"created_at"`
		UpdatedAt  int64  `json:"updated_at"`
	} `json:"data"`
}

// KelingQueryResponse 可灵查询任务响应
type KelingQueryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskID        string `json:"task_id"`
		TaskStatus    string `json:"task_status"`
		TaskStatusMsg string `json:"task_status_msg"`
		CreatedAt     int64  `json:"created_at"`
		UpdatedAt     int64  `json:"updated_at"`
		TaskResult    struct {
			Images []struct {
				Index int    `json:"index"`
				URL   string `json:"url"`
			} `json:"images,omitempty"`
			Videos []struct {
				ID       string `json:"id"`
				URL      string `json:"url"`
				Duration string `json:"duration"`
			} `json:"videos,omitempty"`
		} `json:"task_result"`
	} `json:"data"`
}

// CreateTask 创建视频生成任务
func (a *KelingAdapter) CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error) {
	// 解析任务参数
	paramsMap, ok := task.Params.(map[string]interface{})
	if !ok {
		return CreateTaskResponse{}, errors.New("invalid params type for KeLing video task")
	}

	// 构建请求参数
	payload := KelingCreateRequest{
		Prompt: task.Prompt,
	}

	// 从 params 中提取参数
	if modelName, ok := paramsMap["model_name"].(string); ok {
		payload.ModelName = modelName
	}
	if prompt, ok := paramsMap["prompt"].(string); ok {
		payload.Prompt = prompt
	}
	if negativePrompt, ok := paramsMap["negative_prompt"].(string); ok {
		payload.NegativePrompt = negativePrompt
	}
	if cfgScale, ok := paramsMap["cfg_scale"].(float64); ok {
		payload.CfgScale = cfgScale
	}
	if mode, ok := paramsMap["mode"].(string); ok {
		payload.Mode = mode
	}
	if aspectRatio, ok := paramsMap["aspect_ratio"].(string); ok {
		payload.AspectRatio = aspectRatio
	}
	if duration, ok := paramsMap["duration"].(string); ok {
		payload.Duration = duration
	}

	if sound, ok := paramsMap["sound"].(bool); ok {
		payload.Sound = sound
	}

	// 处理图生视频
	taskType, ok := paramsMap["task_type"].(string)
	if ok && taskType == "image2video" {
		if image, ok := paramsMap["image"].(string); ok {
			payload.Image = image
		}
		if imageTail, ok := paramsMap["image_tail"].(string); ok {
			payload.ImageTail = imageTail
		}
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("failed to marshal payload: %v", err)
	}
	logger.Debugf("KelingCreateRequest: %+v", string(jsonPayload))

	// 发送请求
	url := fmt.Sprintf("%s/kling/v1/videos/%s", videoConfig.ApiURL, taskType)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+videoConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: time.Duration(30) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return CreateTaskResponse{}, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResponse KelingCreateResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return CreateTaskResponse{}, fmt.Errorf("failed to parse response: %v", err)
	}

	if apiResponse.Code != 0 {
		return CreateTaskResponse{}, fmt.Errorf("API error: %s", apiResponse.Message)
	}

	return CreateTaskResponse{
		TaskId:    apiResponse.Data.TaskID,
		Channel:   videoConfig.ApiURL,
		Prompt:    task.Prompt,
		State:     types.VideoStatusPending,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// QueryTask 查询任务状态
func (a *KelingAdapter) QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error) {
	// 从 taskId 中提取 action（可灵的 taskId 格式可能包含 action 信息）
	// 这里需要从任务信息中获取 task_type，暂时使用 text2video 作为默认值
	action := "text2video"

	// 尝试从 channel 或其他地方获取 action，这里简化处理
	// 实际应该从任务信息中获取

	url := fmt.Sprintf("%s/kling/v1/videos/%s/%s", videoConfig.ApiURL, action, taskId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+videoConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Duration(30) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return QueryTaskResponse{}, fmt.Errorf("unexpected status code: %d, %s", res.StatusCode, string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response KelingQueryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return QueryTaskResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Code != 0 {
		return QueryTaskResponse{}, fmt.Errorf("API error: %s", response.Message)
	}

	// 转换状态
	state := response.Data.TaskStatus
	status := state
	switch state {
	case "in_progress", "processing":
		status = types.VideoStatusInProgress
	case "completed", "succeed", "success":
		status = types.VideoStatusSuccess
	case "failed":
		status = types.VideoStatusFailed
	default:
		status = types.VideoStatusPending
	}

	// 构建响应
	result := QueryTaskResponse{
		TaskId:    response.Data.TaskID,
		Status:    status,
		ErrMsg:    response.Data.TaskStatusMsg,
		StatusMsg: response.Data.TaskStatusMsg,
		Output:    string(body),
	}

	// 提取视频URL
	if len(response.Data.TaskResult.Videos) > 0 {
		result.VideoURL = response.Data.TaskResult.Videos[0].URL
	}

	return result, nil
}
