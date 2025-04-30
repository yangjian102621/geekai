package model

type Config struct {
	Id         uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Key     string `gorm:"column:marker;type:varchar(20);uniqueIndex;not null;comment:标识" json:"marker"`
	Config string `gorm:"column:config_json;type:text;not null" json:"config_json"`
}

func (m *Config) TableName() string {
	return "chatgpt_configs"
}
