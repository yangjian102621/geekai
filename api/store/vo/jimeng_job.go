package vo

// JimengJob 即梦AI任务VO
type JimengJob struct {
	Id         uint   `json:"id"`
	UserId     uint   `json:"user_id"`
	TaskId     string `json:"task_id"`
	Type       string `json:"type"`
	ReqKey     string `json:"req_key"`
	Prompt     string `json:"prompt"`
	TaskParams string `json:"task_params"`
	ImgURL     string `json:"img_url"`
	VideoURL   string `json:"video_url"`
	RawData    string `json:"raw_data"`
	Progress   int    `json:"progress"`
	Status     string `json:"status"`
	ErrMsg     string `json:"err_msg"`
	Power      int    `json:"power"`
	CreatedAt  int64  `json:"created_at"` // 时间戳
	UpdatedAt  int64  `json:"updated_at"` // 时间戳
}
