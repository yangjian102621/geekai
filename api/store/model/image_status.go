package model

// 图片任务状态常量（普通图片、MidJourney、Suno 共用）。
const (
	ImageStatusPending     = "pending"
	ImageStatusInProgress  = "in_progress"
	ImageStatusDownloading = "downloading"
	ImageStatusSuccess     = "success"
	ImageStatusFailed      = "failed"
)
