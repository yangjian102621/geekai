package model

type AdminUser struct {
	BaseModel
	Username    string
	Password    string
	Salt        string // 密码盐
	Status      bool   `gorm:"default:true"` // 当前状态
	LastLoginAt int64  // 最后登录时间
	LastLoginIp string // 最后登录 IP
}
