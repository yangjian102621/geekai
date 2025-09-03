package ai3d

import (
	"fmt"
	"geekai/core/types"

	tencent3d "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ai3d/v20250513"
	tencentcloud "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type Tencent3DClient struct {
	client *tencent3d.Client
	config types.Tencent3DConfig
}

type Tencent3DParams struct {
	Prompt          string      `json:"prompt"`                      // 文本提示词
	ImageURL        string      `json:"image_url"`                   // 输入图片URL
	ResultFormat    string      `json:"result_format"`               // 输出格式
	EnablePBR       bool        `json:"enable_pbr"`                  // 是否开启PBR材质
	MultiViewImages []ViewImage `json:"multi_view_images,omitempty"` // 多视角图片
}

type ViewImage struct {
	ViewType     string `json:"view_type"`      // 视角类型 (left/right/back)
	ViewImageURL string `json:"view_image_url"` // 图片URL
}

func NewTencent3DClient(sysConfig *types.SystemConfig) (*Tencent3DClient, error) {
	config := sysConfig.AI3D.Tencent
	credential := tencentcloud.NewCredential(config.SecretId, config.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ai3d.tencentcloudapi.com"

	client, err := tencent3d.NewClient(credential, config.Region, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create tencent 3D client: %v", err)
	}

	return &Tencent3DClient{
		client: client,
		config: config,
	}, nil
}

func (c *Tencent3DClient) UpdateConfig(config types.Tencent3DConfig) error {
	c.config = config
	credential := tencentcloud.NewCredential(config.SecretId, config.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ai3d.tencentcloudapi.com"

	client, err := tencent3d.NewClient(credential, config.Region, cpf)
	if err != nil {
		return fmt.Errorf("failed to create tencent 3D client: %v", err)
	}
	c.client = client
	return nil
}

// SubmitJob 提交3D生成任务
func (c *Tencent3DClient) SubmitJob(params Tencent3DParams) (string, error) {
	request := tencent3d.NewSubmitHunyuanTo3DJobRequest()

	if params.Prompt != "" {
		request.Prompt = tencentcloud.StringPtr(params.Prompt)
	}

	if params.ImageURL != "" {
		request.ImageUrl = tencentcloud.StringPtr(params.ImageURL)
	}

	if params.ResultFormat != "" {
		request.ResultFormat = tencentcloud.StringPtr(params.ResultFormat)
	}

	request.EnablePBR = tencentcloud.BoolPtr(params.EnablePBR)

	if len(params.MultiViewImages) > 0 {
		var viewImages []*tencent3d.ViewImage
		for _, img := range params.MultiViewImages {
			viewImage := &tencent3d.ViewImage{
				ViewType:     tencentcloud.StringPtr(img.ViewType),
				ViewImageUrl: tencentcloud.StringPtr(img.ViewImageURL),
			}
			viewImages = append(viewImages, viewImage)
		}
		request.MultiViewImages = viewImages
	}

	response, err := c.client.SubmitHunyuanTo3DJob(request)
	if err != nil {
		return "", fmt.Errorf("failed to submit tencent 3D job: %v", err)
	}

	if response.Response.JobId == nil {
		return "", fmt.Errorf("no job ID returned from tencent 3D API")
	}

	return *response.Response.JobId, nil
}

// QueryJob 查询任务状态
func (c *Tencent3DClient) QueryJob(jobId string) (*types.AI3DJobResult, error) {
	request := tencent3d.NewQueryHunyuanTo3DJobRequest()
	request.JobId = tencentcloud.StringPtr(jobId)

	response, err := c.client.QueryHunyuanTo3DJob(request)
	if err != nil {
		return nil, fmt.Errorf("failed to query tencent 3D job: %v", err)
	}

	result := &types.AI3DJobResult{
		JobId:    jobId,
		Status:   *response.Response.Status,
		Progress: 0,
	}

	// 根据状态设置进度
	switch *response.Response.Status {
	case "WAIT":
		result.Status = "pending"
		result.Progress = 10
	case "RUN":
		result.Status = "processing"
		result.Progress = 50
	case "DONE":
		result.Status = "completed"
		result.Progress = 100
		// 处理结果文件
		if len(response.Response.ResultFile3Ds) > 0 {
			for _, file := range response.Response.ResultFile3Ds {
				if file.Url != nil {
					result.FileURL = *file.Url
				}
				if file.PreviewImageUrl != nil {
					result.PreviewURL = *file.PreviewImageUrl
				}
				// TODO 取第一个文件
			}
		}
	case "FAIL":
		result.Status = "failed"
		result.Progress = 0
		if response.Response.ErrorMessage != nil {
			result.ErrorMsg = *response.Response.ErrorMessage
		}
	}

	return result, nil
}

// GetSupportedModels 获取支持的模型列表
func (c *Tencent3DClient) GetSupportedModels() []types.AI3DModel {
	return []types.AI3DModel{
		{Name: "Hunyuan3D-3", Power: 500, Formats: []string{"GLB", "OBJ", "STL", "USDZ", "FBX", "MP4"}, Desc: "Hunyuan3D 是腾讯混元团队推出的高质量 3D 生成模型，具备高保真度、细节丰富和高效生成的特点，可快速将文本或图像转换为逼真的 3D 物体。"},
	}
}
