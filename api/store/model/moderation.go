package model

import "time"

type Moderation struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint      `gorm:"column:user_id;type:int(11);not null;comment:用户ID" json:"user_id"`
	Input     string    `gorm:"column:prompt;type:text;not null;comment:用户输入" json:"input"`
	Output    string    `gorm:"column:output;type:text;not null;comment:AI 输出" json:"output"`
	Result    string    `gorm:"column:result;type:text;not null;comment:鉴别结果" json:"result"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *Moderation) TableName() string {
	return "geekai_moderation"
}
