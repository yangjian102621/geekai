package types

import "os"

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

// CaptchaConfig 行为验证码配置
type CaptchaConfig struct {
	ApiKey  string `json:"api_key,omitempty"`
	Type    string `json:"type,omitempty"` // 验证码类型, 可选值: "dot" 或 "slide"
	Enabled bool   `json:"enabled,omitempty"`
}

// WxLoginConfig 微信扫码登录配置
type WxLoginConfig struct {
	ApiKey    string `json:"api_key,omitempty"`
	NotifyURL string `json:"notify_url,omitempty"` // 登录成功回调 URL
	Enabled   bool   `json:"enabled,omitempty"`    // 是否启用微信登录
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
