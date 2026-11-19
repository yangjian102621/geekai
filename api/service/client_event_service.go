package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/store/model"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var ErrInvalidClientEvents = errors.New("invalid client events")

type ClientEventValidationError struct {
	EventNumber int
	Field       string
	Reason      string
}

func (e *ClientEventValidationError) Error() string {
	return fmt.Sprintf("invalid client events: event=%d field=%s reason=%s", e.EventNumber, e.Field, e.Reason)
}

func (e *ClientEventValidationError) Unwrap() error { return ErrInvalidClientEvents }

// ClientEventInput matches the desktop protocol. Identity and IP come from HTTP.
type ClientEventInput struct {
	Name       string          `json:"name"`
	TS         int64           `json:"ts"`
	Props      json.RawMessage `json:"props"`
	MachineID  string          `json:"machine_id"`
	SessionID  string          `json:"session_id"`
	AppVersion string          `json:"app_ver"`
	OS         string          `json:"os"`
}

type ClientEventService struct {
	db *gorm.DB
}

func NewClientEventService(db *gorm.DB) *ClientEventService {
	return &ClientEventService{db: db}
}

func (s *ClientEventService) Save(ctx context.Context, events []ClientEventInput, userID *string, ip string) error {
	if len(events) == 0 || len(events) > 100 {
		return &ClientEventValidationError{Field: "events", Reason: "invalid_batch_size"}
	}
	now := time.Now().UTC()
	rows := make([]model.ClientEvent, 0, len(events))
	for i, event := range events {
		for _, field := range []struct {
			name  string
			value string
			limit int
		}{
			{"name", event.Name, 64}, {"machine_id", event.MachineID, 64}, {"session_id", event.SessionID, 64},
			{"app_ver", event.AppVersion, 64}, {"os", event.OS, 255},
		} {
			reason := ""
			switch {
			case strings.TrimSpace(field.value) == "":
				reason = "required_field"
			case !utf8.ValidString(field.value):
				reason = "invalid_utf8"
			case utf8.RuneCountInString(field.value) > field.limit:
				reason = "field_too_long"
			}
			if reason != "" {
				return &ClientEventValidationError{EventNumber: i + 1, Field: field.name, Reason: reason}
			}
		}
		props := bytes.TrimSpace(event.Props)
		if !json.Valid(props) || len(props) == 0 || props[0] != '{' {
			return &ClientEventValidationError{EventNumber: i + 1, Field: "props", Reason: "invalid_props"}
		}
		if event.TS <= 0 || event.TS > 253402300799 {
			return &ClientEventValidationError{EventNumber: i + 1, Field: "ts", Reason: "invalid_timestamp"}
		}
		if event.Name == "login" {
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(props, &fields); err != nil {
				return &ClientEventValidationError{EventNumber: i + 1, Field: "props", Reason: "invalid_props"}
			}
			if _, exists := fields["first_use_gift"]; exists {
				return &ClientEventValidationError{EventNumber: i + 1, Field: "props.first_use_gift", Reason: "server_managed_field"}
			}
		}
		rows = append(rows, model.ClientEvent{
			UserID: userID, MachineID: event.MachineID, SessionID: event.SessionID,
			ClientVersion: event.AppVersion, OS: event.OS, IP: ip, EventType: event.Name,
			Props: string(event.Props), EventAt: time.Unix(event.TS, 0).UTC(), ReceivedAt: now,
		})
	}
	// Validate the entire batch before GORM inserts it in a transaction.
	// Suppress SQL payloads; the handler logs only safe error codes and metadata.
	db := s.db.Session(&gorm.Session{Logger: s.db.Logger.LogMode(gormlogger.Silent)})
	return db.WithContext(ctx).Create(&rows).Error
}
