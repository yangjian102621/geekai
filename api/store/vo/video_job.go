package vo

type VideoJob struct {
	Id        uint           `json:"id"`
	UserId    uint           `json:"user_id"`
	Channel   string         `json:"channel"`
	Type      string         `json:"type"`
	TaskId    string         `json:"task_id"`
	Prompt    string         `json:"prompt"`    // 提示词
	VideoURL  string         `json:"video_url"` // 视频 URL
	Status    string         `json:"status"`    // 任务状态
	Progress  int            `json:"progress"`  // 任务进度(0-100)
	Publish   bool           `json:"publish"`   // 是否发布
	ErrMsg    string         `json:"err_msg"`   // 错误信息
	Params    map[string]any `json:"params"`    // 任务参数 json（用于前端显示）
	Output    map[string]any `json:"output"`    // 任务输出信息 json（管理后台使用）
	Power     int            `json:"power"`     // 消耗算力
	CreatedAt int64          `json:"created_at"`
}
