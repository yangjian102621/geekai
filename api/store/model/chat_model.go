package model

import (
	"time"
)

type ChatModel struct {
	Id          uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Description string    `gorm:"column:description;type:varchar(1024);not null;default:'';comment:模型类型描述" json:"description"`
	Category    string    `gorm:"column:category;type:varchar(1024);not null;default:'';comment:模型类别" json:"category"`
	Type        string    `gorm:"column:type;type:varchar(10);not null;default:chat;comment:模型类型（chat,img）" json:"type"`
	Name        string    `gorm:"column:name;type:varchar(255);not null;comment:模型名称" json:"name"`
	Value       string    `gorm:"column:value;type:varchar(255);not null;comment:模型值" json:"value"`
	SortNum     int       `gorm:"column:sort_num;type:tinyint(1);not null;comment:排序数字" json:"sort_num"`
	Enabled     bool       `gorm:"column:enabled;type:tinyint(1);not null;default:0;comment:是否启用模型" json:"enabled"`
	Power       int       `gorm:"column:power;type:smallint;not null;comment:消耗算力点数" json:"power"`
	Temperature float32   `gorm:"column:temperature;type:float(3,1);not null;default:1.0;comment:模型创意度" json:"temperature"`
	MaxTokens   int       `gorm:"column:max_tokens;type:int;not null;default:1024;comment:最大响应长度" json:"max_tokens"`
	MaxContext  int       `gorm:"column:max_context;type:int;not null;default:4096;comment:最大上下文长度" json:"max_context"`
	Open        bool       `gorm:"column:open;type:tinyint(1);not null;comment:是否开放模型" json:"open"`
	KeyId       uint       `gorm:"column:key_id;type:int;not null;comment:绑定API KEY ID" json:"key_id"`
	Options     string    `gorm:"column:options;type:text;not null;comment:模型自定义选项" json:"options"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func (m *ChatModel) TableName() string {
	return "chatgpt_chat_models"
}
