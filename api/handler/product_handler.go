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

type ProductHandler struct {
	BaseHandler
}

func NewProductHandler(app *core.AppServer, db *gorm.DB) *ProductHandler {
	return &ProductHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *ProductHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/product/")
	group.GET("list", h.List)
}

// List 模型列表
func (h *ProductHandler) List(c *gin.Context) {
	var items []model.Product
	var list = make([]vo.Product, 0)
	res := h.DB.Where("enabled", true).Order("sort_num ASC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var product vo.Product
			err := utils.CopyObject(item, &product)
			if err == nil {
				product.Id = item.Id
				product.CreatedAt = item.CreatedAt.Unix()
				product.UpdatedAt = item.UpdatedAt.Unix()
				list = append(list, product)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, list)
}
