package vo

type AdminUser struct {
	BaseVo
	Username    string `json:"username"`
	Salt        string `json:"salt"`          // 密码盐
	ExpiredTime int64  `json:"expired_time"`  // 账户到期时间
	Status      bool   `json:"status"`        // 当前状态
	LastLoginAt int64  `json:"last_login_at"` // 最后登录时间
	LastLoginIp string `json:"last_login_ip"` // 最后登录 IP
}
