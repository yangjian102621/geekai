package model

import (
	"time"
)

type ChatRole struct {
	Id          uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(30);not null;comment:角色名称" json:"name"`
	Tid         uint       `gorm:"column:tid;type:int;not null;comment:分类ID" json:"tid"`
	Marker      string    `gorm:"column:marker;type:varchar(30);uniqueIndex;not null;comment:角色标识" json:"marker"`
	Context string    `gorm:"column:context_json;type:text;not null;comment:角色语料 json" json:"context_json"`
	HelloMsg    string    `gorm:"column:hello_msg;type:varchar(255);not null;comment:打招呼信息" json:"hello_msg"`
	Icon        string    `gorm:"column:icon;type:varchar(255);not null;comment:角色图标" json:"icon"`
	Enable      bool       `gorm:"column:enable;type:tinyint(1);not null;comment:是否被启用" json:"enable"`
	SortNum     int       `gorm:"column:sort_num;type:smallint;not null;default:0;comment:角色排序" json:"sort_num"`
	ModelId     uint       `gorm:"column:model_id;type:int;not null;default:0;comment:绑定模型ID" json:"model_id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *ChatRole) TableName() string {
	return "chatgpt_chat_roles"
}
