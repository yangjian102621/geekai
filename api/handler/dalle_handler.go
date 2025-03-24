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
	"geekai/service"
	"geekai/service/dalle"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type DallJobHandler struct {
	BaseHandler
	redis       *redis.Client
	dallService *dalle.Service
	uploader    *oss.UploaderManager
	userService *service.UserService
}

func NewDallJobHandler(app *core.AppServer, db *gorm.DB, service *dalle.Service, manager *oss.UploaderManager, userService *service.UserService) *DallJobHandler {
	return &DallJobHandler{
		dallService: service,
		uploader:    manager,
		userService: userService,
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

func (h *DallJobHandler) preCheck(c *gin.Context) bool {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return false
	}
	if user.Power < h.App.SysConfig.DallPower {
		resp.ERROR(c, "当前用户剩余算力不足以完成本次绘画！")
		return false
	}

	return true

}

// Image 创建一个绘画任务
func (h *DallJobHandler) Image(c *gin.Context) {
	if !h.preCheck(c) {
		return
	}

	var data types.DallTask
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	task := types.DallTask{
		ClientId:         data.ClientId,
		UserId:           uint(userId),
		Prompt:           data.Prompt,
		Quality:          data.Quality,
		Size:             data.Size,
		Style:            data.Style,
		Power:            h.App.SysConfig.DallPower,
		TranslateModelId: h.App.SysConfig.TranslateModelId,
	}
	job := model.DallJob{
		UserId:   uint(userId),
		Prompt:   data.Prompt,
		Power:    task.Power,
		TaskInfo: utils.JsonEncode(task),
	}
	res := h.DB.Create(&job)
	if res.Error != nil {
		resp.ERROR(c, "error with save job: "+res.Error.Error())
		return
	}

	task.Id = job.Id
	h.dallService.PushTask(task)
	resp.SUCCESS(c)
}

// ImgWall 照片墙
func (h *DallJobHandler) ImgWall(c *gin.Context) {
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)
	err, jobs := h.getData(true, 0, page, pageSize, true)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, jobs)
}

// JobList 获取 SD 任务列表
func (h *DallJobHandler) JobList(c *gin.Context) {
	finish := h.GetBool(c, "finish")
	userId := h.GetLoginUserId(c)
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)
	publish := h.GetBool(c, "publish")

	err, jobs := h.getData(finish, userId, page, pageSize, publish)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, jobs)
}

// JobList 获取任务列表
func (h *DallJobHandler) getData(finish bool, userId uint, page int, pageSize int, publish bool) (error, vo.Page) {

	session := h.DB.Session(&gorm.Session{})
	if finish {
		session = session.Where("progress >= ?", 100).Order("id DESC")
	} else {
		session = session.Where("progress < ?", 100).Order("id ASC")
	}
	if userId > 0 {
		session = session.Where("user_id = ?", userId)
	}
	if publish {
		session = session.Where("publish", publish)
	}
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}
	// 统计总数
	var total int64
	session.Model(&model.DallJob{}).Count(&total)

	var items []model.DallJob
	res := session.Find(&items)
	if res.Error != nil {
		return res.Error, vo.Page{}
	}

	var jobs = make([]vo.DallJob, 0)
	for _, item := range items {
		var job vo.DallJob
		err := utils.CopyObject(item, &job)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}

	return nil, vo.NewPage(total, page, pageSize, jobs)
}

// Remove remove task image
func (h *DallJobHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var job model.DallJob
	if res := h.DB.Where("id = ? AND user_id = ?", id, userId).First(&job); res.Error != nil {
		resp.ERROR(c, "记录不存在")
		return
	}

	// 删除任务
	err := h.DB.Delete(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// remove image
	err = h.uploader.GetUploadHandler().Delete(job.ImgURL)
	if err != nil {
		logger.Error("remove image failed: ", err)
	}

	resp.SUCCESS(c)
}

// Publish 发布/取消发布图片到画廊显示
func (h *DallJobHandler) Publish(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	action := h.GetBool(c, "action") // 发布动作，true => 发布，false => 取消分享

	err := h.DB.Model(&model.DallJob{Id: uint(id), UserId: userId}).UpdateColumn("publish", action).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
