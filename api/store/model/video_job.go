package model

import "time"

type VideoJob struct {
	Id        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint      `gorm:"column:user_id;type:int(11);not null;comment:用户 ID" json:"user_id"`
	Channel   string    `gorm:"column:channel;type:varchar(100);not null;comment:渠道" json:"channel"`
	TaskId    string    `gorm:"column:task_id;type:varchar(100);not null;comment:任务 ID" json:"task_id"`
	Params    string    `gorm:"column:params;type:text;comment:视频任务参数 JSON" json:"params"`
	Type      string    `gorm:"column:type;type:varchar(20);comment:任务类型,luma,runway,cogvideo" json:"type"`
	Prompt    string    `gorm:"column:prompt;type:text;not null;comment:提示词" json:"prompt"`
	VideoURL  string    `gorm:"column:video_url;type:varchar(2000);not null;comment:视频地址" json:"video_url"`
	Status    string    `gorm:"column:status;type:varchar(20);default:pending;comment:任务状态:pending,in_progress,downloading,success,failed" json:"status"`
	Progress  int       `gorm:"column:progress;type:smallint;default:0;comment:任务进度(0-100)" json:"progress"`
	Publish   int       `gorm:"column:publish;type:tinyint(1);not null;comment:是否发布" json:"publish"`
	ErrMsg    string    `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Output    string    `gorm:"column:output;type:text;comment:任务输出的原始信息" json:"output"`
	Power     int       `gorm:"column:power;type:smallint;not null;default:0;comment:消耗算力" json:"power"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *VideoJob) TableName() string {
	return "geekai_video_jobs"
}
