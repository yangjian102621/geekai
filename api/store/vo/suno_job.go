package vo

type SunoJob struct {
	Id           uint                   `json:"id"`
	UserId       int                    `json:"user_id"`
	Channel      string                 `json:"channel"`
	Title        string                 `json:"title"`
	Type         int                    `json:"type"`
	TaskId       string                 `json:"task_id"`
	RefTaskId    string                 `json:"ref_task_id"`  // 续写的任务id
	Tags         string                 `json:"tags"`         // 歌曲风格和标签
	Instrumental bool                   `json:"instrumental"` // 是否生成纯音乐
	ExtendSecs   int                    `json:"extend_secs"`  // 续写秒数
	SongId       string                 `json:"song_id"`      // 续写的歌曲id
	RefSongId    string                 `json:"ref_song_id"`  // 续写的歌曲id
	Prompt       string                 `json:"prompt"`       // 提示词
	CoverURL     string                 `json:"cover_url"`    // 封面图 URL
	AudioURL     string                 `json:"audio_url"`    // 音频 URL
	ModelName    string                 `json:"model_name"`   // 模型名称
	Progress     int                    `json:"progress"`     // 任务进度
	Duration     int                    `json:"duration"`     // 银屏时长，秒
	Publish      bool                   `json:"publish"`      // 是否发布
	ErrMsg       string                 `json:"err_msg"`      // 错误信息
	RawData      map[string]interface{} `json:"raw_data"`     // 原始数据 json
	Power        int                    `json:"power"`        // 消耗算力
	RefSong      map[string]interface{} `json:"ref_song,omitempty"`
	User         map[string]interface{} `json:"user,omitempty"` //关联用户信息
	PlayTimes    int                    `json:"play_times"`     // 播放次数
	CreatedAt    int64                  `json:"created_at"`
}
