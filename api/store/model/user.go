package model

import (
	"time"
)

type User struct {
	Id             uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username       string    `gorm:"column:username;type:varchar(30);uniqueIndex;not null;comment:用户名" json:"username"`
	Mobile         string    `gorm:"column:mobile;type:char(11);comment:手机号" json:"mobile"`
	Email          string    `gorm:"column:email;type:varchar(50);comment:邮箱地址" json:"email"`
	Nickname       string    `gorm:"column:nickname;type:varchar(30);not null;comment:昵称" json:"nickname"`
	Password       string    `gorm:"column:password;type:char(64);not null;comment:密码" json:"password"`
	Avatar         string    `gorm:"column:avatar;type:varchar(255);not null;comment:头像" json:"avatar"`
	Salt           string    `gorm:"column:salt;type:char(12);not null;comment:密码盐" json:"salt"`
	Power          int       `gorm:"column:power;type:int;not null;default:0;comment:剩余算力" json:"power"`
	ExpiredTime    int64       `gorm:"column:expired_time;type:int;not null;comment:用户过期时间" json:"expired_time"`
	Status         bool       `gorm:"column:status;type:tinyint(1);not null;comment:当前状态" json:"status"`
	ChatConfig string    `gorm:"column:chat_config;type:text;not null;comment:聊天配置json" json:"chat_config"`
	ChatRoles  string    `gorm:"column:chat_roles_json;type:text;not null;comment:聊天角色 json" json:"chat_roles_json"`
	ChatModels string    `gorm:"column:chat_models_json;type:text;not null;comment:AI模型 json" json:"chat_models_json"`
	LastLoginAt    int64       `gorm:"column:last_login_at;type:int;not null;comment:最后登录时间" json:"last_login_at"`
	Vip            bool       `gorm:"column:vip;type:tinyint(1);not null;default:0;comment:是否会员" json:"vip"`
	LastLoginIp    string    `gorm:"column:last_login_ip;type:char(16);not null;comment:最后登录 IP" json:"last_login_ip"`
	OpenId         string    `gorm:"column:openid;type:varchar(100);comment:第三方登录账号ID" json:"openid"`
	Platform       string    `gorm:"column:platform;type:varchar(30);comment:登录平台" json:"platform"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *User) TableName() string {
	return "chatgpt_users"
}
