package vo

type ChatItem struct {
	BaseVo
	UserId   uint   `json:"user_id"`
	Icon     string `json:"icon"`
	RoleId   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
	ChatId   string `json:"chat_id"`
	ModelId  uint   `json:"model_id"`
	Model    string `json:"model"`
	Title    string `json:"title"`
}
