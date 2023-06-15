package model

// ApiKey OpenAI API 模型
type ApiKey struct {
	BaseModel
	UserId     uint   //用户ID，系统添加的用户 ID 为 0
	Value      string // API Key 的值
	LastUsedAt int64  // 最后使用时间
}
