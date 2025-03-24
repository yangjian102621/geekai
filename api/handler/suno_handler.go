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
	sunoService *suno.Service
	uploader    *oss.UploaderManager
	userService *service.UserService
}

func NewSunoHandler(app *core.AppServer, db *gorm.DB, service *suno.Service, uploader *oss.UploaderManager, userService *service.UserService) *SunoHandler {
	return &SunoHandler{
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
		sunoService: service,
		uploader:    uploader,
		userService: userService,
	}
}

func (h *SunoHandler) Create(c *gin.Context) {

	var data struct {
		ClientId     string `json:"client_id"`
		Prompt       string `json:"prompt"`
		Instrumental bool   `json:"instrumental"`
		Lyrics       string `json:"lyrics"`
		Model        string `json:"model"`
		Tags         string `json:"tags"`
		Title        string `json:"title"`
		Type         int    `json:"type"`
		RefTaskId    string `json:"ref_task_id"`         // 续写的任务id
		ExtendSecs   int    `json:"extend_secs"`         // 续写秒数
		RefSongId    string `json:"ref_song_id"`         // 续写的歌曲id
		SongId       string `json:"song_id,omitempty"`   // 要拼接的歌曲id
		AudioURL     string `json:"audio_url,omitempty"` // 上传自己创作的歌曲
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

	if user.Power < h.App.SysConfig.SunoPower {
		resp.ERROR(c, "您的算力不足，请充值后再试！")
		return
	}

	// 歌曲拼接
	if data.SongId != "" && data.Type == 3 {
		var song model.SunoJob
		if err := h.DB.Where("song_id = ?", data.SongId).First(&song).Error; err == nil {
			data.Instrumental = song.Instrumental
			data.Model = song.ModelName
			data.Tags = song.Tags
		}
		// 拼接歌词
		var refSong model.SunoJob
		if err := h.DB.Where("song_id = ?", data.RefSongId).First(&refSong).Error; err == nil {
			data.Prompt = fmt.Sprintf("%s\n%s", song.Prompt, refSong.Prompt)
		}
	}
	task := types.SunoTask{
		ClientId:     data.ClientId,
		UserId:       int(h.GetLoginUserId(c)),
		Type:         data.Type,
		Title:        data.Title,
		RefTaskId:    data.RefTaskId,
		RefSongId:    data.RefSongId,
		ExtendSecs:   data.ExtendSecs,
		Prompt:       data.Prompt,
		Tags:         data.Tags,
		Model:        data.Model,
		Instrumental: data.Instrumental,
		SongId:       data.SongId,
		AudioURL:     data.AudioURL,
	}

	// 插入数据库
	job := model.SunoJob{
		UserId:       task.UserId,
		Prompt:       data.Prompt,
		Instrumental: data.Instrumental,
		ModelName:    data.Model,
		TaskInfo:     utils.JsonEncode(task),
		Tags:         data.Tags,
		Title:        data.Title,
		Type:         data.Type,
		RefSongId:    data.RefSongId,
		RefTaskId:    data.RefTaskId,
		ExtendSecs:   data.ExtendSecs,
		Power:        h.App.SysConfig.SunoPower,
		SongId:       utils.RandString(32),
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
	task.Id = job.Id
	h.sunoService.PushTask(task)

	// update user's power
	err = h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:      types.PowerConsume,
		Model:     job.ModelName,
		Remark:    fmt.Sprintf("Suno 文生歌曲，%s", job.ModelName),
		CreatedAt: time.Now(),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

func (h *SunoHandler) List(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
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

	// 只有失败，或者超时的任务才能删除
	if job.Progress != service.FailTaskProgress || time.Now().Before(job.CreatedAt.Add(time.Minute*10)) {
		resp.ERROR(c, "只有失败和超时(10分钟)的任务才能删除！")
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
