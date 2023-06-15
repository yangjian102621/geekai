package vo

type ChatItem struct {
	BaseVo
	UserId uint   `json:"user_id"`
	Icon   string `json:"icon"`
	RoleId uint   `json:"role_id"`
	ChatId string `json:"chat_id"`
	Model  string `json:"model"`
	Title  string `json:"title"`
}
