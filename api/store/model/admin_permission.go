package model

import "time"

type AdminPermission struct {
	Id        int `gorm:"primarykey;column:id"`
	Name      string
	Slug      string
	Sort      int
	Pid       int
	CreatedAt time.Time
	UpdatedAt time.Time
}
