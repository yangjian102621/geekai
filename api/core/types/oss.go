package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type OSSConfig struct {
	Active  string             `json:"active,omitempty"`
	Local   LocalStorageConfig `json:"local,omitempty"`
	Minio   MiniOssConfig      `json:"minio,omitempty"`
	QiNiu   QiNiuOssConfig     `json:"qiniu,omitempty"`
	AliYun  AliYunOssConfig    `json:"aliyun,omitempty"`
	Tencent TencentOssConfig   `json:"tencent,omitempty"`
}

type MiniOssConfig struct {
	Endpoint      string `json:"endpoint,omitempty"`
	AccessKey     string `json:"access_key,omitempty"`
	AccessSecret  string `json:"access_secret,omitempty"`
	Bucket        string `json:"bucket,omitempty"`
	UseSSL        bool   `json:"use_ssl,omitempty"`
	Domain        string `json:"domain,omitempty"`
	ThumbTemplate string `json:"thumb_template,omitempty"` // 缩略图模板，使用{width}和{height}作为变量占位符
}

type QiNiuOssConfig struct {
	Zone          string `json:"zone,omitempty"`
	AccessKey     string `json:"access_key,omitempty"`
	AccessSecret  string `json:"access_secret,omitempty"`
	Bucket        string `json:"bucket,omitempty"`
	Domain        string `json:"domain,omitempty"`
	ThumbTemplate string `json:"thumb_template,omitempty"` // 缩略图模板，使用{width}和{height}作为变量占位符，默认：?imageView2/4/w/{width}/h/{height}/q/75
}

type AliYunOssConfig struct {
	Endpoint      string `json:"endpoint,omitempty"`
	AccessKey     string `json:"access_key,omitempty"`
	AccessSecret  string `json:"access_secret,omitempty"`
	Bucket        string `json:"bucket,omitempty"`
	Domain        string `json:"domain,omitempty"`
	ThumbTemplate string `json:"thumb_template,omitempty"` // 缩略图模板，使用{width}和{height}作为变量占位符，默认：?x-oss-process=image/resize,m_lfit,w_{width},h_{height}
}

type LocalStorageConfig struct {
	BasePath      string `json:"base_path,omitempty"`
	BaseURL       string `json:"base_url,omitempty"`
	ThumbTemplate string `json:"thumb_template,omitempty"` // 缩略图模板，使用{width}和{height}作为变量占位符，默认：?imageView2/4/w/{width}/h/{height}/q/75
}

type TencentOssConfig struct {
	Region        string `json:"region,omitempty"`
	SecretId      string `json:"secret_id,omitempty"`
	SecretKey     string `json:"secret_key,omitempty"`
	Bucket        string `json:"bucket,omitempty"`
	Domain        string `json:"domain,omitempty"`
	ThumbTemplate string `json:"thumb_template,omitempty"` // 缩略图模板，使用{width}和{height}作为变量占位符，默认：?imageView2/1/w/{width}/h/{height}/format/jpg
}
