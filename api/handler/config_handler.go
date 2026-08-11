package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"geekai/core"
	"geekai/core/types"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	BaseHandler
	uploaderManager *oss.UploaderManager
	sysConfig       *types.SystemConfig
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB, uploaderManager *oss.UploaderManager, sysConfig *types.SystemConfig) *ConfigHandler {
	return &ConfigHandler{
		BaseHandler:     BaseHandler{App: app, DB: db},
		uploaderManager: uploaderManager,
		sysConfig:       sysConfig,
	}
}

// RegisterRoutes 注册路由
func (h *ConfigHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/config/")

	// 无需授权的接口
	group.GET("get", h.Get)
	group.GET("oss/thumb", h.GetOssThumbTemplate)
}

// Get 获取指定的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	key := c.Query("key")
	var config model.Config
	err := h.DB.Where("name", key).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.SUCCESS(c, nil)
		return
	}

	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var value map[string]any
	err = utils.JsonDecode(config.Value, &value)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if key == types.ConfigKeyWxGzh {
		delete(value, "secret")
		delete(value, "token")
		delete(value, "encoding_aes_key")
	}
	resp.SUCCESS(c, value)
}

// GetOssThumbTemplate 获取当前存储引擎的缩略图模板
func (h *ConfigHandler) GetOssThumbTemplate(c *gin.Context) {
	template := h.uploaderManager.GetThumbTemplate()
	resp.SUCCESS(c, gin.H{
		"template": template,
		"active":   h.sysConfig.OSS.Active,
	})
}
