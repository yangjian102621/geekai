package vo

type InviteCode struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Code      string `json:"code"`
	Hits      int    `json:"hits"`
	RegNum    int    `json:"reg_num"`
	CreatedAt int64  `json:"created_at"`
}
