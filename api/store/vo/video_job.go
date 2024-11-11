package vo

type VideoJob struct {
	Id        uint                   `json:"id"`
	UserId    int                    `json:"user_id"`
	Channel   string                 `json:"channel"`
	Type      string                 `json:"type"`
	TaskId    string                 `json:"task_id"`
	Prompt    string                 `json:"prompt"`     // 提示词
	PromptExt string                 `json:"prompt_ext"` // 提示词
	CoverURL  string                 `json:"cover_url"`  // 封面图 URL
	VideoURL  string                 `json:"video_url"`  // 无水印视频 URL
	WaterURL  string                 `json:"water_url"`  // 有水印视频 URL
	Progress  int                    `json:"progress"`   // 任务进度
	Publish   bool                   `json:"publish"`    // 是否发布
	ErrMsg    string                 `json:"err_msg"`    // 错误信息
	RawData   map[string]interface{} `json:"raw_data"`   // 原始数据 json
	Power     int                    `json:"power"`      // 消耗算力
	CreatedAt int64                  `json:"created_at"`
}
