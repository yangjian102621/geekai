package types

// ApiRequest API 请求实体
type ApiRequest struct {
	Model       string        `json:"model,omitempty"` // 兼容百度文心一言
	Temperature float32       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens,omitempty"` // 兼容百度文心一言
	Stream      bool          `json:"stream"`
	Messages    []interface{} `json:"messages,omitempty"`
	Prompt      []interface{} `json:"prompt,omitempty"` // 兼容 ChatGLM
	Tools       []interface{} `json:"tools,omitempty"`
	ToolChoice  string        `json:"tool_choice,omitempty"`
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
	Role      string      `json:"role"`
	Name      string      `json:"name"`
	Content   interface{} `json:"content"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
}

// ChatSession 聊天会话对象
type ChatSession struct {
	SessionId string    `json:"session_id"`
	ClientIP  string    `json:"client_ip"` // 客户端 IP
	Username  string    `json:"username"`  // 当前登录的 username
	UserId    uint      `json:"user_id"`   // 当前登录的 user ID
	ChatId    string    `json:"chat_id"`   // 客户端聊天会话 ID, 多会话模式专用字段
	Model     ChatModel `json:"model"`     // GPT 模型
}

type ChatModel struct {
	Id       uint     `json:"id"`
	Platform Platform `json:"platform"`
	Value    string   `json:"value"`
	Weight   int      `json:"weight"`
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

var ModelToTokens = map[string]int{
	"gpt-3.5-turbo":     4096,
	"gpt-3.5-turbo-16k": 16384,
	"gpt-4":             8192,
	"gpt-4-32k":         32768,
	"chatglm_pro":       32768, // 清华智普
	"chatglm_std":       16384,
	"chatglm_lite":      4096,
	"ernie_bot_turbo":   8192, // 文心一言
	"general":           8192, // 科大讯飞
	"general2":          8192,
	"general3":          8192,
}

func GetModelMaxToken(model string) int {
	if token, ok := ModelToTokens[model]; ok {
		return token
	}
	return 4096
}
