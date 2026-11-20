package handler_test

// 本文件是"鉴权单点收敛"的运行时回归：
// 它用真实 handler + 真实 Router 组装路由表，断言改造前漏挂鉴权的 7 个 admin 组已经 401，
// 同时断言 login / logout 仍然公开 —— 这两件事是本次重构里最容易被搞反的一对。

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/handler/admin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func adminTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	app := core.NewServerWithEngine(engine, &types.AppConfig{}, nil)
	app.SysConfig = &types.SystemConfig{}
	// 指向一个必定不可达的端口且禁用重试：handler 内部的 redis 调用会快速失败并走错误分支，
	// 从而在"没有真实 Redis"的前提下也能走完 login / logout 的完整路由链。
	deadRedis := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", MaxRetries: -1})

	// DB 为 nil 是安全的：未鉴权请求会在进入 handler 之前被中间件 401 掉。
	admin.NewAdminHandler(app, nil, deadRedis, nil).RegisterRoutes(app.Router)
	admin.NewDashboardHandler(app, nil).RegisterRoutes(app.Router)
	admin.NewOrderHandler(app, nil).RegisterRoutes(app.Router)
	admin.NewProductHandler(app, nil).RegisterRoutes(app.Router)
	admin.NewRedeemHandler(app, nil).RegisterRoutes(app.Router)
	return engine
}

func adminCall(engine *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

// TestAdminLoginLogoutStayPublic 对应验收 F 组：管理员入口本身绝不能被鉴权拦住。
func TestAdminLoginLogoutStayPublic(t *testing.T) {
	engine := adminTestEngine(t)

	// 非法 body 让 Login 在参数校验阶段就返回 400，不触碰 DB。重点不是状态码本身，
	// 而是它绝不可能是 401（401 意味着路由被误收进鉴权组）。
	if w := adminCall(engine, http.MethodPost, "/api/admin/login", `not-json`); w.Code != http.StatusBadRequest {
		t.Fatalf("POST /api/admin/login = %d, 期望 400；若为 401 说明管理后台入口被误收进鉴权组", w.Code)
	}
	if w := adminCall(engine, http.MethodGet, "/api/admin/logout", ""); w.Code == http.StatusUnauthorized {
		t.Fatalf("GET /api/admin/logout 返回 401：登出被误收进鉴权组")
	}
}

// TestAdminProtectedGroupsNowReturn401 对应验收 C 组 + C-破坏性：
// 改造前这 7 个组没有任何鉴权中间件，其中 order/clear 能清空订单表、redeem/create 能匿名造钱。
func TestAdminProtectedGroupsNowReturn401(t *testing.T) {
	engine := adminTestEngine(t)
	cases := []struct{ method, path string }{
		{http.MethodGet, "/api/admin/session"},
		{http.MethodGet, "/api/admin/dashboard/stats"},
		{http.MethodPost, "/api/admin/order/list"},
		{http.MethodGet, "/api/admin/order/clear"}, // 破坏性：清空订单表
		{http.MethodGet, "/api/admin/product/list"},
		{http.MethodPost, "/api/admin/product/save"},
		{http.MethodGet, "/api/admin/redeem/list"},
		{http.MethodPost, "/api/admin/redeem/create"}, // 破坏性：匿名造钱
		{http.MethodPost, "/api/admin/redeem/export"},
	}
	for _, tc := range cases {
		if w := adminCall(engine, tc.method, tc.path, `{}`); w.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s = %d, 期望 401（该路由没有落在管理员鉴权链上）", tc.method, tc.path, w.Code)
		}
	}
}
