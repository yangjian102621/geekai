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
	"geekai/service/suno"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type SunoHandler struct {
	BaseHandler
	service *suno.Service
}

func NewSunoHandler(app *core.AppServer, db *gorm.DB) *SunoHandler {
	return &SunoHandler{
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *SunoHandler) Client(c *gin.Context) {
	//ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	//if err != nil {
	//	logger.Error(err)
	//	c.Abort()
	//	return
	//}
	//
	//userId := h.GetInt(c, "user_id", 0)
	//if userId == 0 {
	//	logger.Info("Invalid user ID")
	//	c.Abort()
	//	return
	//}
	//
	////client := types.NewWsClient(ws)
	//logger.Infof("New websocket connected, IP: %s", c.RemoteIP())
}

func (h *SunoHandler) Create(c *gin.Context) {

	var data struct {
		Prompt       string `json:"prompt"`
		Instrumental bool   `json:"instrumental"`
		Lyrics       string `json:"lyrics"`
		Model        string `json:"model"`
		Tags         string `json:"tags"`
		Title        string `json:"title"`
		Type         int    `json:"type"`
		RefTaskId    string `json:"ref_task_id"` // 续写的任务id
		ExtendSecs   int    `json:"extend_secs"` // 续写秒数
		RefSongId    string `json:"ref_song_id"` // 续写的歌曲id
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 插入数据库
	job := model.SunoJob{
		UserId:       int(h.GetLoginUserId(c)),
		Prompt:       data.Prompt,
		Instrumental: data.Instrumental,
		ModelName:    data.Model,
		Tags:         data.Tags,
		Title:        data.Title,
		Type:         data.Type,
		RefSongId:    data.RefSongId,
		RefTaskId:    data.RefTaskId,
		ExtendSecs:   data.ExtendSecs,
		Power:        h.App.SysConfig.SunoPower,
	}
	tx := h.DB.Create(&job)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	// 创建任务
	h.service.PushTask(types.SunoTask{
		Id:           job.Id,
		UserId:       job.UserId,
		Type:         job.Type,
		Title:        job.Title,
		Lyrics:       data.Lyrics,
		RefTaskId:    data.RefTaskId,
		RefSongId:    data.RefSongId,
		ExtendSecs:   data.ExtendSecs,
		Prompt:       data.Prompt,
		Tags:         data.Tags,
		Model:        data.Model,
		Instrumental: data.Instrumental,
	})

	// update user's power
	tx = h.DB.Model(&model.User{}).Where("id = ?", job.UserId).UpdateColumn("power", gorm.Expr("power - ?", job.Power))
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
			Model:     job.ModelName,
			Remark:    fmt.Sprintf("Suno 文生歌曲，%s", job.ModelName),
			CreatedAt: time.Now(),
		})
	}

	var itemVo vo.SunoJob
	_ = utils.CopyObject(job, &itemVo)
	resp.SUCCESS(c, itemVo)
}

func (h *SunoHandler) List(c *gin.Context) {

}

func (h *SunoHandler) Remove(c *gin.Context) {

}

func (h *SunoHandler) Publish(c *gin.Context) {

}
