package model

import "time"

type File struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	Name      string    `gorm:"column:name;type:varchar(255);not null;comment:文件名" json:"name"`
	ObjKey    string    `gorm:"column:obj_key;type:varchar(100);comment:文件标识" json:"obj_key"`
	URL       string    `gorm:"column:url;type:varchar(255);not null;comment:文件地址" json:"url"`
	Ext       string    `gorm:"column:ext;type:varchar(10);not null;comment:文件后缀" json:"ext"`
	Size      int64     `gorm:"column:size;type:bigint;not null;default:0;comment:文件大小" json:"size"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
}

func (m *File) TableName() string {
	return "chatgpt_files"
}
