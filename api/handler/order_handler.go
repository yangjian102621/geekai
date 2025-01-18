package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	BaseHandler
}

func NewOrderHandler(app *core.AppServer, db *gorm.DB) *OrderHandler {
	return &OrderHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// List 订单列表
func (h *OrderHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	userId := h.GetLoginUserId(c)
	session := h.DB.Session(&gorm.Session{}).Where("user_id = ? AND status = ?", userId, types.OrderPaidSuccess)
	var total int64
	session.Model(&model.Order{}).Count(&total)
	var items []model.Order
	var list = make([]vo.Order, 0)
	offset := (page - 1) * pageSize
	res := session.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var order vo.Order
			err := utils.CopyObject(item, &order)
			if err == nil {
				order.Id = item.Id
				order.CreatedAt = item.CreatedAt.Unix()
				order.UpdatedAt = item.UpdatedAt.Unix()
				payMethod, ok := types.PayMethods[item.PayWay]
				if !ok {
					payMethod = item.PayWay
				}
				payName, ok := types.PayNames[item.PayType]
				if !ok {
					payName = item.PayWay
				}
				order.PayMethod = payMethod
				order.PayName = payName
				list = append(list, order)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, list))
}

// Query 查询订单状态
func (h *OrderHandler) Query(c *gin.Context) {
	orderNo := h.GetTrim(c, "order_no")
	var order model.Order
	res := h.DB.Where("order_no = ?", orderNo).First(&order)
	if res.Error != nil {
		resp.ERROR(c, "Order not found")
		return
	}

	if order.Status == types.OrderPaidSuccess {
		resp.SUCCESS(c, gin.H{"status": order.Status})
		return
	}

	counter := 0
	for {
		time.Sleep(time.Second)
		var item model.Order
		h.DB.Where("order_no = ?", orderNo).First(&item)
		if counter >= 15 || item.Status == types.OrderPaidSuccess || item.Status != order.Status {
			order.Status = item.Status
			break
		}
		counter++
	}

	resp.SUCCESS(c, gin.H{"status": order.Status})
}
