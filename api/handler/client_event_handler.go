package handler

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"geekai/core"
	"geekai/service"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const defaultClientEventMaxBytes int64 = 256 * 1024

type ClientEventHandler struct {
	BaseHandler
	events   *service.ClientEventService
	oauth    *service.TokengineOAuth
	gift     *service.FirstUseGiftService
	maxBytes int64
}

func NewClientEventHandler(app *core.AppServer, db *gorm.DB) *ClientEventHandler {
	maxBytes := app.Config.ClientEventMaxBytes
	if maxBytes <= 0 {
		maxBytes = defaultClientEventMaxBytes
	}
	oauth := service.NewTokengineOAuth(app.Config.TokengineOAuth)
	return &ClientEventHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
		events:      service.NewClientEventService(db),
		oauth:       oauth,
		gift:        service.NewFirstUseGiftService(db, oauth, app.Config.TokengineManagementToken()),
		maxBytes:    maxBytes,
	}
}

func (h *ClientEventHandler) RegisterRoutes(r *core.Router) {
	public := r.PublicGroup("")
	public.POST("/api/events", h.Report)
	public.POST("/api/v1/events", h.Report)
	// Existing desktop installations post to this path.
	public.POST("/v1/events", h.Report)
}

func (h *ClientEventHandler) Report(c *gin.Context) {
	c.Header("X-Request-ID", rand.Text())
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxBytes)
	var payload struct {
		Events []json.RawMessage `json:"events"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&payload)
	if err == nil {
		var extra any
		if trailing := decoder.Decode(&extra); trailing != io.EOF {
			err = trailing
			if err == nil {
				err = errMultipleClientEventJSON
			}
		}
	}
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			rejectClientEvent(c, http.StatusRequestEntityTooLarge, "report_too_large", "decode", err)
		} else {
			rejectClientEvent(c, http.StatusBadRequest, "invalid_report_json", "decode", err)
		}
		return
	}
	if len(payload.Events) == 0 || len(payload.Events) > 100 {
		rejectClientEvent(c, http.StatusBadRequest, "invalid_report", "validate",
			&service.ClientEventValidationError{Field: "events", Reason: "invalid_batch_size"})
		return
	}
	events := make([]service.ClientEventInput, len(payload.Events))
	c.Set(clientEventBatchKey, events)
	for i, raw := range payload.Events {
		if err := json.Unmarshal(raw, &events[i]); err != nil {
			failure := &service.ClientEventValidationError{EventNumber: i + 1, Field: "events", Reason: "json_type_mismatch"}
			var typeError *json.UnmarshalTypeError
			if errors.As(err, &typeError) && typeError.Field != "" {
				failure.Field = typeError.Field
			}
			rejectClientEvent(c, http.StatusBadRequest, "invalid_report_json", "decode", failure)
			return
		}
	}

	userID, ok := h.reportUserID(c)
	if !ok {
		return
	}
	if err := h.events.Save(c.Request.Context(), events, userID, getRealClientIP(c)); err != nil {
		if errors.Is(err, service.ErrInvalidClientEvents) {
			rejectClientEvent(c, http.StatusBadRequest, "invalid_report", "validate", err)
		} else {
			rejectClientEvent(c, http.StatusInternalServerError, "report_save_failed", "save", err)
		}
		return
	}
	if userID != nil {
		for _, event := range events {
			if event.Name != "login" {
				continue
			}
			ctx, cancel := context.WithTimeout(context.WithoutCancel(c.Request.Context()), 20*time.Second)
			result, err := h.gift.Claim(ctx, *userID)
			cancel()
			if err != nil {
				fields := []any{"request_id", c.Writer.Header().Get("X-Request-ID"), "user_id", *userID, "status", result.Status}
				reason := clientEventDatabaseReason(err)
				var failure *service.FirstUseGiftError
				if errors.As(err, &failure) {
					reason = failure.Reason
					fields = append(fields, "upstream_status", failure.StatusCode)
					if errors.Unwrap(err) != nil {
						fields = append(fields, "database_reason", clientEventDatabaseReason(err))
					}
				}
				fields = append(fields, "reason", reason)
				logger.Errorw("first_use_gift_failed", fields...)
			}
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{"accepted": len(events)})
}

func (h *ClientEventHandler) reportUserID(c *gin.Context) (*string, bool) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return nil, true
	}
	fields := strings.Fields(header)
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		rejectClientEvent(c, http.StatusUnauthorized, "unauthorized", "authenticate", nil)
		return nil, false
	}
	identity, err := h.oauth.UserInfo(c.Request.Context(), fields[1])
	if err != nil {
		if errors.Is(err, service.ErrTokengineUnauthorized) {
			rejectClientEvent(c, http.StatusUnauthorized, "unauthorized", "authenticate", err)
		} else {
			rejectClientEvent(c, http.StatusBadGateway, "identity_unavailable", "authenticate", err)
		}
		return nil, false
	}
	return &identity.Subject, true
}

// Matches ThinkBuddy's header precedence; the ingress must sanitize these headers.
func getRealClientIP(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}

	for _, header := range []string{
		c.GetHeader("CF-Connecting-IP"),
		c.GetHeader("True-Client-IP"),
		c.GetHeader("X-Real-IP"),
	} {
		if ip := normalizeClientIP(header); ip != "" {
			return ip
		}
	}

	for _, part := range strings.Split(c.GetHeader("X-Forwarded-For"), ",") {
		if ip := normalizeClientIP(part); ip != "" {
			return ip
		}
	}

	if ip := normalizeClientIP(c.ClientIP()); ip != "" {
		return ip
	}

	return normalizeClientIP(c.Request.RemoteAddr)
}

func normalizeClientIP(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.Trim(value, "\"'")
	if ip, _, err := net.SplitHostPort(value); err == nil {
		value = ip
	}
	value = strings.Trim(value, "[]")
	if net.ParseIP(value) == nil {
		return ""
	}
	return value
}
