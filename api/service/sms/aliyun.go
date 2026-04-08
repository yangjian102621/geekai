package sms

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/types"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type AliYunSmsService struct {
	config types.SmsConfigAli
	client *dysmsapi.Client
	domain string
	zoneId string
}

func NewAliYunSmsService(sysConfig *types.SystemConfig) (*AliYunSmsService, error) {
	config := sysConfig.SMS.Ali
	domain := "dysmsapi.aliyuncs.com"
	zoneId := "cn-hangzhou"

	s := AliYunSmsService{
		config: config,
		domain: domain,
		zoneId: zoneId,
	}
	if sysConfig.SMS.Active == Ali {
		err := s.UpdateConfig(config)
		if err != nil {
			logger.Errorf("阿里云短信初始化失败: %v", err)
		}
	}
	return &s, nil
}

func (s *AliYunSmsService) UpdateConfig(config types.SmsConfigAli) error {
	client, err := dysmsapi.NewClientWithAccessKey(
		s.zoneId,
		config.AccessKey,
		config.AccessSecret)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	s.client = client
	s.config = config
	return nil
}

func (s *AliYunSmsService) SendVerifyCode(mobile string, code int) error {
	if s.client == nil {
		return fmt.Errorf("阿里云短信服务未初始化")
	}
	// 创建短信请求并设置参数
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.Domain = s.domain
	request.PhoneNumbers = mobile
	request.SignName = s.config.Sign
	request.TemplateCode = s.config.CodeTempId
	request.TemplateParam = fmt.Sprintf("{\"code\":\"%d\"}", code) // 短信模板中的参数

	// 发送短信
	response, err := s.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("failed to send SMS:%v", err)
	}

	if response.Code != "OK" {
		return fmt.Errorf("failed to send SMS:%v", response.Message)
	}
	return nil
}

var _ Service = &AliYunSmsService{}
