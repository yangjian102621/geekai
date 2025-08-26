package oss

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/types"
	"strings"

	logger2 "geekai/logger"
)

var logger = logger2.GetLogger()

type UploaderManager struct {
	local  *LocalStorage
	aliyun *AliYunOss
	mini   *MiniOss
	qiniu  *QiNiuOss
	config *types.OSSConfig
}

func NewUploaderManager(sysConfig *types.SystemConfig, local *LocalStorage, aliyun *AliYunOss, mini *MiniOss, qiniu *QiNiuOss) (*UploaderManager, error) {
	if sysConfig.OSS.Active == "" {
		sysConfig.OSS.Active = Local
	}
	sysConfig.OSS.Active = strings.ToLower(sysConfig.OSS.Active)

	return &UploaderManager{
		config: &sysConfig.OSS,
		local:  local,
		aliyun: aliyun,
		mini:   mini,
		qiniu:  qiniu,
	}, nil
}

func (m *UploaderManager) GetUploadHandler() Uploader {
	switch m.config.Active {
	case Local:
		return m.local
	case AliYun:
		return m.aliyun
	case Minio:
		return m.mini
	case QiNiu:
		return m.qiniu
	}
	return m.local
}
