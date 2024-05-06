package sms

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core/types"
	logger2 "geekai/logger"
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
