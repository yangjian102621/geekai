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
	"geekai/core/middleware"
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

type SunoHandler struct {
	handler.BaseHandler
	userService *service.UserService
	uploader    *oss.UploaderManager
}

func NewSunoHandler(app *core.AppServer, db *gorm.DB, userService *service.UserService, manager *oss.UploaderManager) *SunoHandler {
	return &SunoHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}, userService: userService, uploader: manager}
}

// RegisterRoutes 注册路由
func (h *SunoHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/admin/suno/")

	// 需要管理员授权的接口
	group.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		group.POST("list", h.SunoList)
		group.GET("remove", h.Remove)
	}
}

type sunoQuery struct {
	Title     string   `json:"title"`
	Prompt    string   `json:"prompt"`
	CreatedAt []string `json:"created_at"`
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
}

// SunoList Suno 任务列表
func (h *SunoHandler) SunoList(c *gin.Context) {
	var data sunoQuery
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Title != "" {
		session = session.Where("title LIKE ?", "%"+data.Title+"%")
	}
	if data.Prompt != "" {
		// 同时查询 prompt 字段和 params JSON 字段中的 prompt
		session = session.Where("prompt LIKE ? OR JSON_EXTRACT(params, '$.prompt') LIKE ?", "%"+data.Prompt+"%", "%"+data.Prompt+"%")
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

func (h *SunoHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	tx := h.DB.Begin()
	var job model.SunoJob
	if err := h.DB.Where("id", id).First(&job).Error; err != nil {
		resp.ERROR(c, "记录不存在")
		return
	}

	// 删除任务
	tx.Delete(&job)
	md := "suno"
	power := job.Power
	userId := int(job.UserId)
	remark := fmt.Sprintf("SUNO 任务失败，退回算力。任务ID：%d，Err: %s", job.Id, job.ErrMsg)
	needRefund := job.Progress != 100
	fileURL := job.AudioURL

	if needRefund {
		err := h.userService.IncreasePower(uint(userId), power, model.PowerLog{
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
	// remove file
	err := h.uploader.GetUploadHandler().Delete(fileURL)
	if err != nil {
		logger.Error("remove file failed: ", err)
	}

	resp.SUCCESS(c)
}
