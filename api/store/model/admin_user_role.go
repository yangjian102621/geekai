package model

import "time"

type AdminUserRole struct {
	Id        int `gorm:"primarykey;column:id"`
	AdminId   uint
	RoleId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
