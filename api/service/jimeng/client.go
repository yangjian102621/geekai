package jimeng

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"geekai/logger"
)

var clientLogger = logger.GetLogger()

// Client 即梦API客户端
type Client struct {
	accessKey  string
	secretKey  string
	region     string
	service    string
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建即梦API客户端
func NewClient(accessKey, secretKey string) *Client {
	return &Client{
		accessKey: accessKey,
		secretKey: secretKey,
		region:    "cn-north-1",
		service:   "cv",
		baseURL:   "https://visual.volcengineapi.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SubmitTask 提交任务
func (c *Client) SubmitTask(req *SubmitTaskRequest) (*SubmitTaskResponse, error) {
	// 构建请求URL
	queryParams := map[string]string{
		"Action":  "CVSync2AsyncSubmitTask",
		"Version": "2022-08-31",
	}

	reqURL := c.buildURL(queryParams)

	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create http request failed: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 签名请求
	if err := c.signRequest(httpReq, reqBody); err != nil {
		return nil, fmt.Errorf("sign request failed: %w", err)
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send http request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	clientLogger.Infof("Jimeng SubmitTask Response: %s", string(respBody))

	// 解析响应
	var result SubmitTaskResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &result, nil
}

// QueryTask 查询任务
func (c *Client) QueryTask(req *QueryTaskRequest) (*QueryTaskResponse, error) {
	// 构建请求URL
	queryParams := map[string]string{
		"Action":  "CVSync2AsyncGetResult",
		"Version": "2022-08-31",
	}

	reqURL := c.buildURL(queryParams)

	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create http request failed: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 签名请求
	if err := c.signRequest(httpReq, reqBody); err != nil {
		return nil, fmt.Errorf("sign request failed: %w", err)
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send http request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	clientLogger.Infof("Jimeng QueryTask Response: %s", string(respBody))

	// 解析响应
	var result QueryTaskResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &result, nil
}

// SubmitSyncTask 提交同步任务（仅用于文生图）
func (c *Client) SubmitSyncTask(req *SubmitTaskRequest) (*QueryTaskResponse, error) {
	// 构建请求URL
	queryParams := map[string]string{
		"Action":  "CVProcess",
		"Version": "2022-08-31",
	}

	reqURL := c.buildURL(queryParams)

	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create http request failed: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 签名请求
	if err := c.signRequest(httpReq, reqBody); err != nil {
		return nil, fmt.Errorf("sign request failed: %w", err)
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send http request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	clientLogger.Infof("Jimeng SubmitSyncTask Response: %s", string(respBody))

	// 解析响应
	var result QueryTaskResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &result, nil
}

// buildURL 构建请求URL
func (c *Client) buildURL(queryParams map[string]string) string {
	u, _ := url.Parse(c.baseURL)
	q := u.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// signRequest 签名请求
func (c *Client) signRequest(req *http.Request, body []byte) error {
	now := time.Now().UTC()

	// 设置基本头部
	req.Header.Set("X-Date", now.Format("20060102T150405Z"))
	req.Header.Set("Host", req.URL.Host)

	// 计算内容哈希
	contentHash := sha256.Sum256(body)
	req.Header.Set("X-Content-Sha256", hex.EncodeToString(contentHash[:]))

	// 构建签名字符串
	canonicalRequest := c.buildCanonicalRequest(req)
	credentialScope := fmt.Sprintf("%s/%s/%s/request", now.Format("20060102"), c.region, c.service)
	stringToSign := fmt.Sprintf("HMAC-SHA256\n%s\n%s\n%s",
		now.Format("20060102T150405Z"), credentialScope, sha256Hash(canonicalRequest))

	// 计算签名
	signature := c.calculateSignature(stringToSign, now)

	// 设置Authorization头部
	authorization := fmt.Sprintf("HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		c.accessKey, credentialScope, c.getSignedHeaders(req), signature)
	req.Header.Set("Authorization", authorization)

	return nil
}

// buildCanonicalRequest 构建规范请求
func (c *Client) buildCanonicalRequest(req *http.Request) string {
	// HTTP方法
	method := req.Method

	// 规范URI
	uri := req.URL.Path
	if uri == "" {
		uri = "/"
	}

	// 规范查询字符串
	query := req.URL.Query()
	var queryParts []string
	for k, v := range query {
		for _, val := range v {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(val)))
		}
	}
	sort.Strings(queryParts)
	canonicalQuery := strings.Join(queryParts, "&")

	// 规范头部
	var headerParts []string
	headers := make(map[string]string)
	for k, v := range req.Header {
		key := strings.ToLower(k)
		if len(v) > 0 {
			headers[key] = strings.TrimSpace(v[0])
		}
	}

	var headerKeys []string
	for k := range headers {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)

	for _, k := range headerKeys {
		headerParts = append(headerParts, fmt.Sprintf("%s:%s", k, headers[k]))
	}
	canonicalHeaders := strings.Join(headerParts, "\n") + "\n"

	// 签名头部
	signedHeaders := c.getSignedHeaders(req)

	// 载荷哈希
	payloadHash := req.Header.Get("X-Content-Sha256")

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method, uri, canonicalQuery, canonicalHeaders, signedHeaders, payloadHash)
}

// getSignedHeaders 获取签名头部
func (c *Client) getSignedHeaders(req *http.Request) string {
	var headers []string
	for k := range req.Header {
		headers = append(headers, strings.ToLower(k))
	}
	sort.Strings(headers)
	return strings.Join(headers, ";")
}

// calculateSignature 计算签名
func (c *Client) calculateSignature(stringToSign string, t time.Time) string {
	kDate := hmacSha256([]byte("HMAC-SHA256"+c.secretKey), []byte(t.Format("20060102")))
	kRegion := hmacSha256(kDate, []byte(c.region))
	kService := hmacSha256(kRegion, []byte(c.service))
	kSigning := hmacSha256(kService, []byte("request"))
	signature := hmacSha256(kSigning, []byte(stringToSign))
	return hex.EncodeToString(signature)
}

// hmacSha256 计算HMAC-SHA256
func hmacSha256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// sha256Hash 计算SHA256哈希
func sha256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
