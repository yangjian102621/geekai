package model

import "time"

type SdJob struct {
	Id       uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId   uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	Type     string    `gorm:"column:type;type:varchar(20);default:txt2img;comment:任务类别" json:"type"`
	TaskId   string    `gorm:"column:task_id;type:char(30);uniqueIndex;not null;comment:任务 ID" json:"task_id"`
	TaskInfo string    `gorm:"column:task_info;type:text;not null;comment:任务详情" json:"task_info"`
	Prompt   string    `gorm:"column:prompt;type:text;not null;comment:会话提示词" json:"prompt"`
	ImgURL   string    `gorm:"column:img_url;type:varchar(255);comment:图片URL" json:"img_url"`
	Params   string    `gorm:"column:params;type:text;comment:绘画参数json" json:"params"`
	Progress int       `gorm:"column:progress;type:smallint;default:0;comment:任务进度" json:"progress"`
	Publish  int       `gorm:"column:publish;type:tinyint(1);not null;comment:是否发布" json:"publish"`
	ErrMsg   string    `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Power    int       `gorm:"column:power;type:smallint;not null;default:0;comment:消耗算力" json:"power"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *SdJob) TableName() string {
	return "chatgpt_sd_jobs"
}
