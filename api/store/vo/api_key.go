package vo

// ApiKey OpenAI API 模型
type ApiKey struct {
	BaseVo
	Platform   string `json:"platform"`
	Type       string `json:"type"`
	Value      string `json:"value"`        // API Key 的值
	LastUsedAt int64  `json:"last_used_at"` // 最后使用时间
}
