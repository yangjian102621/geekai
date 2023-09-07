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

func NewUploaderManager(config *types.AppConfig) (*UploaderManager, error) {
	active := Local
	if config.OSS.Active != "" {
		active = strings.ToUpper(config.OSS.Active)
	}
	var handler Uploader
	switch active {
	case Local:
		handler = NewLocalStorageService(config)
		break
	case Minio:
		service, err := NewMinioService(config)
		if err != nil {
			return nil, err
		}
		handler = service
		break
	case QiNiu:
		handler = NewQiNiuService(config)
	}

	return &UploaderManager{handler: handler}, nil
}

func (m *UploaderManager) GetUploadHandler() Uploader {
	return m.handler
}
