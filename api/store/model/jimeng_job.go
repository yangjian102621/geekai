package model

import (
	"geekai/core/types"
	"time"
)

// JimengJob 即梦AI任务模型
type JimengJob struct {
	Id        uint               `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint               `gorm:"column:user_id;type:int(11);not null;index;comment:用户ID" json:"user_id"`
	TaskId    string             `gorm:"column:task_id;type:varchar(100);not null;index;comment:任务ID" json:"task_id"`
	Type      types.JMTaskType   `gorm:"column:type;type:varchar(50);not null;comment:任务类型" json:"type"`
	ReqKey    string             `gorm:"column:req_key;type:varchar(100);comment:请求Key" json:"req_key"`
	Prompt    string             `gorm:"column:prompt;type:text;comment:提示词" json:"prompt"`
	Params    string             `gorm:"column:params;type:text;comment:任务参数JSON" json:"params"`
	ImgURL    string             `gorm:"column:img_url;type:varchar(1024);comment:图片或封面URL" json:"img_url"`
	VideoURL  string             `gorm:"column:video_url;type:varchar(1024);comment:视频URL" json:"video_url"`
	RawData   string             `gorm:"column:raw_data;type:text;comment:原始API响应" json:"raw_data"`
	Progress  int                `gorm:"column:progress;type:int;default:0;comment:进度百分比" json:"progress"`
	Status    types.JMTaskStatus `gorm:"column:status;type:varchar(20);default:'pending';comment:任务状态" json:"status"`
	ErrMsg    string             `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Power     int                `gorm:"column:power;type:int(11);default:0;comment:消耗算力" json:"power"`
	CreatedAt time.Time          `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time          `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at"`
}

// TableName 返回数据表名称
func (JimengJob) TableName() string {
	return "geekai_jimeng_jobs"
}
