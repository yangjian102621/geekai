package service

import (
	"context"
	"encoding/json"
	"errors"
	"geekai/core/types"
	"geekai/store/model"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormlogger "gorm.io/gorm/logger"
)

const (
	QuotaWaiting    = "waiting_binding"
	QuotaPending    = "pending"
	QuotaProcessing = "processing"
	QuotaSucceeded  = "succeeded"
	QuotaFailed     = "failed"
	QuotaUnknown    = "unknown"
)

func prepareRewardQuota(tx *gorm.DB, user model.User, entry *model.PowerLog) error {
	if entry.Amount <= 0 || (entry.Type != types.PowerSignIn && entry.Type != types.PowerInvite) {
		return nil
	}
	var config model.Config
	err := tx.Where("name = ?", types.ConfigKeyRewardQuota).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	var rule types.RewardQuotaRule
	if json.Unmarshal([]byte(config.Value), &rule) != nil {
		return errors.New("额度换算配置无效")
	}
	if !rule.Enabled {
		return nil
	}
	if rule.QuotaPerPoint <= 0 || int64(entry.Amount) > math.MaxInt64/rule.QuotaPerPoint {
		return errors.New("奖励额度超出有效范围")
	}
	entry.QuotaAmount = int64(entry.Amount) * rule.QuotaPerPoint
	entry.QuotaStatus = QuotaWaiting
	if user.TokengineSubject != "" && user.TokengineID != nil && *user.TokengineID != "" {
		entry.QuotaSubject, entry.QuotaIdentity = user.TokengineSubject, *user.TokengineID
		entry.QuotaStatus = QuotaPending
	}
	return nil
}

type RewardQuotaService struct {
	db     *gorm.DB
	config *types.AppConfig
	oauth  *TokengineOAuth
}

func NewRewardQuotaService(db *gorm.DB, config *types.AppConfig) *RewardQuotaService {
	return &RewardQuotaService{
		db:     db.Session(&gorm.Session{Logger: db.Logger.LogMode(gormlogger.Silent)}),
		config: config, oauth: NewTokengineOAuth(config.TokengineOAuth),
	}
}

func (s *RewardQuotaService) Validate() error {
	base, err := oauthConfigURL(s.oauth.config.BaseURL)
	token := s.config.TokengineManagementToken()
	if !s.config.TokengineOAuth.Enabled || err != nil || base.Path != "" || token == "" || strings.ContainsAny(token, "\r\n") {
		return errors.New("请先配置 Token 平台登录地址和统一管理令牌")
	}
	return nil
}

func (s *RewardQuotaService) Run(ctx context.Context) {
	for ctx.Err() == nil {
		if err := s.processOne(ctx); err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) && ctx.Err() == nil {
			logger.Errorw("reward_quota_worker_failed")
		}
		timer := time.NewTimer(5 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (s *RewardQuotaService) processOne(ctx context.Context) error {
	if s.Validate() != nil {
		return gorm.ErrRecordNotFound
	}
	// An interrupted request may have reached the provider; never automatically resend it.
	if err := s.db.WithContext(ctx).Model(&model.PowerLog{}).
		Where("quota_status = ? AND quota_attempted_at < ?", QuotaProcessing, time.Now().Add(-time.Minute)).
		Updates(map[string]any{"quota_status": QuotaUnknown, "quota_error": "grant_interrupted"}).Error; err != nil {
		return err
	}
	var reward model.PowerLog
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bound := tx.Model(&model.User{}).Select("id").Where("tokengine_subject <> '' AND tokengine_id IS NOT NULL AND tokengine_id <> ''")
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("quota_status = ? OR (quota_status = ? AND user_id IN (?))", QuotaPending, QuotaWaiting, bound).
			Order("id ASC").First(&reward).Error; err != nil {
			return err
		}
		if reward.QuotaSubject == "" {
			var user model.User
			if err := tx.First(&user, reward.UserId).Error; err != nil {
				return err
			}
			if user.TokengineID == nil || user.TokengineSubject == "" {
				return errors.New("reward identity unavailable")
			}
			reward.QuotaSubject, reward.QuotaIdentity = user.TokengineSubject, *user.TokengineID
		}
		claimed := tx.Model(&reward).Where("quota_status = ?", reward.QuotaStatus).Updates(map[string]any{
			"quota_status": QuotaProcessing, "quota_subject": reward.QuotaSubject,
			"quota_identity": reward.QuotaIdentity, "quota_attempted_at": time.Now(), "quota_error": "",
		})
		if claimed.Error != nil {
			return claimed.Error
		}
		if claimed.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		return err
	}
	status, reason := s.grant(ctx, reward)
	finalCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	result := s.db.WithContext(finalCtx).Model(&model.PowerLog{}).
		Where("id = ? AND quota_status = ?", reward.Id, QuotaProcessing).
		Updates(map[string]any{"quota_status": status, "quota_error": reason})
	if result.Error != nil || result.RowsAffected != 1 {
		logger.Errorw("reward_quota_outcome_save_failed", "reward_id", reward.Id, "provider_outcome", status)
		return errors.New("save quota outcome failed")
	}
	if status != QuotaSucceeded {
		logger.Warnw("reward_quota_grant_failed", "reward_id", reward.Id, "user_id", reward.UserId, "status", status, "reason", reason)
	}
	return nil
}

func (s *RewardQuotaService) grant(ctx context.Context, reward model.PowerLog) (string, string) {
	id, err := strconv.ParseInt(strings.TrimPrefix(reward.QuotaSubject, "u_"), 10, 64)
	if err != nil || id <= 0 || reward.QuotaSubject != "u_"+strconv.FormatInt(id, 10) || reward.QuotaAmount <= 0 {
		return QuotaFailed, "invalid_reward"
	}
	if reward.QuotaIdentity != s.oauth.IdentityKey(reward.QuotaSubject) {
		return QuotaFailed, "identity_provider_mismatch"
	}
	body, _ := json.Marshal(map[string]any{"id": id, "action": "add_quota", "mode": "add", "value": reward.QuotaAmount})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.oauth.config.BaseURL+"/api/user/manage", strings.NewReader(string(body)))
	if err != nil {
		return QuotaFailed, "invalid_grant_request"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.TokengineManagementToken())
	var payload struct {
		Success *bool `json:"success"`
	}
	if err := s.oauth.readJSON(req, &payload); err != nil {
		var failure *TokengineOAuthError
		if errors.As(err, &failure) {
			return QuotaUnknown, failure.Reason
		}
		return QuotaUnknown, "grant_failed"
	}
	if payload.Success == nil {
		return QuotaUnknown, "invalid_grant_response"
	}
	if !*payload.Success {
		return QuotaFailed, "provider_rejected"
	}
	return QuotaSucceeded, ""
}
