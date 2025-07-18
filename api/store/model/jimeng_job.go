package model

import (
	"time"
)

// JimengJob 即梦AI任务模型
type JimengJob struct {
	Id         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     uint      `gorm:"column:user_id;type:int;not null;index;comment:用户ID" json:"user_id"`
	TaskId     string    `gorm:"column:task_id;type:varchar(100);not null;index;comment:任务ID" json:"task_id"`
	Type       string    `gorm:"column:type;type:varchar(50);not null;comment:任务类型" json:"type"`
	ReqKey     string    `gorm:"column:req_key;type:varchar(100);comment:请求Key" json:"req_key"`
	Prompt     string    `gorm:"column:prompt;type:text;comment:提示词" json:"prompt"`
	TaskParams string    `gorm:"column:task_params;type:text;comment:任务参数JSON" json:"task_params"`
	ImgURL     string    `gorm:"column:img_url;type:varchar(1024);comment:图片或封面URL" json:"img_url"`
	VideoURL   string    `gorm:"column:video_url;type:varchar(1024);comment:视频URL" json:"video_url"`
	RawData    string    `gorm:"column:raw_data;type:text;comment:原始API响应" json:"raw_data"`
	Progress   int       `gorm:"column:progress;type:int;default:0;comment:进度百分比" json:"progress"`
	Status     string    `gorm:"column:status;type:varchar(20);default:'pending';comment:任务状态" json:"status"`
	ErrMsg     string    `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Power      int       `gorm:"column:power;type:int;default:0;comment:消耗算力" json:"power"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at"`
}

// JimengJobStatus 即梦任务状态常量
const (
	JimengJobStatusPending    = "pending"
	JimengJobStatusProcessing = "processing"
	JimengJobStatusCompleted  = "completed"
	JimengJobStatusFailed     = "failed"
)

// JimengJobType 即梦任务类型常量
const (
	JimengJobTypeTextToImage          = "text_to_image"           // 文生图
	JimengJobTypeImageToImagePortrait = "image_to_image_portrait" // 图生图人像写真
	JimengJobTypeImageEdit            = "image_edit"              // 图像编辑
	JimengJobTypeImageEffects         = "image_effects"           // 图像特效
	JimengJobTypeTextToVideo          = "text_to_video"           // 文生视频
	JimengJobTypeImageToVideo         = "image_to_video"          // 图生视频
)

// ReqKey 常量定义
const (
	ReqKeyTextToImage          = "high_aes_general_v30l_zt2i" // 文生图
	ReqKeyImageToImagePortrait = "i2i_portrait_photo"         // 图生图人像写真
	ReqKeyImageEdit            = "seededit_v3.0"              // 图像编辑
	ReqKeyImageEffects         = "i2i_multi_style_zx2x"       // 图像特效
	ReqKeyTextToVideo          = "jimeng_vgfm_t2v_l20"        // 文生视频
	ReqKeyImageToVideo         = "jimeng_vgfm_i2v_l20"        // 图生视频
)

// TableName 返回数据表名称
func (JimengJob) TableName() string {
	return "chatgpt_jimeng_jobs"
}
