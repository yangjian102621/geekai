package model

import "time"

// LibrarySource is a server-managed connection; credentials never enter JSON responses.
type LibrarySource struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:120;not null"`
	Provider   string    `json:"provider" gorm:"size:32;not null"`
	BaseURL    string    `json:"base_url" gorm:"size:1024;not null"`
	Credential string    `json:"-" gorm:"type:text"`
	Enabled    bool      `json:"enabled" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (LibrarySource) TableName() string { return "geekai_library_sources" }

type LibraryResource struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	SourceID    uint      `json:"source_id" gorm:"not null;uniqueIndex:ux_library_remote"`
	RemoteID    string    `json:"remote_id" gorm:"size:192;not null;uniqueIndex:ux_library_remote"`
	Name        string    `json:"name" gorm:"size:120;not null"`
	Description string    `json:"description" gorm:"type:text"`
	SortNum     int       `json:"sort_num" gorm:"not null"`
	Published   bool      `json:"published" gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (LibraryResource) TableName() string { return "geekai_library_resources" }
