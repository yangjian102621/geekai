package model

import (
	"geekai/store/vo"
	"time"
)

// PPTJob PPT 生成任务（持久化到数据库）
type PPTJob struct {
	Id              uint         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskId          string       `gorm:"column:task_id;type:varchar(64);uniqueIndex;not null;comment:任务 ID" json:"task_id"`
	UserId          uint         `gorm:"column:user_id;type:int(11);not null;index;comment:用户 ID" json:"user_id"`
	Status          string       `gorm:"column:status;type:varchar(32);not null;default:pending;comment:任务状态 pending,processing,completed,failed" json:"status"`
	ErrMsg          string       `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Prompt          string       `gorm:"column:prompt;type:text;comment:用户补充提示" json:"prompt"`
	Title           string       `gorm:"column:title;type:varchar(256);default:'';comment:列表展示标题(LLM)" json:"title"`
	Thumb           string       `gorm:"column:thumb;type:varchar(1024);default:'';comment:列表缩略图 URL" json:"thumb"`
	Content         string       `gorm:"column:content;type:text;not null;comment:PPT 内容大纲" json:"content"`
	Params          vo.PPTParams `gorm:"column:params;type:text;comment:生成参数 JSON(语言,生成模式,页数)" json:"params"`
	Slides          vo.PPTSlides `gorm:"column:slides;type:text;comment:生成的幻灯片列表 JSON" json:"slides"`
	TotalSlides     int          `gorm:"column:total_slides;type:int(11);default:0;comment:总页数" json:"total_slides"`
	CompletedSlides int          `gorm:"column:completed_slides;type:int(11);default:0;comment:已完成页数" json:"completed_slides"`
	CreatedAt       time.Time    `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *PPTJob) TableName() string {
	return "geekai_ppt_jobs"
}
