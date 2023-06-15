package model

type UserLoginLog struct {
	BaseModel
	UserId       uint
	Username     string
	LoginIp      string
	LoginAddress string
}
