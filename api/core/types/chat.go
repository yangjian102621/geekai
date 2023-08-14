package types

// ApiRequest API 请求实体
type ApiRequest struct {
	Model       string        `json:"model"`
	Temperature float32       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
	Stream      bool          `json:"stream"`
	Messages    []interface{} `json:"messages"`
	Functions   []Function    `json:"functions"`
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
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Role         string       `json:"role"`
	Name         string       `json:"name"`
	Content      interface{}  `json:"content"`
	FunctionCall FunctionCall `json:"function_call,omitempty"`
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

type MjTask struct {
	ChatId      string
	MessageId   string
	MessageHash string
	UserId      uint
	RoleId      uint
	Icon        string
}

type ApiError struct {
	Error struct {
		Message string
		Type    string
		Param   interface{}
		Code    string
	}
}

const PromptMsg = "prompt" // prompt message
const ReplyMsg = "reply"   // reply message
const MjMsg = "mj"

var ModelToTokens = map[string]int{
	"gpt-3.5-turbo":     4096,
	"gpt-3.5-turbo-16k": 16384,
	"gpt-4":             8192,
	"gpt-4-32k":         32768,
}

const TaskStorePrefix = "/tasks/"
