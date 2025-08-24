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
)

type SmsManager struct {
	aliyun *AliYunSmsService
	bao    *BaoSmsService
	active string
}

var logger = logger2.GetLogger()

func NewSmsManager(sysConfig *types.SystemConfig, aliyun *AliYunSmsService, bao *BaoSmsService) (*SmsManager, error) {

	return &SmsManager{
		active: sysConfig.SMS.Active,
		aliyun: aliyun,
		bao:    bao,
	}, nil
}

func (m *SmsManager) GetService() Service {
	if m.active == Ali {
		return m.aliyun
	}
	return m.bao
}

func (m *SmsManager) SetActive(active string) {
	m.active = active
}
