package vo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// SunoParam Suno 任务参数
type SunoParam struct {
	Prompt       string `json:"prompt"`           // 提示词
	Instrumental bool   `json:"instrumental"`     // 是否为纯音乐
	Tags         string `json:"tags"`             // 歌曲风格和标签
	ExtendSecs   int    `json:"extend_secs"`      // 续写秒数
	Lyrics       string `json:"lyrics,omitempty"` // 歌词（可选）
	Model        string `json:"model"`            // 模型名称
}

// Value 实现 driver.Valuer 接口，用于将 SunoParam 序列化为 JSON 字符串存储到数据库
func (p SunoParam) Value() (driver.Value, error) {
	if p.Prompt == "" && !p.Instrumental && p.Tags == "" && p.ExtendSecs == 0 && p.Lyrics == "" && p.Model == "" {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取 JSON 字符串并反序列化为 SunoParam
func (p *SunoParam) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal SunoParam value: %v", value)
	}
	return json.Unmarshal(bytes, p)
}

type SunoJob struct {
	Id           uint                   `json:"id"`
	UserId       uint                   `json:"user_id"`
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
	Output       map[string]interface{} `json:"output"`       // 原始输出数据 json
	Params       SunoParam              `json:"params"`       // 任务参数
	Power        int                    `json:"power"`        // 消耗算力
	RefSong      map[string]interface{} `json:"ref_song,omitempty"`
	User         map[string]interface{} `json:"user,omitempty"` //关联用户信息
	PlayTimes    int                    `json:"play_times"`     // 播放次数
	CreatedAt    int64                  `json:"created_at"`
}
