package vo

type HistoryMessage struct {
	BaseVo
	ChatId     string `json:"chat_id"`
	UserId     uint   `json:"user_id"`
	RoleId     uint   `json:"role_id"`
	Model      string `json:"model"`
	Type       string `json:"type"`
	Icon       string `json:"icon"`
	Tokens     int    `json:"tokens"`
	Content    string `json:"content"`
	UseContext bool   `json:"use_context"`
}
