package model

import "time"

type AppType struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(50);not null;comment:名称" json:"name"`
	Icon      string    `gorm:"column:icon;type:varchar(255);not null;comment:图标URL" json:"icon"`
	SortNum   int       `gorm:"column:sort_num;type:tinyint;not null;comment:排序" json:"sort_num"`
	Enabled   bool       `gorm:"column:enabled;type:tinyint(1);not null;comment:是否启用" json:"enabled"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

// TableName 表名
func (m *AppType) TableName() string {
	return "chatgpt_app_types"
}
