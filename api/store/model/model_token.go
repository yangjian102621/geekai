package model

import "time"

// ModelToken 每个用户保留一个独立的模型调用令牌，数据库只保存摘要。
type ModelToken struct {
	UserId    uint      `gorm:"primaryKey;autoIncrement:false" json:"-"`
	TokenHash string    `gorm:"type:char(64);not null;uniqueIndex" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (ModelToken) TableName() string { return "thinkbuddy_model_tokens" }
