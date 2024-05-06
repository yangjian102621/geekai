package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type OSSConfig struct {
	Active string
	Local  LocalStorageConfig
	Minio  MiniOssConfig
	QiNiu  QiNiuOssConfig
	AliYun AliYunOssConfig
}
type MiniOssConfig struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	Bucket       string
	SubDir       string
	UseSSL       bool
	Domain       string
}

type QiNiuOssConfig struct {
	Zone         string
	AccessKey    string
	AccessSecret string
	Bucket       string
	SubDir       string
	Domain       string
}

type AliYunOssConfig struct {
	Endpoint     string
	AccessKey    string
	AccessSecret string
	Bucket       string
	SubDir       string
	Domain       string
}

type LocalStorageConfig struct {
	BasePath string
	BaseURL  string
}
