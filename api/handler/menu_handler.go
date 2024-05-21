package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MenuHandler struct {
	BaseHandler
}

func NewMenuHandler(app *core.AppServer, db *gorm.DB) *MenuHandler {
	return &MenuHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// List 数据列表
func (h *MenuHandler) List(c *gin.Context) {
	var items []model.Menu
	var list = make([]vo.Menu, 0)
	res := h.DB.Where("enabled", true).Order("sort_num ASC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var product vo.Menu
			err := utils.CopyObject(item, &product)
			if err == nil {
				list = append(list, product)
			}
		}
	}
	resp.SUCCESS(c, list)
}
