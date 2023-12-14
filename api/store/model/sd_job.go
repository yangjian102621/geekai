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
	CreatedAt time.Time
}

func (SdJob) TableName() string {
	return "chatgpt_sd_jobs"
}
