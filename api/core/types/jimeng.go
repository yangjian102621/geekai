package types

// JimengConfig 即梦AI配置
type JimengConfig struct {
	AccessKey string      `json:"access_key"`
	SecretKey string      `json:"secret_key"`
	Power     JimengPower `json:"power"`
}

// JimengPower 即梦AI算力配置
type JimengPower struct {
	Image          int `json:"image"`           // 图片生成算力，单位：积分/张
	Video          int `json:"video"`           // 视频生成算力，单位：积分/秒
	VirtualHuman   int `json:"virtual_human"`   // 数字人视频生成算力，单位：积分/秒
	ActionTransfer int `json:"action_transfer"` // 视频动作迁移算力，单位：积分/秒
}

// JMTaskStatus 任务状态
type JMTaskStatus string

const (
	JMTaskStatusInQueue    = JMTaskStatus("in_queue")   // 任务已提交
	JMTaskStatusGenerating = JMTaskStatus("generating") // 任务处理中
	JMTaskStatusDone       = JMTaskStatus("done")       // 处理完成
	JMTaskStatusNotFound   = JMTaskStatus("not_found")  // 任务未找到
	JMTaskStatusSuccess    = JMTaskStatus("success")    // 任务成功
	JMTaskStatusFailed     = JMTaskStatus("failed")     // 任务失败
	JMTaskStatusExpired    = JMTaskStatus("expired")    // 任务过期
)

// JMTaskType 任务类型
type JMTaskType string

const (
	JMTaskTypeImage          = JMTaskType("image")           // 文生图
	JMTaskTypeVideo          = JMTaskType("video")           // 图生图
	JMTaskTypeVirtualHuman   = JMTaskType("virtual_human")   // 图像编辑
	JMTaskTypeActionTransfer = JMTaskType("action_transfer") // 图像特效
)

// JimengTaskRequest 即梦AI任务请求
type JimengTaskRequest struct {
	ReqKey string `json:"req_key"` // 请求Key
	// 公共参数
	Prompt    string   `json:"prompt,omitempty"`
	ImageUrls []string `json:"image_urls,omitempty"`

	// 图片生成参数
	Size      string `json:"size,omitempty"`
	UsePreLLM bool   `json:"use_pre_llm,omitempty"`

	// 视频生成参数
	Duration       string `json:"duration,omitempty"`    // 视频时长
	TemplateId     string `json:"template_id,omitempty"` // 运镜模板ID
	AspectRatio    string `json:"aspect_ratio,omitempty"`
	CameraStrength string `json:"camera_strength,omitempty"` // 运镜强度

	// 数字人视频生成参数
	AudioURL string `json:"audio_url,omitempty"` // 音频URL

	// 视频动作迁移参数
	VideoURL string `json:"video_url,omitempty"` // 动作视频URL
}
