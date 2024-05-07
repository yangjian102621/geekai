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
	config types.ApiConfig
	levelDB *store.LevelDB
	license types.License
	machineId string
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
		config: server.Config.ApiConfig,
		levelDB: levelDB,
		license: license,
		machineId: machineId,
	}
}

// ActiveLicense 激活 License
func (s *LicenseService) ActiveLicense(license string, machineId string) error {
	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    struct {
			Name      string `json:"name"`
			License   string `json:"license"`
			Mid       string `json:"mid"`
			ExpiredAt int64  `json:"expired_at"`
			UserNum   int    `json:"user_num"`
		}
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
		return fmt.Errorf( "激活失败：%v", res.Message)
	}

	err = s.levelDB.Put(types.LicenseKey, types.License{
		Key:       license,
		MachineId: machineId,
		UserNum:   res.Data.UserNum,
		ExpiredAt: res.Data.ExpiredAt,
		IsActive: true,
	})
	if err != nil {
		return fmt.Errorf("保存许可证书失败：%v", err)
	}

	return nil
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

	if !strings.HasPrefix(uri, "https://gpt.bemore.lol") &&
		!strings.HasPrefix(uri, "https://api.openai.com") &&
		!strings.HasPrefix(uri, "http://cdn.chat-plus.net") &&
		!strings.HasPrefix(uri, "https://api.chat-plus.net") {
		return fmt.Errorf("当前 API 地址 %s 不在白名单列表当中。",uri)
	}

	return nil
}