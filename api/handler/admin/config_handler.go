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
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	handler.BaseHandler
	levelDB *store.LevelDB
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB, levelDB *store.LevelDB) *ConfigHandler {
	return &ConfigHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}, levelDB: levelDB}
}

func (h *ConfigHandler) Update(c *gin.Context) {
	var data struct {
		Key    string `json:"key"`
		Config struct {
			types.SystemConfig
			Content string `json:"content,omitempty"`
			Updated bool   `json:"updated,omitempty"`
		} `json:"config"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
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
