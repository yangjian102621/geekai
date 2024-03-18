package vo

type User struct {
	BaseVo
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar"`
	Salt        string   `json:"salt"`          // 密码盐
	Power       int      `json:"power"`         // 剩余算力
	ChatRoles   []string `json:"chat_roles"`    // 聊天角色集合
	ChatModels  []int    `json:"chat_models"`   // AI模型集合
	ExpiredTime int64    `json:"expired_time"`  // 账户到期时间
	Status      bool     `json:"status"`        // 当前状态
	LastLoginAt int64    `json:"last_login_at"` // 最后登录时间
	LastLoginIp string   `json:"last_login_ip"` // 最后登录 IP
	Vip         bool     `json:"vip"`
}
