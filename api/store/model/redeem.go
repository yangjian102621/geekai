package model

import "time"

// 兑换码

type Redeem struct {
	Id         uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	Name       string    `gorm:"column:name;type:varchar(30);not null;comment:兑换码名称" json:"name"`
	Power      int       `gorm:"column:power;type:int;not null;comment:算力" json:"power"`
	Code       string    `gorm:"column:code;type:varchar(100);uniqueIndex;not null;comment:兑换码" json:"code"`
	Enabled    bool       `gorm:"column:enabled;type:tinyint(1);not null;comment:是否启用" json:"enabled"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	RedeemedAt int64       `gorm:"column:redeemed_at;type:int;not null;comment:兑换时间" json:"redeemed_at"`
}

func (m *Redeem) TableName() string {
	return "chatgpt_redeems"
}
