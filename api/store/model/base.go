package model

import "time"

type BaseModel struct {
	Id        uint `gorm:"primarykey;column:id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
