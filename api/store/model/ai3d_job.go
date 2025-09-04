package model

import (
	"geekai/core/types"
	"time"
)

type AI3DJob struct {
	Id         uint               `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     uint               `gorm:"column:user_id;type:int(11);not null;comment:用户ID" json:"user_id"`
	Type       types.AI3DTaskType `gorm:"column:type;type:varchar(20);not null;comment:API类型 (tencent/gitee)" json:"type"`
	Power      int                `gorm:"column:power;type:int(11);not null;comment:消耗算力" json:"power"`
	TaskId     string             `gorm:"column:task_id;type:varchar(100);comment:第三方任务ID" json:"task_id"`
	FileURL    string             `gorm:"column:file_url;type:varchar(1024);comment:生成的3D模型文件地址" json:"file_url"`
	PreviewURL string             `gorm:"column:preview_url;type:varchar(1024);comment:预览图片地址" json:"preview_url"`
	Model      string             `gorm:"column:model;type:varchar(50);comment:使用的3D模型类型" json:"model"`
	Status     string             `gorm:"column:status;type:varchar(20);not null;default:pending;comment:任务状态" json:"status"`
	ErrMsg     string             `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Params     string             `gorm:"column:params;type:text;comment:任务参数(JSON格式)" json:"params"`
	RawData    string             `gorm:"column:raw_data;type:text;comment:API返回的原始数据" json:"raw_data"`
	CreatedAt  time.Time          `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt  time.Time          `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (m *AI3DJob) TableName() string {
	return "geekai_3d_jobs"
}
