package model

type ChatRole struct {
	BaseModel
	Key      string `gorm:"column:marker;unique"` // 角色唯一标识
	Name     string // 角色名称
	Context  string `gorm:"column:context_json"` // 角色语料信息 json
	HelloMsg string // 打招呼的消息
	Icon     string // 角色聊天图标
	Enable   bool   // 是否启用被启用
	SortNum  int    //排序数字
}
