package service

import (
	"bytes"
	"context"
	"errors"
	"geekai/store/model"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RequestChatCompletion 共用后台的模型密钥配置，用户模型令牌不会发给上游。
func RequestChatCompletion(ctx context.Context, db *gorm.DB, keyID uint, body []byte, key *model.ApiKey) (*http.Response, error) {
	query := db.Where("enabled = ? AND type = ?", true, "chat")
	if keyID > 0 {
		query = query.Where("id = ?", keyID)
	} else {
		query = query.Order("last_used_at ASC, id ASC")
	}
	if err := query.First(key).Error; err != nil || key.Value == "" {
		return nil, errors.New("没有可用的模型密钥")
	}
	endpoint, err := url.Parse(key.ApiURL)
	if err != nil || endpoint.Host == "" || (endpoint.Scheme != "http" && endpoint.Scheme != "https") {
		return nil, errors.New("模型请求地址无效")
	}
	switch strings.TrimRight(endpoint.Path, "/") {
	case "":
		endpoint.Path = "/v1/chat/completions"
	case "/v1":
		endpoint.Path = "/v1/chat/completions"
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+key.Value)
	client := http.DefaultClient
	if key.ProxyURL != "" {
		proxy, err := url.Parse(key.ProxyURL)
		if err != nil || proxy.Host == "" {
			return nil, errors.New("模型代理地址无效")
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy), DisableKeepAlives: true}}
	}
	if err := db.Model(&model.ApiKey{}).Where("id = ?", key.Id).UpdateColumn("last_used_at", time.Now().Unix()).Error; err != nil {
		return nil, err
	}
	return client.Do(request)
}
