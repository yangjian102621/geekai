package vo

import "time"

type SunoJob struct {
	Id           uint   `json:"id"`
	UserId       int    `json:"user_id"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	TaskId       string `json:"task_id"`
	ReferenceId  string `json:"reference_id"`  // 续写的任务id
	Tags         string `json:"tags"`          // 歌曲风格和标签
	Instrumental bool   `json:"instrumental"`  // 是否生成纯音乐
	ExtendSecs   int    `json:"extend_secs"`   // 续写秒数
	SongId       int    `json:"song_id"`       // 续写的歌曲id
	Prompt       string `json:"prompt"`        // 提示词
	ThumbImgURL  string `json:"thumb_img_url"` // 缩略图 URL
	CoverImgURL  string `json:"cover_img_url"` // 封面图 URL
	AudioURL     string `json:"audio_url"`     // 音频 URL
	ModelName    string `json:"model_name"`    // 模型名称
	Progress     int    `json:"progress"`      // 任务进度
	Duration     int    `json:"duration"`      // 银屏时长，秒
	Publish      bool   `json:"publish"`       // 是否发布
	ErrMsg       string `json:"err_msg"`       // 错误信息
	RawData      string `json:"raw_data"`      // 原始数据 json
	Power        int    `json:"power"`         // 消耗算力
	CreatedAt    time.Time
}

func (SunoJob) TableName() string {
	return "chatgpt_suno_jobs"
}
