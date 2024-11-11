package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// TaskType 任务类别
type TaskType string

func (t TaskType) String() string {
	return string(t)
}

const (
	TaskImage     = TaskType("image")
	TaskBlend     = TaskType("blend")
	TaskSwapFace  = TaskType("swapFace")
	TaskUpscale   = TaskType("upscale")
	TaskVariation = TaskType("variation")
)

// MjTask MidJourney 任务
type MjTask struct {
	Id               uint     `json:"id"`      // 任务ID
	TaskId           string   `json:"task_id"` // 中转任务ID
	ClientId         string   `json:"client_id"`
	ImgArr           []string `json:"img_arr"`
	Type             TaskType `json:"type"`
	UserId           int      `json:"user_id"`
	Prompt           string   `json:"prompt,omitempty"`
	NegPrompt        string   `json:"neg_prompt,omitempty"`
	Params           string   `json:"full_prompt"`
	Index            int      `json:"index,omitempty"`
	MessageId        string   `json:"message_id,omitempty"`
	MessageHash      string   `json:"message_hash,omitempty"`
	ChannelId        string   `json:"channel_id"`         // 渠道ID，用来区分是哪个渠道创建的任务，一个任务的 create 和 action 操作必须要再同一个渠道
	Mode             string   `json:"mode"`               // 绘画模式，relax, fast, turbo
	TranslateModelId int      `json:"translate_model_id"` // 提示词翻译模型ID
}

type SdTask struct {
	Id               int          `json:"id"` // job 数据库ID
	Type             TaskType     `json:"type"`
	ClientId         string       `json:"client_id"`
	UserId           int          `json:"user_id"`
	Params           SdTaskParams `json:"params"`
	RetryCount       int          `json:"retry_count"`
	TranslateModelId int          `json:"translate_model_id"` // 提示词翻译模型ID
}

type SdTaskParams struct {
	ClientId     string  `json:"client_id"` // 客户端ID
	TaskId       string  `json:"task_id"`
	Prompt       string  `json:"prompt"`     // 提示词
	NegPrompt    string  `json:"neg_prompt"` // 反向提示词
	Steps        int     `json:"steps"`      // 迭代步数，默认20
	Sampler      string  `json:"sampler"`    // 采样器
	Scheduler    string  `json:"scheduler"`  // 采样调度
	FaceFix      bool    `json:"face_fix"`   // 面部修复
	CfgScale     float32 `json:"cfg_scale"`  //引导系数，默认 7
	Seed         int64   `json:"seed"`       // 随机数种子
	Height       int     `json:"height"`
	Width        int     `json:"width"`
	HdFix        bool    `json:"hd_fix"`         // 启用高清修复
	HdRedrawRate float32 `json:"hd_redraw_rate"` // 高清修复重绘幅度
	HdScale      int     `json:"hd_scale"`       // 放大倍数
	HdScaleAlg   string  `json:"hd_scale_alg"`   // 放大算法
	HdSteps      int     `json:"hd_steps"`       // 高清修复迭代步数
}

// DallTask DALL-E task
type DallTask struct {
	ClientId string `json:"client_id"`
	Id       uint   `json:"id"`
	UserId   uint   `json:"user_id"`
	Prompt   string `json:"prompt"`
	N        int    `json:"n"`
	Quality  string `json:"quality"`
	Size     string `json:"size"`
	Style    string `json:"style"`

	Power            int `json:"power"`
	TranslateModelId int `json:"translate_model_id"` // 提示词翻译模型ID
}

type SunoTask struct {
	ClientId     string `json:"client_id"`
	Id           uint   `json:"id"`
	Channel      string `json:"channel"`
	UserId       int    `json:"user_id"`
	Type         int    `json:"type"`
	Title        string `json:"title"`
	RefTaskId    string `json:"ref_task_id,omitempty"`
	RefSongId    string `json:"ref_song_id,omitempty"`
	Prompt       string `json:"prompt"` // 提示词/歌词
	Tags         string `json:"tags"`
	Model        string `json:"model"`
	Instrumental bool   `json:"instrumental"`          // 是否纯音乐
	ExtendSecs   int    `json:"extend_secs,omitempty"` // 延长秒杀
	SongId       string `json:"song_id,omitempty"`     // 合并歌曲ID
	AudioURL     string `json:"audio_url"`             // 用户上传音频地址
}

const (
	VideoLuma   = "luma"
	VideoRunway = "runway"
	VideoCog    = "cog"
)

type VideoTask struct {
	ClientId         string      `json:"client_id"`
	Id               uint        `json:"id"`
	Channel          string      `json:"channel"`
	UserId           int         `json:"user_id"`
	Type             string      `json:"type"`
	TaskId           string      `json:"task_id"`
	Prompt           string      `json:"prompt"` // 提示词
	Params           VideoParams `json:"params"`
	TranslateModelId int         `json:"translate_model_id"` // 提示词翻译模型ID
}

type VideoParams struct {
	PromptOptimize bool   `json:"prompt_optimize"` // 是否优化提示词
	Loop           bool   `json:"loop"`            // 是否循环参考图
	StartImgURL    string `json:"start_img_url"`   // 第一帧参考图地址
	EndImgURL      string `json:"end_img_url"`     // 最后一帧参考图地址
	Model          string `json:"model"`           // 使用哪个模型生成视频
	Radio          string `json:"radio"`           // 视频尺寸
	Style          string `json:"style"`           // 风格
	Duration       int    `json:"duration"`        // 视频时长（秒）
}
