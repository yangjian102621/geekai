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
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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
	if errors.Is(err, gorm.ErrRecordNotFound) && h.App.SysConfig.Base.EnabledRegister {
		if h.App.SysConfig.Base.DailyPower <= 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "checkin_disabled"})
			return
		}
		user, err = h.provisionTokengineCheckinUser(ctx, subject)
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "account_not_bound"})
		return
	case errors.Is(err, errCheckinIdentityUnavailable):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "identity_unavailable"})
		return
	case errors.Is(err, errCheckinUsernameConflict):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "username_conflict"})
		return
	case errors.Is(err, errCheckinMobileConflict):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "mobile_conflict"})
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

var (
	errCheckinIdentityUnavailable = errors.New("check-in identity unavailable")
	errCheckinUsernameConflict    = errors.New("check-in username conflicts with an existing account")
	errCheckinMobileConflict      = errors.New("check-in mobile conflicts with an existing account")
)

func (h *UserHandler) provisionTokengineCheckinUser(ctx context.Context, subject string) (model.User, error) {
	profile, err := h.tokengineOAuth.UserProfile(ctx, subject, h.App.Config.TokengineManagementToken())
	if err != nil {
		var failure *service.TokengineOAuthError
		if errors.As(err, &failure) && failure.StatusCode == http.StatusNotFound {
			return model.User{}, gorm.ErrRecordNotFound
		}
		logger.Warnw("checkin_tokengine_profile_unavailable", "error", err)
		return model.User{}, errCheckinIdentityUnavailable
	}
	if profile.Status != 1 {
		return model.User{}, nil
	}
	return h.tokengineUser(ctx, service.TokengineIdentity{Subject: subject, Username: profile.Username}, profile, "", tokengineDesktopCheckin)
}

func (h *UserHandler) createTokengineCheckinUser(ctx context.Context, subject, key, mobile string, profile *service.TokengineUserProfile) (model.User, error) {
	if strings.TrimSpace(profile.Username) == "" || utf8.RuneCountInString(profile.Username) > 30 || utf8.RuneCountInString(profile.DisplayName) > 30 {
		logger.Warnw("checkin_tokengine_profile_invalid_name")
		return model.User{}, errCheckinIdentityUnavailable
	}
	// The legacy request carries only a subject. Provision its own account without linking it by contact or issuing a login session.
	user, err := h.createTokengineUser(ctx, model.User{
		Username:         profile.Username,
		Nickname:         profile.DisplayName,
		Mobile:           mobile,
		Platform:         "tokengine",
		TokengineID:      &key,
		TokengineSubject: subject,
	}, "")
	var duplicate *mysql.MySQLError
	if errors.Is(err, gorm.ErrDuplicatedKey) || errors.As(err, &duplicate) && duplicate.Number == 1062 {
		return model.User{}, errCheckinUsernameConflict
	}
	return user, err
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
