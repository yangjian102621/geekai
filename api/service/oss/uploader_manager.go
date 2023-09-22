package oss

import (
	"chatplus/core/types"
	"strings"
)

type UploaderManager struct {
	handler Uploader
}

const Local = "LOCAL"
const Minio = "MINIO"
const QiNiu = "QINIU"
const AliYun = "ALIYUN"

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
