package service

import (
	"context"
	"crypto/tls"
	"errors"
	"geekai/core/types"
	"geekai/store/model"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestTokengineExchangeStillValidatesTokenAndIdentity(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			if r.Method != http.MethodPost || r.ParseForm() != nil || r.Form.Get("code") != "test-code" || r.Form.Get("code_verifier") != "test-verifier" || r.Form.Get("client_id") != "geekai-web" {
				t.Error("OAuth exchange request changed")
			}
			_, _ = io.WriteString(w, `{"access_token":"provider-secret","token_type":"Bearer"}`)
		case "/oauth/userinfo":
			if r.Header.Get("Authorization") != "Bearer provider-secret" {
				t.Error("userinfo request lost its access token")
			}
			_, _ = io.WriteString(w, `{"sub":"u_42","username":"desktop-user"}`)
		default:
			t.Error("unexpected OAuth endpoint")
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer upstream.Close()
	oauth := NewTokengineOAuth(types.TokengineOAuthConfig{
		Enabled: true, BaseURL: upstream.URL, InternalBaseURL: upstream.URL, ClientID: "geekai-web",
		RedirectURL: "http://localhost" + TokengineCallbackPath,
	})
	identity, err := oauth.Exchange(context.Background(), "test-code", "test-verifier")
	if err != nil || identity.Subject != "u_42" || identity.Username != "desktop-user" {
		t.Fatalf("OAuth login regression: %+v %v", identity, err)
	}
}

func TestTokengineUsesInternalBaseURLForBackendRequests(t *testing.T) {
	config := types.TokengineOAuthConfig{
		Enabled: true, BaseURL: "https://provider.example", ClientID: "thinkbuddy-web",
		RedirectURL: "https://thinkbuddy.example" + TokengineCallbackPath,
	}
	publicOAuth := NewTokengineOAuth(config)
	config.InternalBaseURL = "http://provider.test.svc.cluster.local:3000/"
	oauth := NewTokengineOAuth(config)

	authorizeURL, err := oauth.AuthorizeURL("test-state", "test-verifier")
	if err != nil || !strings.HasPrefix(authorizeURL, config.BaseURL+"/oauth/authorize?") {
		t.Fatalf("authorization must use the public address: %q, %v", authorizeURL, err)
	}
	if key := oauth.IdentityKey("u_42"); key == "" || key != publicOAuth.IdentityKey("u_42") {
		t.Fatal("the internal address changed the account identity key")
	}

	oauth.client.Transport = tokengineTestTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Scheme != "http" || r.URL.Host != "provider.test.svc.cluster.local:3000" {
			t.Fatalf("backend request used %s://%s", r.URL.Scheme, r.URL.Host)
		}
		var body string
		switch r.URL.Path {
		case "/oauth/token":
			if r.Method != http.MethodPost || r.ParseForm() != nil || r.Form.Get("redirect_uri") != config.RedirectURL ||
				r.Form.Get("client_id") != config.ClientID || r.Form.Get("code_verifier") != "test-verifier" {
				t.Error("internal token exchange changed the OAuth parameters")
			}
			body = `{"access_token":"test-token","token_type":"Bearer"}`
		case "/oauth/userinfo":
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Error("missing OAuth access token")
			}
			body = `{"sub":"u_42"}`
		case "/api/user/42":
			if r.Header.Get("Authorization") != "Bearer test-service-token" {
				t.Error("missing management token")
			}
			body = `{"success":true,"data":{"id":42,"role":1,"status":1}}`
		case "/api/subscription/admin/users/42/subscriptions", "/api/user/manage":
			if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer test-service-token" {
				t.Errorf("invalid grant request for %s", r.URL.Path)
			}
			body = `{"success":true}`
		default:
			t.Fatalf("unexpected internal endpoint: %s", r.URL.Path)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
	})

	identity, err := oauth.Exchange(context.Background(), "test-code", "test-verifier")
	if err != nil || identity.Subject != "u_42" {
		t.Fatalf("internal token exchange failed: %+v %v", identity, err)
	}
	profile, err := oauth.UserProfile(context.Background(), identity.Subject, "test-service-token")
	if err != nil || profile == nil || profile.ID != 42 {
		t.Fatalf("internal profile lookup failed: %+v %v", profile, err)
	}

	db, _ := clientEventTestDB(t)
	gift := NewFirstUseGiftService(db, oauth, "test-service-token")
	if err := gift.grantSubscription(context.Background(), "42", 5); err != nil {
		t.Fatalf("internal subscription grant failed: %v", err)
	}
	quota := NewRewardQuotaService(db, &types.AppConfig{TokengineOAuth: config, TokengineServiceToken: "test-service-token"})
	quota.oauth.client.Transport = oauth.client.Transport
	if err := quota.Validate(); err != nil {
		t.Fatal(err)
	}
	reward := model.PowerLog{QuotaSubject: identity.Subject, QuotaIdentity: oauth.IdentityKey(identity.Subject), QuotaAmount: 100}
	if status, reason := quota.grant(context.Background(), reward); status != QuotaSucceeded {
		t.Fatalf("internal quota grant failed: %s %s", status, reason)
	}
}

func TestTokengineOAuthRejectsInvalidInternalBaseURL(t *testing.T) {
	for _, address := range []string{
		"", "ftp://provider.test", "http://", "/", "//provider.test",
		"http://user:password@provider.test", "http://provider.test/api/",
		"http://provider.test/?route=api", "http://provider.test?", "http://provider.test#",
	} {
		config := types.TokengineOAuthConfig{
			Enabled: true, BaseURL: "https://provider.example", ClientID: "thinkbuddy-web",
			InternalBaseURL: address, RedirectURL: "https://thinkbuddy.example" + TokengineCallbackPath,
		}
		if err := NewTokengineOAuth(config).Validate(); err == nil {
			t.Errorf("accepted invalid internal address %q", address)
		}
	}
}

func TestTokengineExchangeClassifiesTokenRejections(t *testing.T) {
	for _, tc := range []struct {
		name   string
		body   string
		reason string
	}{
		{"redirect mismatch", `{"error":"invalid_grant","error_description":"redirect_uri mismatch"}`, "redirect_uri_mismatch"},
		{"wrong verifier", `{"error":"invalid_grant","error_description":"PKCE verification failed"}`, "pkce_verification_failed"},
		{"expired or used code", `{"error":"invalid_grant","error_description":"code is invalid, expired or already used"}`, "invalid_authorization_code"},
		{"unknown description", `{"error":"invalid_grant","error_description":"provider-secret"}`, "invalid_grant"},
		{"unknown error", `{"error":"provider-secret"}`, "upstream_http_error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oauth := NewTokengineOAuth(types.TokengineOAuthConfig{
				Enabled: true, BaseURL: "https://provider.example", ClientID: "thinkbuddy-web",
				InternalBaseURL: "http://provider.test", RedirectURL: "http://localhost" + TokengineCallbackPath,
			})
			oauth.client.Transport = tokengineTestTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/oauth/token" {
					t.Fatalf("unexpected request to %s", r.URL.Path)
				}
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader(tc.body)),
				}, nil
			})
			_, err := oauth.Exchange(context.Background(), "test-code", "test-verifier")
			var failure *TokengineOAuthError
			if !errors.As(err, &failure) || failure.Reason != tc.reason || failure.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected %s with HTTP 400, got %v", tc.reason, err)
			}
		})
	}
}

func TestTokengineUserInfoNeverForwardsCredentialsOnRedirect(t *testing.T) {
	var forwarded atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwarded.Store(true)
		_, _ = io.WriteString(w, `{"sub":"forged-user"}`)
	}))
	defer target.Close()
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer upstream.Close()
	oauth := NewTokengineOAuth(types.TokengineOAuthConfig{Enabled: true, BaseURL: "https://provider.example", InternalBaseURL: upstream.URL})
	_, err := oauth.UserInfo(context.Background(), "provider-secret")
	var failure *TokengineOAuthError
	if !errors.As(err, &failure) || failure.Reason != "redirect_rejected" || forwarded.Load() {
		t.Fatal("redirect accepted or desktop credentials forwarded")
	}
}

type tokengineTestTransport func(*http.Request) (*http.Response, error)

func (f tokengineTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestTokengineUserInfoClassifiesFailuresWithoutLeakingContents(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
		body   string
		reason string
	}{
		{"timeout", context.DeadlineExceeded, 0, "", "timeout"},
		{"canceled", context.Canceled, 0, "", "request_canceled"},
		{"DNS", &net.DNSError{Err: "private-provider", Name: "private-provider", IsNotFound: true}, 0, "", "dns_failed"},
		{"connection", &net.OpError{Op: "dial", Err: errors.New("private-provider")}, 0, "", "connection_failed"},
		{"certificate", &tls.CertificateVerificationError{Err: errors.New("private-provider")}, 0, "", "tls_failed"},
		{"unknown transport", errors.New("private-provider"), 0, "", "request_failed"},
		{"unauthorized", nil, 401, `{"token":"provider-secret"}`, "invalid_token"},
		{"forbidden", nil, 403, `{}`, "invalid_token"},
		{"upstream unavailable", nil, 503, `private-provider`, "upstream_http_error"},
		{"invalid JSON", nil, 200, `private-provider`, "invalid_response_json"},
		{"missing identity", nil, 200, `{"username":"private-provider"}`, "invalid_identity"},
		{"oversized response", nil, 200, strings.Repeat("x", (1<<20)+1), "response_too_large"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oauth := NewTokengineOAuth(types.TokengineOAuthConfig{Enabled: true, BaseURL: "https://provider.example", InternalBaseURL: "http://private-provider"})
			oauth.client.Transport = tokengineTestTransport(func(r *http.Request) (*http.Response, error) {
				if r.Header.Get("Authorization") != "Bearer provider-secret" {
					t.Error("missing access token")
				}
				if tc.err != nil {
					return nil, tc.err
				}
				return &http.Response{StatusCode: tc.status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(tc.body))}, nil
			})
			_, err := oauth.UserInfo(context.Background(), "provider-secret")
			var failure *TokengineOAuthError
			if !errors.As(err, &failure) || failure.Reason != tc.reason {
				t.Fatalf("unexpected diagnostic: %v", err)
			}
			if (tc.status == 401 || tc.status == 403) && !errors.Is(err, ErrTokengineUnauthorized) {
				t.Fatal("unauthorized sentinel lost")
			}
			if tc.status >= 400 && failure.StatusCode != tc.status {
				t.Fatal("upstream status lost")
			}
			if strings.Contains(err.Error(), "provider-secret") || strings.Contains(err.Error(), "private-provider") {
				t.Fatalf("unsafe diagnostic: %v", err)
			}
		})
	}
}
