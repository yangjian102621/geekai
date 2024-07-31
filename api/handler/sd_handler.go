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
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SdJobHandler struct {
	BaseHandler
	redis     *redis.Client
	pool      *sd.ServicePool
	uploader  *oss.UploaderManager
	snowflake *service.Snowflake
	leveldb   *store.LevelDB
}

func NewSdJobHandler(app *core.AppServer, db *gorm.DB, pool *sd.ServicePool, manager *oss.UploaderManager, snowflake *service.Snowflake, levelDB *store.LevelDB) *SdJobHandler {
	return &SdJobHandler{
		pool:      pool,
		uploader:  manager,
		snowflake: snowflake,
		leveldb:   levelDB,
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *SdJobHandler) Client(c *gin.Context) {
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
	h.pool.Clients.Put(uint(userId), client)
	logger.Infof("New websocket connected, IP: %s", c.RemoteIP())
}

func (h *SdJobHandler) preCheck(c *gin.Context) bool {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if !h.pool.HasAvailableService() {
		resp.ERROR(c, "Stable-Diffusion 池子中没有没有可用的服务！")
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
	params := types.SdTaskParams{
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
	}

	job := model.SdJob{
		UserId:    userId,
		Type:      types.TaskImage.String(),
		TaskId:    params.TaskId,
		Params:    utils.JsonEncode(params),
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

	h.pool.PushTask(types.SdTask{
		Id:     int(job.Id),
		Type:   types.TaskImage,
		Params: params,
		UserId: userId,
	})

	client := h.pool.Clients.Get(uint(job.UserId))
	if client != nil {
		_ = client.Send([]byte("Task Updated"))
	}

	// update user's power
	tx := h.DB.Model(&model.User{}).Where("id = ?", job.UserId).UpdateColumn("power", gorm.Expr("power - ?", job.Power))
	// 记录算力变化日志
	if tx.Error == nil && tx.RowsAffected > 0 {
		user, _ := h.GetLoginUser(c)
		h.DB.Create(&model.PowerLog{
			UserId:    user.Id,
			Username:  user.Username,
			Type:      types.PowerConsume,
			Amount:    job.Power,
			Balance:   user.Power - job.Power,
			Mark:      types.PowerSub,
			Model:     "stable-diffusion",
			Remark:    fmt.Sprintf("绘图操作，任务ID：%s", job.TaskId),
			CreatedAt: time.Now(),
		})
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
func (h *SdJobHandler) getData(finish bool, userId uint, page int, pageSize int, publish bool) (error, []vo.SdJob) {

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

	var items []model.SdJob
	res := session.Find(&items)
	if res.Error != nil {
		return res.Error, nil
	}

	var jobs = make([]vo.SdJob, 0)
	for _, item := range items {
		var job vo.SdJob
		err := utils.CopyObject(item, &job)
		if err != nil {
			continue
		}

		if item.Progress < 100 {
			// 从 leveldb 中获取图片预览数据
			var imageData string
			err = h.leveldb.Get(item.TaskId, &imageData)
			if err == nil {
				job.ImgURL = "data:image/png;base64," + imageData
			}
		}
		jobs = append(jobs, job)
	}

	return nil, jobs
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
			Model:     "stable-diffusion",
			Remark:    fmt.Sprintf("任务失败，退回算力。任务ID：%s， Err: %s", job.TaskId, job.ErrMsg),
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

	client := h.pool.Clients.Get(uint(job.UserId))
	if client != nil {
		_ = client.Send([]byte(sd.Finished))
	}

	resp.SUCCESS(c)
}

// Publish 发布/取消发布图片到画廊显示
func (h *SdJobHandler) Publish(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	action := h.GetBool(c, "action") // 发布动作，true => 发布，false => 取消分享

	res := h.DB.Model(&model.SdJob{Id: uint(id), UserId: int(userId)}).UpdateColumn("publish", action)
	if res.Error != nil {
		logger.Error("error with update database：", res.Error)
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c)
}
