package model

type ChatModel struct {
	BaseModel
	Platform    string
	Name        string
	Value       string // API Key 的值
	SortNum     int
	Enabled     bool
	Power       int     // 每次对话消耗算力
	Open        bool    // 是否开放模型给所有人使用
	MaxTokens   int     // 最大响应长度
	MaxContext  int     // 最大上下文长度
	Temperature float32 // 模型温度
}
