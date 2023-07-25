package service

import (
	"chatplus/core/types"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"time"
)

type CaptchaService struct {
	config types.ChatPlusApiConfig
	client *req.Client
}

func NewCaptchaService(config types.ChatPlusApiConfig) *CaptchaService {
	return &CaptchaService{
		config: config,
		client: req.C().SetTimeout(10 * time.Second),
	}
}

func (s *CaptchaService) Get() (interface{}, error) {
	if s.config.Token == "" {
		return nil, errors.New("无效的 API Token")
	}

	url := fmt.Sprintf("%s/api/captcha/get", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("AppId", s.config.AppId).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.config.Token)).
		SetSuccessResult(&res).Get(url)
	if err != nil || r.IsErrorState() {
		return nil, fmt.Errorf("请求 API 失败：%v", err)
	}

	if res.Code != types.Success {
		return nil, fmt.Errorf("请求 API 失败：%s", res.Message)
	}

	return res.Data, nil
}

func (s *CaptchaService) Check(data interface{}) bool {
	url := fmt.Sprintf("%s/api/captcha/check", s.config.ApiURL)
	var res types.BizVo
	r, err := s.client.R().
		SetHeader("AppId", s.config.AppId).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", s.config.Token)).
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
