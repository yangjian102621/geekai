package types

// ApiRequest API 请求实体
type ApiRequest struct {
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
	Messages    []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiResponse struct {
	Choices []ChoiceItem `json:"choices"`
}

// ChoiceItem API 响应实体
type ChoiceItem struct {
	Delta        Message `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}
