package handler

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	mysqldriver "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const clientEventProps = `{ "message": " 测试 ", "nested": {"zero":0,"flag":false,"nil":null}, "integer":9007199254740993 }`
const clientEventJSON = `{"name":"error_backend","ts":1792576514,"machine_id":"device-1","session_id":"session-1","app_ver":"1.2.3","os":"Windows 11","props":` + clientEventProps + `,"user_id":"u_forged","ip":"1.2.3.4"}`
const clientEventBody = `{"events":[` + clientEventJSON + `]}`

func clientEventRouter(t *testing.T, cfg types.AppConfig) (*gin.Engine, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{
		DisableAutomaticPing: true, Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(middleware.ParameterHandlerMiddleware())
	app := &core.AppServer{Config: &cfg}
	NewClientEventHandler(app, db).RegisterRoutes(core.NewRouter(engine, &cfg, nil))
	return engine, mock
}

func postClientEvent(engine http.Handler, path, body, credential string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if credential != "" {
		req.Header.Set("Authorization", credential)
	}
	req.Header.Set("X-Forwarded-For", "203.0.113.7")
	req.RemoteAddr = "192.0.2.5:1234"
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func expectClientEventWrite(mock sqlmock.Sqlmock, userID driver.Value) {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `client_event`").WithArgs(
		userID, "device-1", "session-1", "1.2.3", "Windows 11", "203.0.113.7", "error_backend",
		clientEventProps, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func clientEventFailureResponse(t *testing.T, w *httptest.ResponseRecorder, code, reason string) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != code || body["reason"] != reason {
		t.Fatalf("unexpected failure response: %s", w.Body.String())
	}
	if id := w.Header().Get("X-Request-ID"); id == "" || body["request_id"] != id {
		t.Fatalf("missing correlation ID: %s", w.Body.String())
	}
	return body
}

func TestClientEventRoutesPreserveAnonymousIdentityAndRawJSON(t *testing.T) {
	engine, mock := clientEventRouter(t, types.AppConfig{})
	for _, path := range []string{"/api/events", "/api/v1/events", "/v1/events"} {
		expectClientEventWrite(mock, nil)
		w := postClientEvent(engine, path, clientEventBody, "")
		if w.Code != http.StatusOK || w.Body.String() != `{"accepted":1}` {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
	}
}

func TestClientEventRoutesStoreNormalizedProxyIP(t *testing.T) {
	engine, mock := clientEventRouter(t, types.AppConfig{})
	for _, path := range []string{"/api/v1/events", "/v1/events"} {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `client_event`").WithArgs(
			nil, "device-1", "session-1", "1.2.3", "Windows 11", "2001:db8::1", "error_backend",
			clientEventProps, sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(clientEventBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("CF-Connecting-IP", `"[2001:db8::1]:443"`)
		req.Header.Set("True-Client-IP", "203.0.113.2")
		req.Header.Set("X-Real-IP", "203.0.113.3")
		req.Header.Set("X-Forwarded-For", "203.0.113.4, 10.0.0.1")
		req.RemoteAddr = "192.0.2.5:1234"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK || w.Body.String() != `{"accepted":1}` {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
	}
}

func TestClientEventRejectsInvalidJSONAndEnforcesLimitBeforeMiddlewareReads(t *testing.T) {
	engine, _ := clientEventRouter(t, types.AppConfig{})
	for _, body := range []string{"null", `{}`, `{"events":[]}`, `{"events":[`, clientEventBody + `{}`, clientEventBody + `x`, `{"events":[` + clientEventJSON + `,{}]}`} {
		w := postClientEvent(engine, "/api/v1/events", body, "")
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected invalid batch rejection, got %d %s", w.Code, w.Body.String())
		}
	}
	for _, path := range []string{"/api/events", "/api/v1/events", "/v1/events"} {
		for _, body := range []string{
			strings.Replace(clientEventBody, `" 测试 "`, `"`+strings.Repeat("a", 256*1024)+`"`, 1),
			clientEventBody + strings.Repeat(" ", 256*1024),
		} {
			w := postClientEvent(engine, path, body, "")
			if w.Code != http.StatusRequestEntityTooLarge {
				t.Fatalf("body cap bypassed on %s: %d %s", path, w.Code, w.Body.String())
			}
		}
	}
	engine, _ = clientEventRouter(t, types.AppConfig{ClientEventMaxBytes: 32})
	if w := postClientEvent(engine, "/api/v1/events", clientEventBody, ""); w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("configured limit ignored: %d", w.Code)
	}
}

func TestClientEventVerifiesDesktopOAuthWithoutCreatingAccountsOrRewards(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   int
		reason string
	}{
		{"valid", 200, `{"sub":"u_148","username":"desktop-user"}`, 200, ""},
		{"wrapped identity", 200, `{"data":{"sub":"u_148","username":"desktop-user"}}`, 200, ""},
		{"expired", 401, `{"error":"expired"}`, 401, "invalid_token"},
		{"forbidden", 403, `{}`, 401, "invalid_token"},
		{"offline", 503, `{}`, 502, "upstream_http_error"},
		{"missing identity", 200, `{"username":"desktop-user"}`, 502, "invalid_identity"},
		{"invalid response", 200, `{`, 502, "invalid_response_json"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/oauth/userinfo" || r.Header.Get("Authorization") != "Bearer desktop-token" {
					t.Error("desktop token was not verified through userinfo")
				}
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer upstream.Close()
			engine, mock := clientEventRouter(t, types.AppConfig{TokengineOAuth: types.TokengineOAuthConfig{
				Enabled: true, BaseURL: "https://provider.example", InternalBaseURL: upstream.URL,
			}})
			if tc.want == http.StatusOK {
				expectClientEventWrite(mock, "u_148")
			}
			w := postClientEvent(engine, "/v1/events", clientEventBody, "Bearer desktop-token")
			if w.Code != tc.want || strings.Contains(w.Body.String(), "desktop-token") {
				t.Fatalf("unexpected identity result: %d %s", w.Code, w.Body.String())
			}
			if tc.want == 401 {
				clientEventFailureResponse(t, w, "unauthorized", tc.reason)
			} else if tc.want == 502 {
				clientEventFailureResponse(t, w, "identity_unavailable", tc.reason)
			}
		})
	}
	engine, _ := clientEventRouter(t, types.AppConfig{})
	for _, token := range []string{"Bearer", "Bearer a b", "invalid-token"} {
		if w := postClientEvent(engine, "/api/v1/events", clientEventBody, token); w.Code != http.StatusUnauthorized {
			t.Fatalf("invalid token accepted: %d %s", w.Code, w.Body.String())
		}
	}
	if w := postClientEvent(engine, "/api/v1/events", clientEventBody, "Bearer desktop-token"); w.Code != http.StatusBadGateway {
		t.Fatalf("unconfigured verification accepted: %d", w.Code)
	} else {
		clientEventFailureResponse(t, w, "identity_unavailable", "identity_not_configured")
	}
}

func TestClientEventHidesDatabaseErrors(t *testing.T) {
	for _, tc := range []struct {
		err    error
		reason string
	}{
		{errors.New("private database credentials"), "database_error"},
		{&mysqldriver.MySQLError{Number: 1146, Message: "private database credentials"}, "database_table_missing"},
		{&mysqldriver.MySQLError{Number: 1062, Message: "private database credentials"}, "database_duplicate_key"},
		{&mysqldriver.MySQLError{Number: 1406, Message: "private database credentials"}, "database_value_too_long"},
		{mysqldriver.ErrInvalidConn, "database_connection_failed"},
		{context.DeadlineExceeded, "database_timeout"},
	} {
		t.Run(tc.reason, func(t *testing.T) {
			engine, mock := clientEventRouter(t, types.AppConfig{})
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `client_event`").WillReturnError(tc.err)
			mock.ExpectRollback()
			w := postClientEvent(engine, "/api/v1/events", clientEventBody, "")
			if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), "private database credentials") {
				t.Fatalf("database failure response: %d %s", w.Code, w.Body.String())
			}
			clientEventFailureResponse(t, w, "report_save_failed", tc.reason)
		})
	}
}

func TestClientEventFailuresIdentifyInvalidEventWithoutEchoingValues(t *testing.T) {
	engine, _ := clientEventRouter(t, types.AppConfig{})
	for _, tc := range []struct {
		name, body, code, reason, field string
		eventNumber                     float64
	}{
		{"empty body", "", "invalid_report_json", "empty_body", "", 0},
		{"incomplete JSON", `{"events":[`, "invalid_report_json", "incomplete_json", "", 0},
		{"multiple JSON values", clientEventBody + `{}`, "invalid_report_json", "multiple_json_values", "", 0},
		{"wrong type", strings.Replace(clientEventBody, `1792576514`, `"private-value"`, 1), "invalid_report_json", "json_type_mismatch", "ts", 1},
		{"second event wrong type", `{"events":[` + clientEventJSON + `,` + strings.Replace(clientEventJSON, `1792576514`, `"private-value"`, 1) + `]}`, "invalid_report_json", "json_type_mismatch", "ts", 2},
		{"empty batch", `{"events":[]}`, "invalid_report", "invalid_batch_size", "events", 0},
		{"second event timestamp", `{"events":[` + clientEventJSON + `,` + strings.Replace(clientEventJSON, `1792576514`, `0`, 1) + `]}`, "invalid_report", "invalid_timestamp", "ts", 2},
		{"missing device", strings.Replace(clientEventBody, `"device-1"`, `""`, 1), "invalid_report", "required_field", "machine_id", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := postClientEvent(engine, "/v1/events", tc.body, "")
			if w.Code != http.StatusBadRequest || strings.Contains(w.Body.String(), "private-value") || strings.Contains(w.Body.String(), "测试") {
				t.Fatalf("unsafe failure response: %d %s", w.Code, w.Body.String())
			}
			body := clientEventFailureResponse(t, w, tc.code, tc.reason)
			if tc.field != "" && body["field"] != tc.field {
				t.Fatalf("wrong field: %s", w.Body.String())
			}
			if tc.eventNumber > 0 && body["event_number"] != tc.eventNumber {
				t.Fatalf("wrong event number: %s", w.Body.String())
			}
		})
	}
}

func TestClientEventRequestIDIsGeneratedByServer(t *testing.T) {
	engine, _ := clientEventRouter(t, types.AppConfig{})
	ids := map[string]bool{}
	for range 2 {
		req := httptest.NewRequest(http.MethodPost, "/v1/events", strings.NewReader(`{}`))
		req.Header.Set("X-Request-ID", "client-controlled-id")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		body := clientEventFailureResponse(t, w, "invalid_report", "invalid_batch_size")
		id := body["request_id"].(string)
		if id == "client-controlled-id" || ids[id] {
			t.Fatal("request ID reused or controlled by client")
		}
		ids[id] = true
	}
}

func TestRealClientIPMatchesThinkBuddy(t *testing.T) {
	for _, tc := range []struct {
		name    string
		headers map[string]string
		remote  string
		want    string
	}{
		{"cloudflare first", map[string]string{"CF-Connecting-IP": "203.0.113.1", "True-Client-IP": "203.0.113.2", "X-Real-IP": "203.0.113.3", "X-Forwarded-For": "203.0.113.4, 10.0.0.1"}, "10.0.0.2:8080", "203.0.113.1"},
		{"true client before real ip", map[string]string{"CF-Connecting-IP": "invalid", "True-Client-IP": "203.0.113.2", "X-Real-IP": "203.0.113.3"}, "10.0.0.2:8080", "203.0.113.2"},
		{"real ip before forwarded", map[string]string{"X-Real-IP": "203.0.113.3", "X-Forwarded-For": "203.0.113.4, 10.0.0.1"}, "10.0.0.2:8080", "203.0.113.3"},
		{"first valid forwarded ip", map[string]string{"X-Forwarded-For": "invalid, 203.0.113.4, 10.0.0.1"}, "10.0.0.2:8080", "203.0.113.4"},
		{"quoted ipv6 with port", map[string]string{"X-Real-IP": `"[2001:db8::1]:443"`}, "10.0.0.2:8080", "2001:db8::1"},
		{"ipv4 with port", map[string]string{"X-Real-IP": "203.0.113.5:1234"}, "10.0.0.2:8080", "203.0.113.5"},
		{"whitespace around header", map[string]string{"X-Real-IP": " 203.0.113.5 "}, "10.0.0.2:8080", "203.0.113.5"},
		{"invalid headers fall back", map[string]string{"CF-Connecting-IP": "invalid", "True-Client-IP": "invalid", "X-Real-IP": "invalid", "X-Forwarded-For": "also-invalid"}, "192.0.2.5:1234", "192.0.2.5"},
		{"direct ipv6 connection", nil, "[2001:db8::2]:1234", "2001:db8::2"},
		{"bare remote ip", nil, "192.0.2.5", "192.0.2.5"},
		{"invalid connection", nil, "invalid", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("POST", "/v1/events", nil)
			c.Request.RemoteAddr = tc.remote
			for name, value := range tc.headers {
				c.Request.Header.Set(name, value)
			}
			if got := getRealClientIP(c); got != tc.want {
				t.Fatalf("got IP %q, want %q", got, tc.want)
			}
		})
	}
	if got := getRealClientIP(nil); got != "" {
		t.Fatalf("nil context returned IP %q", got)
	}
	if got := getRealClientIP(&gin.Context{}); got != "" {
		t.Fatalf("missing request returned IP %q", got)
	}
}

func TestClientEventDiagnosticsCorrelateLogsWithoutSensitiveContents(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	previous := logger
	logger = zap.New(core).Sugar()
	t.Cleanup(func() { logger = previous })
	for _, tc := range []struct {
		name, stage, code, reason string
		status                    int
		err                       error
	}{
		{"field", "validate", "invalid_report", "invalid_timestamp", 400, &service.ClientEventValidationError{EventNumber: 2, Field: "ts", Reason: "invalid_timestamp"}},
		{"database", "save", "report_save_failed", "database_table_missing", 500, &mysqldriver.MySQLError{Number: 1146, Message: "private-secret SQL contents"}},
		{"unknown database", "save", "report_save_failed", "database_error", 500, errors.New("private-secret database password")},
		{"upstream", "authenticate", "identity_unavailable", "upstream_http_error", 502, &service.TokengineOAuthError{Reason: "upstream_http_error", StatusCode: 503}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logs.TakeAll()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/v1/events", strings.NewReader("private-secret body"))
			c.Request.Header.Set("Authorization", "Bearer private-secret")
			c.Header("X-Request-ID", "request-1")
			c.Set(clientEventBatchKey, []service.ClientEventInput{
				{MachineID: "first-device", SessionID: "first-session"},
				{MachineID: "second-device", SessionID: "second-session", Props: json.RawMessage(`{"message":"private-secret"}`)},
			})
			rejectClientEvent(c, tc.status, tc.code, tc.stage, tc.err)
			entries := logs.All()
			if len(entries) != 1 {
				t.Fatalf("expected one diagnostic log, got %d", len(entries))
			}
			fields := entries[0].ContextMap()
			if fields["request_id"] != "request-1" || fields["reason"] != tc.reason || fields["stage"] != tc.stage {
				t.Fatalf("missing diagnostic context: %+v", fields)
			}
			wantDevice := "first-device"
			if tc.stage == "validate" {
				wantDevice = "second-device"
				if fields["event_number"] != int64(2) || fields["field"] != "ts" {
					t.Fatalf("missing invalid event details: %+v", fields)
				}
			}
			if fields["machine_id"] != wantDevice {
				t.Fatalf("wrong correlation metadata: %+v", fields)
			}
			if tc.name == "database" && fields["database_code"] != uint16(1146) {
				t.Fatalf("missing database code: %+v", fields)
			}
			if tc.name == "upstream" && fields["upstream_status"] != int64(503) {
				t.Fatalf("missing upstream status: %+v", fields)
			}
			encoded, err := json.Marshal(fields)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(encoded), "private-secret") || strings.Contains(w.Body.String(), "private-secret") {
				t.Fatal("diagnostics leaked request contents or raw errors")
			}
		})
	}
}
