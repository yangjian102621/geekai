package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/store"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/shirou/gopsutil/host"
)

type LicenseService struct {
	config       types.ApiConfig
	levelDB      *store.LevelDB
	license      *types.License
	urlWhiteList []string
	machineId    string
}

func NewLicenseService(server *core.AppServer, levelDB *store.LevelDB) *LicenseService {
	var license types.License
	var machineId string
	_ = levelDB.Get(types.LicenseKey, &license)
	info, err := host.Info()
	if err == nil {
		machineId = info.HostID
	}
	logger.Infof("License: %+v", license)
	return &LicenseService{
		config:    server.Config.ApiConfig,
		levelDB:   levelDB,
		license:   &license,
		machineId: machineId,
	}
}

type License struct {
	Name      string              `json:"name"`
	License   string              `json:"license"`
	MachineId string              `json:"mid"`
	ActiveAt  int64               `json:"active_at"`
	ExpiredAt int64               `json:"expired_at"`
	UserNum   int                 `json:"user_num"`
	Configs   types.LicenseConfig `json:"configs"`
}

// ActiveLicense 激活 License
func (s *LicenseService) ActiveLicense(license string, machineId string) error {
	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    License       `json:"data"`
	}
	apiURL := fmt.Sprintf("%s/%s", s.config.ApiURL, "api/license/active")
	response, err := req.C().R().
		SetBody(map[string]string{"license": license, "machine_id": machineId}).
		SetSuccessResult(&res).Post(apiURL)
	if err != nil {
		return fmt.Errorf("发送激活请求失败: %v", err)
	}

	if response.IsErrorState() {
		return fmt.Errorf("发送激活请求失败：%v", response.Status)
	}

	if res.Code != types.Success {
		return fmt.Errorf("激活失败：%v", res.Message)
	}

	s.license = &types.License{
		Key:       license,
		MachineId: machineId,
		Configs:   res.Data.Configs,
		ExpiredAt: res.Data.ExpiredAt,
		IsActive:  true,
	}
	err = s.levelDB.Put(types.LicenseKey, s.license)
	if err != nil {
		return fmt.Errorf("保存许可证书失败：%v", err)
	}
	return nil
}

// SyncLicense 定期同步 License
func (s *LicenseService) SyncLicense() {
	go func() {
		retryCounter := 0
		for {
			license, err := s.fetchLicense()
			if err != nil {
				retryCounter++
				if retryCounter < 5 {
					logger.Warn(err)
				}
				s.license.IsActive = false
			} else {
				s.license = license
			}

			urls, err := s.fetchUrlWhiteList()
			if err == nil {
				s.urlWhiteList = urls
			}

			time.Sleep(time.Second * 10)
		}
	}()
}

func (s *LicenseService) fetchLicense() (*types.License, error) {
	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    License       `json:"data"`
	}
	apiURL := fmt.Sprintf("%s/%s", s.config.ApiURL, "api/license/check")
	response, err := req.C().R().
		SetBody(map[string]string{"license": s.license.Key, "machine_id": s.machineId}).
		SetSuccessResult(&res).Post(apiURL)
	if err != nil {
		return nil, fmt.Errorf("发送激活请求失败: %v", err)
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("激活失败：%v", response.Status)
	}
	if res.Code != types.Success {
		return nil, fmt.Errorf("激活失败：%v", res.Message)
	}

	return &types.License{
		Key:       res.Data.License,
		MachineId: res.Data.MachineId,
		Configs:   res.Data.Configs,
		ExpiredAt: res.Data.ExpiredAt,
		IsActive:  true,
	}, nil
}

func (s *LicenseService) fetchUrlWhiteList() ([]string, error) {
	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    []string      `json:"data"`
	}
	apiURL := fmt.Sprintf("%s/%s", s.config.ApiURL, "api/license/urls")
	response, err := req.C().R().SetSuccessResult(&res).Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	if response.IsErrorState() {
		return nil, fmt.Errorf("发送请求失败：%v", response.Status)
	}
	if res.Code != types.Success {
		return nil, fmt.Errorf("获取白名单失败：%v", res.Message)
	}

	return res.Data, nil
}

// GetLicense 获取许可信息
func (s *LicenseService) GetLicense() *types.License {
	return s.license
}

// IsValidApiURL 判断是否合法的中转 URL
func (s *LicenseService) IsValidApiURL(uri string) error {
	// 获得许可授权的直接放行
	if s.license.IsActive {
		if s.license.MachineId != s.machineId {
			return errors.New("系统使用了盗版的许可证书")
		}

		if time.Now().Unix() > s.license.ExpiredAt {
			return errors.New("系统许可证书已经过期")
		}
		return nil
	}

	if len(s.urlWhiteList) == 0 {
		urls, err := s.fetchUrlWhiteList()
		if err == nil {
			s.urlWhiteList = urls
		}
	}

	for _, v := range s.urlWhiteList {
		if strings.HasPrefix(uri, v) {
			return nil
		}
	}
	return fmt.Errorf("当前 API 地址 %s 不在白名单列表当中。", uri)
}
