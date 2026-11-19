package handler

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/service"
	"io"
	"net"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

const clientEventBatchKey = "client_event_batch"

var errMultipleClientEventJSON = errors.New("multiple JSON values")

func rejectClientEvent(c *gin.Context, status int, code, stage string, err error) {
	reason := code
	response := gin.H{"error": code, "request_id": c.Writer.Header().Get("X-Request-ID")}
	fields := []any{
		"stage", stage, "status", status, "path", c.Request.URL.Path,
	}
	var validation *service.ClientEventValidationError
	var oauth *service.TokengineOAuthError
	var tooLarge *http.MaxBytesError
	var typeError *json.UnmarshalTypeError
	var mysqlError *mysql.MySQLError
	eventNumber := 0
	switch {
	case errors.As(err, &validation):
		reason = validation.Reason
		response["field"] = validation.Field
		eventNumber = validation.EventNumber
		if eventNumber > 0 {
			response["event_number"] = eventNumber
		}
	case errors.As(err, &oauth):
		reason = oauth.Reason
		if oauth.StatusCode > 0 {
			fields = append(fields, "upstream_status", oauth.StatusCode)
		}
	case stage == "decode":
		switch {
		case errors.As(err, &tooLarge):
			reason = "request_too_large"
		case errors.As(err, &typeError):
			reason = "json_type_mismatch"
			response["field"] = typeError.Field
			fields = append(fields, "offset", typeError.Offset)
		case errors.Is(err, errMultipleClientEventJSON):
			reason = "multiple_json_values"
		case errors.Is(err, io.EOF):
			reason = "empty_body"
		case errors.Is(err, io.ErrUnexpectedEOF):
			reason = "incomplete_json"
		default:
			var syntax *json.SyntaxError
			if errors.As(err, &syntax) {
				reason = "json_syntax_error"
				fields = append(fields, "offset", syntax.Offset)
			} else {
				reason = "request_body_read_failed"
			}
		}
	case stage == "authenticate":
		reason = "invalid_authorization"
	case stage == "save":
		reason = clientEventDatabaseReason(err)
		fields = append(fields, "error_type", fmt.Sprintf("%T", err))
		if errors.As(err, &mysqlError) {
			fields = append(fields, "database_code", mysqlError.Number)
		}
	}
	response["reason"] = reason
	for key, value := range response {
		fields = append(fields, key, value)
	}
	if value, exists := c.Get(clientEventBatchKey); exists {
		events := value.([]service.ClientEventInput)
		fields = append(fields, "batch_size", len(events))
		index := 0
		if eventNumber > 0 {
			index = eventNumber - 1
		}
		if index < len(events) {
			// Only bounded correlation fields; never log props, tokens or raw errors.
			for key, value := range map[string]string{
				"machine_id": events[index].MachineID, "session_id": events[index].SessionID,
			} {
				if value != "" && utf8.ValidString(value) && utf8.RuneCountInString(value) <= 64 {
					fields = append(fields, key, value)
				}
			}
		}
	}
	if status >= http.StatusInternalServerError {
		logger.Errorw("client_event_report_failed", fields...)
	} else {
		logger.Warnw("client_event_report_failed", fields...)
	}
	c.JSON(status, response)
}

func clientEventDatabaseReason(err error) string {
	var mysqlError *mysql.MySQLError
	var network net.Error
	switch {
	case errors.Is(err, context.Canceled):
		return "request_canceled"
	case errors.Is(err, context.DeadlineExceeded), errors.As(err, &network) && network.Timeout():
		return "database_timeout"
	case errors.Is(err, driver.ErrBadConn), errors.Is(err, mysql.ErrInvalidConn):
		return "database_connection_failed"
	case errors.As(err, &mysqlError):
		if reason := map[uint16]string{
			1044: "database_access_denied", 1045: "database_access_denied",
			1049: "database_not_found", 1146: "database_table_missing",
			1054: "database_column_missing", 1062: "database_duplicate_key",
			1153: "database_packet_too_large", 1205: "database_lock_timeout",
			1213: "database_deadlock", 1290: "database_write_restricted",
			1366: "database_invalid_value", 1406: "database_value_too_long",
		}[mysqlError.Number]; reason != "" {
			return reason
		}
	}
	var connection *net.OpError
	if errors.As(err, &connection) {
		return "database_connection_failed"
	}
	return "database_error"
}
