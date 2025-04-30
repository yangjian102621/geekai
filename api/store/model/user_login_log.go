package model

import (
	"time"
)

type UserLoginLog struct {
	Id           uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId       uint       `gorm:"column:user_id;type:int;not null;comment:用户ID" json:"user_id"`
	Username     string    `gorm:"column:username;type:varchar(30);not null;comment:用户名" json:"username"`
	LoginIp      string    `gorm:"column:login_ip;type:char(16);not null;comment:登录IP" json:"login_ip"`
	LoginAddress string    `gorm:"column:login_address;type:varchar(30);not null;comment:登录地址" json:"login_address"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *UserLoginLog) TableName() string {
	return "chatgpt_user_login_logs"
}
