package model

import (
	"time"
)

type AdminUser struct {
	Id           uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"column:username;type:varchar(30);uniqueIndex;not null;comment:用户名" json:"username"`
	Password     string    `gorm:"column:password;type:char(64);not null;comment:密码" json:"password"`
	Salt         string    `gorm:"column:salt;type:char(12);not null;comment:密码盐" json:"salt"`
	Status       bool       `gorm:"column:status;type:tinyint(1);not null;comment:当前状态" json:"status"`
	LastLoginAt  int64       `gorm:"column:last_login_at;type:int;not null;comment:最后登录时间" json:"last_login_at"`
	LastLoginIp  string    `gorm:"column:last_login_ip;type:char(16);not null;comment:最后登录 IP" json:"last_login_ip"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at"`
}

// TableName 表名
func (m *AdminUser) TableName() string {
	return "chatgpt_admin_users"
}
