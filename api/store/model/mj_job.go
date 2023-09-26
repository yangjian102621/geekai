package model

import "time"

type MidJourneyJob struct {
	Id          uint `gorm:"primarykey;column:id"`
	Type        string
	UserId      int
	MessageId   string
	ReferenceId string
	ImgURL      string
	Hash        string // message hash
	Progress    int
	Prompt      string
	Started     bool
	CreatedAt   time.Time
}

func (MidJourneyJob) TableName() string {
	return "chatgpt_mj_jobs"
}
