package model

import (
	"chatplus/core/types"
	"time"
)

// PowerLog 算力消费日志
type PowerLog struct {
	Id        uint `gorm:"primarykey;column:id"`
	UserId    uint
	Username  string
	Type      types.PowerType
	Amount    int
	Balance   int
	Model     string          // 模型
	Remark    string          // 备注
	Mark      types.PowerMark // 资金类型
	CreatedAt time.Time
}
