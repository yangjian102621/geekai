package model

import (
	"time"
)

type InviteLog struct {
	Id         uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	InviterId  uint       `gorm:"column:inviter_id;type:int;not null;comment:邀请人ID" json:"inviter_id"`
	UserId     uint       `gorm:"column:user_id;type:int;not null;comment:注册用户ID" json:"user_id"`
	Username   string    `gorm:"column:username;type:varchar(30);not null;comment:用户名" json:"username"`
	InviteCode string    `gorm:"column:invite_code;type:char(8);not null;comment:邀请码" json:"invite_code"`
	Remark     string    `gorm:"column:remark;type:varchar(255);not null;comment:备注" json:"remark"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *InviteLog) TableName() string {
	return "chatgpt_invite_logs"
}
