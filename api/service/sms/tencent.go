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

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentSmsService struct {
	config types.SmsConfigTencent
	client *sms.Client
	region string
}

func NewTencentSmsService(sysConfig *types.SystemConfig) (*TencentSmsService, error) {
	config := sysConfig.SMS.Tencent
	region := config.Region
	if region == "" {
		region = "ap-guangzhou" // 默认使用广州地区
	}

	s := TencentSmsService{
		config: config,
		region: region,
	}
	if sysConfig.SMS.Active == Tencent {
		err := s.UpdateConfig(config)
		if err != nil {
			logger.Errorf("腾讯云短信初始化失败: %v", err)
		}
	}
	return &s, nil
}

func (s *TencentSmsService) UpdateConfig(config types.SmsConfigTencent) error {
	if config.SecretId == "" || config.SecretKey == "" {
		// 配置不完整时不初始化客户端
		s.config = config
		if config.Region != "" {
			s.region = config.Region
		} else {
			s.region = "ap-guangzhou"
		}
		return nil
	}

	region := config.Region
	if region == "" {
		region = "ap-guangzhou"
	}

	// 创建凭证
	credential := common.NewCredential(
		config.SecretId,
		config.SecretKey,
	)

	// 创建客户端配置
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	// 创建客户端
	client, err := sms.NewClient(credential, region, cpf)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	s.client = client
	s.config = config
	s.region = region
	return nil
}

// SendVerifyCode 发送验证码短信
// 注意：腾讯云后台配置的短信模板内容应与配置中的 code_template 一致
// 模板只需要1个参数：{1} 表示验证码，例如：{1}为您的验证码，请于5分钟内填写，如非本人操作，请忽略本短信。
func (s *TencentSmsService) SendVerifyCode(mobile string, code int) error {
	if s.client == nil {
		return fmt.Errorf("腾讯云短信服务未初始化")
	}

	// 创建发送短信请求
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(s.config.SmsSdkAppId)
	request.SignName = common.StringPtr(s.config.Sign)
	request.TemplateId = common.StringPtr(s.config.CodeTempId)
	request.PhoneNumberSet = common.StringPtrs([]string{mobile})
	request.TemplateParamSet = common.StringPtrs([]string{fmt.Sprintf("%d", code), "5"})

	// 发送短信
	response, err := s.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	// 检查响应
	if response.Response == nil {
		return fmt.Errorf("failed to send SMS: response is nil")
	}

	if len(response.Response.SendStatusSet) == 0 {
		return fmt.Errorf("failed to send SMS: no send status")
	}

	sendStatus := response.Response.SendStatusSet[0]
	if sendStatus.Code == nil || *sendStatus.Code != "Ok" {
		message := "unknown error"
		if sendStatus.Message != nil {
			message = *sendStatus.Message
		}
		return fmt.Errorf("failed to send SMS: %s", message)
	}

	return nil
}

var _ Service = &TencentSmsService{}
