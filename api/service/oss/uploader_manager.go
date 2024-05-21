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
)

type UploaderManager struct {
	handler Uploader
}

func NewUploaderManager(config *types.AppConfig) (*UploaderManager, error) {
	active := Local
	if config.OSS.Active != "" {
		active = strings.ToUpper(config.OSS.Active)
	}
	var handler Uploader
	switch active {
	case Local:
		handler = NewLocalStorage(config)
		break
	case Minio:
		client, err := NewMiniOss(config)
		if err != nil {
			return nil, err
		}
		handler = client
		break
	case QiNiu:
		handler = NewQiNiuOss(config)
		break
	case AliYun:
		client, err := NewAliYunOss(config)
		if err != nil {
			return nil, err
		}
		handler = client
		break
	}

	return &UploaderManager{handler: handler}, nil
}

func (m *UploaderManager) GetUploadHandler() Uploader {
	return m.handler
}
