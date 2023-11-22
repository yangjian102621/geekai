package model

type ChatModel struct {
	BaseModel
	Platform string
	Name     string
	Value    string // API Key 的值
	SortNum  int
	Enabled  bool
	Weight   int  // 对话权重，每次对话扣减多少次对话额度
	Open     bool // 是否开放模型给所有人使用
}
