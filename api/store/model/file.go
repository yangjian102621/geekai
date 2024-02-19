package model

import "time"

type File struct {
	Id        uint `gorm:"primarykey;column:id"`
	UserId    uint
	Name      string
	ObjKey    string
	URL       string
	Ext       string
	Size      int64
	CreatedAt time.Time
}
