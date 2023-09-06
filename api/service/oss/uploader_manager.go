package oss

import (
	"chatplus/core/types"
	"strings"
)

type UploaderManager struct {
	active         string
	uploadServices map[string]Uploader
}

const Local = "LOCAL"
const Minio = "MINIO"
const QiNiu = "QINIU"

func NewUploaderManager(config *types.AppConfig) (*UploaderManager, error) {
	services := make(map[string]Uploader)
	if config.OSS.Minio.AccessKey != "" {
		minioService, err := NewMinioService(config)
		if err != nil {
			return nil, err
		}
		services[Minio] = minioService
	}
	if config.OSS.Local.BasePath != "" {
		services[Local] = NewLocalStorageService(config)
	}
	if config.OSS.QiNiu.AccessKey != "" {
		services[QiNiu] = NewQiNiuService(config)
	}
	active := Local
	if config.OSS.Active != "" {
		active = strings.ToUpper(config.OSS.Active)
	}
	return &UploaderManager{uploadServices: services, active: active}, nil
}

func (m *UploaderManager) GetActiveService() Uploader {
	return m.uploadServices[m.active]
}
