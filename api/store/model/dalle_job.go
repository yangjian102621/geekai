package model

import "time"

type DallJob struct {
	Id        uint `gorm:"primarykey;column:id"`
	UserId    uint
	Prompt    string
	ImgURL    string
	OrgURL    string
	Publish   bool
	Power     int
	Progress  int
	ErrMsg    string
	CreatedAt time.Time
}
