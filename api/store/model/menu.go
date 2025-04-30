package model

// Menu 系统菜单
type Menu struct {
	Id       uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;type:varchar(30);not null;comment:菜单名称" json:"name"`
	Icon     string `gorm:"column:icon;type:varchar(150);not null;comment:菜单图标" json:"icon"`
	URL      string `gorm:"column:url;type:varchar(100);not null;comment:地址" json:"url"`
	SortNum  int    `gorm:"column:sort_num;type:smallint;not null;comment:排序" json:"sort_num"`
	Enabled  bool    `gorm:"column:enabled;type:tinyint(1);not null;comment:是否启用" json:"enabled"`
}

func (m *Menu) TableName() string {
	return "chatgpt_menus"
}
