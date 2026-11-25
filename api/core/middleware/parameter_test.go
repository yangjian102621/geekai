package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestTrimJSONBody 是参数中间件的精度回归。
// 判据是二值的：既能修掉 float64 中转的精度丢失，又不能让"整数值浮点写法"绑定不进 int 字段。
func TestTrimJSONBody(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{"int64 保真", `{"id":1234567890123456789}`, `{"id":1234567890123456789}`, true},
		{"2^53+1 保真", `{"id":9007199254740993}`, `{"id":9007199254740993}`, true},
		{"整数值浮点写法归一", `{"f":1.0}`, `{"f":1}`, true},
		{"科学计数法归一", `{"n":2e3}`, `{"n":2000}`, true},
		{"真小数保持", `{"p":9.99}`, `{"p":9.99}`, true},
		{"字符串去两端空格", `{"s":"  hi  "}`, `{"s":"hi"}`, true},
		{"嵌套数组内字符串与数字", `{"a":["  x  ",1.0]}`, `{"a":["x",1]}`, true},
		{"非法 JSON 不改写", `{bad`, ``, false},
		{"顶层数组原样透传", `[1.0,2]`, ``, false},
		{"顶层标量原样透传", `123`, ``, false},
	}
	for _, tc := range cases {
		got, ok := trimJSONBody([]byte(tc.in))
		if ok != tc.ok {
			t.Fatalf("%s: ok = %v, 期望 %v", tc.name, ok, tc.ok)
		}
		if !ok {
			continue
		}
		if string(got) != tc.want {
			t.Fatalf("%s: trimJSONBody(%s) = %s, 期望 %s", tc.name, tc.in, got, tc.want)
		}
	}
}

// TestParameterMiddlewareRestoresOriginalBodyOnFailure 断言解析失败时 handler 能看到原始字节（而不是空 body）。
func TestParameterMiddlewareRestoresOriginalBodyOnFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ParameterHandlerMiddleware())

	const broken = `{"a": bad`
	var seen string
	engine.POST("/api/echo", func(c *gin.Context) {
		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Errorf("读取 body 失败: %v", err)
		}
		seen = string(raw)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/echo", strings.NewReader(broken))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(httptest.NewRecorder(), req)

	if seen != broken {
		t.Fatalf("handler 读到的 body = %q, 期望原始字节 %q", seen, broken)
	}
}

// TestParameterMiddlewareSyncsContentLength 断言回写 body 后 ContentLength 同步（旧实现未同步）。
func TestParameterMiddlewareSyncsContentLength(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ParameterHandlerMiddleware())

	const body = `{"s":"  hi  ","id":1234567890123456789}`
	var gotBody string
	var gotLen int64
	engine.POST("/api/echo", func(c *gin.Context) {
		raw, _ := io.ReadAll(c.Request.Body)
		gotBody = string(raw)
		gotLen = c.Request.ContentLength
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(httptest.NewRecorder(), req)

	want := `{"id":1234567890123456789,"s":"hi"}`
	if gotBody != want {
		t.Fatalf("body = %q, 期望 %q", gotBody, want)
	}
	if gotLen != int64(len(want)) {
		t.Fatalf("ContentLength = %d, 期望 %d（必须与回写后的 body 长度一致）", gotLen, len(want))
	}
}

// TestParameterMiddlewareSkipsMessageOriginalPaths 断言跳过清单里的路径连 query 都不 trim。
func TestParameterMiddlewareSkipsMessageOriginalPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ParameterHandlerMiddleware())

	const body = `{"content":"  keep me  "}`
	for _, path := range []string{
		"/api/events", "/api/v1/events", "/v1/events",
		"/api/checkin", "/api/v1/checkin",
		"/api/library", "/api/v1/library", "/api/admin/library/sources",
	} {
		var gotQuery, gotBody string
		engine.POST(path, func(c *gin.Context) {
			raw, _ := io.ReadAll(c.Request.Body)
			gotBody = string(raw)
			gotQuery = c.Request.URL.RawQuery
			c.Status(http.StatusOK)
		})
		req := httptest.NewRequest(http.MethodPost, path+"?q=%20x%20", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		engine.ServeHTTP(httptest.NewRecorder(), req)

		if gotBody != body {
			t.Fatalf("%s: body 被改写了: %q", path, gotBody)
		}
		if gotQuery != "q=%20x%20" {
			t.Fatalf("%s: query 被 trim 了: %q", path, gotQuery)
		}
	}
}
