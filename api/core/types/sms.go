package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type SMSConfig struct {
	Active string       `json:"active"`
	Ali    SmsConfigAli `json:"aliyun"`
	Bao    SmsConfigBao `json:"bao"`
}

// SmsConfigAli 阿里云短信平台配置
type SmsConfigAli struct {
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Sign         string `json:"sign"`         // 短信签名
	CodeTempId   string `json:"code_temp_id"` // 验证码短信模板 ID
}

// SmsConfigBao 短信宝平台配置
type SmsConfigBao struct {
	Username     string `json:"username"`      //短信宝平台注册的用户名
	Password     string `json:"password"`      //短信宝平台注册的密码
	Sign         string `json:"sign"`          // 短信签名
	CodeTemplate string `json:"code_template"` // 验证码短信模板 匹配
}
