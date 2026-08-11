package adapters

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"io"
	"time"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

// LumaAdapter Luma 视频生成适配器
type LumaAdapter struct {
	db         *gorm.DB
	httpClient *req.Client
}

// NewLumaAdapter 创建 Luma 适配器
func NewLumaAdapter(db *gorm.DB) *LumaAdapter {
	return &LumaAdapter{
		db:         db,
		httpClient: req.C().SetTimeout(time.Minute * 3),
	}
}

// GetProvider 获取服务提供商名称
func (a *LumaAdapter) GetProvider() string {
	return "luma"
}

// LumaCreateRequest Luma 创建任务请求
type LumaCreateRequest struct {
	ModelName    string `json:"model_name"`
	UserPrompt   string `json:"user_prompt"`
	ExpandPrompt bool   `json:"expand_prompt,omitempty"`
	Loop         bool   `json:"loop,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`     // 图生视频
	ImageEndURL  string `json:"image_end_url,omitempty"` // 图生视频
	Duration     string `json:"duration,omitempty"`      // 视频时长
	Resolution   string `json:"resolution,omitempty"`    // 视频分辨率
}

// LumaCreateResponse Luma 创建任务响应
type LumaCreateResponse struct {
	Id        string `json:"id"`
	Prompt    string `json:"prompt"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	Channel   string `json:"channel,omitempty"`
}

// LumaQueryResponse Luma 查询任务响应
type LumaQueryResponse struct {
	Id    string `json:"id"`
	State string `json:"state"`
	Video struct {
		URL         string `json:"url"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		Thumbnail   string `json:"thumbnail"`
		DownloadURL string `json:"download_url"`
	} `json:"video"`
	Prompt    string `json:"prompt"`
	Thumbnail struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnail"`
}

// CreateTask 创建视频生成任务
func (a *LumaAdapter) CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error) {
	// 解析任务参数
	paramsMap, ok := task.Params.(map[string]any)
	if !ok {
		return CreateTaskResponse{}, errors.New("invalid params type for Luma video task")
	}

	// 构建请求参数
	reqBody := LumaCreateRequest{
		UserPrompt: task.Prompt,
	}

	// 从 params 中提取参数
	if expandPrompt, ok := paramsMap["expand_prompt"].(bool); ok {
		reqBody.ExpandPrompt = expandPrompt
	}
	if model, ok := paramsMap["model"].(string); ok {
		reqBody.ModelName = model
	}
	if loop, ok := paramsMap["loop"].(bool); ok {
		reqBody.Loop = loop
	}
	if imageURL, ok := paramsMap["image_url"].(string); ok {
		reqBody.ImageURL = imageURL
	}
	if imageEndURL, ok := paramsMap["image_end_url"].(string); ok {
		reqBody.ImageEndURL = imageEndURL
	}
	if duration, ok := paramsMap["duration"].(string); ok {
		reqBody.Duration = duration
	}
	if resolution, ok := paramsMap["resolution"].(string); ok {
		reqBody.Resolution = resolution
	}
	// 发送请求
	apiURL := fmt.Sprintf("%s/luma/generations", videoConfig.ApiURL)
	r, err := a.httpClient.R().
		SetHeader("Authorization", "Bearer "+videoConfig.ApiKey).
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
	var res LumaCreateResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	return CreateTaskResponse{
		TaskId:    res.Id,
		Channel:   videoConfig.ApiURL,
		Prompt:    res.Prompt,
		State:     types.VideoStatusPending,
		CreatedAt: res.CreatedAt,
	}, nil
}

// QueryTask 查询任务状态
func (a *LumaAdapter) QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error) {
	apiURL := fmt.Sprintf("%s/luma/generations/%s", videoConfig.ApiURL, taskId)
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
	var res LumaQueryResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	switch res.State {
	case "completed", "succeed", "success":
		res.State = types.VideoStatusSuccess
	case "in_progress", "running":
		res.State = types.VideoStatusInProgress
	case "failed":
		res.State = types.VideoStatusFailed
	default:
		res.State = types.VideoStatusPending
	}

	// 构建响应
	response := QueryTaskResponse{
		TaskId:   res.Id,
		Status:   res.State,
		VideoURL: res.Video.DownloadURL,
		Prompt:   res.Prompt,
	}

	// 如果有原始数据，转换为 JSON 字符串
	if len(body) > 0 {
		response.Output = string(body)
	}

	return response, nil
}
