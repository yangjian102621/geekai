package model

type ChatModel struct {
	BaseModel
	Platform string
	Name     string
	Value    string // API Key 的值
	SortNum  int
	Enabled  bool
}
