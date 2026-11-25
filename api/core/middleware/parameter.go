package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// skipPathPrefixes 是"消息原文必须逐字保留"的路径，既不 trim query 也不改写 body。
//   - 资料库接口自行限制并解析请求体。
var skipPathPrefixes = []string{"/api/library", "/api/v1/library", "/api/admin/library"}

// ParameterHandlerMiddleware 统一参数处理。
//
// 相对旧实现修正 4 点，其余语义（query trim / 跳过清单 / multipart 跳过）保持不变：
//  1. JSON body 解析改用 json.Decoder + UseNumber()：杜绝 float64 中转导致的 int64 精度丢失
//     （旧实现把 1234567890123456789 写成 1234567890123456800）；
//  2. 数字规范化：1.0 → 1、2e3 → 2000，与旧 float64 行为保持一致，避免把整数值浮点写法弄坏；
//  3. 解析失败还原**原始字节**（旧实现里 ShouldBindJSON 失败后 body 已被读空，handler 拿到空 body）；
//  4. 回写 body 后同步 ContentLength（旧实现未同步，下游按旧长度读取可能截断或多读）。
//
// 刻意保留：JSON 字符串值两端 TrimSpace。去掉它属于用户可见的行为变更（例如聊天框里粘贴的
// 缩进代码首尾空白会被吞），需要产品与前端一起确认，本轮不动。
func ParameterHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过清单必须放在最前面：旧实现对这些路径连 query 都不 trim。
		if shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		trimQueryParams(c)

		if !shouldRewriteJSONBody(c) || c.Request.Body == nil {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		rewritten, ok := trimJSONBody(body)
		if !ok { // 解析失败：原样还原原始字节，不吞 body
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(rewritten))
		c.Request.ContentLength = int64(len(rewritten))
		c.Next()
	}
}

// shouldSkipPath 判断路径是否属于"原文保留"清单。
func shouldSkipPath(path string) bool {
	// 客户端事件自行限制请求体大小，并保留事件属性中的原始 JSON。
	if path == "/api/events" || path == "/api/v1/events" || path == "/v1/events" {
		return true
	}
	if path == "/api/checkin" || path == "/api/v1/checkin" {
		return true // 签到接口自行限制请求体大小。
	}
	for _, prefix := range skipPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// trimQueryParams 对 GET 参数做两端去空格，与原实现逐字一致。
func trimQueryParams(c *gin.Context) {
	params := c.Request.URL.Query()
	for key, values := range params {
		for i, value := range values {
			params[key][i] = strings.TrimSpace(value)
		}
	}
	c.Request.URL.RawQuery = params.Encode()
}

// shouldRewriteJSONBody 判断是否需要改写请求体：只处理 application/json，跳过文件上传。
func shouldRewriteJSONBody(c *gin.Context) bool {
	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		return false
	}
	return strings.Contains(contentType, "application/json")
}

// trimJSONBody 解析、规范化并重新序列化 JSON 对象。
// 返回 ok=false 表示"不该改 / 改不了"，调用方必须还原原始字节。
func trimJSONBody(body []byte) ([]byte, bool) {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var data any
	if err := dec.Decode(&data); err != nil {
		return nil, false
	}
	if _, isObject := data.(map[string]any); !isObject {
		return nil, false // 只处理 JSON 对象，数组/标量原样透传
	}
	out, err := json.Marshal(sanitizeJSON(data))
	if err != nil {
		return nil, false
	}
	return out, true
}

// sanitizeJSON 一次遍历同时完成两件事（两者互不影响，故合并以省一次递归）：
//  1. 字符串值 TrimSpace —— 保留旧语义；
//  2. json.Number 规范化 —— 让 1.0 → 1、2e3 → 2000 与旧 float64 行为一致，
//     同时靠 Int64() 保住 >2^53 的整数精度（旧实现中转 float64 会丢精度）。
//
// 规范化结果仍然是 json.Number，因此 json.Marshal 会原样输出数字字面量，不会被二次格式化。
func sanitizeJSON(data any) any {
	switch v := data.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		// 整数（含超出 float64 精度范围的大整数）直接原样保真。
		if i, err := v.Int64(); err == nil {
			return json.Number(strconv.FormatInt(i, 10))
		}
		// 整数值的浮点写法（1.0 / 2e3）归一成整数字面量，保持对 int 字段的可绑定性。
		if f, err := v.Float64(); err == nil && f == math.Trunc(f) && math.Abs(f) < 1e18 {
			return json.Number(strconv.FormatInt(int64(f), 10))
		}
		return v
	case map[string]any:
		for key, value := range v {
			v[key] = sanitizeJSON(value)
		}
		return v
	case []any:
		for i, value := range v {
			v[i] = sanitizeJSON(value)
		}
		return v
	default:
		return data
	}
}
