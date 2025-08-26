package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"errors"
	"geekai/core"
	"geekai/core/middleware"
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
	configService  *service.ConfigService
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB, levelDB *store.LevelDB, licenseService *service.LicenseService, configService *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		BaseHandler:    handler.BaseHandler{App: app, DB: db},
		levelDB:        levelDB,
		licenseService: licenseService,
		configService:  configService,
	}
}

// RegisterRoutes 注册路由
func (h *ConfigHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/admin/config/")

	// 需要管理员授权的接口
	group.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		group.POST("update", h.Update)
		group.GET("get", h.Get)
		group.POST("active", h.Active)
		group.POST("test", h.Test)
		group.GET("license", h.GetLicense)
	}
}

func (h *ConfigHandler) Update(c *gin.Context) {
	var payload struct {
		Key       string             `json:"key"`
		Config    json.RawMessage    `json:"config"`
		ConfigBak types.SystemConfig `json:"config_bak,omitempty"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("Update config failed: %v", err)
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if payload.Key == "system" {
		var sys types.SystemConfig
		if err := json.Unmarshal(payload.Config, &sys); err != nil {
			resp.ERROR(c, "系统配置解析失败: "+err.Error())
			return
		}
		if (sys.Base.Copyright != payload.ConfigBak.Base.Copyright) && !h.licenseService.GetLicense().Configs.DeCopy {
			resp.ERROR(c, "您无权修改版权信息，请先联系作者获取授权")
			return
		}

	}

	// 使用统一配置服务写入与广播
	if err := h.configService.Set(payload.Key, payload.Config); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	if payload.Key == "system" {
		var sys types.SystemConfig
		if err := json.Unmarshal(payload.Config, &sys); err == nil {
			h.App.SysConfig = &sys
		}
	}
	resp.SUCCESS(c)
}

// Get 获取指定的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	key := c.Query("key")
	var config model.Config
	res := h.DB.Where("name", key).First(&config)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			resp.SUCCESS(c, map[string]interface{}{})
			return
		}
		resp.ERROR(c, res.Error.Error())
		return
	}

	var value map[string]interface{}
	err := utils.JsonDecode(config.Value, &value)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, value)
}

// Test 配置测试（占位）
func (h *ConfigHandler) Test(c *gin.Context) {
	var data struct {
		Key string `json:"key"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	msg, err := h.configService.Test(data.Key)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, msg)
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

	resp.SUCCESS(c)

}

// GetLicense 获取 License 信息
func (h *ConfigHandler) GetLicense(c *gin.Context) {
	license := h.licenseService.GetLicense()
	resp.SUCCESS(c, license)
}
