package model

// ApiKey OpenAI API 模型
type ApiKey struct {
	BaseModel
	Platform   string
	Value      string // API Key 的值
	LastUsedAt int64  // 最后使用时间
}
