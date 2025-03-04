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
	"geekai/service/video"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoHandler struct {
	BaseHandler
	videoService *video.Service
	uploader     *oss.UploaderManager
	userService  *service.UserService
}

func NewVideoHandler(app *core.AppServer, db *gorm.DB, service *video.Service, uploader *oss.UploaderManager, userService *service.UserService) *VideoHandler {
	return &VideoHandler{
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
		videoService: service,
		uploader:     uploader,
		userService:  userService,
	}
}

func (h *VideoHandler) LumaCreate(c *gin.Context) {

	var data struct {
		Prompt        string `json:"prompt"`
		FirstFrameImg string `json:"first_frame_img,omitempty"`
		EndFrameImg   string `json:"end_frame_img,omitempty"`
		ExpandPrompt  bool   `json:"expand_prompt,omitempty"`
		Loop          bool   `json:"loop,omitempty"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 检查 Prompt 长度
	if data.Prompt == "" {
		resp.ERROR(c, "prompt is needed")
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	if user.Power < h.App.SysConfig.LumaPower {
		resp.ERROR(c, "您的算力不足，请充值后再试！")
		return
	}

	userId := int(h.GetLoginUserId(c))
	params := types.LumaVideoParams{
		PromptOptimize: data.ExpandPrompt,
		Loop:           data.Loop,
		StartImgURL:    data.FirstFrameImg,
		EndImgURL:      data.EndFrameImg,
	}
	task := types.VideoTask{
		UserId:           userId,
		Type:             types.VideoLuma,
		Prompt:           data.Prompt,
		Params:           params,
		TranslateModelId: h.App.SysConfig.TranslateModelId,
	}
	// 插入数据库
	job := model.VideoJob{
		UserId:   userId,
		Type:     types.VideoLuma,
		Prompt:   data.Prompt,
		Power:    h.App.SysConfig.LumaPower,
		TaskInfo: utils.JsonEncode(task),
	}
	tx := h.DB.Create(&job)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	// 创建任务
	task.Id = job.Id
	h.videoService.PushTask(task)

	// update user's power
	err = h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "luma",
		Remark: fmt.Sprintf("Luma 文生视频，任务ID：%d", job.Id),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *VideoHandler) KeLingCreate(c *gin.Context) {

	var data struct {
		Channel       string              `json:"channel"`
		TaskType      string              `json:"task_type"`       // 任务类型: text2video/image2video
		Model         string              `json:"model"`           // 模型: kling-v1-5,kling-v1-6
		Prompt        string              `json:"prompt"`          // 视频描述
		NegPrompt     string              `json:"negative_prompt"` // 负面提示词
		CfgScale      float64             `json:"cfg_scale"`       // 相关性系数(0-1)
		Mode          string              `json:"mode"`            // 生成模式: std/pro
		AspectRatio   string              `json:"aspect_ratio"`    // 画面比例: 16:9/9:16/1:1
		Duration      string              `json:"duration"`        // 视频时长: 5/10
		CameraControl types.CameraControl `json:"camera_control"`  // 摄像机控制
		Image         string              `json:"image"`           // 参考图片URL(image2video)
		ImageTail     string              `json:"image_tail"`      // 尾帧图片URL(image2video)
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 计算当前任务所需算力
	key := fmt.Sprintf("%s_%s_%s", data.Model, data.Mode, data.Duration)
	power := h.App.SysConfig.KeLingPowers[key]
	if power == 0 {
		resp.ERROR(c, "当前模型暂不支持")
		return
	}
	if user.Power < power {
		resp.ERROR(c, "您的算力不足，请充值后再试！")
		return
	}

	if data.Prompt == "" {
		resp.ERROR(c, "prompt is needed")
		return
	}

	userId := int(h.GetLoginUserId(c))
	params := types.KeLingVideoParams{
		TaskType:      data.TaskType,
		Model:         data.Model,
		Prompt:        data.Prompt,
		NegPrompt:     data.NegPrompt,
		CfgScale:      data.CfgScale,
		Mode:          data.Mode,
		AspectRatio:   data.AspectRatio,
		Duration:      data.Duration,
		CameraControl: data.CameraControl,
		Image:         data.Image,
		ImageTail:     data.ImageTail,
	}
	task := types.VideoTask{
		UserId:           userId,
		Type:             types.VideoKeLing,
		Prompt:           data.Prompt,
		Params:           params,
		TranslateModelId: h.App.SysConfig.TranslateModelId,
		Channel:          data.Channel,
	}
	// 插入数据库
	job := model.VideoJob{
		UserId:   userId,
		Type:     types.VideoKeLing,
		Prompt:   data.Prompt,
		Power:    power,
		TaskInfo: utils.JsonEncode(task),
	}
	tx := h.DB.Create(&job)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	// 创建任务
	task.Id = job.Id
	h.videoService.PushTask(task)

	// update user's power
	err = h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "keling",
		Remark: fmt.Sprintf("keling 文生视频，任务ID：%d", job.Id),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *VideoHandler) List(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	t := c.Query("type")
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	all := h.GetBool(c, "all")
	session := h.DB.Session(&gorm.Session{})
	if t != "" {
		session = session.Where("type", t)
	}
	if all {
		session = session.Where("publish", 0).Where("progress", 100)
	} else {
		session = session.Where("user_id", userId)
	}
	// 统计总数
	var total int64
	session.Model(&model.VideoJob{}).Count(&total)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}
	var list []model.VideoJob
	err := session.Order("id desc").Find(&list).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 转换为 VO
	items := make([]vo.VideoJob, 0)
	for _, v := range list {
		var item vo.VideoJob
		err = utils.CopyObject(v, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = v.CreatedAt.Unix()
		if item.VideoURL == "" {
			item.VideoURL = v.WaterURL
		}
		items = append(items, item)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, items))
}

func (h *VideoHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var job model.VideoJob
	err := h.DB.Where("id = ?", id).Where("user_id", userId).First(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 只有失败或者超时的任务才能删除
	if !(job.Progress == service.FailTaskProgress || time.Now().After(job.CreatedAt.Add(time.Minute*30))) {
		resp.ERROR(c, "只有失败和超时(30分钟)的任务才能删除！")
		return
	}

	// 删除任务
	err = h.DB.Delete(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 删除文件
	_ = h.uploader.GetUploadHandler().Delete(job.CoverURL)
	_ = h.uploader.GetUploadHandler().Delete(job.VideoURL)

	resp.SUCCESS(c)
}

func (h *VideoHandler) Publish(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	publish := h.GetBool(c, "publish")
	var job model.VideoJob
	err := h.DB.Where("id = ?", id).Where("user_id", userId).First(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	err = h.DB.Model(&job).UpdateColumn("publish", publish).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
