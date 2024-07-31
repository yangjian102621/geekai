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
	"geekai/service/oss"
	"geekai/service/suno"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type SunoHandler struct {
	BaseHandler
	service  *suno.Service
	uploader *oss.UploaderManager
}

func NewSunoHandler(app *core.AppServer, db *gorm.DB, service *suno.Service, uploader *oss.UploaderManager) *SunoHandler {
	return &SunoHandler{
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
		service:  service,
		uploader: uploader,
	}
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *SunoHandler) Client(c *gin.Context) {
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
	if data.Lyrics != "" {
		job.Prompt = data.Lyrics
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
		RefTaskId:    data.RefTaskId,
		RefSongId:    data.RefSongId,
		ExtendSecs:   data.ExtendSecs,
		Prompt:       job.Prompt,
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

	client := h.service.Clients.Get(uint(job.UserId))
	if client != nil {
		_ = client.Send([]byte("Task Updated"))
	}
	resp.SUCCESS(c)
}

func (h *SunoHandler) List(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)
	session := h.DB.Session(&gorm.Session{}).Where("user_id", userId)

	// 统计总数
	var total int64
	session.Model(&model.SunoJob{}).Count(&total)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}
	var list []model.SunoJob
	err := session.Order("id desc").Find(&list).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 初始化续写关系
	songIds := make([]string, 0)
	for _, v := range list {
		if v.RefTaskId != "" {
			songIds = append(songIds, v.RefSongId)
		}
	}
	var tasks []model.SunoJob
	h.DB.Where("song_id IN ?", songIds).Find(&tasks)
	songMap := make(map[string]model.SunoJob)
	for _, t := range tasks {
		songMap[t.SongId] = t
	}
	// 转换为 VO
	items := make([]vo.SunoJob, 0)
	for _, v := range list {
		var item vo.SunoJob
		err = utils.CopyObject(v, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = v.CreatedAt.Unix()
		if s, ok := songMap[v.RefSongId]; ok {
			item.RefSong = map[string]interface{}{
				"id":    s.Id,
				"title": s.Title,
				"cover": s.CoverURL,
				"audio": s.AudioURL,
			}
		}
		items = append(items, item)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, items))
}

func (h *SunoHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var job model.SunoJob
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
			Model:     job.ModelName,
			Remark:    fmt.Sprintf("Suno 任务失败，退回算力。任务ID：%s，Err:%s", job.TaskId, job.ErrMsg),
			CreatedAt: time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			resp.ERROR(c, err.Error())
			return
		}
	}
	tx.Commit()

	// 删除文件
	_ = h.uploader.GetUploadHandler().Delete(job.CoverURL)
	_ = h.uploader.GetUploadHandler().Delete(job.AudioURL)
}

func (h *SunoHandler) Publish(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	publish := h.GetBool(c, "publish")
	err := h.DB.Model(&model.SunoJob{}).Where("id", id).Where("user_id", userId).UpdateColumn("publish", publish).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

func (h *SunoHandler) Update(c *gin.Context) {
	var data struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Cover string `json:"cover"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Id == 0 || data.Title == "" || data.Cover == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	userId := h.GetLoginUserId(c)
	var item model.SunoJob
	if err := h.DB.Where("id", data.Id).Where("user_id", userId).First(&item).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	item.Title = data.Title
	item.CoverURL = data.Cover

	if err := h.DB.Updates(&item).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// Detail 歌曲详情
func (h *SunoHandler) Detail(c *gin.Context) {
	songId := c.Query("song_id")
	if songId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var item model.SunoJob
	if err := h.DB.Where("song_id", songId).First(&item).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 读取用户信息
	var user model.User
	if err := h.DB.Where("id", item.UserId).First(&user).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	var itemVo vo.SunoJob
	if err := utils.CopyObject(item, &itemVo); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	itemVo.CreatedAt = item.CreatedAt.Unix()
	itemVo.User = map[string]interface{}{
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
	}

	resp.SUCCESS(c, itemVo)
}

// Play 增加歌曲播放次数
func (h *SunoHandler) Play(c *gin.Context) {
	songId := c.Query("song_id")
	if songId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	h.DB.Model(&model.SunoJob{}).Where("song_id", songId).UpdateColumn("play_times", gorm.Expr("play_times + ?", 1))
}

const genLyricTemplate = `
你是一位才华横溢的作曲家，拥有丰富的情感和细腻的笔触，你对文字有着独特的感悟力，能将各种情感和意境巧妙地融入歌词中。
请以【%s】为主题创作一首歌曲，歌曲时间不要太短，3分钟左右，不要输出任何解释性的内容。
输出格式如下：
歌曲名称
第一节：
{{歌词内容}}
副歌：
{{歌词内容}}

第二节：
{{歌词内容}}
副歌：
{{歌词内容}}

尾声：
{{歌词内容}}
`

// Lyric 生成歌词
func (h *SunoHandler) Lyric(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	content, err := utils.OpenAIRequest(h.DB, fmt.Sprintf(genLyricTemplate, data.Prompt), "gpt-4o-mini")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, content)
}
