package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"geekai/core/types"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const TokengineCallbackPath = "/api/user/login/tokengine/callback"

var ErrTokengineUnauthorized = errors.New("invalid Tokengine access token")
var errTokengineRedirect = errors.New("OAuth endpoint redirected")

// Keep credentials, response contents and request URLs out of error messages.
type TokengineOAuthError struct {
	Reason     string
	StatusCode int
	cause      error
}

func (e *TokengineOAuthError) Error() string { return "OAuth failed: " + e.Reason }
func (e *TokengineOAuthError) Unwrap() error { return e.cause }

type TokengineOAuth struct {
	config types.TokengineOAuthConfig
	client *http.Client
}

// Only decode identity fields; provider credentials and AI tokens stay out of local sessions.
type TokengineIdentity struct {
	Subject  string `json:"sub"`
	Username string `json:"username"`
	Phone    string `json:"phone,omitempty"`
	Balance  string `json:"balance,omitempty"`
}

// Explicit fields exclude upstream passwords, management credentials and AI tokens.
type TokengineUserProfile struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	Telephone    string `json:"telephone"`
	Role         int    `json:"role"`
	Status       int    `json:"status"`
	Group        string `json:"group"`
	Quota        int64  `json:"quota"`
	UsedQuota    int64  `json:"used_quota"`
	RequestCount int64  `json:"request_count"`
	AffCode      string `json:"aff_code"`
	InviterID    int64  `json:"inviter_id"`
}

func NewTokengineOAuth(config types.TokengineOAuthConfig) *TokengineOAuth {
	config.BaseURL = strings.TrimRight(config.BaseURL, "/")
	return &TokengineOAuth{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errTokengineRedirect
			},
		},
	}
}

func (s *TokengineOAuth) Validate() error {
	if !s.config.Enabled {
		return errors.New("Token 平台登录未启用")
	}
	if strings.TrimSpace(s.config.ClientID) == "" || len(s.config.ClientID) > 64 {
		return errors.New("Token 平台客户端 ID 配置不正确")
	}
	base, err := oauthConfigURL(s.config.BaseURL)
	if err != nil || base.Path != "" {
		return errors.New("Token 平台地址配置不正确")
	}
	callback, err := oauthConfigURL(s.config.RedirectURL)
	if err != nil || callback.Path != TokengineCallbackPath || callback.RawPath != "" {
		return errors.New("Token 平台回调地址配置不正确")
	}
	return nil
}

func oauthConfigURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return nil, errors.New("invalid OAuth URL")
	}
	loopback := u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" || u.Hostname() == "::1"
	if u.Scheme != "https" && !(u.Scheme == "http" && loopback) {
		return nil, errors.New("OAuth requires HTTPS outside localhost")
	}
	return u, nil
}

func OAuthRandom() (string, error) {
	var data [32]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

func (s *TokengineOAuth) AuthorizeURL(state, verifier string) (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	challenge := sha256.Sum256([]byte(verifier))
	query := url.Values{
		"response_type":         {"code"},
		"client_id":             {s.config.ClientID},
		"redirect_uri":          {s.config.RedirectURL},
		"scope":                 {"openid"},
		"state":                 {state},
		"code_challenge":        {base64.RawURLEncoding.EncodeToString(challenge[:])},
		"code_challenge_method": {"S256"},
	}
	return s.config.BaseURL + "/oauth/authorize?" + query.Encode(), nil
}

// Include the provider origin so test and production identities cannot share an account.
func (s *TokengineOAuth) IdentityKey(subject string) string {
	u, _ := url.Parse(s.config.BaseURL)
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	if u.Scheme == "https" && u.Port() == "443" {
		u.Host = u.Hostname()
		if strings.Contains(u.Host, ":") {
			u.Host = "[" + u.Host + "]"
		}
	}
	digest := sha256.Sum256([]byte(u.String() + "\x00" + subject))
	return hex.EncodeToString(digest[:])
}

func (s *TokengineOAuth) Exchange(ctx context.Context, code, verifier string) (TokengineIdentity, error) {
	var identity TokengineIdentity
	if err := s.Validate(); err != nil {
		return identity, err
	}
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {s.config.ClientID},
		"redirect_uri":  {s.config.RedirectURL},
		"code":          {code},
		"code_verifier": {verifier},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.BaseURL+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return identity, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	if err = s.readJSON(req, &token); err != nil {
		return identity, err
	}
	if token.AccessToken == "" || len(token.AccessToken) > 16384 || !strings.EqualFold(token.TokenType, "Bearer") {
		return identity, errors.New("invalid OAuth token response")
	}
	return s.UserInfo(ctx, token.AccessToken)
}

// UserInfo verifies desktop access tokens without creating users or rewards.
func (s *TokengineOAuth) UserInfo(ctx context.Context, accessToken string) (TokengineIdentity, error) {
	var identity TokengineIdentity
	if strings.TrimSpace(accessToken) == "" || len(accessToken) > 16384 {
		return identity, &TokengineOAuthError{Reason: "invalid_token", cause: ErrTokengineUnauthorized}
	}
	base, err := oauthConfigURL(s.config.BaseURL)
	if !s.config.Enabled || err != nil || base.Path != "" {
		return identity, &TokengineOAuthError{Reason: "identity_not_configured"}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.config.BaseURL+"/oauth/userinfo", nil)
	if err != nil {
		return identity, &TokengineOAuthError{Reason: "invalid_identity_request", cause: err}
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	var payload struct {
		TokengineIdentity
		Data *TokengineIdentity `json:"data"`
	}
	if err = s.readJSON(req, &payload); err != nil {
		return identity, err
	}
	identity = payload.TokengineIdentity
	if identity.Subject == "" && payload.Data != nil {
		identity = *payload.Data
	}
	if strings.TrimSpace(identity.Subject) == "" || len(identity.Subject) > 256 ||
		utf8.RuneCountInString(identity.Username) > 64 || utf8.RuneCountInString(identity.Phone) > 32 || len(identity.Balance) > 128 {
		return identity, &TokengineOAuthError{Reason: "invalid_identity"}
	}
	return identity, nil
}

// UserProfile accepts only a subject already verified by UserInfo.
func (s *TokengineOAuth) UserProfile(ctx context.Context, subject, serviceToken string) (*TokengineUserProfile, error) {
	base, err := oauthConfigURL(s.config.BaseURL)
	serviceToken = strings.TrimSpace(serviceToken)
	if !s.config.Enabled || err != nil || base.Path != "" || serviceToken == "" {
		return nil, &TokengineOAuthError{Reason: "profile_not_configured"}
	}
	id, err := strconv.ParseInt(strings.TrimPrefix(subject, "u_"), 10, 64)
	if err != nil || id <= 0 || subject != "u_"+strconv.FormatInt(id, 10) {
		return nil, &TokengineOAuthError{Reason: "invalid_profile_subject"}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.config.BaseURL+"/api/user/"+strconv.FormatInt(id, 10), nil)
	if err != nil {
		return nil, &TokengineOAuthError{Reason: "invalid_profile_request", cause: err}
	}
	req.Header.Set("Authorization", "Bearer "+serviceToken)
	req.Header.Set("Accept", "application/json")
	var payload struct {
		Success bool                  `json:"success"`
		Data    *TokengineUserProfile `json:"data"`
	}
	if err := s.readJSON(req, &payload); err != nil {
		return nil, err
	}
	if !payload.Success || payload.Data == nil || payload.Data.ID != id || payload.Data.Role <= 0 || payload.Data.Status <= 0 {
		return nil, &TokengineOAuthError{Reason: "invalid_profile_response"}
	}
	for _, field := range []struct {
		value string
		limit int
	}{
		{payload.Data.Username, 64}, {payload.Data.DisplayName, 128}, {payload.Data.Email, 255},
		{payload.Data.Telephone, 32}, {payload.Data.Group, 64}, {payload.Data.AffCode, 32},
	} {
		if utf8.RuneCountInString(field.value) > field.limit {
			return nil, &TokengineOAuthError{Reason: "invalid_profile_response"}
		}
	}
	return payload.Data, nil
}

func (s *TokengineOAuth) readJSON(req *http.Request, target interface{}) error {
	response, err := s.client.Do(req)
	if err != nil {
		return &TokengineOAuthError{Reason: tokengineNetworkFailureReason(err), cause: err}
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return &TokengineOAuthError{Reason: "invalid_token", StatusCode: response.StatusCode, cause: ErrTokengineUnauthorized}
	}
	if response.StatusCode != http.StatusOK {
		return &TokengineOAuthError{Reason: "upstream_http_error", StatusCode: response.StatusCode}
	}
	const limit = 1 << 20
	body, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		reason := tokengineNetworkFailureReason(err)
		if reason == "request_failed" {
			reason = "response_read_failed"
		}
		return &TokengineOAuthError{Reason: reason, StatusCode: response.StatusCode, cause: err}
	}
	if len(body) > limit {
		return &TokengineOAuthError{Reason: "response_too_large", StatusCode: response.StatusCode}
	}
	if err := json.Unmarshal(body, target); err != nil {
		return &TokengineOAuthError{Reason: "invalid_response_json", StatusCode: response.StatusCode, cause: err}
	}
	return nil
}

func tokengineNetworkFailureReason(err error) string {
	var timeout net.Error
	var dns *net.DNSError
	var connection *net.OpError
	var certificate *tls.CertificateVerificationError
	var tlsHeader tls.RecordHeaderError
	switch {
	case errors.Is(err, errTokengineRedirect):
		return "redirect_rejected"
	case errors.Is(err, context.Canceled):
		return "request_canceled"
	case errors.Is(err, context.DeadlineExceeded), errors.As(err, &timeout) && timeout.Timeout():
		return "timeout"
	case errors.As(err, &dns):
		return "dns_failed"
	case errors.As(err, &certificate), errors.As(err, &tlsHeader):
		return "tls_failed"
	case errors.As(err, &connection):
		return "connection_failed"
	default:
		return "request_failed"
	}
}
