package model

import "time"

type MidJourneyJob struct {
	Id          uint `gorm:"primarykey;column:id"`
	Type        string
	UserId      int
	TaskId      string
	ChannelId   string
	MessageId   string
	ReferenceId string
	ImgURL      string
	OrgURL      string // 原图地址
	Hash        string // message hash
	Progress    int
	Prompt      string
	UseProxy    bool // 是否使用反代加载图片
	CreatedAt   time.Time
}

func (MidJourneyJob) TableName() string {
	return "chatgpt_mj_jobs"
}
