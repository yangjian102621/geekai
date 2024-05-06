package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type SMSConfig struct {
	Active string
	Ali    SmsConfigAli
	Bao    SmsConfigBao
}

// SmsConfigAli 阿里云短信平台配置
type SmsConfigAli struct {
	AccessKey    string
	AccessSecret string
	Product      string
	Domain       string
	Sign         string // 短信签名
	CodeTempId   string // 验证码短信模板 ID
}

// SmsConfigBao 短信宝平台配置
type SmsConfigBao struct {
	Username     string //短信宝平台注册的用户名
	Password     string //短信宝平台注册的密码
	Domain       string //域名
	Sign         string // 短信签名
	CodeTemplate string // 验证码短信模板 匹配
}
