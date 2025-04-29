package model

import "time"

type InviteCode struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户ID" json:"user_id"`
	Code      string    `gorm:"column:code;type:char(8);uniqueIndex;not null;comment:邀请码" json:"code"`
	Hits      int       `gorm:"column:hits;type:int;not null;comment:点击次数" json:"hits"`
	RegNum    int       `gorm:"column:reg_num;type:smallint;not null;comment:注册数量" json:"reg_num"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *InviteCode) TableName() string {
	return "chatgpt_invite_codes"
}
