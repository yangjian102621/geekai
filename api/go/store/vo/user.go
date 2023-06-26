package vo

import "chatplus/core/types"

type User struct {
	BaseVo
	Username    string           `json:"username,omitempty"`
	Nickname    string           `json:"nickname,omitempty"`
	Avatar      string           `json:"avatar,omitempty"`
	Salt        string           `json:"salt,omitempty"`          // 密码盐
	Tokens      int64            `json:"tokens,omitempty"`        // 剩余tokens
	Calls       int              `json:"calls,omitempty"`         // 剩余对话次数
	ChatConfig  types.ChatConfig `json:"chat_config,omitempty"`   // 聊天配置
	ChatRoles   []string         `json:"chat_roles,omitempty"`    // 聊天角色集合
	ExpiredTime int64            `json:"expired_time,omitempty"`  // 账户到期时间
	Status      bool             `json:"status,omitempty"`        // 当前状态
	LastLoginAt int64            `json:"last_login_at,omitempty"` // 最后登录时间
	LastLoginIp string           `json:"last_login_ip,omitempty"` // 最后登录 IP
}
