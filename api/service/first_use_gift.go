package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type FirstUseGiftService struct {
	db     *gorm.DB
	config types.FirstUseGiftConfig
	client *http.Client
}

type FirstUseGiftResult struct {
	Status string `json:"status"`
	PlanID int    `json:"plan_id,omitempty"`
}

// Keep upstream bodies, credentials and SQL parameters out of diagnostics.
type FirstUseGiftError struct {
	Reason     string
	StatusCode int
	retryable  bool
	cause      error
}

func (e *FirstUseGiftError) Error() string { return "first use gift: " + e.Reason }
func (e *FirstUseGiftError) Unwrap() error { return e.cause }

type firstUseGiftState struct {
	PlanID    int       `json:"plan_id"`
	Status    string    `json:"status"`
	LastError string    `json:"last_error,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewFirstUseGiftService(db *gorm.DB, config types.FirstUseGiftConfig, oauthBaseURL, serviceToken string) *FirstUseGiftService {
	if strings.TrimSpace(config.BaseURL) == "" {
		config.BaseURL = oauthBaseURL
	}
	config.BaseURL = strings.TrimRight(strings.TrimSpace(config.BaseURL), "/")
	config.ServiceToken = strings.TrimSpace(serviceToken)
	return &FirstUseGiftService{
		db: db.Session(&gorm.Session{Logger: db.Logger.LogMode(gormlogger.Silent)}), config: config,
		client: &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
			return errTokengineRedirect
		}},
	}
}

// Reuses ThinkBuddy web's first authenticated login and first_use_gift state.
// Processing or uncertain outcomes must never automatically grant again.
func (s *FirstUseGiftService) Claim(ctx context.Context, userID string) (FirstUseGiftResult, error) {
	numericID := strings.TrimPrefix(userID, "u_")
	id, err := strconv.Atoi(numericID)
	if err != nil || id <= 0 || userID != "u_"+strconv.Itoa(id) {
		return FirstUseGiftResult{}, &FirstUseGiftError{Reason: "invalid_user_id"}
	}

	var login model.ClientEvent
	var props map[string]json.RawMessage
	result := FirstUseGiftResult{}
	// The database lock also coordinates separate server instances.
	digest := sha256.Sum256([]byte(userID))
	lockName := fmt.Sprintf("first_use_gift:%x", digest[:24])
	err = s.db.WithContext(ctx).Connection(func(db *gorm.DB) error {
		db = db.Session(&gorm.Session{})
		var acquired sql.NullInt64
		if err := db.Raw("SELECT GET_LOCK(?, 5)", lockName).Scan(&acquired).Error; err != nil {
			return &FirstUseGiftError{Reason: "lock_failed", cause: err}
		}
		if !acquired.Valid || acquired.Int64 != 1 {
			return &FirstUseGiftError{Reason: "lock_unavailable"}
		}
		defer func() {
			cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			defer cancel()
			if err := db.WithContext(cleanup).Exec("DO RELEASE_LOCK(?)", lockName).Error; err != nil {
				// Discard the locked connection; keep the committed reservation eligible for this grant.
				_ = db.Statement.ConnPool.(*sql.Conn).Raw(func(any) error { return driver.ErrBadConn })
				logger.Warnw("first_use_gift_lock_release_failed", "user_id", userID)
			}
		}()
		// Prefer the existing marker: concurrent inserts can commit out of ID order.
		err := db.Where("user_id = ? AND event_type = ?", userID, "login").
			Order("JSON_CONTAINS_PATH(props, 'one', '$.first_use_gift') DESC, id ASC").First(&login).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Status = "awaiting_login"
			return nil
		}
		if err != nil {
			return &FirstUseGiftError{Reason: "login_lookup_failed", cause: err}
		}
		if err := json.Unmarshal([]byte(login.Props), &props); err != nil || props == nil {
			return &FirstUseGiftError{Reason: "invalid_login_props"}
		}
		var gift firstUseGiftState
		if raw, exists := props["first_use_gift"]; exists {
			if err := json.Unmarshal(raw, &gift); err != nil {
				return &FirstUseGiftError{Reason: "invalid_gift_state"}
			}
			result = FirstUseGiftResult{Status: gift.Status, PlanID: gift.PlanID}
			switch gift.Status {
			case "succeeded":
				result.Status = "already_received"
				return nil
			case "processing", "unknown":
				return nil
			case "failed":
				// Only confirmed rejections and configuration failures can be retried.
			default:
				return &FirstUseGiftError{Reason: "invalid_gift_state"}
			}
		}
		var config model.Config
		err = db.Where("name = ?", types.ConfigKeyFirstUseGift).First(&config).Error
		var rule types.FirstUseGiftRule
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return &FirstUseGiftError{Reason: "gift_config_load_failed", cause: err}
		}
		if err == nil {
			if json.Unmarshal([]byte(config.Value), &rule) != nil || rule.PlanID < 0 || rule.Enabled && rule.PlanID == 0 {
				return &FirstUseGiftError{Reason: "invalid_gift_config"}
			}
		}
		if !rule.Enabled {
			result.Status = "disabled"
			return nil
		}
		result = FirstUseGiftResult{PlanID: rule.PlanID}
		return saveFirstUseGiftState(db, &login, props, firstUseGiftState{PlanID: result.PlanID, Status: "processing"})
	})
	if err != nil || result.Status != "" {
		return result, err
	}

	grantErr := s.grantSubscription(ctx, numericID, result.PlanID)
	state := firstUseGiftState{PlanID: result.PlanID, Status: "succeeded"}
	if grantErr != nil {
		state.Status, state.LastError = "unknown", grantErr.Reason
		if grantErr.retryable {
			state.Status = "failed"
		}
	}
	result.Status = state.Status
	// Persist the outcome even if the caller disconnected during the request.
	finalCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	if err := saveFirstUseGiftState(s.db.WithContext(finalCtx), &login, props, state); err != nil {
		return result, err
	}
	if grantErr != nil {
		return result, grantErr
	}
	result.Status = "granted"
	return result, nil
}

func saveFirstUseGiftState(db *gorm.DB, login *model.ClientEvent, props map[string]json.RawMessage, gift firstUseGiftState) error {
	gift.UpdatedAt = time.Now().UTC()
	encoded, err := json.Marshal(gift)
	if err != nil {
		return &FirstUseGiftError{Reason: "encode_gift_state_failed"}
	}
	props["first_use_gift"] = encoded
	encoded, err = json.Marshal(props)
	if err != nil {
		return &FirstUseGiftError{Reason: "encode_login_props_failed"}
	}
	update := db.Model(login).Update("props", string(encoded))
	if update.Error != nil || update.RowsAffected != 1 {
		return &FirstUseGiftError{Reason: "save_gift_state_failed", cause: update.Error}
	}
	return nil
}

func (s *FirstUseGiftService) grantSubscription(ctx context.Context, numericID string, planID int) *FirstUseGiftError {
	base, err := oauthConfigURL(s.config.BaseURL)
	if err != nil || base.Path != "" || s.config.ServiceToken == "" || strings.ContainsAny(s.config.ServiceToken, "\r\n") {
		return &FirstUseGiftError{Reason: "subscription_not_configured", retryable: true}
	}
	body, _ := json.Marshal(map[string]int{"plan_id": planID})
	endpoint := s.config.BaseURL + "/api/subscription/admin/users/" + numericID + "/subscriptions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return &FirstUseGiftError{Reason: "invalid_subscription_request", retryable: true}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.ServiceToken)
	response, err := s.client.Do(req)
	if err != nil {
		return &FirstUseGiftError{Reason: tokengineNetworkFailureReason(err)}
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return &FirstUseGiftError{Reason: "upstream_http_error", StatusCode: response.StatusCode}
	}
	const limit = 1 << 20
	data, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return &FirstUseGiftError{Reason: "response_read_failed", StatusCode: response.StatusCode}
	}
	if len(data) > limit {
		return &FirstUseGiftError{Reason: "response_too_large", StatusCode: response.StatusCode}
	}
	var payload struct {
		Success *bool `json:"success"`
	}
	if err := json.Unmarshal(data, &payload); err != nil || payload.Success == nil {
		return &FirstUseGiftError{Reason: "invalid_subscription_response", StatusCode: response.StatusCode}
	}
	if !*payload.Success {
		return &FirstUseGiftError{Reason: "subscription_rejected", StatusCode: response.StatusCode, retryable: true}
	}
	return nil
}
