package oss

import (
	"chatplus/core/types"
	"strings"
)

type UploaderManager struct {
	active         string
	uploadServices map[string]Uploader
}

const uploaderLocal = "LOCAL"
const uploaderMinio = "MINIO"

func NewUploaderManager(config *types.AppConfig) (*UploaderManager, error) {
	services := make(map[string]Uploader)
	if config.OSS.Minio.AccessKey != "" {
		minioService, err := NewMinioService(config)
		if err != nil {
			return nil, err
		}
		services[uploaderMinio] = minioService
	}
	if config.OSS.Local.BasePath != "" {
		services[uploaderLocal] = NewLocalStorageService(config)
	}
	active := uploaderLocal
	if config.OSS.Active != "" {
		active = strings.ToUpper(config.OSS.Active)
	}
	return &UploaderManager{uploadServices: services, active: active}, nil
}

func (m *UploaderManager) GetActiveService() Uploader {
	return m.uploadServices[m.active]
}
