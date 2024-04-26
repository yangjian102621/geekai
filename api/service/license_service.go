package service

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/shirou/gopsutil/host"
	"strings"
	"time"
)

type LicenseService struct {
	config       types.ApiConfig
	levelDB      *store.LevelDB
	license      types.License
	urlWhiteList []string
	machineId    string
}

func NewLicenseService(server *core.AppServer, levelDB *store.LevelDB) * LicenseService {
	var license types.License
	var machineId string
	_ = levelDB.Get(types.LicenseKey, &license)
	info, err := host.Info()
	if err == nil {
		machineId = info.HostID
	}
	return &LicenseService{
		config:    server.Config.ApiConfig,
		levelDB:   levelDB,
		license:   license,
		machineId: machineId,
	}
}

type License struct {
	Name      string `json:"name"`
	Value     string `json:"license"`
	Mid       string `json:"mid"`
	ExpiredAt int64  `json:"expired_at"`
	UserNum   int    `json:"user_num"`
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
		return fmt.Errorf( "发送激活请求失败：%v", response.Status)
	}

	if res.Code != types.Success {
		return fmt.Errorf("激活失败：%v", res.Message)
	}

	s.license = types.License{
		Key:       license,
		MachineId: machineId,
		UserNum:   res.Data.UserNum,
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
		for {
			var res struct {
				Code    types.BizCode `json:"code"`
				Message string        `json:"message"`
				Data    struct {
					License License  `json:"license"`
					Urls    []string `json:"urls"`
				}
			}
			apiURL := fmt.Sprintf("%s/%s", s.config.ApiURL, "api/license/check")
			response, err := req.C().R().
				SetBody(map[string]string{"license": s.license.Key, "machine_id": s.machineId}).
				SetSuccessResult(&res).Post(apiURL)
			if err != nil {
				logger.Errorf("发送激活请求失败: %v", err)
				goto next
			}
			if response.IsErrorState() {
				logger.Errorf("激活失败：%v", response.Status)
				goto next
			}
			if res.Code != types.Success {
				logger.Errorf("激活失败：%v", res.Message)
				s.license.IsActive = false
				goto next
			}

			s.license = types.License{
				Key:       res.Data.License.Value,
				MachineId: res.Data.License.Mid,
				UserNum:   res.Data.License.UserNum,
				ExpiredAt: res.Data.License.ExpiredAt,
				IsActive:  true,
			}
			s.urlWhiteList = res.Data.Urls
			logger.Debugf("同步 License 成功：%v\n%v", s.license, s.urlWhiteList)
		next:
			time.Sleep(time.Second * 10)
		}
	}()
}

// GetLicense 获取许可信息
func (s *LicenseService) GetLicense() types.License {
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

	for _, v := range s.urlWhiteList {
		if strings.HasPrefix(uri, v) {
			return nil
		}
	}
	return fmt.Errorf("当前 API 地址 %s 不在白名单列表当中。", uri)
}