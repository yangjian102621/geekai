package model

import "time"

type AppType struct {
	Id        uint `gorm:"primarykey"`
	Name      string
	Icon      string
	SortNum   int
	CreatedAt time.Time
}
