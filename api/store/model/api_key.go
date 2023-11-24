package model

// ApiKey OpenAI API 模型
type ApiKey struct {
	BaseModel
	Platform   string
	Type       string // 用途 chat => 聊天，img => 绘图
	Value      string // API Key 的值
	LastUsedAt int64  // 最后使用时间
}
