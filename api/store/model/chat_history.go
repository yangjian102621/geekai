package model

import "gorm.io/gorm"

type ChatMessage struct {
	BaseModel
	ChatId     string // 会话 ID
	UserId     uint   // 用户 ID
	RoleId     uint   // 角色 ID
	Model      string // AI模型
	Type       string
	Icon       string
	Tokens     int
	Content    string
	UseContext bool // 是否可以作为聊天上下文
	DeletedAt  gorm.DeletedAt
}

func (ChatMessage) TableName() string {
	return "chatgpt_chat_history"
}
