package model

import "time"

type DallJob struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户ID" json:"user_id"`
	Prompt    string    `gorm:"column:prompt;type:text;not null;comment:提示词" json:"prompt"`
	TaskInfo  string    `gorm:"column:task_info;type:text;not null;comment:任务详情" json:"task_info"`
	ImgURL    string    `gorm:"column:img_url;type:varchar(255);not null;comment:图片地址" json:"img_url"`
	OrgURL    string    `gorm:"column:org_url;type:varchar(1024);comment:原图地址" json:"org_url"`
	Publish   int       `gorm:"column:publish;type:tinyint(1);not null;comment:是否发布" json:"publish"`
	Power     int       `gorm:"column:power;type:smallint;not null;comment:消耗算力" json:"power"`
	Progress  int       `gorm:"column:progress;type:smallint;not null;comment:任务进度" json:"progress"`
	ErrMsg    string    `gorm:"column:err_msg;type:varchar(1024);not null;comment:错误信息" json:"err_msg"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *DallJob) TableName() string {
	return "chatgpt_dall_jobs"
}
