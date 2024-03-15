package model

// ApiKey OpenAI API 模型
type ApiKey struct {
	BaseModel
	Platform   string
	Name       string
	Type       string // 用途 chat => 聊天，img => 绘图
	Value      string // API Key 的值
	ApiURL     string // 当前 KEY 的 API 地址
	Enabled    bool   // 是否启用
	ProxyURL   string // 代理地址
	LastUsedAt int64  // 最后使用时间
}
