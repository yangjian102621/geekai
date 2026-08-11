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

// MiniMaxAdapter MiniMax 视频生成适配器
type MiniMaxAdapter struct {
	db         *gorm.DB
	httpClient *req.Client
}

// NewMiniMaxAdapter 创建 MiniMax 适配器
func NewMiniMaxAdapter(db *gorm.DB) *MiniMaxAdapter {
	return &MiniMaxAdapter{
		db:         db,
		httpClient: req.C().SetTimeout(time.Minute * 3),
	}
}

// GetProvider 获取服务提供商名称
func (a *MiniMaxAdapter) GetProvider() string {
	return "minimax"
}

// MiniMaxCreateRequest MiniMax 创建任务请求
type MiniMaxCreateRequest struct {
	Model           string `json:"model"`
	Prompt          string `json:"prompt"`
	Duration        int    `json:"duration,omitempty"`
	Resolution      string `json:"resolution,omitempty"`
	FirstFrameImage string `json:"first_frame_image,omitempty"`
	LastFrameImage  string `json:"last_frame_image,omitempty"`
	PromptOptimizer bool   `json:"prompt_optimizer,omitempty"`
}

// MiniMaxCreateResponse MiniMax 创建任务响应
type MiniMaxCreateResponse struct {
	TaskId   string `json:"task_id"`
	BaseResp struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

// MiniMaxFile MiniMax 文件信息
type MiniMaxFile struct {
	Bytes       int    `json:"bytes"`
	CreatedAt   int64  `json:"created_at"`
	DownloadURL string `json:"download_url"`
	FileId      int64  `json:"file_id"`
	Filename    string `json:"filename"`
	Purpose     string `json:"purpose"`
}

// MiniMaxQueryResponse MiniMax 查询任务响应
type MiniMaxQueryResponse struct {
	TaskId      string       `json:"task_id"`
	Status      string       `json:"status"`
	FileId      string       `json:"file_id,omitempty"` // 顶层 file_id 可能是字符串
	File        *MiniMaxFile `json:"file,omitempty"`    // file 对象包含详细信息
	VideoWidth  int          `json:"video_width,omitempty"`
	VideoHeight int          `json:"video_height,omitempty"`
	VideoURL    string       `json:"video_url,omitempty"`
	Prompt      string       `json:"prompt,omitempty"`
	ErrMsg      string       `json:"err_msg,omitempty"`
	StatusMsg   string       `json:"status_msg,omitempty"`
	BaseResp    struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

// CreateTask 创建视频生成任务
func (a *MiniMaxAdapter) CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error) {
	// 解析任务参数
	paramsMap, ok := task.Params.(map[string]interface{})
	if !ok {
		return CreateTaskResponse{}, fmt.Errorf("invalid params type for MiniMax video task")
	}

	// 构建请求参数
	reqBody := MiniMaxCreateRequest{
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
	if firstFrameImage, ok := paramsMap["first_frame_image"].(string); ok {
		reqBody.FirstFrameImage = firstFrameImage
	}
	if lastFrameImage, ok := paramsMap["last_frame_image"].(string); ok {
		reqBody.LastFrameImage = lastFrameImage
	}
	if promptOptimizer, ok := paramsMap["prompt_optimizer"].(bool); ok {
		reqBody.PromptOptimizer = promptOptimizer
	} else {
		reqBody.PromptOptimizer = true // 默认值
	}

	// 发送请求
	apiURL := fmt.Sprintf("%s/minimax/v1/video_generation", videoConfig.ApiURL)
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
	var res MiniMaxCreateResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	if res.BaseResp.StatusCode != 0 {
		return CreateTaskResponse{}, fmt.Errorf("API 返回错误：%s", res.BaseResp.StatusMsg)
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
func (a *MiniMaxAdapter) QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error) {
	// MiniMax 查询接口
	apiURL := fmt.Sprintf("%s/minimax/v1/query/video_generation?task_id=%s", videoConfig.ApiURL, taskId)
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
	var res MiniMaxQueryResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	if res.BaseResp.StatusCode != 0 {
		return QueryTaskResponse{}, fmt.Errorf("API 返回错误：%s", res.BaseResp.StatusMsg)
	}

	// 转换状态（处理大小写）
	state := strings.ToLower(res.Status)
	switch state {
	case "completed", "succeed", "success":
		state = types.VideoStatusSuccess
	case "in_progress", "running":
		state = types.VideoStatusInProgress
	case "failed":
		state = types.VideoStatusFailed
	default:
		state = types.VideoStatusPending
	}

	// 获取视频URL，优先从 file.download_url 获取
	videoURL := res.VideoURL
	if res.File != nil && res.File.DownloadURL != "" {
		videoURL = res.File.DownloadURL
	}

	// 构建响应
	response := QueryTaskResponse{
		TaskId:    res.TaskId,
		Status:    state,
		VideoURL:  videoURL,
		Prompt:    res.Prompt,
		ErrMsg:    res.ErrMsg,
		StatusMsg: res.StatusMsg,
	}

	// 如果有原始数据，转换为 JSON 字符串
	if len(body) > 0 {
		response.Output = string(body)
	}

	return response, nil
}
