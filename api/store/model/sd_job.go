package model

import "time"

type SdJob struct {
	Id        uint `gorm:"primarykey;column:id"`
	Type      string
	UserId    int
	TaskId    string
	ImgURL    string
	Progress  int
	Prompt    string
	Params    string
	Publish   bool //是否发布图片到画廊
	CreatedAt time.Time
}

func (SdJob) TableName() string {
	return "chatgpt_sd_jobs"
}
