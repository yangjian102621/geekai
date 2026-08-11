package adapters

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"io"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

// WanAdapter Wan（通义万相）视频生成适配器
type WanAdapter struct {
	db         *gorm.DB
	httpClient *req.Client
}

// NewWanAdapter 创建 Wan 适配器
func NewWanAdapter(db *gorm.DB) *WanAdapter {
	return &WanAdapter{
		db:         db,
		httpClient: req.C().SetTimeout(time.Minute * 3),
	}
}

// GetProvider 获取服务提供商名称
func (a *WanAdapter) GetProvider() string {
	return "wan"
}

// WanCreateRequest Wan 创建任务请求
type WanCreateRequest struct {
	Prompt         string   `json:"prompt"`
	Model          string   `json:"model"`
	Duration       int      `json:"duration,omitempty"`
	Resolution     string   `json:"resolution,omitempty"`
	NegativePrompt string   `json:"negative_prompt,omitempty"`
	Images         []string `json:"images,omitempty"`
	PromptExtend   bool     `json:"prompt_extend,omitempty"`
}

// WanCreateResponse Wan 创建任务响应
type WanCreateResponse struct {
	TaskId string `json:"task_id"`
}

// WanQueryResponse Wan 查询任务响应
type WanQueryResponse struct {
	TaskId     string       `json:"task_id"`
	Platform   string       `json:"platform"`
	Action     string       `json:"action"`
	Status     string       `json:"status"`
	FailReason string       `json:"fail_reason"`
	SubmitTime int64        `json:"submit_time"`
	StartTime  int64        `json:"start_time"`
	FinishTime int64        `json:"finish_time"`
	Progress   string       `json:"progress"`
	Data       WanQueryData `json:"data"`
	SearchItem string       `json:"search_item"`
}

// WanQueryData Wan 查询响应中的 data 字段
type WanQueryData struct {
	Output string `json:"output"`
}

// CreateTask 创建视频生成任务
func (a *WanAdapter) CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error) {
	// 解析任务参数
	paramsMap, ok := task.Params.(map[string]any)
	if !ok {
		return CreateTaskResponse{}, fmt.Errorf("invalid params type for Wan video task")
	}

	// 构建请求参数
	reqBody := WanCreateRequest{
		Prompt: task.Prompt,
	}

	// 从 params 中提取参数
	if model, ok := paramsMap["model"].(string); ok {
		reqBody.Model = model
	}
	if duration, ok := paramsMap["duration"].(float64); ok {
		reqBody.Duration = int(duration)
	} else if duration, ok := paramsMap["duration"].(int); ok {
		reqBody.Duration = duration
	}
	if resolution, ok := paramsMap["resolution"].(string); ok {
		reqBody.Resolution = resolution
	}
	if images, ok := paramsMap["images"].([]any); ok {
		imageUrls := make([]string, 0)
		for _, img := range images {
			if imgStr, ok := img.(string); ok {
				imageUrls = append(imageUrls, imgStr)
			}
		}
		if len(imageUrls) > 0 {
			reqBody.Images = imageUrls
		}
	}
	if negativePrompt, ok := paramsMap["negative_prompt"].(string); ok {
		reqBody.NegativePrompt = negativePrompt
	}
	if promptExtend, ok := paramsMap["prompt_extend"].(bool); ok {
		reqBody.PromptExtend = promptExtend
	}

	logger.Debugf("WanCreateRequest: %+v", reqBody)

	// 发送请求
	apiURL := fmt.Sprintf("%s/v2/videos/generations", videoConfig.ApiURL)
	r, err := a.httpClient.R().
		SetHeader("Authorization", "Bearer "+videoConfig.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post(apiURL)

	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("请求 API 出错：%v", err)
	}

	if r.StatusCode != 200 && r.StatusCode != 201 {
		body, _ := io.ReadAll(r.Body)
		return CreateTaskResponse{}, fmt.Errorf("请求 API 出错：%d, %s", r.StatusCode, string(body))
	}

	body, _ := io.ReadAll(r.Body)
	var res WanCreateResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	return CreateTaskResponse{
		TaskId:    res.TaskId,
		Channel:   videoConfig.ApiURL,
		Prompt:    task.Prompt,
		State:     "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// QueryTask 查询任务状态
func (a *WanAdapter) QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error) {
	apiURL := fmt.Sprintf("%s/v2/videos/generations/%s", videoConfig.ApiURL, taskId)
	r, err := a.httpClient.R().
		SetHeader("Authorization", "Bearer "+videoConfig.ApiKey).
		Get(apiURL)

	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("请求 API 失败：%v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		body, _ := io.ReadAll(r.Body)
		return QueryTaskResponse{}, fmt.Errorf("API 返回失败：%d, %s", r.StatusCode, string(body))
	}

	body, _ := io.ReadAll(r.Body)
	var res WanQueryResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	// 转换状态（SUCCESS -> success, FAILED -> failed, 其他保持原样）
	state := strings.ToLower(res.Status)
	switch state {
	case "in_progress", "running":
		state = types.VideoStatusInProgress
	case "completed", "succeed", "success":
		state = types.VideoStatusSuccess
	case "failed", "failure":
		state = types.VideoStatusFailed
	default:
		state = types.VideoStatusPending
	}

	// 解析进度（从 "100%" 转换为 100）
	progress := 0
	if res.Progress != "" {
		// 移除 % 符号并转换为整数
		progressStr := strings.TrimSuffix(res.Progress, "%")
		if p, err := fmt.Sscanf(progressStr, "%d", &progress); err == nil && p == 1 {
			// 成功解析
		}
	}

	// 从 data.output 中提取视频 URL
	videoURL := res.Data.Output

	// 构建响应
	response := QueryTaskResponse{
		TaskId:   res.TaskId,
		Status:   state,
		Progress: progress,
		VideoURL: videoURL,
		ErrMsg:   res.FailReason,
	}

	// 如果有原始数据，转换为 JSON 字符串
	if len(body) > 0 {
		response.Output = string(body)
	}

	return response, nil
}
