package vo

type UserLoginLog struct {
	BaseVo
	UserId       uint   `json:"user_id"`
	Username     string `json:"username"`
	LoginIp      string `json:"login_ip"`
	LoginAddress string `json:"login_address"`
}
