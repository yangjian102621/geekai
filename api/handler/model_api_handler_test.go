package handler

import (
	"database/sql/driver"
	"encoding/json"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func modelAPITestHandler(t *testing.T) (*ModelAPIHandler, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	gin.SetMode(gin.TestMode)
	app := &core.AppServer{Engine: gin.New(), Config: &types.AppConfig{}}
	app.Engine.Use(middleware.ParameterHandlerMiddleware())
	return NewModelAPIHandler(app, db, service.NewUserService(db)), mock
}

func modelAPIRequest(engine *gin.Engine, method, path, token, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func expectModelAPIAuth(mock sqlmock.Sqlmock, key string, status bool, expiry int64) {
	mock.ExpectQuery("SELECT .*thinkbuddy_model_tokens.*token_hash").WithArgs(modelTokenHash(key)).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "token_hash"}).AddRow(7, modelTokenHash(key)))
	mock.ExpectQuery("SELECT .*geekai_users").WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "status", "expired_time", "chat_models_json", "power"}).
			AddRow(7, "alice", status, expiry, "[8]", 10))
}

func TestModelAPIAuthenticationAndCatalog(t *testing.T) {
	h, mock := modelAPITestHandler(t)
	h.RegisterRoutes()
	for _, key := range []string{"", "raw-site-jwt", "Bearer upstream-key", "Bearer sk-thinkbuddy-short"} {
		w := modelAPIRequest(h.App.Engine, "GET", "/v1/models", key, "")
		if w.Code != 401 || !strings.Contains(w.Body.String(), `"error"`) {
			t.Fatalf("invalid credentials accepted: %d %s", w.Code, w.Body.String())
		}
	}
	key := modelTokenPrefix + strings.Repeat("a", 43)
	mock.ExpectQuery("SELECT .*thinkbuddy_model_tokens.*token_hash").WithArgs(modelTokenHash(key)).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}))
	if w := modelAPIRequest(h.App.Engine, "GET", "/v1/models", "Bearer "+key, ""); w.Code != 401 {
		t.Fatalf("revoked token accepted: %d", w.Code)
	}
	for _, account := range []struct {
		status bool
		expiry int64
	}{{false, 0}, {true, time.Now().Unix() - 1}} {
		expectModelAPIAuth(mock, key, account.status, account.expiry)
		if w := modelAPIRequest(h.App.Engine, "GET", "/v1/models", "Bearer "+key, ""); w.Code != 403 {
			t.Fatalf("unavailable account accepted: %d", w.Code)
		}
	}
	expectModelAPIAuth(mock, key, true, 0)
	mock.ExpectQuery("SELECT .*geekai_chat_models.*type = \\? AND enabled = \\?.*\\(open = \\? OR id IN \\(\\?\\)\\)").
		WithArgs("chat", true, true, 8).
		WillReturnRows(sqlmock.NewRows([]string{"id", "value", "created_at", "key_id"}).
			AddRow(1, "test-model", time.Unix(10, 0), 99).AddRow(8, "test-model", time.Unix(10, 0), 99))
	w := modelAPIRequest(h.App.Engine, "GET", "/v1/models", "Bearer "+key, "")
	var catalog struct {
		Object string
		Data   []struct {
			ID      string `json:"id"`
			OwnedBy string `json:"owned_by"`
		}
	}
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &catalog) != nil || catalog.Object != "list" || len(catalog.Data) != 1 || catalog.Data[0].ID != "test-model" || catalog.Data[0].OwnedBy != "thinkbuddy" {
		t.Fatalf("bad standard catalog: %d %s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), "key_id") || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("catalog leaked internal fields or was cacheable: %s", w.Body.String())
	}
}

type tokenHashArgument struct{ value string }

func (a *tokenHashArgument) Match(value driver.Value) bool {
	a.value, _ = value.(string)
	return len(a.value) == 64 && !strings.Contains(a.value, "sk-")
}

func TestModelTokenRotationAndRevocation(t *testing.T) {
	h, mock := modelAPITestHandler(t)
	h.RegisterRoutes()
	for _, method := range []string{"GET", "POST", "DELETE"} {
		if w := modelAPIRequest(h.App.Engine, method, "/api/user/model-token", "", ""); w.Code != 401 {
			t.Fatalf("unprotected token management: %d", w.Code)
		}
	}
	// 管理操作从已认证的用户上下文取身份，忽略请求中的 user_id。
	account := h.App.Engine.Group("/test-token", func(c *gin.Context) {
		c.Set(types.LoginUserCache, model.User{Id: 7, Status: true})
	})
	account.POST("", h.RotateToken)
	account.DELETE("", h.RevokeToken)
	account.GET("", h.TokenStatus)
	for i := 0; i < 2; i++ {
		digest := &tokenHashArgument{}
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `thinkbuddy_model_tokens`.*ON DUPLICATE KEY UPDATE").
			WithArgs(7, digest, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		w := modelAPIRequest(h.App.Engine, "POST", "/test-token", "", `{"user_id":999}`)
		var result struct {
			Code int
			Data struct{ Token string }
		}
		if json.Unmarshal(w.Body.Bytes(), &result) != nil || result.Code != 0 || len(result.Data.Token) != len(modelTokenPrefix)+43 || modelTokenHash(result.Data.Token) != digest.value {
			t.Fatalf("token not generated or stored as hash: %s", w.Body.String())
		}
	}
	mock.ExpectQuery("SELECT .*thinkbuddy_model_tokens.*user_id").WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "token_hash", "created_at"}).AddRow(7, "private-hash", time.Now()))
	if w := modelAPIRequest(h.App.Engine, "GET", "/test-token", "", ""); !strings.Contains(w.Body.String(), `"enabled":true`) || strings.Contains(w.Body.String(), "private-hash") {
		t.Fatalf("unsafe token status: %s", w.Body.String())
	}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `thinkbuddy_model_tokens` WHERE user_id = \\?").WithArgs(7).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if w := modelAPIRequest(h.App.Engine, "DELETE", "/test-token", "", ""); !strings.Contains(w.Body.String(), `"code":0`) {
		t.Fatalf("revoke failed: %s", w.Body.String())
	}
}
