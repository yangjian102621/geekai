package types

// AI3DConfig 3D生成配置
type AI3DConfig struct {
	Tencent Tencent3DConfig `json:"tencent,omitempty"`
	Gitee   Gitee3DConfig   `json:"gitee,omitempty"`
}

// Tencent3DConfig 腾讯云3D配置
type Tencent3DConfig struct {
	SecretId  string      `json:"secret_id,omitempty"`
	SecretKey string      `json:"secret_key,omitempty"`
	Region    string      `json:"region,omitempty"`
	Enabled   bool        `json:"enabled,omitempty"`
	Models    []AI3DModel `json:"models,omitempty"`
}

// Gitee3DConfig Gitee 3D配置
type Gitee3DConfig struct {
	APIKey  string      `json:"api_key,omitempty"`
	Enabled bool        `json:"enabled,omitempty"`
	Models  []AI3DModel `json:"models,omitempty"`
}

type AI3DTaskType string

const (
	AI3DTaskTypeTencent AI3DTaskType = "tencent"
	AI3DTaskTypeGitee   AI3DTaskType = "gitee"
)

// AI3DJobResult 3D任务结果
type AI3DJobResult struct {
	TaskId     string `json:"task_id"`     // 任务ID
	Status     string `json:"status"`      // 任务状态
	FileURL    string `json:"file_url"`    // 3D模型文件URL
	PreviewURL string `json:"preview_url"` // 预览图片URL
	ErrorMsg   string `json:"error_msg"`   // 错误信息
	RawData    string `json:"raw_data"`    // 原始数据
}

// AI3DModel 3D模型配置
type AI3DModel struct {
	Name    string   `json:"name"`    // 模型名称
	Desc    string   `json:"desc"`    // 模型描述
	Power   int      `json:"power"`   // 算力消耗
	Formats []string `json:"formats"` // 支持输出的文件格式
}

// AI3DJobRequest 3D任务请求
type AI3DJobRequest struct {
	Type     string `json:"type"`      // API类型 (tencent/gitee)
	Model    string `json:"model"`     // 3D模型类型
	Prompt   string `json:"prompt"`    // 文本提示词
	ImageURL string `json:"image_url"` // 输入图片URL
	Power    int    `json:"power"`     // 消耗算力
}

// AI3DJobStatus 3D任务状态
const (
	AI3DJobStatusPending    = "pending"    // 等待中
	AI3DJobStatusProcessing = "processing" // 处理中
	AI3DJobStatusSuccess    = "success"    // 已完成
	AI3DJobStatusFailed     = "failed"     // 失败
)
