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
	Model             string  `json:"model"`                         // 模型名称
	FileFormat        string  `json:"file_format,omitempty"`         // 文件格式(Step1X-3D、Hi3DGen模型适用)，支持 glb 和 stl
	Type              string  `json:"type,omitempty"`                // 输出格式(Hunyuan3D-2模型适用)
	ImageURL          string  `json:"image_url"`                     // 输入图片URL
	Texture           bool    `json:"texture,omitempty"`             // 是否开启纹理
	Seed              int     `json:"seed,omitempty"`                // 随机种子
	NumInferenceSteps int     `json:"num_inference_steps,omitempty"` //迭代次数
	GuidanceScale     float64 `json:"guidance_scale,omitempty"`      //引导系数
	OctreeResolution  int     `json:"octree_resolution,omitempty"`   // 3D 渲染精度，越高3D 细节越丰富
}

type Gitee3DResponse struct {
	TaskID string `json:"task_id"`
	Output struct {
		FileURL    string `json:"file_url,omitempty"`
		PreviewURL string `json:"preview_url,omitempty"`
	} `json:"output"`
	Status      string `json:"status"`
	CreatedAt   any    `json:"created_at"`
	StartedAt   any    `json:"started_at"`
	CompletedAt any    `json:"completed_at"`
	Urls        struct {
		Get    string `json:"get"`
		Cancel string `json:"cancel"`
	} `json:"urls"`
}

type GiteeErrorResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

func NewGitee3DClient(sysConfig *types.SystemConfig) *Gitee3DClient {
	return &Gitee3DClient{
		httpClient: req.C().SetTimeout(time.Minute * 3),
		config:     sysConfig.AI3D.Gitee,
		apiURL:     "https://ai.gitee.com/v1",
	}
}

func (c *Gitee3DClient) UpdateConfig(config types.Gitee3DConfig) {
	c.config = config
}

func (c *Gitee3DClient) GetConfig() *types.Gitee3DConfig {
	return &c.config
}

// SubmitJob 提交3D生成任务
func (c *Gitee3DClient) SubmitJob(params Gitee3DParams) (string, error) {

	var giteeResp Gitee3DResponse
	response, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.config.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetSuccessResult(&giteeResp).
		Post(c.apiURL + "/async/image-to-3d")

	if err != nil {
		return "", fmt.Errorf("failed to submit gitee 3D job: %v", err)
	}

	if giteeResp.TaskID == "" {
		var giteeErr GiteeErrorResponse
		_ = json.Unmarshal(response.Bytes(), &giteeErr)
		return "", fmt.Errorf("no task ID returned from gitee 3D API: %s", giteeErr.Message)
	}

	return giteeResp.TaskID, nil
}

// QueryJob 查询任务状态
func (c *Gitee3DClient) QueryJob(taskId string) (*types.AI3DJobResult, error) {
	var giteeResp Gitee3DResponse
	apiURL := fmt.Sprintf("%s/task/%s", c.apiURL, taskId)
	response, err := c.httpClient.R().
		SetHeader("Authorization", "Bearer "+c.config.APIKey).
		SetSuccessResult(&giteeResp).
		Get(apiURL)

	if err != nil {
		return nil, fmt.Errorf("failed to query gitee 3D job: %v", err)
	}

	result := &types.AI3DJobResult{
		TaskId: taskId,
		Status: c.convertStatus(giteeResp.Status),
	}

	if giteeResp.TaskID == "" {
		var giteeErr GiteeErrorResponse
		_ = json.Unmarshal(response.Bytes(), &giteeErr)
		result.ErrorMsg = giteeErr.Message
	} else if giteeResp.Status == "success" {
		result.FileURL = giteeResp.Output.FileURL
	}
	result.RawData = string(response.Bytes())

	logger.Debugf("gitee 3D job response: %+v", result)

	return result, nil
}

// convertStatus 转换Gitee状态到系统状态
func (c *Gitee3DClient) convertStatus(giteeStatus string) string {
	switch giteeStatus {
	case "waiting":
		return types.AI3DJobStatusPending
	case "in_progress":
		return types.AI3DJobStatusProcessing
	case "success":
		return types.AI3DJobStatusSuccess
	case "failure", "cancelled":
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
