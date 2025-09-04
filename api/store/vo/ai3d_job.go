package vo

import "geekai/core/types"

type AI3DJob struct {
	Id         uint               `json:"id"`
	UserId     uint               `json:"user_id"`
	Type       types.AI3DTaskType `json:"type"`
	Power      int                `json:"power"`
	TaskId     string             `json:"task_id"`
	FileURL    string             `json:"file_url"`
	PreviewURL string             `json:"preview_url"`
	Model      string             `json:"model"`
	Status     string             `json:"status"`
	ErrMsg     string             `json:"err_msg"`
	Params     AI3DJobParams      `json:"params"`
	CreatedAt  int64              `json:"created_at"`
	UpdatedAt  int64              `json:"updated_at"`
}

// AI3DJobParams 创建3D任务请求
type AI3DJobParams struct {
	// 通用参数
	JobId      uint               `json:"job_id,omitempty"`      // 任务ID
	Type       types.AI3DTaskType `json:"type,omitempty"`        // API类型 (tencent/gitee)
	Model      string             `json:"model,omitempty"`       // 3D模型类型
	Prompt     string             `json:"prompt,omitempty"`      // 文本提示词
	ImageURL   string             `json:"image_url,omitempty"`   // 输入图片URL
	FileFormat string             `json:"file_format,omitempty"` // 输出文件格式
	Power      int                `json:"power,omitempty"`       // 消耗算力
	// 腾讯3d专有参数
	EnablePBR bool `json:"enable_pbr,omitempty"` // 是否开启PBR材质
	// Gitee3d专有参数
	Texture           bool    `json:"texture,omitempty"`             // 是否开启纹理
	Seed              int     `json:"seed,omitempty"`                // 随机种子
	NumInferenceSteps int     `json:"num_inference_steps,omitempty"` //迭代次数
	GuidanceScale     float64 `json:"guidance_scale,omitempty"`      //引导系数
	OctreeResolution  int     `json:"octree_resolution"`             // 3D 渲染精度，越高3D 细节越丰富
}
