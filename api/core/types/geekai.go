package types

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// GeekAI 增值服务
var GeekAPIURL = "https://sapi.geekai.me"

func init() {
	if os.Getenv("GEEK_API_URL") != "" {
		GeekAPIURL = os.Getenv("GEEK_API_URL")
	}
}

// decodeBool 兼容历史数据中布尔值被存为字符串（"true"/"false"/"1"/"0"）的情况
func decodeBool(b []byte) (bool, error) {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return false, err
	}
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return false, err
		}
		return parsed, nil
	default:
		return false, fmt.Errorf("无法将 %T 解析为 bool", v)
	}
}

// CaptchaConfig 行为验证码配置
type CaptchaConfig struct {
	ApiKey  string `json:"api_key,omitempty"`
	Type    string `json:"type,omitempty"` // 验证码类型, 可选值: "dot" 或 "slide"
	Enabled bool   `json:"enabled,omitempty"`
}

func (c *CaptchaConfig) UnmarshalJSON(b []byte) error {
	var aux struct {
		ApiKey  string          `json:"api_key,omitempty"`
		Type    string          `json:"type,omitempty"`
		Enabled json.RawMessage `json:"enabled,omitempty"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	c.ApiKey = aux.ApiKey
	c.Type = aux.Type
	if len(aux.Enabled) > 0 {
		enabled, err := decodeBool(aux.Enabled)
		if err != nil {
			return err
		}
		c.Enabled = enabled
	}
	return nil
}

// WxLoginConfig 微信扫码登录配置
type WxLoginConfig struct {
	ApiKey    string `json:"api_key,omitempty"`
	NotifyURL string `json:"notify_url,omitempty"` // 登录成功回调 URL
	Enabled   bool   `json:"enabled,omitempty"`    // 是否启用微信登录
}

func (c *WxLoginConfig) UnmarshalJSON(b []byte) error {
	var aux struct {
		ApiKey    string          `json:"api_key,omitempty"`
		NotifyURL string          `json:"notify_url,omitempty"`
		Enabled   json.RawMessage `json:"enabled,omitempty"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	c.ApiKey = aux.ApiKey
	c.NotifyURL = aux.NotifyURL
	if len(aux.Enabled) > 0 {
		enabled, err := decodeBool(aux.Enabled)
		if err != nil {
			return err
		}
		c.Enabled = enabled
	}
	return nil
}

// 微信公众号配置
type WxGzhConfig struct {
	AppId          string `json:"app_id,omitempty"`
	Secret         string `json:"secret,omitempty"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	Enabled        bool   `json:"enabled"`
}

// WxGzhMenuConfig 公众号自定义菜单草稿（与微信公众平台 menu/create 的 button 结构一致）
type WxGzhMenuConfig struct {
	Button []WxGzhMenuButton `json:"button"`
}

// WxGzhMenuButton 单条菜单；含子菜单时仅填 Name 与 SubButton，不填 Type
type WxGzhMenuButton struct {
	Type       string            `json:"type,omitempty"`        // view | click 等
	Name       string            `json:"name"`
	Key        string            `json:"key,omitempty"`         // click
	URL        string            `json:"url,omitempty"`         // view
	AppID      string            `json:"appid,omitempty"`       // miniprogram
	PagePath   string            `json:"pagepath,omitempty"`    // miniprogram
	SubButton  []WxGzhMenuButton `json:"sub_button,omitempty"`  // 子菜单
}
