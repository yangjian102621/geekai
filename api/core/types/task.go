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
	TaskModal     = TaskType("modal") // 局部重绘
)

// MjTask MidJourney 任务
type MjTask struct {
	Id               uint     `json:"id"`      // 任务ID
	TaskId           string   `json:"task_id"` // 中转任务ID
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
	MaskBase64       string   `json:"mask_base64,omitempty"` // 局部重绘蒙版 base64（仅 TaskModal 使用）
}

// ImageTask Image generation task
type ImageTask struct {
	ModelId          uint     `json:"model_id"`
	ModelName        string   `json:"model_name"`
	ModelValue       string   `json:"model_value"`
	Image            []string `json:"image,omitempty"`
	Id               uint     `json:"id"`
	TaskId           string   `json:"task_id"` // Kapon 异步任务 ID，与 API 返回的 task_id 一致
	UserId           uint     `json:"user_id"`
	Prompt           string   `json:"prompt"`
	AspectRatio      string   `json:"aspect_ratio"`
	Size             string   `json:"size"`
	Power            int      `json:"power"`
	TranslateModelId int      `json:"translate_model_id"` // 提示词翻译模型ID
}

type SunoTask struct {
	Id           uint   `json:"id"`
	Channel      string `json:"channel"`
	UserId       int    `json:"user_id"`
	Type         int    `json:"type"`
	Title        string `json:"title"`
	RefTaskId    string `json:"ref_task_id,omitempty"`
	RefSongId    string `json:"ref_song_id,omitempty"`
	Prompt       string `json:"prompt"`           // 提示词
	Lyrics       string `json:"lyrics,omitempty"` // 歌词
	Tags         string `json:"tags"`
	Model        string `json:"model"`
	Instrumental bool   `json:"instrumental"`          // 是否纯音乐
	ExtendSecs   int    `json:"extend_secs,omitempty"` // 延长秒杀
	SongId       string `json:"song_id,omitempty"`     // 合并歌曲ID
	AudioURL     string `json:"audio_url"`             // 用户上传音频地址
}

const (
	VideoLuma    = "luma"
	VideoSora    = "sora"
	VideoVeo     = "veo"
	VideoKeLing  = "keling"
	VideoMiniMax = "minimax"
	VideoWan     = "wan"
	VideoDoubao  = "doubao"
)

type VideoTask struct {
	Id               uint   `json:"id"`
	Channel          string `json:"channel"`
	UserId           int    `json:"user_id"`
	Type             string `json:"type"` // provider（不带版本号：veo, sora, luma）
	TaskId           string `json:"task_id"`
	Prompt           string `json:"prompt"` // 提示词
	Params           any    `json:"params"`
	TranslateModelId int    `json:"translate_model_id"` // 提示词翻译模型ID
}
