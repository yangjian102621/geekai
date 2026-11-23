package service

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func clientEventTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{
		DisableAutomaticPing: true, Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	return db, mock
}

func validClientEvent() ClientEventInput {
	return ClientEventInput{
		Name: "app_start", TS: 1792576514, Props: json.RawMessage(`{"source":"desktop"}`),
		MachineID: "device-1", SessionID: "session-1", AppVersion: "1.2.3", OS: "Windows 11",
	}
}

type eventTime struct{ timestamp int64 }

func (m eventTime) Match(value driver.Value) bool {
	stamp, ok := value.(time.Time)
	return ok && stamp.Unix() == m.timestamp && stamp.Location() == time.UTC
}

func TestClientEventSavePreservesLargePropsAndSameSecondEvents(t *testing.T) {
	db, mock := clientEventTestDB(t)
	event := validClientEvent()
	event.Props = json.RawMessage(`{ "message":"` + strings.Repeat("a", 70*1024) + `", "integer":9007199254740993, "zero":0, "flag":false, "nil":null }`)
	userID := "u_148"
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `client_event`").WithArgs(
		userID, event.MachineID, event.SessionID, event.AppVersion, event.OS, "203.0.113.7", event.Name,
		string(event.Props), eventTime{event.TS}, sqlmock.AnyArg(),
		userID, event.MachineID, event.SessionID, event.AppVersion, event.OS, "203.0.113.7", event.Name,
		string(event.Props), eventTime{event.TS}, sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 2))
	mock.ExpectCommit()
	if err := NewClientEventService(db).Save(context.Background(), []ClientEventInput{event, event}, &userID, "203.0.113.7"); err != nil {
		t.Fatal(err)
	}
}

func TestClientEventSaveRejectsWholeInvalidBatch(t *testing.T) {
	db, _ := clientEventTestDB(t)
	cases := []struct {
		name   string
		change func(*ClientEventInput)
		field  string
		reason string
	}{
		{"missing name", func(e *ClientEventInput) { e.Name = "" }, "name", "required_field"},
		{"blank device", func(e *ClientEventInput) { e.MachineID = " \t" }, "machine_id", "required_field"},
		{"missing session", func(e *ClientEventInput) { e.SessionID = "" }, "session_id", "required_field"},
		{"missing version", func(e *ClientEventInput) { e.AppVersion = "" }, "app_ver", "required_field"},
		{"missing os", func(e *ClientEventInput) { e.OS = "" }, "os", "required_field"},
		{"long name", func(e *ClientEventInput) { e.Name = strings.Repeat("a", 65) }, "name", "field_too_long"},
		{"long device", func(e *ClientEventInput) { e.MachineID = strings.Repeat("a", 65) }, "machine_id", "field_too_long"},
		{"long session", func(e *ClientEventInput) { e.SessionID = strings.Repeat("a", 65) }, "session_id", "field_too_long"},
		{"long version", func(e *ClientEventInput) { e.AppVersion = strings.Repeat("a", 65) }, "app_ver", "field_too_long"},
		{"long unicode os", func(e *ClientEventInput) { e.OS = strings.Repeat("字", 256) }, "os", "field_too_long"},
		{"invalid encoding", func(e *ClientEventInput) { e.Name = string([]byte{0xff}) }, "name", "invalid_utf8"},
		{"missing timestamp", func(e *ClientEventInput) { e.TS = 0 }, "ts", "invalid_timestamp"},
		{"negative timestamp", func(e *ClientEventInput) { e.TS = -1 }, "ts", "invalid_timestamp"},
		{"timestamp overflow", func(e *ClientEventInput) { e.TS = 253402300800 }, "ts", "invalid_timestamp"},
		{"missing props", func(e *ClientEventInput) { e.Props = nil }, "props", "invalid_props"},
		{"null props", func(e *ClientEventInput) { e.Props = json.RawMessage(`null`) }, "props", "invalid_props"},
		{"array props", func(e *ClientEventInput) { e.Props = json.RawMessage(`[]`) }, "props", "invalid_props"},
		{"invalid props", func(e *ClientEventInput) { e.Props = json.RawMessage(`{`) }, "props", "invalid_props"},
		{"reserved gift", func(e *ClientEventInput) { e.Name = "login"; e.Props = json.RawMessage(`{"first_use_gift":null}`) }, "props.first_use_gift", "server_managed_field"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			event := validClientEvent()
			tc.change(&event)
			err := NewClientEventService(db).Save(context.Background(), []ClientEventInput{validClientEvent(), event}, nil, "192.0.2.1")
			var validation *ClientEventValidationError
			if !errors.Is(err, ErrInvalidClientEvents) || !errors.As(err, &validation) {
				t.Fatalf("validation details lost: %v", err)
			}
			if validation.EventNumber != 2 || validation.Field != tc.field || validation.Reason != tc.reason {
				t.Fatalf("wrong validation details: %+v", validation)
			}
		})
	}
	for _, count := range []int{0, 101} {
		events := make([]ClientEventInput, count)
		for i := range events {
			events[i] = validClientEvent()
		}
		if err := NewClientEventService(db).Save(context.Background(), events, nil, "192.0.2.1"); !errors.Is(err, ErrInvalidClientEvents) {
			t.Fatalf("invalid batch length %d was accepted: %v", count, err)
		}
	}
}

func TestClientEventSaveRollsBackDatabaseFailure(t *testing.T) {
	db, mock := clientEventTestDB(t)
	want := errors.New("insert failed")
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `client_event`").WillReturnError(want)
	mock.ExpectRollback()
	err := NewClientEventService(db).Save(context.Background(), []ClientEventInput{validClientEvent(), validClientEvent()}, nil, "192.0.2.1")
	if !errors.Is(err, want) {
		t.Fatalf("database failure lost: %v", err)
	}
}
