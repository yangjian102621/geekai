package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiKeyHandler struct {
	handler.BaseHandler
}

func NewApiKeyHandler(app *core.AppServer, db *gorm.DB) *ApiKeyHandler {
	return &ApiKeyHandler{BaseHandler: handler.BaseHandler{DB: db, App: app}}
}

func (h *ApiKeyHandler) Save(c *gin.Context) {
	var data struct {
		Id       uint   `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Value    string `json:"value"`
		ApiURL   string `json:"api_url"`
		Enabled  bool   `json:"enabled"`
		ProxyURL string `json:"proxy_url"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	apiKey := model.ApiKey{}
	if data.Id > 0 {
		h.DB.Find(&apiKey, data.Id)
	}
	apiKey.Value = data.Value
	apiKey.Type = data.Type
	apiKey.ApiURL = data.ApiURL
	apiKey.Enabled = data.Enabled
	apiKey.ProxyURL = data.ProxyURL
	apiKey.Name = data.Name
	err := h.DB.Save(&apiKey).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var keyVo vo.ApiKey
	err = utils.CopyObject(apiKey, &keyVo)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("拷贝数据失败：%v", err))
		return
	}
	keyVo.Id = apiKey.Id
	keyVo.CreatedAt = apiKey.CreatedAt.Unix()
	resp.SUCCESS(c, keyVo)
}

func (h *ApiKeyHandler) List(c *gin.Context) {
	status := h.GetBool(c, "status")
	t := h.GetTrim(c, "type")

	session := h.DB.Session(&gorm.Session{})
	if status {
		session = session.Where("enabled", true)
	}
	if t != "" {
		session = session.Where("type", t)
	}

	var items []model.ApiKey
	var keys = make([]vo.ApiKey, 0)
	res := session.Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var key vo.ApiKey
			err := utils.CopyObject(item, &key)
			if err == nil {
				key.Id = item.Id
				key.CreatedAt = item.CreatedAt.Unix()
				key.UpdatedAt = item.UpdatedAt.Unix()
				keys = append(keys, key)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, keys)
}

func (h *ApiKeyHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Model(&model.ApiKey{}).Where("id = ?", data.Id).Update(data.Filed, data.Value).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *ApiKeyHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Where("id", id).Delete(&model.ApiKey{}).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}
