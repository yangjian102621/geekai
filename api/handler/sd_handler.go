package handler

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
	"geekai/service"
	"geekai/service/oss"
	"geekai/service/sd"
	"geekai/store"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SdJobHandler struct {
	BaseHandler
	redis       *redis.Client
	sdService   *sd.Service
	uploader    *oss.UploaderManager
	snowflake   *service.Snowflake
	leveldb     *store.LevelDB
	userService *service.UserService
}

func NewSdJobHandler(app *core.AppServer,
	db *gorm.DB,
	service *sd.Service,
	manager *oss.UploaderManager,
	snowflake *service.Snowflake,
	userService *service.UserService,
	levelDB *store.LevelDB) *SdJobHandler {
	return &SdJobHandler{
		sdService:   service,
		uploader:    manager,
		snowflake:   snowflake,
		leveldb:     levelDB,
		userService: userService,
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

func (h *SdJobHandler) preCheck(c *gin.Context) bool {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if user.Power < h.App.SysConfig.SdPower {
		resp.ERROR(c, "当前用户剩余算力不足以完成本次绘画！")
		return false
	}

	return true

}

// Image 创建一个绘画任务
func (h *SdJobHandler) Image(c *gin.Context) {
	if !h.preCheck(c) {
		return
	}

	var data types.SdTaskParams
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Width <= 0 {
		data.Width = 512
	}
	if data.Height <= 0 {
		data.Height = 512
	}
	if data.CfgScale <= 0 {
		data.CfgScale = 7
	}
	if data.Seed == 0 {
		data.Seed = -1
	}
	if data.Steps <= 0 {
		data.Steps = 20
	}
	if data.Sampler == "" {
		data.Sampler = "Euler a"
	}
	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	taskId, err := h.snowflake.Next(true)
	if err != nil {
		resp.ERROR(c, "error with generate task id: "+err.Error())
		return
	}

	task := types.SdTask{
		ClientId: data.ClientId,
		Type:     types.TaskImage,
		Params: types.SdTaskParams{
			TaskId:       taskId,
			Prompt:       data.Prompt,
			NegPrompt:    data.NegPrompt,
			Steps:        data.Steps,
			Sampler:      data.Sampler,
			FaceFix:      data.FaceFix,
			CfgScale:     data.CfgScale,
			Seed:         data.Seed,
			Height:       data.Height,
			Width:        data.Width,
			HdFix:        data.HdFix,
			HdRedrawRate: data.HdRedrawRate,
			HdScale:      data.HdScale,
			HdScaleAlg:   data.HdScaleAlg,
			HdSteps:      data.HdSteps,
		},
		UserId:           userId,
		TranslateModelId: h.App.SysConfig.TranslateModelId,
	}

	job := model.SdJob{
		UserId:    userId,
		Type:      types.TaskImage.String(),
		TaskId:    taskId,
		Params:    utils.JsonEncode(task.Params),
		TaskInfo:  utils.JsonEncode(task),
		Prompt:    data.Prompt,
		Progress:  0,
		Power:     h.App.SysConfig.SdPower,
		CreatedAt: time.Now(),
	}
	res := h.DB.Create(&job)
	if res.Error != nil {
		resp.ERROR(c, "error with save job: "+res.Error.Error())
		return
	}

	task.Id = int(job.Id)
	h.sdService.PushTask(task)

	// update user's power
	err = h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "stable-diffusion",
		Remark: fmt.Sprintf("绘图操作，任务ID：%s", job.TaskId),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// ImgWall 照片墙
func (h *SdJobHandler) ImgWall(c *gin.Context) {
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
func (h *SdJobHandler) JobList(c *gin.Context) {
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

// JobList 获取 MJ 任务列表
func (h *SdJobHandler) getData(finish bool, userId uint, page int, pageSize int, publish bool) (error, vo.Page) {

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
	session.Model(&model.SdJob{}).Count(&total)

	var items []model.SdJob
	res := session.Find(&items)
	if res.Error != nil {
		return res.Error, vo.Page{}
	}

	var jobs = make([]vo.SdJob, 0)
	for _, item := range items {
		var job vo.SdJob
		err := utils.CopyObject(item, &job)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}

	return nil, vo.NewPage(total, page, pageSize, jobs)
}

// Remove remove task image
func (h *SdJobHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var job model.SdJob
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
func (h *SdJobHandler) Publish(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	action := h.GetBool(c, "action") // 发布动作，true => 发布，false => 取消分享

	err := h.DB.Model(&model.SdJob{Id: uint(id), UserId: int(userId)}).UpdateColumn("publish", action).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
