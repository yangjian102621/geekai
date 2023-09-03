package model

import "gorm.io/gorm"

type ChatItem struct {
	BaseModel
	ChatId    string `gorm:"column:chat_id;unique"` // 会话 ID
	UserId    uint   // 用户 ID
	RoleId    uint   // 角色 ID
	ModelId   uint   // 会话模型
	Title     string // 会话标题
	DeletedAt gorm.DeletedAt
}
