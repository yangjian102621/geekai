package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type OSSConfig struct {
	Active string             `json:"active,omitempty"`
	Local  LocalStorageConfig `json:"local,omitempty"`
	Minio  MiniOssConfig      `json:"minio,omitempty"`
	QiNiu  QiNiuOssConfig     `json:"qiniu,omitempty"`
	AliYun AliYunOssConfig    `json:"aliyun,omitempty"`
}

type MiniOssConfig struct {
	Endpoint     string `json:"endpoint,omitempty"`
	AccessKey    string `json:"access_key,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
	Bucket       string `json:"bucket,omitempty"`
	UseSSL       bool   `json:"use_ssl,omitempty"`
	Domain       string `json:"domain,omitempty"`
}

type QiNiuOssConfig struct {
	Zone         string `json:"zone,omitempty"`
	AccessKey    string `json:"access_key,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
	Bucket       string `json:"bucket,omitempty"`
	Domain       string `json:"domain,omitempty"`
}

type AliYunOssConfig struct {
	Endpoint     string `json:"endpoint,omitempty"`
	AccessKey    string `json:"access_key,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
	Bucket       string `json:"bucket,omitempty"`
	Domain       string `json:"domain,omitempty"`
}

type LocalStorageConfig struct {
	BasePath string `json:"base_path,omitempty"`
	BaseURL  string `json:"base_url,omitempty"`
}
