package vo

type AI3DJob struct {
	Id         uint   `json:"id"`
	UserId     uint   `json:"user_id"`
	Type       string `json:"type"`
	Power      int    `json:"power"`
	TaskId     string `json:"task_id"`
	ImgURL     string `json:"img_url"`
	PreviewURL string `json:"preview_url"`
	Model      string `json:"model"`
	Status     string `json:"status"`
	ErrMsg     string `json:"err_msg"`
	Params     string `json:"params"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type AI3DJobCreate struct {
	Type     string `json:"type" binding:"required"`  // API类型 (tencent/gitee)
	Model    string `json:"model" binding:"required"` // 3D模型类型
	Prompt   string `json:"prompt"`                   // 文本提示词
	ImageURL string `json:"image_url"`                // 输入图片URL
	Power    int    `json:"power" binding:"required"` // 消耗算力
}

type ThreeDJobList struct {
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Total    int       `json:"total"`
	List     []AI3DJob `json:"list"`
}
