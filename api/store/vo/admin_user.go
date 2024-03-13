package vo

type AdminUser struct {
	BaseVo
	Username    string      `json:"username"`
	Status      bool        `json:"status"`        // 当前状态
	LastLoginAt int64       `json:"last_login_at"` // 最后登录时间
	LastLoginIp string      `json:"last_login_ip"` // 最后登录 IP
	RoleIds     interface{} `json:"role_ids"`      //角色ids
}
