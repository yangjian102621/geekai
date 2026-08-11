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

// VideoAdapter 视频生成适配器接口
type VideoAdapter interface {
	// CreateTask 创建视频生成任务
	CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error)

	// QueryTask 查询任务状态
	QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error)

	// GetProvider 获取服务提供商名称（不带版本号：veo, sora, luma）
	GetProvider() string
}

// VeoAdapter Veo 视频生成适配器
type VeoAdapter struct {
	db         *gorm.DB
	httpClient *req.Client
}

// NewVeoAdapter 创建 Veo 适配器
func NewVeoAdapter(db *gorm.DB) *VeoAdapter {
	return &VeoAdapter{
		db:         db,
		httpClient: req.C().SetTimeout(time.Minute * 3),
	}
}

// GetProvider 获取服务提供商名称
func (a *VeoAdapter) GetProvider() string {
	return "veo"
}

// VeoCreateRequest Veo 创建任务请求
type VeoCreateRequest struct {
	Prompt         string   `json:"prompt"`
	Model          string   `json:"model"`
	EnhancePrompt  bool     `json:"enhance_prompt,omitempty"`
	EnableUpsample bool     `json:"enable_upsample,omitempty"`
	AspectRatio    string   `json:"aspect_ratio,omitempty"`
	Images         []string `json:"images,omitempty"` // 图生视频时使用
}

// VeoCreateResponse Veo 创建任务响应
type VeoCreateResponse struct {
	TaskId string `json:"task_id"`
}

// VeoQueryResponse Veo 查询任务响应
type VeoQueryResponse struct {
	TaskId     string       `json:"task_id"`
	Platform   string       `json:"platform"`
	Action     string       `json:"action"`
	Status     string       `json:"status"`
	FailReason string       `json:"fail_reason"`
	SubmitTime int64        `json:"submit_time"`
	StartTime  int64        `json:"start_time"`
	FinishTime int64        `json:"finish_time"`
	Progress   string       `json:"progress"`
	Data       VeoQueryData `json:"data"`
	SearchItem string       `json:"search_item"`
}

// VeoQueryData Veo 查询响应中的 data 字段
type VeoQueryData struct {
	Output string `json:"output"`
}

// CreateTask 创建视频生成任务
func (a *VeoAdapter) CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error) {
	// 解析任务参数
	paramsMap, ok := task.Params.(map[string]any)
	if !ok {
		return CreateTaskResponse{}, fmt.Errorf("invalid params type for Veo video task")
	}

	// 构建请求参数
	reqBody := VeoCreateRequest{
		Prompt: task.Prompt,
	}

	// 从 params 中提取参数
	if model, ok := paramsMap["model"].(string); ok {
		reqBody.Model = model
	}
	if enhancePrompt, ok := paramsMap["enhance_prompt"].(bool); ok {
		reqBody.EnhancePrompt = enhancePrompt
	}
	if enableUpsample, ok := paramsMap["enable_upsample"].(bool); ok {
		reqBody.EnableUpsample = enableUpsample
	}
	if aspectRatio, ok := paramsMap["aspect_ratio"].(string); ok {
		reqBody.AspectRatio = aspectRatio
	}

	// 处理图生视频（images 参数）
	if images, ok := paramsMap["images"].([]interface{}); ok {
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
	var res VeoCreateResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	return CreateTaskResponse{
		TaskId:    res.TaskId,
		Channel:   videoConfig.ApiURL,
		Prompt:    task.Prompt,
		State:     types.VideoStatusPending,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// QueryTask 查询任务状态
func (a *VeoAdapter) QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error) {
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
	var res VeoQueryResponse
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
	case "failed":
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
