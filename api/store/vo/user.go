package vo

import "chatplus/core/types"

type User struct {
	BaseVo
	Username    string           `json:"username"`
	Mobile      string           `json:"mobile"`
	Nickname    string           `json:"nickname"`
	Avatar      string           `json:"avatar"`
	Salt        string           `json:"salt"`   // 密码盐
	Tokens      int64            `json:"tokens"` // 剩余tokens
	Calls       int              `json:"calls"`  // 剩余对话次数
	ImgCalls    int              `json:"img_calls"`
	ChatConfig  types.ChatConfig `json:"chat_config"`   // 聊天配置
	ChatRoles   []string         `json:"chat_roles"`    // 聊天角色集合
	ExpiredTime int64            `json:"expired_time"`  // 账户到期时间
	Status      bool             `json:"status"`        // 当前状态
	LastLoginAt int64            `json:"last_login_at"` // 最后登录时间
	LastLoginIp string           `json:"last_login_ip"` // 最后登录 IP
}
