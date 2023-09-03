package vo

type ChatItem struct {
	BaseVo
	UserId  uint   `json:"user_id"`
	Icon    string `json:"icon"`
	RoleId  uint   `json:"role_id"`
	ChatId  string `json:"chat_id"`
	ModelId uint   `json:"model_id"`
	Title   string `json:"title"`
}
