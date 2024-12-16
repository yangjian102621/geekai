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
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
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

// Client WebSocket 客户端，用于通知任务状态变更
func (h *VideoHandler) Client(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Abort()
		return
	}

	userId := h.GetInt(c, "user_id", 0)
	if userId == 0 {
		logger.Info("Invalid user ID")
		c.Abort()
		return
	}

	client := types.NewWsClient(ws)
	h.videoService.Clients.Put(uint(userId), client)
	logger.Infof("New websocket connected, IP: %s", c.RemoteIP())
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
	if data.Prompt == "" {
		resp.ERROR(c, "prompt is needed")
		return
	}

	userId := int(h.GetLoginUserId(c))
	params := types.VideoParams{
		PromptOptimize: data.ExpandPrompt,
		Loop:           data.Loop,
		StartImgURL:    data.FirstFrameImg,
		EndImgURL:      data.EndFrameImg,
	}
	// 插入数据库
	job := model.VideoJob{
		UserId: userId,
		Type:   types.VideoLuma,
		Prompt: data.Prompt,
		Power:  h.App.SysConfig.LumaPower,
		Params: utils.JsonEncode(params),
	}
	tx := h.DB.Create(&job)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	// 创建任务
	h.videoService.PushTask(types.VideoTask{
		Id:     job.Id,
		UserId: userId,
		Type:   types.VideoLuma,
		Prompt: data.Prompt,
		Params: params,
	})

	// update user's power
	err := h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "luma",
		Remark: fmt.Sprintf("Luma 文生视频，任务ID：%d", job.Id),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	client := h.videoService.Clients.Get(uint(job.UserId))
	if client != nil {
		_ = client.Send([]byte("Task Updated"))
	}
	resp.SUCCESS(c)
}

func (h *VideoHandler) List(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	t := c.Query("type")
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	all := h.GetBool(c, "all")
	session := h.DB.Session(&gorm.Session{}).Where("user_id", userId)
	if t != "" {
		session = session.Where("type", t)
	}
	if all {
		session = session.Where("publish", 0).Where("progress", 100)
	} else {
		session = session.Where("user_id", h.GetLoginUserId(c))
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
	// 删除任务
	tx := h.DB.Begin()
	if err := tx.Delete(&job).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, err.Error())
		return
	}

	// 如果任务未完成，或者任务失败，则恢复用户算力
	if job.Progress != 100 {
		err = h.userService.IncreasePower(job.UserId, job.Power, model.PowerLog{
			Type:   types.PowerRefund,
			Model:  "luma",
			Remark: fmt.Sprintf("Luma 任务失败，退回算力。任务ID：%s，Err:%s", job.TaskId, job.ErrMsg),
		})
		if err != nil {
			tx.Rollback()
			resp.ERROR(c, err.Error())
			return
		}
	}
	tx.Commit()

	// 删除文件
	_ = h.uploader.GetUploadHandler().Delete(job.CoverURL)
	_ = h.uploader.GetUploadHandler().Delete(job.VideoURL)
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
