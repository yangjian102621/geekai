package model

import (
	"geekai/core/types"
	"time"
)

// PowerLog 算力消费日志
type PowerLog struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户ID" json:"user_id"`
	Username  string    `gorm:"column:username;type:varchar(30);not null;comment:用户名" json:"username"`
	Type      types.PowerType       `gorm:"column:type;type:tinyint(1);not null;comment:类型（1：充值，2：消费，3：退费）" json:"type"`
	Amount    int       `gorm:"column:amount;type:smallint;not null;comment:算力数值" json:"amount"`
	Balance   int       `gorm:"column:balance;type:int;not null;comment:余额" json:"balance"`
	Model     string    `gorm:"column:model;type:varchar(30);not null;comment:模型" json:"model"`
	Remark    string    `gorm:"column:remark;type:varchar(512);not null;comment:备注" json:"remark"`
	Mark      types.PowerMark       `gorm:"column:mark;type:tinyint(1);not null;comment:资金类型（0：支出，1：收入）" json:"mark"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
}

func (m *PowerLog) TableName() string {
	return "chatgpt_power_logs"
}
