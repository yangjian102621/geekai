package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/types"
	"time"

	"github.com/imroc/req/v3"
)

type CaptchaService struct {
	config types.CaptchaConfig
	client *req.Client
}

func NewCaptchaService(captchaConfig types.CaptchaConfig) *CaptchaService {
	return &CaptchaService{
		config: captchaConfig,
		client: req.C().SetTimeout(10 * time.Second),
	}
}

func (s *CaptchaService) UpdateConfig(config types.CaptchaConfig) {
	s.config = config
}

func (s *CaptchaService) GetConfig() types.CaptchaConfig {
	return s.config
}

func (s *CaptchaService) Get() (interface{}, error) {
	url := fmt.Sprintf("%s/api/captcha/get", types.GeekAPIURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.config.ApiKey)).
		SetSuccessResult(&res).Get(url)
	if err != nil || r.IsErrorState() {
		return nil, fmt.Errorf("请求 API 失败：%v", err)
	}

	if res.Code != types.Success {
		return nil, fmt.Errorf("请求 API 失败：%s", res.Message)
	}

	return res.Data, nil
}

func (s *CaptchaService) Check(data any) bool {
	url := fmt.Sprintf("%s/api/captcha/check", types.GeekAPIURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.config.ApiKey)).
		SetBodyJsonMarshal(data).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return false
	}

	if res.Code != types.Success {
		return false
	}

	return true
}

func (s *CaptchaService) SlideGet() (any, error) {
	url := fmt.Sprintf("%s/api/captcha/slide/get", types.GeekAPIURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.config.ApiKey)).
		SetSuccessResult(&res).Get(url)
	if err != nil || r.IsErrorState() {
		return nil, fmt.Errorf("请求 API 失败：%v", err)
	}

	if res.Code != types.Success {
		return nil, fmt.Errorf("请求 API 失败：%s", res.Message)
	}

	return res.Data, nil
}

func (s *CaptchaService) SlideCheck(data any) bool {
	url := fmt.Sprintf("%s/api/captcha/slide/check", types.GeekAPIURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.config.ApiKey)).
		SetBodyJsonMarshal(data).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return false
	}

	if res.Code != types.Success {
		return false
	}

	return true
}
