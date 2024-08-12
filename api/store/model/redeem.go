package model

import "time"

// 兑换码

type Redeem struct {
	Id         uint   `gorm:"primarykey;column:id"`
	UserId     uint   // 用户 ID
	Name       string // 名称
	Power      int    // 算力
	Code       string // 兑换码
	Enabled    bool   // 启用状态
	RedeemedAt int64  // 兑换时间
	CreatedAt  time.Time
}
