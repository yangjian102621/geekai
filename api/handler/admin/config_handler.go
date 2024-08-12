package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	handler.BaseHandler
	levelDB        *store.LevelDB
	licenseService *service.LicenseService
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB, levelDB *store.LevelDB, licenseService *service.LicenseService) *ConfigHandler {
	return &ConfigHandler{
		BaseHandler:    handler.BaseHandler{App: app, DB: db},
		levelDB:        levelDB,
		licenseService: licenseService,
	}
}

func (h *ConfigHandler) Update(c *gin.Context) {
	var data struct {
		Key    string `json:"key"`
		Config struct {
			types.SystemConfig
			Content string `json:"content,omitempty"`
			Updated bool   `json:"updated,omitempty"`
		} `json:"config"`
		ConfigBak types.SystemConfig `json:"config_bak,omitempty"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// ONLY authorized user can change the copyright
	if (data.Key == "system" && data.Config.Copyright != data.ConfigBak.Copyright) && !h.licenseService.GetLicense().Configs.DeCopy {
		resp.ERROR(c, "您无权修改版权信息，请先联系作者获取授权")
		return
	}

	value := utils.JsonEncode(&data.Config)
	config := model.Config{Key: data.Key, Config: value}
	res := h.DB.FirstOrCreate(&config, model.Config{Key: data.Key})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	if config.Id > 0 {
		config.Config = value
		res := h.DB.Updates(&config)
		if res.Error != nil {
			resp.ERROR(c, res.Error.Error())
			return
		}

		// update config cache for AppServer
		var cfg model.Config
		h.DB.Where("marker", data.Key).First(&cfg)
		var err error
		if data.Key == "system" {
			err = utils.JsonDecode(cfg.Config, &h.App.SysConfig)
		}
		if err != nil {
			resp.ERROR(c, "Failed to update config cache: "+err.Error())
			return
		}
		logger.Infof("Update AppServer's config successfully: %v", config.Config)
	}

	resp.SUCCESS(c, config)
}

// Get 获取指定的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	key := c.Query("key")
	var config model.Config
	res := h.DB.Where("marker", key).First(&config)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var value map[string]interface{}
	err := utils.JsonDecode(config.Config, &value)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, value)
}

// Active 激活系统
func (h *ConfigHandler) Active(c *gin.Context) {
	var data struct {
		License string `json:"license"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	info, err := host.Info()
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	err = h.licenseService.ActiveLicense(data.License, info.HostID)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, info.HostID)
}

// GetLicense 获取 License 信息
func (h *ConfigHandler) GetLicense(c *gin.Context) {
	license := h.licenseService.GetLicense()
	resp.SUCCESS(c, license)
}
