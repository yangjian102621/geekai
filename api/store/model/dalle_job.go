package model

import "time"

type DallJob struct {
	Id        uint `gorm:"primarykey;column:id"`
	UserId    int
	TaskId      string
	Prompt    string
	ImgURL       string
	Publish       bool
	Power      int
	Progress      int
	ErrMsg      string
	CreatedAt time.Time
}
