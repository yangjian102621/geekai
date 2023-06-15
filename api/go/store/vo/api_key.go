package vo

// ApiKey OpenAI API 模型
type ApiKey struct {
	BaseVo
	UserId     uint   `json:"user_id"`      //用户ID，系统添加的用户 ID 为 0
	Value      string `json:"value"`        // API Key 的值
	LastUsedAt int64  `json:"last_used_at"` // 最后使用时间
}
