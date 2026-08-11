package adapters

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"io"
	"mime/multipart"
	"net/http"
	"geekai/utils"
	"time"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

// SoraAdapter Sora 视频生成适配器
type SoraAdapter struct {
	db         *gorm.DB
	httpClient *req.Client
}

// NewSoraAdapter 创建 Sora 适配器
func NewSoraAdapter(db *gorm.DB) *SoraAdapter {
	return &SoraAdapter{
		db:         db,
		httpClient: req.C().SetTimeout(time.Minute * 3),
	}
}

// GetProvider 获取服务提供商名称
func (a *SoraAdapter) GetProvider() string {
	return "sora"
}

// SoraCreateRequest Sora 创建任务请求
type SoraCreateRequest struct {
	Model          string      `json:"model"`                     // 模型名称：sora-2, sora-2-pro
	Prompt         string      `json:"prompt"`                    // 提示词
	Size           string      `json:"size,omitempty"`            // 分辨率：1280x720, 720x1280, 1792x1024, 1024x1792
	InputReference interface{} `json:"input_reference,omitempty"` // 图生视频的参考图片，官方为对象 {"image_url": "..."}，也兼容字符串 URL
	Seconds        string      `json:"seconds,omitempty"`         // 视频时长（秒），默认4秒
	Watermark      bool        `json:"watermark,omitempty"`       // 是否添加水印
}

// SoraCreateResponse Sora 创建任务响应
type SoraCreateResponse struct {
	ID        string     `json:"id"`              // 任务ID
	Object    string     `json:"object"`          // 对象类型，固定为 "video"
	Model     string     `json:"model"`           // 模型名称
	Status    string     `json:"status"`          // 状态：queued, in_progress, completed, failed
	CreatedAt int64      `json:"created_at"`      // 创建时间戳
	Seconds   string     `json:"seconds"`         // 视频时长
	Size      string     `json:"size"`            // 分辨率
	Error     *SoraError `json:"error,omitempty"` // 错误信息（成功时为null）
}

// SoraQueryResponse Sora 查询任务响应
type SoraQueryResponse struct {
	ID        string     `json:"id"`              // 任务ID
	Object    string     `json:"object"`          // 对象类型，固定为 "video"
	Model     string     `json:"model"`           // 模型名称
	Status    string     `json:"status"`          // 状态：queued, in_progress, completed, failed
	Progress  int        `json:"progress"`        // 进度（0-100）
	CreatedAt int64      `json:"created_at"`      // 创建时间戳
	Seconds   string     `json:"seconds"`         // 视频时长
	Size      string     `json:"size"`            // 分辨率
	Error     *SoraError `json:"error,omitempty"` // 错误信息（成功时为null）
	VideoURL  string     `json:"video_url"`       // 视频URL（成功时生成）
}

type SoraError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateTask 创建视频生成任务
func (a *SoraAdapter) CreateTask(task types.VideoTask, videoConfig *types.VideoConfig) (CreateTaskResponse, error) {
	// 解析任务参数
	paramsMap, ok := task.Params.(map[string]any)
	if !ok {
		return CreateTaskResponse{}, fmt.Errorf("invalid params type for Sora video task")
	}

	// 是否调用官方 Sora 接口
	isOfficial := false
	if v, ok := paramsMap["is_official"].(bool); ok {
		isOfficial = v
	}

	// 提取通用参数
	model, ok := paramsMap["model"].(string)
	if !ok || model == "" {
		return CreateTaskResponse{}, fmt.Errorf("model 参数必填")
	}

	size, _ := paramsMap["size"].(string)

	seconds := "10" // 默认 10 秒
	if v, ok := paramsMap["seconds"].(string); ok && v != "" {
		seconds = v
	} else if duration, ok := paramsMap["duration"].(float64); ok {
		seconds = fmt.Sprintf("%.0f", duration)
	} else if duration, ok := paramsMap["duration"].(int); ok {
		seconds = fmt.Sprintf("%d", duration)
	}

	watermark := false
	if v, ok := paramsMap["watermark"].(bool); ok {
		watermark = v
	}

	// 处理图生视频（input_reference 参数）
	// 支持单个字符串或数组的第一个元素
	var imageURL string
	if inputRef, ok := paramsMap["input_reference"].(string); ok && inputRef != "" {
		imageURL = inputRef
	} else if images, ok := paramsMap["images"].([]interface{}); ok && len(images) > 0 {
		// 兼容旧的 images 参数格式
		if imgStr, ok := images[0].(string); ok && imgStr != "" {
			imageURL = imgStr
		}
	} else if image, ok := paramsMap["image"].(string); ok && image != "" {
		// 兼容 image 参数
		imageURL = image
	}

	// 官方 Sora：使用 multipart/form-data 携带文件
	if isOfficial && imageURL != "" {
		return a.createOfficialSoraTask(task, videoConfig, model, size, seconds, imageURL)
	}

	// 其他场景：保持原来的 JSON 调用，input_reference 继续传 URL 字符串
	reqBody := SoraCreateRequest{
		Model:    model,
		Prompt:   task.Prompt,
		Size:     size,
		Seconds:  seconds,
		Watermark: watermark,
	}

	if imageURL != "" {
		reqBody.InputReference = imageURL
	}

	// 发送 JSON 请求
	apiURL := fmt.Sprintf("%s/v1/videos", videoConfig.ApiURL)
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
	var res SoraCreateResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	// 转换状态：queued -> pending
	state := res.Status
	if state == "queued" || state == "in_progress" || state == "" {
		state = "pending"
	}

	return CreateTaskResponse{
		TaskId:    res.ID,
		Channel:   videoConfig.ApiURL,
		Prompt:    task.Prompt,
		State:     state,
		CreatedAt: time.Unix(res.CreatedAt, 0).Format(time.RFC3339),
	}, nil
}

// createOfficialSoraTask 调用官方 Sora API，使用 multipart/form-data 携带图片文件
func (a *SoraAdapter) createOfficialSoraTask(task types.VideoTask, videoConfig *types.VideoConfig, model, size, seconds, imageURL string) (CreateTaskResponse, error) {
	if videoConfig == nil || videoConfig.ApiURL == "" || videoConfig.ApiKey == "" {
		return CreateTaskResponse{}, fmt.Errorf("Sora 视频配置不完整")
	}

	imgData, err := downloadImageBytes(imageURL)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("下载参考图片失败：%v", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 文本字段
	if err = writer.WriteField("prompt", task.Prompt); err != nil {
		return CreateTaskResponse{}, err
	}
	if err = writer.WriteField("model", model); err != nil {
		return CreateTaskResponse{}, err
	}
	if size != "" {
		if err = writer.WriteField("size", size); err != nil {
			return CreateTaskResponse{}, err
		}
	}
	if seconds != "" {
		if err = writer.WriteField("seconds", seconds); err != nil {
			return CreateTaskResponse{}, err
		}
	}

	// 文件字段
	fileWriter, err := writer.CreateFormFile("input_reference", "image")
	if err != nil {
		return CreateTaskResponse{}, err
	}
	if _, err = fileWriter.Write(imgData); err != nil {
		return CreateTaskResponse{}, err
	}

	if err = writer.Close(); err != nil {
		return CreateTaskResponse{}, err
	}

	apiURL := fmt.Sprintf("%s/v1/videos", videoConfig.ApiURL)
	req, err := http.NewRequest(http.MethodPost, apiURL, &buf)
	if err != nil {
		return CreateTaskResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+videoConfig.ApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 3 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return CreateTaskResponse{}, fmt.Errorf("请求官方 Sora API 出错：%v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return CreateTaskResponse{}, fmt.Errorf("请求官方 Sora API 出错：%d, %s", resp.StatusCode, string(body))
	}

	var res SoraCreateResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return CreateTaskResponse{}, fmt.Errorf("解析官方 Sora API 数据失败：%v, %s", err, string(body))
	}

	state := res.Status
	if state == "queued" || state == "in_progress" || state == "" {
		state = "pending"
	}

	return CreateTaskResponse{
		TaskId:    res.ID,
		Channel:   videoConfig.ApiURL,
		Prompt:    task.Prompt,
		State:     state,
		CreatedAt: time.Unix(res.CreatedAt, 0).Format(time.RFC3339),
	}, nil
}

// downloadImageBytes 下载远程图片并返回二进制内容，用于 multipart 文件上传
func downloadImageBytes(imageURL string) ([]byte, error) {
	body, _, err := utils.FetchURLBytes(context.Background(), imageURL, "", 3*time.Minute, 2, 32<<20)
	return body, err
}

// downloadImageAsDataURL 下载远程图片并转为 data URL，避免向官方 Sora 直接传地址
// QueryTask 查询任务状态
func (a *SoraAdapter) QueryTask(taskId string, channel string, videoConfig *types.VideoConfig) (QueryTaskResponse, error) {
	apiURL := fmt.Sprintf("%s/v1/videos/%s", videoConfig.ApiURL, taskId)
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
	var res SoraQueryResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return QueryTaskResponse{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	// 转换状态：queued -> pending, completed -> success
	state := res.Status
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

	// 处理错误信息
	errMsg := ""
	if res.Error != nil {
		errMsg = res.Error.Message
	} else {
		errMsg = fmt.Sprintf("进度: %d%%", res.Progress)
	}

	// 构建响应
	response := QueryTaskResponse{
		TaskId:    res.ID,
		Status:    state,
		Progress:  res.Progress,
		VideoURL:  res.VideoURL,
		Prompt:    "", // Sora API 响应中不包含 prompt 字段
		ErrMsg:    errMsg,
		StatusMsg: fmt.Sprintf("进度: %d%%", res.Progress),
	}

	// 如果有原始数据，转换为 JSON 字符串
	if len(body) > 0 {
		response.Output = string(body)
	}

	return response, nil
}
