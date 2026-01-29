package vo

type ChatModel struct {
	BaseVo
	Name        string            `json:"name"`
	Value       string            `json:"value"`
	Enabled     bool              `json:"enabled"`
	SortNum     int               `json:"sort_num"`
	Power       int               `json:"power"`
	Open        bool              `json:"open"`
	MaxTokens   int               `json:"max_tokens"`  // 最大响应长度
	MaxContext  int               `json:"max_context"` // 最大上下文长度
	Description string            `json:"description"` // 模型描述
	Category    string            `json:"category"`    //模型类别
	Temperature float32           `json:"temperature"` // 模型温度
	KeyId       uint              `json:"key_id,omitempty"`
	KeyName     string            `json:"key_name"`
	Options     map[string]string `json:"options"`
	Type        string            `json:"type"`
}
