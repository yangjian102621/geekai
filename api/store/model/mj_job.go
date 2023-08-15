package model

import "time"

type MidJourneyJob struct {
	Id        uint `gorm:"primarykey;column:id"`
	UserId    uint
	ChatId    string
	MessageId string
	Hash      string
	Content   string
	Prompt    string
	Image     string
	CreatedAt time.Time
}

func (MidJourneyJob) TableName() string {
	return "chatgpt_mj_jobs"
}
