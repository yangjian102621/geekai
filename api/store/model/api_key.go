package model

import (
	"time"
)

// ApiKey OpenAI API 模型
type ApiKey struct {
	Id         uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(30);comment:名称" json:"name"`
	Value      string    `gorm:"column:value;type:varchar(255);not null;comment:API KEY value" json:"value"`
	Type       string    `gorm:"column:type;type:varchar(10);default:chat;not null;comment:用途（chat=>聊天，img=>图片）" json:"type"`
	LastUsedAt int64       `gorm:"column:last_used_at;type:int;not null;comment:最后使用时间" json:"last_used_at"`
	ApiURL     string    `gorm:"column:api_url;type:varchar(255);comment:API 地址" json:"api_url"`
	Enabled    bool       `gorm:"column:enabled;type:tinyint(1);comment:是否启用" json:"enabled"`
	ProxyURL   string    `gorm:"column:proxy_url;type:varchar(100);comment:代理地址" json:"proxy_url"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

// TableName 表名
func (m *ApiKey) TableName() string {
	return "chatgpt_api_keys"
}
