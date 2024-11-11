package model

import "time"

type VideoJob struct {
	Id        uint `gorm:"primarykey;column:id"`
	UserId    int
	Channel   string // 频道
	Type      string // luma,runway,cog
	TaskId    string
	TaskInfo  string // 原始任务信息
	Prompt    string // 提示词
	PromptExt string // 优化后提示词
	CoverURL  string // 封面图 URL
	VideoURL  string // 无水印视频 URL
	WaterURL  string // 有水印视频 URL
	Progress  int    // 任务进度
	Publish   bool   // 是否发布
	ErrMsg    string // 错误信息
	RawData   string // 原始数据 json
	Power     int    // 消耗算力
	CreatedAt time.Time
}

func (VideoJob) TableName() string {
	return "chatgpt_video_jobs"
}
