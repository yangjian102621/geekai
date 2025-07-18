package model

import (
	"time"
)

// JimengJob 即梦AI任务模型
type JimengJob struct {
	Id         uint         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     uint         `gorm:"column:user_id;type:int;not null;index;comment:用户ID" json:"user_id"`
	TaskId     string       `gorm:"column:task_id;type:varchar(100);not null;index;comment:任务ID" json:"task_id"`
	Type       JMTaskType   `gorm:"column:type;type:varchar(50);not null;comment:任务类型" json:"type"`
	ReqKey     string       `gorm:"column:req_key;type:varchar(100);comment:请求Key" json:"req_key"`
	Prompt     string       `gorm:"column:prompt;type:text;comment:提示词" json:"prompt"`
	TaskParams string       `gorm:"column:task_params;type:text;comment:任务参数JSON" json:"task_params"`
	ImgURL     string       `gorm:"column:img_url;type:varchar(1024);comment:图片或封面URL" json:"img_url"`
	VideoURL   string       `gorm:"column:video_url;type:varchar(1024);comment:视频URL" json:"video_url"`
	RawData    string       `gorm:"column:raw_data;type:text;comment:原始API响应" json:"raw_data"`
	Progress   int          `gorm:"column:progress;type:int;default:0;comment:进度百分比" json:"progress"`
	Status     JMTaskStatus `gorm:"column:status;type:varchar(20);default:'pending';comment:任务状态" json:"status"`
	ErrMsg     string       `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Power      int          `gorm:"column:power;type:int(11);default:0;comment:消耗算力" json:"power"`
	CreatedAt  time.Time    `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time    `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at"`
}

// JMTaskStatus 任务状态
type JMTaskStatus string

const (
	JMTaskStatusInQueue    = JMTaskStatus("in_queue")   // 任务已提交
	JMTaskStatusGenerating = JMTaskStatus("generating") // 任务处理中
	JMTaskStatusDone       = JMTaskStatus("done")       // 处理完成
	JMTaskStatusNotFound   = JMTaskStatus("not_found")  // 任务未找到
	JMTaskStatusSuccess    = JMTaskStatus("success")    // 任务成功
	JMTaskStatusFailed     = JMTaskStatus("failed")     // 任务失败
)

// JMTaskType 任务类型
type JMTaskType string

const (
	JMTaskTypeTextToImage  = JMTaskType("text_to_image")  // 文生图
	JMTaskTypeImageToImage = JMTaskType("image_to_image") // 图生图
	JMTaskTypeImageEdit    = JMTaskType("image_edit")     // 图像编辑
	JMTaskTypeImageEffects = JMTaskType("image_effects")  // 图像特效
	JMTaskTypeTextToVideo  = JMTaskType("text_to_video")  // 文生视频
	JMTaskTypeImageToVideo = JMTaskType("image_to_video") // 图生视频
)

// TableName 返回数据表名称
func (JimengJob) TableName() string {
	return "chatgpt_jimeng_jobs"
}
