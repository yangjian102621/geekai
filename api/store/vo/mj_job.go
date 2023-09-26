package vo

import "time"

type MidJourneyJob struct {
	Id          uint      `json:"id"`
	Type        string    `json:"type"`
	UserId      int       `json:"user_id"`
	MessageId   string    `json:"message_id"`
	ReferenceId string    `json:"reference_id"`
	ImgURL      string    `json:"img_url"`
	Hash        string    `json:"hash"`
	Progress    int       `json:"progress"`
	Prompt      string    `json:"prompt"`
	CreatedAt   time.Time `json:"created_at"`
	Started     bool      `json:"started"`
}
