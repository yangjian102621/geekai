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
	active string
}

func NewUploaderManager(sysConfig *types.SystemConfig, local *LocalStorage, aliyun *AliYunOss, mini *MiniOss, qiniu *QiNiuOss) (*UploaderManager, error) {
	if sysConfig.OSS.Active == "" {
		sysConfig.OSS.Active = Local
	}
	sysConfig.OSS.Active = strings.ToLower(sysConfig.OSS.Active)

	return &UploaderManager{
		active: sysConfig.OSS.Active,
		local:  local,
		aliyun: aliyun,
		mini:   mini,
		qiniu:  qiniu,
	}, nil
}

func (m *UploaderManager) GetUploadHandler() Uploader {
	switch m.active {
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

func (m *UploaderManager) UpdateConfig(config types.OSSConfig) {
	switch config.Active {
	case Local:
		m.local.UpdateConfig(config.Local)
	case AliYun:
		m.aliyun.UpdateConfig(config.AliYun)
	case Minio:
		m.mini.UpdateConfig(config.Minio)
	case QiNiu:
		m.qiniu.UpdateConfig(config.QiNiu)
	}
	m.active = config.Active
}
