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
	Publish   bool   //是否发布图片到画廊
	ErrMsg    string // 报错信息
	Power     int    // 消耗算力
	CreatedAt time.Time
}

func (SdJob) TableName() string {
	return "chatgpt_sd_jobs"
}
