package model

import (
	"time"
)

type InviteLog struct {
	Id         uint `gorm:"primarykey;column:id"`
	InviterId  uint
	UserId     uint
	Username   string
	InviteCode string
	Remark     string
	CreatedAt  time.Time
}
