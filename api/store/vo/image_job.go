package vo

type ImageJob struct {
	Id       uint   `json:"id"`
	UserId   int    `json:"user_id"`
	Prompt   string `json:"prompt"`
	Params   string `json:"params"`
	ImgURL   string `json:"img_url"`
	OrgURL   string `json:"org_url"`
	Publish  bool   `json:"publish"`
	Power    int    `json:"power"`
	Progress int    `json:"progress"`
	ErrMsg   string `json:"err_msg"`
	CreatedAt int64  `json:"created_at"`
}
