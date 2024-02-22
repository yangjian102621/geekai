package sms

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"strings"
)

type ServiceManager struct {
	handler Service
}

var logger = logger2.GetLogger()

func NewSendServiceManager(config *types.AppConfig) (*ServiceManager, error) {
	active := Ali
	if config.SMS.Active != "" {
		active = strings.ToUpper(config.SMS.Active)
	}
	var handler Service
	switch active {
	case Ali:
		client, err := NewAliYunSmsService(config)
		if err != nil {
			return nil, err
		}
		handler = client
		break
	case Bao:
		handler = NewSmsBaoSmsService(config)
		break
	}

	return &ServiceManager{handler: handler}, nil
}

func (m *ServiceManager) GetService() Service {
	return m.handler
}
