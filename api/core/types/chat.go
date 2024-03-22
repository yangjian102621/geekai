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
	Functions   []interface{} `json:"functions,omitempty"` // 兼容中转平台

	ToolChoice string `json:"tool_choice,omitempty"`

	Input      map[string]interface{} `json:"input,omitempty"`      //兼容阿里通义千问
	Parameters map[string]interface{} `json:"parameters,omitempty"` //兼容阿里通义千问
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
	Role         string      `json:"role"`
	Name         string      `json:"name"`
	Content      interface{} `json:"content"`
	ToolCalls    []ToolCall  `json:"tool_calls,omitempty"`
	FunctionCall struct {
		Name      string `json:"name,omitempty"`
		Arguments string `json:"arguments,omitempty"`
	} `json:"function_call,omitempty"`
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
	Id          uint     `json:"id"`
	Platform    Platform `json:"platform"`
	Name        string   `json:"name"`
	Value       string   `json:"value"`
	Power       int      `json:"power"`
	MaxTokens   int      `json:"max_tokens"`  // 最大响应长度
	MaxContext  int      `json:"max_context"` // 最大上下文长度
	Temperature float32  `json:"temperature"` // 模型温度
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

// PowerType 算力日志类型
type PowerType int

const (
	PowerRecharge = PowerType(1) // 充值
	PowerConsume  = PowerType(2) // 消费
	PowerRefund   = PowerType(3) // 任务（SD,MJ）执行失败，退款
	PowerInvite   = PowerType(4) // 邀请奖励
	PowerReward   = PowerType(5) // 众筹
	PowerGift     = PowerType(6) // 系统赠送
)

func (t PowerType) String() string {
	switch t {
	case PowerRecharge:
		return "充值"
	case PowerConsume:
		return "消费"
	case PowerRefund:
		return "退款"
	case PowerReward:
		return "众筹"

	}
	return "其他"
}

type PowerMark int

const (
	PowerSub = PowerMark(0)
	PowerAdd = PowerMark(1)
)
