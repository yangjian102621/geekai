package core

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/types"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// newTestRouter 构造一个不依赖 DB / Redis 的 Router。
// Redis 传 nil 是安全的：鉴权中间件只在请求期访问 redis；无 token 的请求在解析阶段就返回 401。
func newTestRouter(t *testing.T) (*Router, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	cfg := &types.AppConfig{}
	cfg.Session.SecretKey = "test-user-secret"
	cfg.AdminSession.SecretKey = "test-admin-secret"
	return NewRouter(engine, cfg, nil), engine
}

// TestRouterHandlerChainLength 断言四个 group 在建组时就已经固化好安全级别。
// B1: admin 组只有 AdminAuthMiddleware 一个 handler。
// B2: user 组只有 UserAuthMiddleware 一个 handler。
// B3: pub 组没有任何 handler（完全公开）。
func TestRouterHandlerChainLength(t *testing.T) {
	r, _ := newTestRouter(t)

	if got, want := len(r.admin.Handlers), 1; got != want {
		t.Fatalf("admin 组 handler 数 = %d, 期望 %d（仅管理员鉴权中间件）", got, want)
	}
	if got, want := len(r.user.Handlers), 1; got != want {
		t.Fatalf("user 组 handler 数 = %d, 期望 %d（仅用户鉴权中间件）", got, want)
	}
	if got, want := len(r.pub.Handlers), 0; got != want {
		t.Fatalf("pub 组 handler 数 = %d, 期望 %d（完全公开）", got, want)
	}
	// adminPub 是 /api/admin 下的公开出口，同样不能带鉴权中间件。
	if got, want := len(r.adminPub.Handlers), 0; got != want {
		t.Fatalf("adminPub 组 handler 数 = %d, 期望 %d（login/logout 必须公开）", got, want)
	}
}

// TestRouterBasePaths 断言前缀拼接结果与改造前逐字一致（搬迁时路径零变化）。
// B4: AdminPublicGroup().BasePath() == "/api/admin"
// B7: AdminGroup("/api/admin/redeem") == "/api/admin/redeem"
func TestRouterBasePaths(t *testing.T) {
	r, _ := newTestRouter(t)

	if got, want := r.AdminPublicGroup().BasePath(), "/api/admin"; got != want {
		t.Fatalf("AdminPublicGroup().BasePath() = %q, 期望 %q", got, want)
	}
	cases := []struct {
		name   string
		prefix string
		want   string
	}{
		{"无尾斜杠", "/api/admin/redeem", "/api/admin/redeem"},
		{"带尾斜杠", "/api/admin/redeem/", "/api/admin/redeem/"},
		{"仅基址", "/api/admin", "/api/admin"},
		{"相对前缀也应可用", "/redeem/", "/api/admin/redeem/"},
	}
	for _, tc := range cases {
		if got := r.AdminGroup(tc.prefix).BasePath(); got != tc.want {
			t.Fatalf("AdminGroup(%q).BasePath() = %q, 期望 %q（%s）", tc.prefix, got, tc.want, tc.name)
		}
	}
	// 拼错的路径绝不可能是 "/api/admin/api/admin/..."
	for _, tc := range cases {
		if got := r.AdminGroup(tc.prefix).BasePath(); strings.Contains(got, "/api/admin/api/admin") {
			t.Fatalf("AdminGroup(%q) 前缀被重复拼接: %q", tc.prefix, got)
		}
	}
}

// TestPublicGroupRejectsAdminPrefix 断言"管理后台命名空间只能由 AdminGroup / AdminPublicGroup 进入"。
// B5 / B6: 试图把 /api/admin 前缀当公开或普通用户路由注册，启动期直接 panic。
func TestPublicGroupRejectsAdminPrefix(t *testing.T) {
	r, _ := newTestRouter(t)

	assertPanic := func(name string, fn func()) {
		t.Helper()
		defer func() {
			rec := recover()
			if rec == nil {
				t.Fatalf("%s 未 panic：/api/admin 前缀被允许了", name)
			}
			msg, ok := rec.(string)
			if !ok || !strings.Contains(msg, "/api/admin") {
				t.Fatalf("%s panic 文案未包含 /api/admin: %v", name, rec)
			}
		}()
		fn()
	}

	assertPanic("PublicGroup(\"/api/admin/x\")", func() { r.PublicGroup("/api/admin/x") })
	assertPanic("PublicGroup(\"/api/admin\")", func() { r.PublicGroup("/api/admin") })
	assertPanic("UserGroup(\"/api/admin/x\")", func() { r.UserGroup("/api/admin/x") })

	// 反向断言：合法的公开前缀不能被误伤。
	if got := r.PublicGroup("/api/public-test/").BasePath(); got != "/api/public-test/" {
		t.Fatalf("PublicGroup 合法前缀被误改: %q", got)
	}
}

// TestStaticFilesInheritsPublicChain 断言静态目录挂且仅挂在公开链上，并保留额外中间件。
func TestStaticFilesInheritsPublicChain(t *testing.T) {
	r, engine := newTestRouter(t)
	extra := func(c *gin.Context) { c.Next() }
	r.StaticFiles("/update", t.TempDir(), extra)

	found := false
	for _, route := range engine.Routes() {
		if route.Path == "/update/*filepath" {
			found = true
		}
	}
	if !found {
		t.Fatalf("StaticFiles 未注册 /update/*filepath，实际路由表: %+v", engine.Routes())
	}
}

// TestRouterExposesNoEngineGetter 是"不给 *gin.Engine 提供任何 getter"这条约束的可执行版本。
// 任何人将来给 Router 加回 Engine() / Gin() / GetEngine()，或者把某个字段导出，都会被这条挡住。
func TestRouterExposesNoEngineGetter(t *testing.T) {
	routerType := reflect.TypeOf(Router{})
	engineType := reflect.TypeOf(&gin.Engine{})

	for i := 0; i < routerType.NumField(); i++ {
		field := routerType.Field(i)
		if field.PkgPath == "" {
			t.Fatalf("Router 字段 %s 被导出：handler 层又能直接拿到路由对象了", field.Name)
		}
	}
	for i := 0; i < routerType.NumMethod(); i++ {
		method := routerType.Method(i)
		for j := 0; j < method.Type.NumOut(); j++ {
			if method.Type.Out(j) == engineType {
				t.Fatalf("Router 方法 %s 返回了 *gin.Engine：等于重新打开了漏挂鉴权的口子", method.Name)
			}
		}
	}
}
