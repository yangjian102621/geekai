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
	ApiKey  string `json:"api_key"`
	Type    string `json:"type"` // 验证码类型, 可选值: "dot" 或 "slide"
	Enabled bool   `json:"enabled"`
}

// WxLoginConfig 微信登录配置
type WxLoginConfig struct {
	ApiKey    string `json:"api_key"`
	NotifyURL string `json:"notify_url"` // 登录成功回调 URL
	Enabled   bool   `json:"enabled"`    // 是否启用微信登录
}
