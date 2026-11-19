package handler

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	tokengineCookiePath   = "/api/user/login/tokengine"
	tokengineStateCookie  = "tokengine_oauth_state"
	tokengineTicketCookie = "tokengine_login_ticket"
)

type tokengineState struct {
	Verifier   string `json:"verifier"`
	InviteCode string `json:"invite_code"`
	ReturnTo   string `json:"return_to"`
}

type tokengineSavedProfile struct {
	service.TokengineIdentity
	Profile *service.TokengineUserProfile `json:"profile,omitempty"`
}

func oauthHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Referrer-Policy", "no-referrer")
}

func (h *UserHandler) tokengineCookie(c *gin.Context, name, value string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, value, maxAge, tokengineCookiePath, "", strings.HasPrefix(h.App.Config.TokengineOAuth.RedirectURL, "https://"), true)
}

func (h *UserHandler) TokengineConfig(c *gin.Context) {
	oauthHeaders(c)
	resp.SUCCESS(c, gin.H{"enabled": h.tokengineOAuth.Validate() == nil})
}

func (h *UserHandler) TokengineStart(c *gin.Context) {
	oauthHeaders(c)
	if h.tokengineOAuth.Validate() != nil {
		h.tokengineFailure(c, "Token 平台登录未启用或配置不正确")
		return
	}
	state, err := service.OAuthRandom()
	if err != nil {
		h.tokengineFailure(c, "无法发起登录，请重试")
		return
	}
	verifier, err := service.OAuthRandom()
	if err != nil {
		h.tokengineFailure(c, "无法发起登录，请重试")
		return
	}
	invite := strings.TrimSpace(c.Query("invite_code"))
	if len(invite) > 128 {
		h.tokengineFailure(c, "邀请码不正确")
		return
	}
	data, err := json.Marshal(tokengineState{
		Verifier: verifier, InviteCode: invite, ReturnTo: oauthReturnPath(c.Query("redirect")),
	})
	if err != nil || h.redis.Set(c, "oauth/tokengine/state/"+state, data, 5*time.Minute).Err() != nil {
		h.tokengineFailure(c, "登录服务暂时不可用，请重试")
		return
	}
	authorizeURL, err := h.tokengineOAuth.AuthorizeURL(state, verifier)
	if err != nil {
		h.tokengineFailure(c, "Token 平台登录配置不正确")
		return
	}
	h.tokengineCookie(c, tokengineTicketCookie, "", -1)
	h.tokengineCookie(c, tokengineStateCookie, state, 300)
	c.Redirect(http.StatusFound, authorizeURL)
}

func (h *UserHandler) TokengineCallback(c *gin.Context) {
	oauthHeaders(c)
	state := c.Query("state")
	cookie, err := c.Cookie(tokengineStateCookie)
	if err != nil || len(state) != 43 || subtle.ConstantTimeCompare([]byte(state), []byte(cookie)) != 1 {
		h.tokengineFailure(c, "登录验证失败，请重新发起登录")
		return
	}
	h.tokengineCookie(c, tokengineStateCookie, "", -1)
	data, err := h.consumeOAuthValue(c, "oauth/tokengine/state/"+state)
	var saved tokengineState
	if err != nil || json.Unmarshal([]byte(data), &saved) != nil || saved.Verifier == "" {
		h.tokengineFailure(c, "登录已过期，请重新发起登录")
		return
	}
	if c.Query("error") != "" {
		h.tokengineFailure(c, "已取消 Token 平台授权")
		return
	}
	code := c.Query("code")
	if code == "" || len(code) > 2048 {
		h.tokengineFailure(c, "缺少授权码，请重新登录")
		return
	}
	identity, err := h.tokengineOAuth.Exchange(c.Request.Context(), code, saved.Verifier)
	if err != nil {
		h.tokengineFailure(c, "Token 平台验证失败，请重试")
		return
	}
	profile, err := h.tokengineOAuth.UserProfile(c.Request.Context(), identity.Subject, h.App.Config.TokengineManagementToken())
	if err != nil {
		var failure *service.TokengineOAuthError
		if errors.As(err, &failure) {
			logger.Warnw("tokengine_profile_sync_failed", "subject", identity.Subject, "reason", failure.Reason, "upstream_status", failure.StatusCode)
		}
	}
	user, err := h.tokengineUser(c.Request.Context(), identity, profile, saved.InviteCode)
	if err != nil {
		h.tokengineFailure(c, err.Error())
		return
	}
	ticket, err := service.OAuthRandom()
	if err != nil || h.redis.Set(c, "oauth/tokengine/ticket/"+ticket, user.Id, time.Minute).Err() != nil {
		h.tokengineFailure(c, "登录服务暂时不可用，请重试")
		return
	}
	// A short-lived HttpOnly cookie carries the one-time handoff, never a token in the URL.
	h.tokengineCookie(c, tokengineTicketCookie, ticket, 60)
	query := url.Values{"tokengine": {"success"}, "redirect": {oauthReturnPath(saved.ReturnTo)}}
	c.Redirect(http.StatusFound, "/login?"+query.Encode())
}

func (h *UserHandler) tokengineUser(ctx context.Context, identity service.TokengineIdentity, profile *service.TokengineUserProfile, invite string) (model.User, error) {
	key := h.tokengineOAuth.IdentityKey(identity.Subject)
	var user model.User
	err := h.DB.Where("tokengine_id = ?", key).First(&user).Error
	newUser := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !newUser {
		return user, errors.New("查询用户失败，请稍后重试")
	}
	lockedMobile := ""
	if newUser {
		if profile == nil {
			return user, errors.New("获取 Token 平台用户资料失败，请稍后重试")
		}
		if registrationMobilePattern.MatchString(profile.Telephone) {
			unlock, err := h.lockMobile(ctx, profile.Telephone)
			if err != nil {
				return user, errors.New("手机号请求正在处理中，请稍后重试")
			}
			defer unlock()
			lockedMobile = profile.Telephone
			var matches []model.User
			if h.DB.Where("mobile = ?", profile.Telephone).Limit(2).Find(&matches).Error != nil {
				return user, errors.New("查询用户失败，请稍后重试")
			}
			if len(matches) > 1 {
				return user, errors.New("该手机号对应多个账号，请联系管理员")
			}
			if len(matches) == 1 {
				user = matches[0]
				if user.TokengineID != nil && *user.TokengineID != "" && *user.TokengineID != key {
					return user, errors.New("该手机号对应的账号已绑定其他 Token 平台用户，请联系管理员")
				}
				newUser = false
			}
		}
	}
	if newUser && !h.App.SysConfig.Base.EnabledRegister {
		return user, errors.New("注册功能已关闭，请使用已注册的账号登录")
	}
	if !newUser && !user.Status {
		return user, errors.New("该用户已被禁止登录，请联系管理员")
	}

	nickname := []rune(strings.TrimSpace(identity.Username))
	if len(nickname) > 30 {
		nickname = nickname[:30]
	}
	username := user.Username
	if (newUser || user.Platform == "tokengine") && len(nickname) > 0 && string(nickname) != user.Username {
		username = string(nickname)
		var occupied int64
		if h.DB.Model(&model.User{}).Where("username = ? AND id <> ?", username, user.Id).Count(&occupied).Error != nil {
			return user, errors.New("查询用户名失败，请稍后重试")
		}
		if occupied > 0 {
			// A matching name never links an unrelated local account to this provider identity.
			const suffixes = "abcdefghijklmnopqrstuvwxyz123456789"
			prefix := nickname
			if len(prefix) > 29 {
				prefix = prefix[:29]
			}
			current := []rune(user.Username)
			if !newUser && len(current) == len(prefix)+1 && strings.HasPrefix(user.Username, string(prefix)) && strings.ContainsRune(suffixes, current[len(current)-1]) {
				username = user.Username
			} else {
				username = ""
				for _, index := range rand.Perm(len(suffixes)) {
					candidate := string(prefix) + string(suffixes[index])
					if h.DB.Model(&model.User{}).Where("username = ? AND id <> ?", candidate, user.Id).Count(&occupied).Error != nil {
						return user, errors.New("查询用户名失败，请稍后重试")
					}
					if occupied == 0 {
						username = candidate
						break
					}
				}
				if username == "" {
					return user, errors.New("该平台用户名的单字符后缀已被占用，请修改平台用户名后重试")
				}
			}
		}
	} else if newUser {
		username = "token_" + utils.RandString(20)
	}

	if profile != nil && strings.TrimSpace(profile.DisplayName) != "" {
		nickname = []rune(strings.TrimSpace(profile.DisplayName))
		if len(nickname) > 30 {
			nickname = nickname[:30]
		}
	}
	if newUser {
		// Keep the verified mobile and identity together if later profile sync fails.
		user, err = h.createNewUser(model.User{
			Username: username, Nickname: string(nickname), Mobile: lockedMobile, Platform: "tokengine", TokengineID: &key,
		}, invite)
		if err != nil {
			// Another callback may have inserted this identity while we were creating it.
			user = model.User{}
			if h.DB.Where("tokengine_id = ?", key).First(&user).Error != nil {
				return user, errors.New("创建用户失败，请检查邀请码或稍后重试")
			}
		}
	}
	if !user.Status {
		return user, errors.New("该用户已被禁止登录，请联系管理员")
	}

	var saved tokengineSavedProfile
	if user.TokengineProfile != "" && json.Unmarshal([]byte(user.TokengineProfile), &saved) != nil {
		return user, errors.New("读取 Token 平台用户资料失败，请稍后重试")
	}
	saved.TokengineIdentity = identity
	if profile != nil {
		saved.Profile = profile
	}
	data, err := json.Marshal(saved)
	if err != nil {
		return user, errors.New("保存 Token 平台用户资料失败，请稍后重试")
	}
	updates := map[string]any{"tokengine_subject": identity.Subject, "tokengine_profile": string(data)}
	bindExisting := user.TokengineID == nil || *user.TokengineID == ""
	if bindExisting {
		updates["tokengine_id"] = key
	}
	if !newUser && username != user.Username {
		updates["username"] = username
	}
	if profile != nil {
		if user.Nickname == "" {
			updates["nickname"] = string(nickname)
		}
		// Contacts may already belong to another local account; never merge by contact.
		contacts := map[string]string{}
		if profile.Email != user.Email && (profile.Email == "" || utils.IsValidEmail(profile.Email) && len([]rune(profile.Email)) <= 50) {
			contacts["email"] = profile.Email
		}
		if profile.Telephone != user.Mobile && (profile.Telephone == "" || registrationMobilePattern.MatchString(profile.Telephone)) {
			contacts["mobile"] = profile.Telephone
			// Hold both numbers so an in-flight SMS login cannot restore the old binding.
			mobiles := []string{user.Mobile, profile.Telephone}
			sort.Strings(mobiles)
			for _, mobile := range mobiles {
				if mobile == "" || mobile == lockedMobile {
					continue
				}
				unlock, err := h.lockMobile(ctx, mobile)
				if err != nil {
					return user, errors.New("手机号请求正在处理中，请稍后重试")
				}
				defer unlock()
			}
		}
		for column, value := range contacts {
			updates[column] = ""
			if value == "" {
				continue
			}
			var occupied int64
			if h.DB.Model(&model.User{}).Where(column+" = ? AND id <> ?", value, user.Id).Count(&occupied).Error != nil {
				return user, errors.New("查询用户联系方式失败，请稍后重试")
			}
			if occupied == 0 {
				updates[column] = value
			} else {
				logger.Warnw("tokengine_contact_conflict", "user_id", user.Id, "field", column)
			}
		}
	}
	update := h.DB.Model(&user)
	if bindExisting {
		update = update.Where("(tokengine_id IS NULL OR tokengine_id = '') AND mobile = ?", lockedMobile)
	}
	result := update.Updates(updates)
	if result.Error != nil || bindExisting && result.RowsAffected == 0 {
		return user, errors.New("同步 Token 平台用户资料失败，请稍后重试")
	}
	if saved.Profile != nil && saved.Profile.Status != 1 {
		return user, errors.New("该 Token 平台用户已被禁止登录，请联系管理员")
	}
	return user, nil
}

func (h *UserHandler) TokengineComplete(c *gin.Context) {
	oauthHeaders(c)
	if h.tokengineOAuth.Validate() != nil {
		resp.ERROR(c, "Token 平台登录未启用")
		return
	}
	ticket, err := c.Cookie(tokengineTicketCookie)
	h.tokengineCookie(c, tokengineTicketCookie, "", -1)
	if err != nil || len(ticket) != 43 {
		resp.ERROR(c, "登录已过期，请重新发起登录")
		return
	}
	userID, err := h.consumeOAuthValue(c, "oauth/tokengine/ticket/"+ticket)
	if err != nil {
		resp.ERROR(c, "登录已过期，请重新发起登录")
		return
	}
	var user model.User
	if h.DB.Where("id = ?", userID).First(&user).Error != nil || !user.Status {
		resp.ERROR(c, "用户不存在或已被禁止登录")
		return
	}
	token, err := h.doLogin(&user, c.ClientIP())
	if err != nil {
		resp.ERROR(c, "登录失败，请稍后重试")
		return
	}
	resp.SUCCESS(c, gin.H{"token": token, "user_id": user.Id, "username": user.Username})
}

func (h *UserHandler) consumeOAuthValue(c *gin.Context, key string) (string, error) {
	return h.redis.Eval(c, `local v = redis.call('GET', KEYS[1]); if v then redis.call('DEL', KEYS[1]) end; return v`, []string{key}).Text()
}

func (h *UserHandler) tokengineFailure(c *gin.Context, message string) {
	h.tokengineCookie(c, tokengineTicketCookie, "", -1)
	c.Redirect(http.StatusFound, "/login?"+url.Values{"tokengine": {"error"}, "message": {message}}.Encode())
}

func oauthReturnPath(raw string) string {
	if len(raw) > 2048 || !strings.HasPrefix(raw, "/") || strings.HasPrefix(raw, "//") || strings.Contains(raw, "\\") || strings.IndexFunc(raw, unicode.IsControl) >= 0 {
		return "/"
	}
	u, err := url.Parse(raw)
	if err != nil || u.IsAbs() || u.Host != "" {
		return "/"
	}
	for _, blocked := range []string{"/admin", "/login", "/register"} {
		if u.Path == blocked || strings.HasPrefix(u.Path, blocked+"/") {
			return "/"
		}
	}
	return raw
}
