package sms

import (
	"chatplus/core/types"
	"strings"
)

type SmsServiceManager struct {
	handler SmsService
}

const Ali = "Ali"
const SmsBao = "SmsBao"

func NewSendServiceManager(config *types.AppConfig) (*SmsServiceManager, error) {
	active := SmsBao
	if config.OSS.Active != "" {
		active = strings.ToUpper(config.SMS.Active)
	}
	var handler SmsService
	switch active {
	case Ali:
		client, err := NewAliYunSmsService(config)
		if err != nil {
			return nil, err
		}
		handler = client
		break
	case SmsBao:
		handler = NewSmsBaoSmsService(config)
		break
	}

	return &SmsServiceManager{handler: handler}, nil
}

func (m *SmsServiceManager) GetUploadHandler() SmsService {
	return m.handler
}
