package model

import "time"

// ClientEvent keeps the existing desktop reporting schema for data migration.
type ClientEvent struct {
	Id            uint      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID        *string   `gorm:"size:256;index:idx_client_user_received,priority:1"`
	MachineID     string    `gorm:"size:64;not null;index:idx_client_machine_received,priority:1"`
	SessionID     string    `gorm:"size:64;not null;index"`
	ClientVersion string    `gorm:"size:64;not null"`
	OS            string    `gorm:"size:255;not null"`
	IP            string    `gorm:"size:45;not null"`
	EventType     string    `gorm:"size:64;not null"`
	Props         string    `gorm:"type:longtext;not null"`
	EventAt       time.Time `gorm:"not null"`
	ReceivedAt    time.Time `gorm:"not null;index;index:idx_client_user_received,priority:2;index:idx_client_machine_received,priority:2"`
}

func (ClientEvent) TableName() string {
	return "client_event"
}
