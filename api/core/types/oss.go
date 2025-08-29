package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type OSSConfig struct {
	Active string             `json:"active"`
	Local  LocalStorageConfig `json:"local"`
	Minio  MiniOssConfig      `json:"minio"`
	QiNiu  QiNiuOssConfig     `json:"qiniu"`
	AliYun AliYunOssConfig    `json:"aliyun"`
}

type MiniOssConfig struct {
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Bucket       string `json:"bucket"`
	UseSSL       bool   `json:"use_ssl"`
	Domain       string `json:"domain"`
}

type QiNiuOssConfig struct {
	Zone         string `json:"zone"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Bucket       string `json:"bucket"`
	Domain       string `json:"domain"`
}

type AliYunOssConfig struct {
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Bucket       string `json:"bucket"`
	Domain       string `json:"domain"`
}

type LocalStorageConfig struct {
	BasePath string `json:"base_path"`
	BaseURL  string `json:"base_url"`
}
