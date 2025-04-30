package model

import (
	"time"
)

type ChatMessage struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	ChatId      string    `gorm:"column:chat_id;type:char(40);not null;index;comment:会话 ID" json:"chat_id"`
	Type        string    `gorm:"column:type;type:varchar(10);not null;comment:类型：prompt|reply" json:"type"`
	Icon        string    `gorm:"column:icon;type:varchar(255);not null;comment:角色图标" json:"icon"`
	RoleId      uint       `gorm:"column:role_id;type:int;not null;comment:角色 ID" json:"role_id"`
	Model       string    `gorm:"column:model;type:varchar(30);comment:模型名称" json:"model"`
	Content     string    `gorm:"column:content;type:text;not null;comment:聊天内容" json:"content"`
	Tokens      int       `gorm:"column:tokens;type:smallint;not null;comment:耗费 token 数量" json:"tokens"`
	TotalTokens int       `gorm:"column:total_tokens;type:int;not null;comment:消耗总Token长度" json:"total_tokens"`
	UseContext  bool       `gorm:"column:use_context;type:tinyint(1);not null;comment:是否允许作为上下文语料" json:"use_context"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *ChatMessage) TableName() string {
	return "chatgpt_chat_history"
}
