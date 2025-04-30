package model

import "time"

type SunoJob struct {
	Id          uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId      uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	Channel     string    `gorm:"column:channel;type:varchar(100);not null;comment:渠道" json:"channel"`
	Title       string    `gorm:"column:title;type:varchar(100);comment:歌曲标题" json:"title"`
	Type        int       `gorm:"column:type;type:tinyint(1);default:0;comment:任务类型,1:灵感创作,2:自定义创作" json:"type"`
	TaskId      string    `gorm:"column:task_id;type:varchar(50);comment:任务 ID" json:"task_id"`
	TaskInfo    string    `gorm:"column:task_info;type:text;not null;comment:任务详情" json:"task_info"`
	RefTaskId   string    `gorm:"column:ref_task_id;type:char(50);comment:引用任务 ID" json:"ref_task_id"`
	Tags        string    `gorm:"column:tags;type:varchar(100);comment:歌曲风格" json:"tags"`
	Instrumental bool      `gorm:"column:instrumental;type:tinyint(1);default:0;comment:是否为纯音乐" json:"instrumental"`
	ExtendSecs  int       `gorm:"column:extend_secs;type:smallint;default:0;comment:延长秒数" json:"extend_secs"`
	SongId      string    `gorm:"column:song_id;type:varchar(50);comment:要续写的歌曲 ID" json:"song_id"`
	RefSongId   string    `gorm:"column:ref_song_id;type:varchar(50);not null;comment:引用的歌曲ID" json:"ref_song_id"`
	Prompt      string    `gorm:"column:prompt;type:varchar(2000);not null;comment:提示词" json:"prompt"`
	CoverURL    string    `gorm:"column:cover_url;type:varchar(512);comment:封面图地址" json:"cover_url"`
	AudioURL    string    `gorm:"column:audio_url;type:varchar(512);comment:音频地址" json:"audio_url"`
	ModelName   string    `gorm:"column:model_name;type:varchar(30);comment:模型地址" json:"model_name"`
	Progress    int       `gorm:"column:progress;type:smallint;default:0;comment:任务进度" json:"progress"`
	Duration    int       `gorm:"column:duration;type:smallint;not null;default:0;comment:歌曲时长" json:"duration"`
	Publish     int       `gorm:"column:publish;type:tinyint(1);not null;comment:是否发布" json:"publish"`
	ErrMsg      string    `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	RawData     string    `gorm:"column:raw_data;type:text;comment:原始数据" json:"raw_data"`
	Power       int       `gorm:"column:power;type:smallint;not null;default:0;comment:消耗算力" json:"power"`
	PlayTimes   int       `gorm:"column:play_times;type:int;comment:播放次数" json:"play_times"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *SunoJob) TableName() string {
	return "chatgpt_suno_jobs"
}
