package model

type User struct {
	BaseModel
	Username    string
	Nickname    string
	Password    string
	Avatar      string
	Salt        string // 密码盐
	Power       int    // 剩余算力
	ChatConfig  string `gorm:"column:chat_config_json"` // 聊天配置 json
	ChatRoles   string `gorm:"column:chat_roles_json"`  // 聊天角色
	ChatModels  string `gorm:"column:chat_models_json"` // AI 模型，不同的用户拥有不同的聊天模型
	ExpiredTime int64  // 账户到期时间
	Status      bool   `gorm:"default:true"` // 当前状态
	LastLoginAt int64  // 最后登录时间
	LastLoginIp string // 最后登录 IP
	Vip         bool   // 是否 VIP 会员
}
