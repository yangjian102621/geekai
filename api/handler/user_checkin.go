package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *UserHandler) TokengineCheckin(c *gin.Context) {
	oauthHeaders(c)
	var request struct {
		ProviderUserID string `json:"provider_user_id"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(c.Writer, c.Request.Body, 1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_provider_user_id"})
		return
	}
	subject := strings.TrimSpace(request.ProviderUserID)
	id, err := strconv.ParseInt(strings.TrimPrefix(subject, "u_"), 10, 64)
	if err != nil || id <= 0 || subject != "u_"+strconv.FormatInt(id, 10) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_provider_user_id"})
		return
	}
	key := h.tokengineOAuth.IdentityKey(subject)
	if key == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "identity_unavailable"})
		return
	}

	ctx := c.Request.Context()
	var user model.User
	err = h.DB.WithContext(ctx).Select("id", "status").Where("tokengine_id = ?", key).Take(&user).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "account_not_bound"})
		return
	case err != nil:
		logger.Errorw("checkin_account_lookup_failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "checkin_failed"})
		return
	case !user.Status:
		c.JSON(http.StatusForbidden, gin.H{"error": "user_blocked"})
		return
	}

	points, alreadyCheckedIn, err := h.checkIn(ctx, user.Id, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "checkin_failed"})
		return
	}
	// The original desktop treats a zero-point success as an invalid response.
	if points <= 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "checkin_disabled"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"already_checked_in": alreadyCheckedIn,
		"reward": gin.H{
			"points_awarded": points,
			"grant_status":   "succeeded",
		},
	})
}

func (h *UserHandler) checkIn(ctx context.Context, userID uint, localOnly bool) (points int, alreadyCheckedIn bool, err error) {
	date := time.Now().Format("2006-01-02")
	key := fmt.Sprintf("signin/%d/%s", userID, date)
	power := max(0, h.App.SysConfig.Base.DailyPower)
	var signedIn bool
	_ = h.levelDB.Get(key, &signedIn)
	grant := h.ledgerService.Grant
	if localOnly {
		points, alreadyCheckedIn, err := h.existingCheckin(ctx, userID, key, date)
		if err != nil || alreadyCheckedIn || signedIn || power == 0 {
			return points, alreadyCheckedIn, err
		}
		grant = h.ledgerService.GrantLocalPoints
	} else if signedIn {
		return 0, true, nil
	}
	var granted bool
	err = h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var grantErr error
		granted, grantErr = grant(tx, userID, power, service.LedgerRef{Type: "signin:" + date, ID: userID}, model.PowerLog{
			RewardKey: &key,
			Type:      types.PowerSignIn,
			Model:     "SignIn",
			Remark:    fmt.Sprintf("每日签到奖励，金额：%d", power),
		})
		return grantErr
	})
	if err != nil {
		logger.Errorw("signin_reward_failed", "user_id", userID)
		return 0, false, err
	}
	if !granted {
		if localOnly {
			return h.existingCheckin(ctx, userID, key, date)
		}
		return 0, true, nil
	}
	h.levelDB.Put(key, true)
	return power, false, nil
}

func (h *UserHandler) existingCheckin(ctx context.Context, userID uint, key, date string) (int, bool, error) {
	var log model.PowerLog
	db := h.DB.WithContext(ctx)
	err := db.Select("amount").Where("reward_key = ?", key).Take(&log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Older GeekAI sign-ins can have a daily log without a reward key.
		err = db.Select("amount").Where("user_id = ? AND type = ? AND DATE(created_at) = ?", userID, types.PowerSignIn, date).
			Order("id ASC").Take(&log).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, false, nil
	}
	return log.Amount, err == nil, err
}
