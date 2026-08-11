package vo

type User struct {
	BaseVo
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	Salt        string `json:"salt"`         // 密码盐
	Power       int    `json:"power"`        // 剩余算力
	ChatModels  []int  `json:"chat_models"`  // AI模型集合
	ChatRoles   []uint `json:"chat_roles"`   // 工作区应用 ID 列表
	ExpiredTime int64  `json:"expired_time"` // 账户到期时间
	Status      bool   `json:"status"`       // 当前状态
	LastLoginAt int64  `json:"last_login_at"`
	LastLoginIp string `json:"last_login_ip"`
	Vip         bool   `json:"vip"`
	OpenId      string `json:"openid"`   // 第三方登录 OpenID
	Platform    string `json:"platform"` // 第三方登录平台
}
