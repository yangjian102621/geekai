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
	"geekai/service"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MediaHandler struct {
	handler.BaseHandler
	userService *service.UserService
	uploader    *oss.UploaderManager
}

func NewMediaHandler(app *core.AppServer, db *gorm.DB, userService *service.UserService, manager *oss.UploaderManager) *MediaHandler {
	return &MediaHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}, userService: userService, uploader: manager}
}

type mediaQuery struct {
	Prompt    string   `json:"prompt"`
	Username  string   `json:"username"`
	CreatedAt []string `json:"created_at"`
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
}

// SunoList Suno 任务列表
func (h *MediaHandler) SunoList(c *gin.Context) {
	var data mediaQuery
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
		session = session.Where("created_at >= ? AND created_at <= ?", data.CreatedAt[0], data.CreatedAt[1])
	}
	var total int64
	session.Model(&model.SunoJob{}).Count(&total)
	var list []model.SunoJob
	var items = make([]vo.SunoJob, 0)
	offset := (data.Page - 1) * data.PageSize
	err := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&list).Error
	if err == nil {
		// 填充数据
		for _, item := range list {
			var job vo.SunoJob
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

// LumaList Luma 视频任务列表
func (h *MediaHandler) LumaList(c *gin.Context) {
	var data mediaQuery
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
		session = session.Where("created_at >= ? AND created_at <= ?", data.CreatedAt[0], data.CreatedAt[1])
	}
	var total int64
	session.Model(&model.VideoJob{}).Count(&total)
	var list []model.VideoJob
	var items = make([]vo.VideoJob, 0)
	offset := (data.Page - 1) * data.PageSize
	err := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&list).Error
	if err == nil {
		// 填充数据
		for _, item := range list {
			var job vo.VideoJob
			err = utils.CopyObject(item, &job)
			if err != nil {
				continue
			}
			job.CreatedAt = item.CreatedAt.Unix()
			if job.VideoURL == "" {
				job.VideoURL = job.WaterURL
			}
			items = append(items, job)
		}
	}

	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}

func (h *MediaHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	tab := c.Query("tab")

	tx := h.DB.Begin()
	var md, remark, fileURL string
	var power, userId, progress int
	switch tab {
	case "suno":
		var job model.SunoJob
		if err := h.DB.Where("id", id).First(&job).Error; err != nil {
			resp.ERROR(c, "记录不存在")
			return
		}
		tx.Delete(&job)
		md = "suno"
		power = job.Power
		userId = job.UserId
		remark = fmt.Sprintf("SUNO 任务失败，退回算力。任务ID：%d，Err: %s", job.Id, job.ErrMsg)
		progress = job.Progress
		fileURL = job.AudioURL
		break
	case "luma":
		var job model.VideoJob
		if res := h.DB.Where("id", id).First(&job); res.Error != nil {
			resp.ERROR(c, "记录不存在")
			return
		}

		// 删除任务
		tx.Delete(&job)
		md = job.Type
		power = job.Power
		userId = job.UserId
		remark = fmt.Sprintf("LUMA 任务失败，退回算力。任务ID：%d，Err: %s", job.Id, job.ErrMsg)
		progress = job.Progress
		fileURL = job.VideoURL
		if fileURL == "" {
			fileURL = job.WaterURL
		}
		break
	default:
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if progress != 100 {
		err := h.userService.IncreasePower(userId, power, model.PowerLog{
			Type:   types.PowerRefund,
			Model:  md,
			Remark: remark,
		})
		if err != nil {
			tx.Rollback()
			resp.ERROR(c, err.Error())
			return
		}
	}
	tx.Commit()
	// remove image
	err := h.uploader.GetUploadHandler().Delete(fileURL)
	if err != nil {
		logger.Error("remove image failed: ", err)
	}

	resp.SUCCESS(c)
}
