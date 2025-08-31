package model

import (
	"time"
)

// Product 充值产品
type Product struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(30);not null;comment:名称" json:"name"`
	Price     float64   `gorm:"column:price;type:decimal(10,2);not null;default:0.00;comment:价格" json:"price"`
	Power     int       `gorm:"column:power;type:int;not null;default:0;comment:增加算力值" json:"power"`
	Enabled   bool      `gorm:"column:enabled;type:tinyint(1);not null;default:0;comment:是否启动" json:"enabled"`
	Sales     int       `gorm:"column:sales;type:int;not null;default:0;comment:销量" json:"sales"`
	SortNum   int       `gorm:"column:sort_num;type:tinyint;not null;default:0;comment:排序" json:"sort_num"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *Product) TableName() string {
	return "geekai_products"
}
