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
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PowerLogHandler struct {
	handler.BaseHandler
}

func NewPowerLogHandler(app *core.AppServer, db *gorm.DB) *PowerLogHandler {
	return &PowerLogHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

func (h *PowerLogHandler) List(c *gin.Context) {
	var data struct {
		Username string   `json:"username"`
		Type     int      `json:"type"`
		Model    string   `json:"model"`
		Date     []string `json:"date"`
		Page     int      `json:"page"`
		PageSize int      `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Model != "" {
		session = session.Where("model", data.Model)
	}
	if data.Type > 0 {
		session = session.Where("type", data.Type)
	}
	if len(data.Date) == 2 {
		start := data.Date[0] + " 00:00:00"
		end := data.Date[1] + " 00:00:00"
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}

	var total int64
	session.Model(&model.PowerLog{}).Count(&total)
	var items []model.PowerLog
	var list = make([]vo.PowerLog, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var log vo.PowerLog
			err := utils.CopyObject(item, &log)
			if err != nil {
				continue
			}
			log.Id = item.Id
			log.CreatedAt = item.CreatedAt.Unix()
			log.TypeStr = item.Type.String()
			list = append(list, log)
		}
	}

	// 统计消费算力总和
	var totalPower float64
	if len(data.Date) == 2 {
		session.Where("mark", 0).Select("SUM(amount) as total_sum").Scan(&totalPower)
	}
	resp.SUCCESS(c, gin.H{"data": vo.NewPage(total, data.Page, data.PageSize, list), "stat": totalPower})
}
