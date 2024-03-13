package model

import "time"

type AdminRolePermission struct {
	Id           int `gorm:"primarykey;column:id"`
	RoleId       int
	PermissionId int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
