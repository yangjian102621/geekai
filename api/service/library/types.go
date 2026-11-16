// Package library manages WeKnora connections distributed to ThinkBuddy.
package library

import (
	"errors"
	"regexp"
)

var (
	ErrNotFound = errors.New("资料不存在或尚未上架")
	ErrUpstream = errors.New("资料源暂时不可用")
	ErrTooLarge = errors.New("资料超过传输大小限制")
)

type Collection struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Provider    string             `json:"provider,omitempty"`
	Connection  *WeKnoraConnection `json:"connection,omitempty"`
}

type WeKnoraConnection struct {
	ServerURL       string `json:"server_url"`
	APIKey          string `json:"api_key"`
	KnowledgeBaseID string `json:"knowledge_base_id"`
}

var identifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,191}$`)

func ValidID(id string) bool { return identifier.MatchString(id) }
