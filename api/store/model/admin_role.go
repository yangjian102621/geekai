package model

import "time"

type AdminRole struct {
	Id          int `gorm:"primarykey;column:id"`
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
