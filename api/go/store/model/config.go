package model

type Config struct {
	Id     uint   `gorm:"primarykey;column:id"`
	Key    string `gorm:"column:marker;unique"`
	Config string `gorm:"column:config_json"`
}
