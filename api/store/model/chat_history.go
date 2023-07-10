package model

type HistoryMessage struct {
	BaseModel
	ChatId  string // 会话 ID
	UserId  uint   // 用户 ID
	RoleId  uint   // 角色 ID
	Type    string
	Icon    string
	Tokens  int
	Content string
}

func (HistoryMessage) TableName() string {
	return "chatgpt_chat_history"
}
