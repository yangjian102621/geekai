package model

type Function struct {
	Id          uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"column:name;type:varchar(30);uniqueIndex;not null;comment:函数名称" json:"name"`
	Label       string `gorm:"column:label;type:varchar(30);comment:函数标签" json:"label"`
	Description string `gorm:"column:description;type:varchar(255);comment:函数描述" json:"description"`
	Parameters  string `gorm:"column:parameters;type:text;comment:函数参数（JSON）" json:"parameters"`
	Token       string `gorm:"column:token;type:varchar(255);comment:API授权token" json:"token"`
	Action      string `gorm:"column:action;type:varchar(255);comment:函数处理 API" json:"action"`
	Enabled     bool    `gorm:"column:enabled;type:tinyint(1);not null;default:0;comment:是否启用" json:"enabled"`
}

func (m *Function) TableName() string {
	return "chatgpt_functions"
}
