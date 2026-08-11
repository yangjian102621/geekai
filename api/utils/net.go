package utils

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/log"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

var logger = log.GetLogger()

const (
	defaultHTTPTimeout          = 45 * time.Second
	defaultTLSHandshakeTimeout = 10 * time.Second
	defaultDialTimeout          = 10 * time.Second
	defaultResponseHeaderTimeout = 30 * time.Second
	defaultIdleConnTimeout         = 90 * time.Second
	defaultMaxIdleConns            = 100
	defaultMaxIdleConnsPerHost    = 10
)

func newHTTPClient(timeout time.Duration, proxy string) *http.Client {
	if timeout <= 0 {
		timeout = defaultHTTPTimeout
	}

	var proxyFn func(*http.Request) (*url.URL, error)
	if strings.TrimSpace(proxy) != "" {
		if proxyURL, err := url.Parse(proxy); err == nil && proxyURL != nil {
			proxyFn = http.ProxyURL(proxyURL)
		}
	}

	transport := &http.Transport{
		Proxy: proxyFn,
		DialContext: (&net.Dialer{
			Timeout:   defaultDialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ResponseHeaderTimeout: defaultResponseHeaderTimeout,
		IdleConnTimeout:       defaultIdleConnTimeout,
		MaxIdleConns:          defaultMaxIdleConns,
		MaxIdleConnsPerHost:  defaultMaxIdleConnsPerHost,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var ne net.Error
	if errors.As(err, &ne) {
		if ne.Timeout() || ne.Temporary() {
			return true
		}
	}

	// 兼容不同 Go/系统对握手超时的错误文本
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "tls handshake timeout") ||
		strings.Contains(msg, "handshake timeout") ||
		strings.Contains(msg, "connection reset by peer") ||
		strings.Contains(msg, "connection timed out")
}

func retryDelay(attempt int) time.Duration {
	// attempt 从 0 开始：1s, 2s, 4s...
	d := time.Second * time.Duration(1<<attempt)
	if d > 8*time.Second {
		d = 8 * time.Second
	}
	return d
}

func readAllLimit(r io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes > 0 {
		return io.ReadAll(io.LimitReader(r, maxBytes))
	}
	return io.ReadAll(r)
}

// FetchURLBytes 发起 GET 并返回响应体（只在 2xx 认为成功）。
// 为了降低“偶发 HTTPS 握手超时”，对可重试错误/状态码会做有限重试。
func FetchURLBytes(ctx context.Context, rawURL string, proxy string, timeout time.Duration, retries int, maxBytes int64) ([]byte, int, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		timeout = defaultHTTPTimeout
	}
	if retries < 0 {
		retries = 0
	}
	if maxBytes <= 0 {
		maxBytes = 8 << 20 // 默认最多读取 8MiB
	}

	var lastErr error
	for attempt := 0; attempt <= retries; attempt++ {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}

		client := newHTTPClient(timeout, proxy)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, 0, err
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < retries && isRetryableError(err) {
				time.Sleep(retryDelay(attempt))
				continue
			}
			return nil, 0, err
		}

		body, readErr := readAllLimit(resp.Body, maxBytes)
		_ = resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			if attempt < retries && isRetryableError(readErr) {
				time.Sleep(retryDelay(attempt))
				continue
			}
			return nil, resp.StatusCode, readErr
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			// 429/5xx 可能是短暂错误，允许有限重试
			status := resp.StatusCode
			lastErr = fmt.Errorf("request failed: status=%d", status)
			if attempt < retries && (status == http.StatusTooManyRequests || status >= 500 && status <= 599) {
				time.Sleep(retryDelay(attempt))
				continue
			}

			// 避免把大 body 全塞进日志/错误里
			preview := strings.TrimSpace(string(body))
			if len(preview) > 256 {
				preview = preview[:256]
			}
			if preview != "" {
				return body, status, fmt.Errorf("request failed: status=%d body=%s", status, preview)
			}
			return body, status, lastErr
		}

		return body, resp.StatusCode, nil
	}

	return nil, 0, lastErr
}

// SendMsg 回复客户片段端消息
func SendMsg(client *types.WsClient, message types.ReplyMessage) {
	message.ClientId = client.Id
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error for decoding json data: %v", err.Error())
		return
	}
	err = client.Send(msg)
	if err != nil {
		logger.Errorf("Error for reply message: %v", err.Error())
	}
}

// SendAndFlush 回复客户端一条完整的消息
func SendAndFlush(ws *types.WsClient, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeText, Body: message})
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeEnd})
}

func SendChunkMsg(ws *types.WsClient, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeText, Body: message})
}

// SendErrMsg 向客户端发送错误消息
func SendErrMsg(ws *types.WsClient, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeErr, Body: message})
}

func SendChannelMsg(ws *types.WsClient, channel types.WsChannel, message interface{}) {
	SendMsg(ws, types.ReplyMessage{Channel: channel, Type: types.MsgTypeText, Body: message})
}

func DownloadImage(imageURL string, proxy string) ([]byte, error) {
	body, _, err := FetchURLBytes(context.Background(), imageURL, proxy, defaultHTTPTimeout, 2, 32<<20)
	return body, err
}

func GetBaseURL(strURL string) string {
	u, err := url.Parse(strURL)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

func GetImgExt(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".png"
	}
	return ext
}

func GetFileSize(url string) (int64, error) {
	// 优先走 HEAD 获取 Content-Length，减少下载成本；若对方不支持，再 fallback 到 GET 统计字节数。
	client := newHTTPClient(defaultHTTPTimeout, "")
	if req, err := http.NewRequest(http.MethodHead, url, nil); err == nil {
		if resp, err := client.Do(req); err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode <= 299 && resp.ContentLength > 0 {
				return resp.ContentLength, nil
			}
		}
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 尽量不把大 body 都读进来，只返回状态码
		return 0, fmt.Errorf("get file size failed: status=%d", resp.StatusCode)
	}
	n, err := io.Copy(io.Discard, resp.Body)
	return n, err
}
