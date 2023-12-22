package model

type Function struct {
	Id          uint `gorm:"primarykey;column:id"`
	Name        string
	Label       string
	Description string
	Parameters  string
	Required    string
	Action      string
	Enabled     bool
}
