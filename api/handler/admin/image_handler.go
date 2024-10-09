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

type ImageHandler struct {
	handler.BaseHandler
}

func NewImageHandler(app *core.AppServer, db *gorm.DB) *ImageHandler {
	return &ImageHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

type query struct {
	Prompt    string   `json:"prompt"`
	Username  string   `json:"username"`
	CreatedAt []string `json:"created_time"`
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
}

// MjList Midjourney 任务列表
func (h *ImageHandler) MjList(c *gin.Context) {
	var data query
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Username != "" {
		var user model.User
		err := h.DB.Where("username", data.Username).First(&user).Error
		if err == nil {
			session = session.Where("user_id", user.Id)
		}
	}
	if data.Prompt != "" {
		session = session.Where("prompt LIKE ?", "%"+data.Prompt+"%")
	}
	if len(data.CreatedAt) == 2 {
		start := utils.Str2stamp(data.CreatedAt[0] + " 00:00:00")
		end := utils.Str2stamp(data.CreatedAt[1] + " 00:00:00")
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}
	var total int64
	session.Model(&model.MidJourneyJob{}).Count(&total)
	var list []model.MidJourneyJob
	var items = make([]vo.MidJourneyJob, 0)
	offset := (data.Page - 1) * data.PageSize
	err := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&list).Error
	if err == nil {
		// 填充数据
		for _, item := range list {
			var job vo.MidJourneyJob
			err = utils.CopyObject(item, &job)
			if err != nil {
				continue
			}
			job.CreatedAt = item.CreatedAt.Unix()
			items = append(items, job)
		}
	}

	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}

// SdList Stable Diffusion 任务列表
func (h *ImageHandler) SdList(c *gin.Context) {
	var data query
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Username != "" {
		var user model.User
		err := h.DB.Where("username", data.Username).First(&user).Error
		if err == nil {
			session = session.Where("user_id", user.Id)
		}
	}
	if data.Prompt != "" {
		session = session.Where("prompt LIKE ?", "%"+data.Prompt+"%")
	}
	if len(data.CreatedAt) == 2 {
		start := utils.Str2stamp(data.CreatedAt[0] + " 00:00:00")
		end := utils.Str2stamp(data.CreatedAt[1] + " 00:00:00")
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}
	var total int64
	session.Model(&model.SdJob{}).Count(&total)
	var list []model.SdJob
	var items = make([]vo.SdJob, 0)
	offset := (data.Page - 1) * data.PageSize
	err := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&list).Error
	if err == nil {
		// 填充数据
		for _, item := range list {
			var job vo.SdJob
			err = utils.CopyObject(item, &job)
			if err != nil {
				continue
			}
			job.CreatedAt = item.CreatedAt.Unix()
			items = append(items, job)
		}
	}

	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}

// DallList DALL-E 任务列表
func (h *ImageHandler) DallList(c *gin.Context) {
	var data query
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Username != "" {
		var user model.User
		err := h.DB.Where("username", data.Username).First(&user).Error
		if err == nil {
			session = session.Where("user_id", user.Id)
		}
	}
	if data.Prompt != "" {
		session = session.Where("prompt LIKE ?", "%"+data.Prompt+"%")
	}
	if len(data.CreatedAt) == 2 {
		start := utils.Str2stamp(data.CreatedAt[0] + " 00:00:00")
		end := utils.Str2stamp(data.CreatedAt[1] + " 00:00:00")
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}
	var total int64
	session.Model(&model.DallJob{}).Count(&total)
	var list []model.DallJob
	var items = make([]vo.DallJob, 0)
	offset := (data.Page - 1) * data.PageSize
	err := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&list).Error
	if err == nil {
		// 填充数据
		for _, item := range list {
			var job vo.DallJob
			err = utils.CopyObject(item, &job)
			if err != nil {
				continue
			}
			job.CreatedAt = item.CreatedAt.Unix()
			items = append(items, job)
		}
	}

	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}
