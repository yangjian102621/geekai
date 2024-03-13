package vo

import (
	"chatplus/core/types"
	"time"
)

type SdJob struct {
	Id        uint               `json:"id"`
	Type      string             `json:"type"`
	UserId    int                `json:"user_id"`
	TaskId    string             `json:"task_id"`
	ImgURL    string             `json:"img_url"`
	Params    types.SdTaskParams `json:"params"`
	Progress  int                `json:"progress"`
	Prompt    string             `json:"prompt"`
	Publish   bool               `json:"publish"`
	ErrMsg    string             `json:"err_msg"`
	Power     int                `json:"power"`
	CreatedAt time.Time          `json:"created_at"`
}
