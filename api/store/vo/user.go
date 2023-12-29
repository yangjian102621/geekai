package vo

import "chatplus/core/types"

type User struct {
	BaseVo
	Mobile      string               `json:"mobile"`
	Nickname    string               `json:"nickname"`
	Avatar      string               `json:"avatar"`
	Salt        string               `json:"salt"`         // 密码盐
	TotalTokens int64                `json:"total_tokens"` // 总消耗tokens
	Calls       int                  `json:"calls"`        // 剩余对话次数
	ImgCalls    int                  `json:"img_calls"`
	ChatConfig  types.UserChatConfig `json:"chat_config"`   // 聊天配置
	ChatRoles   []string             `json:"chat_roles"`    // 聊天角色集合
	ChatModels  []string             `json:"chat_models"`   // AI模型集合
	ExpiredTime int64                `json:"expired_time"`  // 账户到期时间
	Status      bool                 `json:"status"`        // 当前状态
	LastLoginAt int64                `json:"last_login_at"` // 最后登录时间
	LastLoginIp string               `json:"last_login_ip"` // 最后登录 IP
	Vip         bool                 `json:"vip"`
	Tokens      int                  `json:"token"` // 当月消耗的 fee
}
