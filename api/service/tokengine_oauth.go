package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const TokengineCallbackPath = "/api/user/login/tokengine/callback"

type TokengineOAuth struct {
	config types.TokengineOAuthConfig
	client *http.Client
}

// Only decode identity fields; provider credentials and AI tokens stay out of local sessions.
type TokengineIdentity struct {
	Subject  string `json:"sub"`
	Username string `json:"username"`
}

func NewTokengineOAuth(config types.TokengineOAuthConfig) *TokengineOAuth {
	config.BaseURL = strings.TrimRight(config.BaseURL, "/")
	return &TokengineOAuth{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errors.New("OAuth endpoint redirected")
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
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, s.config.BaseURL+"/oauth/userinfo", nil)
	if err != nil {
		return identity, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	if err = s.readJSON(req, &identity); err != nil {
		return identity, err
	}
	if strings.TrimSpace(identity.Subject) == "" || len(identity.Subject) > 256 {
		return identity, errors.New("invalid OAuth identity")
	}
	return identity, nil
}

func (s *TokengineOAuth) readJSON(req *http.Request, target interface{}) error {
	response, err := s.client.Do(req)
	if err != nil {
		return errors.New("OAuth request failed")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("OAuth endpoint returned status %d", response.StatusCode)
	}
	const limit = 1 << 20
	body, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil || len(body) > limit {
		return errors.New("invalid OAuth response size")
	}
	if err := json.Unmarshal(body, target); err != nil {
		return errors.New("invalid OAuth JSON response")
	}
	return nil
}
