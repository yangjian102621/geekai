package model

type Config struct {
	Id    uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name  string `gorm:"column:name;type:varchar(20);uniqueIndex;not null;comment:配置名称"`
	Value string `gorm:"column:value;type:text;not null"`
}

func (m *Config) TableName() string {
	return "chatgpt_configs"
}
