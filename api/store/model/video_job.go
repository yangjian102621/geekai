package model

import "time"

type VideoJob struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	Channel   string    `gorm:"column:channel;type:varchar(100);not null;comment:渠道" json:"channel"`
	TaskId    string    `gorm:"column:task_id;type:varchar(100);not null;comment:任务 ID" json:"task_id"`
	TaskInfo  string    `gorm:"column:task_info;type:text;comment:原始任务信息" json:"task_info"`
	Type      string    `gorm:"column:type;type:varchar(20);comment:任务类型,luma,runway,cogvideo" json:"type"`
	Prompt    string    `gorm:"column:prompt;type:text;not null;comment:提示词" json:"prompt"`
	PromptExt string    `gorm:"column:prompt_ext;type:text;comment:优化后提示词" json:"prompt_ext"`
	CoverURL  string    `gorm:"column:cover_url;type:varchar(512);comment:封面图地址" json:"cover_url"`
	VideoURL  string    `gorm:"column:video_url;type:varchar(512);comment:视频地址" json:"video_url"`
	WaterURL  string    `gorm:"column:water_url;type:varchar(512);comment:带水印的视频地址" json:"water_url"`
	Progress  int       `gorm:"column:progress;type:smallint;default:0;comment:任务进度" json:"progress"`
	Publish   int       `gorm:"column:publish;type:tinyint(1);not null;comment:是否发布" json:"publish"`
	ErrMsg    string    `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	RawData   string    `gorm:"column:raw_data;type:text;comment:原始数据" json:"raw_data"`
	Power     int       `gorm:"column:power;type:smallint;not null;default:0;comment:消耗算力" json:"power"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *VideoJob) TableName() string {
	return "chatgpt_video_jobs"
}
