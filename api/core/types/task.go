package types

// TaskType 任务类别
type TaskType string

func (t TaskType) String() string {
	return string(t)
}

const (
	TaskImage     = TaskType("image")
	TaskUpscale   = TaskType("upscale")
	TaskVariation = TaskType("variation")
	TaskTxt2Img   = TaskType("text2img")
)

// TaskSrc 任务来源
type TaskSrc string

const (
	TaskSrcChat = TaskSrc("chat") // 来自聊天页面
	TaskSrcImg  = TaskSrc("img")  // 专业绘画页面
)

// MjTask MidJourney 任务
type MjTask struct {
	Id          int      `json:"id"`
	SessionId   string   `json:"session_id"`
	Src         TaskSrc  `json:"src"`
	Type        TaskType `json:"type"`
	UserId      int      `json:"user_id"`
	Prompt      string   `json:"prompt,omitempty"`
	ChatId      string   `json:"chat_id,omitempty"`
	RoleId      int      `json:"role_id,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Index       int32    `json:"index,omitempty"`
	MessageId   string   `json:"message_id,omitempty"`
	MessageHash string   `json:"message_hash,omitempty"`
	RetryCount  int      `json:"retry_count"`
}

// SdParams stable diffusion 绘画参数
type SdParams struct {
	TaskId         string  `json:"task_id"`
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt"`
	Steps          int     `json:"steps"`
	Sampler        string  `json:"sampler"`
	FaceFix        bool    `json:"face_fix"`
	CfgScale       float32 `json:"cfg_scale"`
	Seed           int64   `json:"seed"`
	Height         int     `json:"height"`
	Width          int     `json:"width"`
	HdFix          bool    `json:"hd_fix"`
	HdRedrawRate   float32 `json:"hd_redraw_rate"`
	HdScale        int     `json:"hd_scale"`
	HdScaleAlg     string  `json:"hd_scale_alg"`
	HdSampleNum    int     `json:"hd_sample_num"`
}

type SdTask struct {
	Id         int            `json:"id"`
	SessionId  string         `json:"session_id"`
	Src        types.TaskSrc  `json:"src"`
	Type       types.TaskType `json:"type"`
	UserId     int            `json:"user_id"`
	Prompt     string         `json:"prompt,omitempty"`
	Params     types.SdParams `json:"params"`
	RetryCount int            `json:"retry_count"`
}
