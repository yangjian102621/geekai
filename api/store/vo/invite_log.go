package vo

type InviteLog struct {
	Id         uint   `json:"id"`
	InviterId  uint   `json:"inviter_id"`
	UserId     uint   `json:"user_id"`
	Username   string `json:"username"`
	InviteCode string `json:"invite_code"`
	Remark     string `json:"remark"`
	CreatedAt  int64  `json:"created_at"`
}
