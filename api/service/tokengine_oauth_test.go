package service

import (
	"context"
	"crypto/tls"
	"errors"
	"geekai/core/types"
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
		Enabled: true, BaseURL: upstream.URL, ClientID: "geekai-web",
		RedirectURL: "http://localhost" + TokengineCallbackPath,
	})
	identity, err := oauth.Exchange(context.Background(), "test-code", "test-verifier")
	if err != nil || identity.Subject != "u_42" || identity.Username != "desktop-user" {
		t.Fatalf("OAuth login regression: %+v %v", identity, err)
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
	oauth := NewTokengineOAuth(types.TokengineOAuthConfig{Enabled: true, BaseURL: upstream.URL})
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
			oauth := NewTokengineOAuth(types.TokengineOAuthConfig{Enabled: true, BaseURL: "https://private-provider"})
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
