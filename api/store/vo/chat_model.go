package vo

type ChatModel struct {
	BaseVo
	Name        string  `json:"name"`
	Value       string  `json:"value"`
	Enabled     bool    `json:"enabled"`
	SortNum     int     `json:"sort_num"`
	Power       int     `json:"power"`
	Open        bool    `json:"open"`
	MaxTokens   int     `json:"max_tokens"`  // 最大响应长度
	MaxContext  int     `json:"max_context"` // 最大上下文长度
	Temperature float32 `json:"temperature"` // 模型温度
	KeyId       int     `json:"key_id,omitempty"`
	KeyName     string  `json:"key_name"`
}
