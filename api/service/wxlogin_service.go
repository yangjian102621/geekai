package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
)

type WxLoginService struct {
	config      types.WxLoginConfig
	client      *req.Client
	redisClient *redis.Client
}

const loginStateKeyPrefix = "wx_login_state/"

type LoginStatus struct {
	Status string `json:"status"`
	OpenID string `json:"openid"`
	Token  string `json:"token"`
}

const (
	LoginStatusPending = "pending"
	LoginStatusSuccess = "success"
	LoginStatusExpired = "expired" // 登录失效，需要重新登录
)

func NewWxLoginService(config types.WxLoginConfig, redisClient *redis.Client) *WxLoginService {
	return &WxLoginService{
		config:      config,
		client:      req.C().SetTimeout(10 * time.Second),
		redisClient: redisClient,
	}
}

func (s *WxLoginService) UpdateConfig(config types.WxLoginConfig) {
	s.config = config
}

func (s *WxLoginService) GetConfig() types.WxLoginConfig {
	return s.config
}

func (s *WxLoginService) SetConfig(config types.WxLoginConfig) {
	s.config = config
}

func (s *WxLoginService) GetLoginQrCodeUrl(state string) (string, error) {
	if s.config.ApiKey == "" {
		return "", errors.New("无效的 API Key")
	}

	url := fmt.Sprintf("%s/api/auth/wechat/login", types.GeekAPIURL)
	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    struct {
			Ticket string `json:"ticket"`
			Url    string `json:"url"`
		} `json:"data"`
	}
	r, err := s.client.R().
		SetHeader("Authorization", s.config.ApiKey).
		SetBody(map[string]string{
			"notify_url": s.config.NotifyURL,
			"state":      state,
		}).
		SetSuccessResult(&res).Post(url)
	if err != nil || r.IsErrorState() {
		return "", fmt.Errorf("请求 API 失败：%v", err)
	}

	if res.Code != types.Success {
		return "", fmt.Errorf("请求 API 失败：%s", res.Message)
	}

	status := LoginStatus{
		Status: LoginStatusPending,
		OpenID: "",
	}
	s.redisClient.Set(context.Background(), loginStateKeyPrefix+state, utils.JsonEncode(status), time.Hour)

	return res.Data.Url, nil
}

func (s *WxLoginService) GetLoginStatus(state string) (*LoginStatus, error) {
	result, err := s.redisClient.Get(context.Background(), loginStateKeyPrefix+state).Result()
	if err != nil {
		return nil, errors.New("登录失败")
	}

	var status LoginStatus
	err = utils.JsonDecode(result, &status)
	if err != nil {
		return nil, errors.New("登录失败")
	}

	return &status, nil
}

func (s *WxLoginService) SetLoginStatus(state string, status LoginStatus) {
	s.redisClient.Set(context.Background(), loginStateKeyPrefix+state, utils.JsonEncode(status), time.Hour)
}
