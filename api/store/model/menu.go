package model

// Menu 系统菜单
type Menu struct {
	Id      uint   `gorm:"primarykey;column:id"`
	Name    string // 菜单名称
	Icon    string // 菜单图标
	URL     string // 菜单跳转地址
	SortNum int    // 排序
	Enabled bool   // 启用状态
}
