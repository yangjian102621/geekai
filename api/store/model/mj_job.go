package model

import "time"

type MidJourneyJob struct {
	Id        uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    uint       `gorm:"column:user_id;type:int;not null;comment:用户 ID" json:"user_id"`
	TaskId    string    `gorm:"column:task_id;type:varchar(20);uniqueIndex;comment:任务 ID" json:"task_id"`
	TaskInfo  string    `gorm:"column:task_info;type:text;not null;comment:任务详情" json:"task_info"`
	Type      string    `gorm:"column:type;type:varchar(20);default:image;comment:任务类别" json:"type"`
	MessageId string    `gorm:"column:message_id;type:char(40);not null;index;comment:消息 ID" json:"message_id"`
	ChannelId string    `gorm:"column:channel_id;type:varchar(100);comment:频道ID" json:"channel_id"`
	RefId     string    `gorm:"column:reference_id;type:char(40);comment:引用消息 ID" json:"reference_id"`
	Prompt    string    `gorm:"column:prompt;type:text;not null;comment:会话提示词" json:"prompt"`
	ImgURL    string    `gorm:"column:img_url;type:varchar(400);comment:图片URL" json:"img_url"`
	OrgURL    string    `gorm:"column:org_url;type:varchar(400);comment:原始图片地址" json:"org_url"`
	Hash      string    `gorm:"column:hash;type:varchar(100);comment:message hash" json:"hash"`
	Progress  int       `gorm:"column:progress;type:smallint;default:0;comment:任务进度" json:"progress"`
	UseProxy  int       `gorm:"column:use_proxy;type:tinyint(1);not null;default:0;comment:是否使用反代" json:"use_proxy"`
	Publish   int       `gorm:"column:publish;type:tinyint(1);not null;comment:是否发布" json:"publish"`
	ErrMsg    string    `gorm:"column:err_msg;type:varchar(1024);comment:错误信息" json:"err_msg"`
	Power     int       `gorm:"column:power;type:smallint;not null;default:0;comment:消耗算力" json:"power"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

func (m *MidJourneyJob) TableName() string {
	return "chatgpt_mj_jobs"
}
