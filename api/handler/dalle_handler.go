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
	"geekai/service/dalle"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type DallJobHandler struct {
	BaseHandler
	redis    *redis.Client
	service  *dalle.Service
	uploader *oss.UploaderManager
}

func NewDallJobHandler(app *core.AppServer, db *gorm.DB, service *dalle.Service, manager *oss.UploaderManager) *DallJobHandler {
	return &DallJobHandler{
		service:  service,
		uploader: manager,
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *DallJobHandler) Client(c *gin.Context) {
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
	h.service.Clients.Put(uint(userId), client)
	logger.Infof("New websocket connected, IP: %s", c.RemoteIP())
	go func() {
		for {
			_, msg, err := client.Receive()
			if err != nil {
				client.Close()
				h.service.Clients.Delete(uint(userId))
				return
			}

			var message types.WsMessage
			err = utils.JsonDecode(string(msg), &message)
			if err != nil {
				continue
			}

			// 心跳消息
			if message.Type == "heartbeat" {
				logger.Debug("收到 DallE 心跳消息：", message.Content)
				continue
			}
		}
	}()
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
	job := model.DallJob{
		UserId: uint(userId),
		Prompt: data.Prompt,
		Power:  h.App.SysConfig.DallPower,
	}
	res := h.DB.Create(&job)
	if res.Error != nil {
		resp.ERROR(c, "error with save job: "+res.Error.Error())
		return
	}

	h.service.PushTask(types.DallTask{
		JobId:   job.Id,
		UserId:  uint(userId),
		Prompt:  data.Prompt,
		Quality: data.Quality,
		Size:    data.Size,
		Style:   data.Style,
		Power:   job.Power,
	})

	client := h.service.Clients.Get(job.UserId)
	if client != nil {
		_ = client.Send([]byte("Task Updated"))
	}
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
func (h *DallJobHandler) getData(finish bool, userId uint, page int, pageSize int, publish bool) (error, []vo.DallJob) {

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

	var items []model.DallJob
	res := session.Find(&items)
	if res.Error != nil {
		return res.Error, nil
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

	return nil, jobs
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
	tx := h.DB.Begin()
	if err := tx.Delete(&job).Error; err != nil {
		tx.Rollback()
		resp.ERROR(c, err.Error())
		return
	}

	// 如果任务未完成，或者任务失败，则恢复用户算力
	if job.Progress != 100 {
		err := tx.Model(&model.User{}).Where("id = ?", job.UserId).UpdateColumn("power", gorm.Expr("power + ?", job.Power)).Error
		if err != nil {
			tx.Rollback()
			resp.ERROR(c, err.Error())
			return
		}

		var user model.User
		h.DB.Where("id = ?", job.UserId).First(&user)
		err = tx.Create(&model.PowerLog{
			UserId:    user.Id,
			Username:  user.Username,
			Type:      types.PowerConsume,
			Amount:    job.Power,
			Balance:   user.Power,
			Mark:      types.PowerAdd,
			Model:     "dall-e-3",
			Remark:    fmt.Sprintf("任务失败，退回算力。任务ID：%d，Err: %s", job.Id, job.ErrMsg),
			CreatedAt: time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			resp.ERROR(c, err.Error())
			return
		}
	}
	tx.Commit()

	// remove image
	err := h.uploader.GetUploadHandler().Delete(job.ImgURL)
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

	res := h.DB.Model(&model.DallJob{Id: uint(id), UserId: userId}).UpdateColumn("publish", action)
	if res.Error != nil {
		logger.Error("error with update database：", res.Error)
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c)
}
