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
	Reward     string `gorm:"column:reward_json"` // 邀请奖励
	CreatedAt  time.Time
}
