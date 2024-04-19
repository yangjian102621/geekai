package vo

type DallJob struct {
	Id        uint `json:"id"`
	UserId    int `json:"user_id"`
	TaskId      string `json:"task_id"`
	Prompt    string `json:"prompt"`
	ImgURL       string `json:"img_url"`
	Publish       bool `json:"publish"`
	Power      int `json:"power"`
	Progress      int `json:"progress"`
	ErrMsg      string `json:"err_msg"`
	CreatedAt int64 `json:"created_at"`
}
