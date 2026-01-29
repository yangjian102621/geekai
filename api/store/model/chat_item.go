package model

import (
	"time"
)

type ChatItem struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ChatId    string    `gorm:"column:chat_id;type:char(40);uniqueIndex;not null;comment:会话 ID" json:"chat_id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	RoleId    uint       `gorm:"column:role_id;type:int;not null;comment:角色 ID" json:"role_id"`
	Title     string    `gorm:"column:title;type:varchar(100);not null;comment:会话标题" json:"title"`
	ModelId   uint       `gorm:"column:model_id;type:int;not null;default:0;comment:模型 ID" json:"model_id"`
	Model     string    `gorm:"column:model;type:varchar(30);comment:模型名称" json:"model"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at"`
}

func (m *ChatItem) TableName() string {
	return "chatgpt_chat_items"
}
