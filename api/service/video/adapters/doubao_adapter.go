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

// DoubaoAdapter 豆包 Seedance 视频生成适配器（通过 Kapon VolcArk 接入）
type DoubaoAdapter struct {
	db *gorm.DB
}

// NewDoubaoAdapter 创建 Doubao 适配器
func NewDoubaoAdapter(db *gorm.DB) *DoubaoAdapter {
	return &DoubaoAdapter{
		db: db,
	}
}

// GetProvider 获取服务提供商名称
func (a *DoubaoAdapter) GetProvider() string {
	return types.VideoDoubao
}

// doubaoContentItem 请求体中的 content 子项
type doubaoContentItem struct {
	Type     string                 `json:"type"`
	Text     string                 `json:"text,omitempty"`
	ImageURL map[string]string      `json:"image_url,omitempty"`
	Extra    map[string]interface{} `json:"extra,omitempty"` // 预留扩展
}

// doubaoCreateRequest 创建任务请求体
type doubaoCreateRequest struct {
	Model      string              `json:"model"`
	Content    []doubaoContentItem `json:"content"`
	Duration   int                 `json:"duration,omitempty"`
	Frames     int                 `json:"frames,omitempty"`
	Ratio      string              `json:"ratio,omitempty"`
	Resolution string              `json:"resolution,omitempty"`
	Seed       int64               `json:"seed,omitempty"`
	// 其他官方支持的字段，按需追加
}

// doubaoCreateResponse 创建任务响应
type doubaoCreateResponse struct {
	Id         string `json:"id"`
	PlatformId string `json:"platform_id"`
	// 其余字段目前用不到，先不展开
}

// doubaoQueryContent 查询任务 content 字段
type doubaoQueryContent struct {
	VideoURL     string `json:"video_url"`
	LastFrameURL string `json:"last_frame_url"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// doubaoQueryUsage 查询任务 usage 字段
type doubaoQueryUsage struct {
	VideoTokens int `json:"video_tokens"`
}

// doubaoQueryResponse 查询任务响应
type doubaoQueryResponse struct {
	Id         string             `json:"id"`
	PlatformId string             `json:"platform_id"`
	Model      string             `json:"model"`
	Status     string             `json:"status"`
	Content    doubaoQueryContent `json:"content"`
	Duration   int                `json:"duration"`
	Frames     int                `json:"framespersecond"`
	Usage      doubaoQueryUsage   `json:"usage"`
	Error      string             `json:"error,omitempty"`
}

// CreateTask 创建豆包 Seedance 视频任务
func (a *DoubaoAdapter) CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error) {
	if videoConfig == nil {
		return CreateTaskResponse{}, errors.New("视频配置为空")
	}
	if videoConfig.ApiURL == "" || videoConfig.ApiKey == "" {
		return CreateTaskResponse{}, errors.New("豆包视频未配置 ApiURL 或 ApiKey")
	}

	paramsMap, ok := task.Params.(map[string]interface{})
	if !ok {
		return CreateTaskResponse{}, errors.New("invalid params type for Doubao video task")
	}

	// 模型名称：优先从 params.model 读取，否则默认 doubao-seedance-1-5-pro
	modelName := "doubao-seedance-1-5-pro"
	if v, ok := paramsMap["model"].(string); ok && v != "" {
		modelName = v
	}

	// 解析基础参数
	duration := 0
	if v, ok := paramsMap["duration"].(float64); ok {
		duration = int(v)
	}
	if v, ok := paramsMap["duration"].(int); ok {
		duration = v
	}

	ratio := ""
	if v, ok := paramsMap["aspect_ratio"].(string); ok {
		ratio = v
	}

	resolution := ""
	if v, ok := paramsMap["resolution"].(string); ok {
		resolution = v
	}

	var seed int64
	switch v := paramsMap["seed"].(type) {
	case float64:
		seed = int64(v)
	case int:
		seed = int64(v)
	case int64:
		seed = v
	}

	// 构建 content 数组：文本提示词为必填
	content := []doubaoContentItem{
		{
			Type: "text",
			Text: task.Prompt,
		},
	}

	// 如果存在 input_reference（图片 URL），则追加 image_url 项，用于 I2V
	if ref, ok := paramsMap["input_reference"].(string); ok && ref != "" {
		content = append(content, doubaoContentItem{
			Type: "image_url",
			ImageURL: map[string]string{
				"url": ref,
			},
		})
	}

	reqBody := doubaoCreateRequest{
		Model:      modelName,
		Content:    content,
		Duration:   duration,
		Ratio:      ratio,
		Resolution: resolution,
		Seed:       seed,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("序列化豆包请求失败: %v", err)
	}
	logger.Debugf("DoubaoCreateRequest: %s", string(payload))

	url := fmt.Sprintf("%s/seedance/v3/contents/generations/tasks", videoConfig.ApiURL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("创建豆包请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+videoConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("调用豆包接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("读取豆包响应失败: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return CreateTaskResponse{}, fmt.Errorf("豆包接口返回错误状态码: %d, %s", resp.StatusCode, string(body))
	}

	var apiResp doubaoCreateResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return CreateTaskResponse{}, fmt.Errorf("解析豆包创建任务响应失败: %v, body=%s", err, string(body))
	}

	taskId := apiResp.PlatformId
	if taskId == "" {
		taskId = apiResp.Id
	}
	if taskId == "" {
		return CreateTaskResponse{}, fmt.Errorf("豆包创建任务响应缺少任务 ID, body=%s", string(body))
	}

	return CreateTaskResponse{
		TaskId:    taskId,
		Channel:   videoConfig.ApiURL,
		Prompt:    task.Prompt,
		State:     types.VideoStatusPending,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// QueryTask 查询豆包 Seedance 视频任务状态
func (a *DoubaoAdapter) QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error) {
	if videoConfig == nil {
		return QueryTaskResponse{}, errors.New("视频配置为空")
	}
	if videoConfig.ApiURL == "" || videoConfig.ApiKey == "" {
		return QueryTaskResponse{}, errors.New("豆包视频未配置 ApiURL 或 ApiKey")
	}

	url := fmt.Sprintf("%s/seedance/v3/contents/generations/tasks/%s", videoConfig.ApiURL, taskId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("创建查询请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+videoConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("调用豆包查询接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("读取豆包查询响应失败: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return QueryTaskResponse{}, fmt.Errorf("豆包查询接口返回错误状态码: %d, %s", resp.StatusCode, string(body))
	}

	var apiResp doubaoQueryResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return QueryTaskResponse{}, fmt.Errorf("解析豆包查询任务响应失败: %v, body=%s", err, string(body))
	}

	status := apiResp.Status
	progress := 0

	switch status {
	case "queued":
		status = types.VideoStatusPending
		progress = 10
	case "running":
		status = types.VideoStatusInProgress
		progress = 60
	case "succeeded":
		status = types.VideoStatusSuccess
		progress = 100
	case "failed", "cancelled":
		status = types.VideoStatusFailed
	default:
		// 保持原样或视为 pending
		status = types.VideoStatusPending
	}

	errMsg := apiResp.Error
	if errMsg == "" && status == types.VideoStatusFailed {
		errMsg = "doubao task failed"
	}

	result := QueryTaskResponse{
		TaskId:    apiResp.PlatformId,
		Status:    status,
		Progress:  progress,
		VideoURL:  apiResp.Content.VideoURL,
		Prompt:    "",
		ErrMsg:    errMsg,
		StatusMsg: status,
		Output:    string(body),
	}

	if result.TaskId == "" {
		result.TaskId = apiResp.Id
	}

	return result, nil
}
