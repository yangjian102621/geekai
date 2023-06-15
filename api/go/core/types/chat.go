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

// ChatSession 聊天会话对象
type ChatSession struct {
	SessionId string `json:"session_id"`
	ClientIP  string `json:"client_ip"` // 客户端 IP
	Username  string `json:"username"`  // 当前登录的 username
	UserId    uint   `json:"user_id"`   // 当前登录的 user ID
	ChatId    string `json:"chat_id"`   // 客户端聊天会话 ID, 多会话模式专用字段
	Model     string `json:"model"`     // GPT 模型
}

type ApiError struct {
	Error struct {
		Message string
		Type    string
		Param   interface{}
		Code    string
	}
}

const PROMPT_MSG = "prompt" // prompt message
const REPLY_MSG = "reply"   // reply message
