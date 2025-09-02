package ai3d

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"time"

	"github.com/imroc/req/v3"
)

type Gitee3DClient struct {
	httpClient *req.Client
	config     types.Gitee3DConfig
	apiURL     string
}

type Gitee3DParams struct {
	Prompt       string `json:"prompt"`        // 文本提示词
	ImageURL     string `json:"image_url"`     // 输入图片URL
	ResultFormat string `json:"result_format"` // 输出格式
}

type Gitee3DResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

type Gitee3DQueryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status     string `json:"status"`
		Progress   int    `json:"progress"`
		ResultURL  string `json:"result_url"`
		PreviewURL string `json:"preview_url"`
		ErrorMsg   string `json:"error_msg"`
	} `json:"data"`
}

func NewGitee3DClient(sysConfig *types.SystemConfig) *Gitee3DClient {
	return &Gitee3DClient{
		httpClient: req.C().SetTimeout(time.Minute * 3),
		config:     sysConfig.AI3D.Gitee,
		apiURL:     "https://ai.gitee.com/v1/async/image-to-3d",
	}
}

func (c *Gitee3DClient) UpdateConfig(config types.Gitee3DConfig) {
	c.config = config
}

// SubmitJob 提交3D生成任务
func (c *Gitee3DClient) SubmitJob(params Gitee3DParams) (string, error) {
	requestBody := map[string]any{
		"prompt":        params.Prompt,
		"image_url":     params.ImageURL,
		"result_format": params.ResultFormat,
	}

	response, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.config.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(c.apiURL + "/async/image-to-3d")

	if err != nil {
		return "", fmt.Errorf("failed to submit gitee 3D job: %v", err)
	}

	var giteeResp Gitee3DResponse
	if err := json.Unmarshal(response.Bytes(), &giteeResp); err != nil {
		return "", fmt.Errorf("failed to parse gitee response: %v", err)
	}

	if giteeResp.Code != 0 {
		return "", fmt.Errorf("gitee API error: %s", giteeResp.Message)
	}

	if giteeResp.Data.TaskID == "" {
		return "", fmt.Errorf("no task ID returned from gitee 3D API")
	}

	return giteeResp.Data.TaskID, nil
}

// QueryJob 查询任务状态
func (c *Gitee3DClient) QueryJob(taskId string) (*types.AI3DJobResult, error) {
	response, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.config.APIKey).
		Get(fmt.Sprintf("%s/task/%s/get", c.apiURL, taskId))

	if err != nil {
		return nil, fmt.Errorf("failed to query gitee 3D job: %v", err)
	}

	var giteeResp Gitee3DQueryResponse
	if err := json.Unmarshal(response.Bytes(), &giteeResp); err != nil {
		return nil, fmt.Errorf("failed to parse gitee query response: %v", err)
	}

	if giteeResp.Code != 0 {
		return nil, fmt.Errorf("gitee API error: %s", giteeResp.Message)
	}

	result := &types.AI3DJobResult{
		JobId:    taskId,
		Status:   c.convertStatus(giteeResp.Data.Status),
		Progress: giteeResp.Data.Progress,
	}

	// 根据状态设置结果
	switch giteeResp.Data.Status {
	case "completed":
		result.FileURL = giteeResp.Data.ResultURL
		result.PreviewURL = giteeResp.Data.PreviewURL
	case "failed":
		result.ErrorMsg = giteeResp.Data.ErrorMsg
	}

	return result, nil
}

// convertStatus 转换Gitee状态到系统状态
func (c *Gitee3DClient) convertStatus(giteeStatus string) string {
	switch giteeStatus {
	case "pending":
		return types.AI3DJobStatusPending
	case "processing":
		return types.AI3DJobStatusProcessing
	case "completed":
		return types.AI3DJobStatusCompleted
	case "failed":
		return types.AI3DJobStatusFailed
	default:
		return types.AI3DJobStatusPending
	}
}

// GetSupportedModels 获取支持的模型列表
func (c *Gitee3DClient) GetSupportedModels() []types.AI3DModel {
	return []types.AI3DModel{
		{Name: "Hunyuan3D-2", Power: 100, Formats: []string{"GLB"}, Desc: "Hunyuan3D-2 是腾讯混元团队推出的高质量 3D 生成模型，具备高保真度、细节丰富和高效生成的特点，可快速将文本或图像转换为逼真的 3D 物体。"},
		{Name: "Step1X-3D", Power: 55, Formats: []string{"GLB", "STL"}, Desc: "Step1X-3D 是一款由阶跃星辰（StepFun）与光影焕像（LightIllusions）联合研发并开源的高保真 3D 生成模型，专为高质量、可控的 3D 内容创作而设计。"},
		{Name: "Hi3DGen", Power: 35, Formats: []string{"GLB", "STL"}, Desc: "Hi3DGen 是一个 AI 工具，它可以把你上传的普通图片，智能转换成有“立体感”的图片（法线图），常用于制作 3D 效果，比如游戏建模、虚拟现实、动画制作等。"},
	}
}
