package vo

import "time"

type MidJourneyJob struct {
	Id          uint      `json:"id"`
	Type        string    `json:"type"`
	UserId      int       `json:"user_id"`
	ChannelId   string    `json:"channel_id"`
	TaskId      string    `json:"task_id"`
	MessageId   string    `json:"message_id"`
	ReferenceId string    `json:"reference_id"`
	ImgURL      string    `json:"img_url"`
	OrgURL      string    `json:"org_url"`
	Hash        string    `json:"hash"`
	Progress    int       `json:"progress"`
	Prompt      string    `json:"prompt"`
	UseProxy    bool      `json:"use_proxy"`
	Publish     bool      `json:"publish"`
	ErrMsg      string    `json:"err_msg"`
	Power       int       `json:"power"`
	CreatedAt   time.Time `json:"created_at"`
}
