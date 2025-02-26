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
	"geekai/service/mj"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MidJourneyHandler struct {
	BaseHandler
	mjService   *mj.Service
	snowflake   *service.Snowflake
	uploader    *oss.UploaderManager
	userService *service.UserService
}

func NewMidJourneyHandler(app *core.AppServer, db *gorm.DB, snowflake *service.Snowflake, service *mj.Service, manager *oss.UploaderManager, userService *service.UserService) *MidJourneyHandler {
	return &MidJourneyHandler{
		snowflake:   snowflake,
		mjService:   service,
		uploader:    manager,
		userService: userService,
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

func (h *MidJourneyHandler) preCheck(c *gin.Context) bool {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if user.Power < h.App.SysConfig.MjPower {
		resp.ERROR(c, "当前用户剩余算力不足以完成本次绘画！")
		return false
	}

	return true

}

// Image 创建一个绘画任务
func (h *MidJourneyHandler) Image(c *gin.Context) {
	var data struct {
		TaskType  string   `json:"task_type"`
		ClientId  string   `json:"client_id"`
		Prompt    string   `json:"prompt"`
		NegPrompt string   `json:"neg_prompt"`
		Rate      string   `json:"rate"`
		Model     string   `json:"model"`   // 模型
		Chaos     int      `json:"chaos"`   // 创意度取值范围: 0-100
		Raw       bool     `json:"raw"`     // 是否开启原始模型
		Seed      int64    `json:"seed"`    // 随机数
		Stylize   int      `json:"stylize"` // 风格化
		ImgArr    []string `json:"img_arr"`
		Tile      bool     `json:"tile"`    // 重复平铺
		Quality   float32  `json:"quality"` // 画质
		Iw        float32  `json:"iw"`
		CRef      string   `json:"cref"` //生成角色一致的图像
		SRef      string   `json:"sref"` //生成风格一致的图像
		Cw        int      `json:"cw"`   // 参考程度
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if !h.preCheck(c) {
		return
	}

	var params = ""
	if data.Rate != "" && !strings.Contains(params, "--ar") {
		params += " --ar " + data.Rate
	}
	if data.Seed > 0 && !strings.Contains(params, "--seed") {
		params += fmt.Sprintf(" --seed %d", data.Seed)
	}
	if data.Stylize > 0 && !strings.Contains(params, "--s") && !strings.Contains(params, "--stylize") {
		params += fmt.Sprintf(" --s %d", data.Stylize)
	}
	if data.Chaos > 0 && !strings.Contains(params, "--c") && !strings.Contains(params, "--chaos") {
		params += fmt.Sprintf(" --c %d", data.Chaos)
	}
	if len(data.ImgArr) > 0 && data.Iw > 0 {
		params += fmt.Sprintf(" --iw %.2f", data.Iw)
	}
	if data.Raw {
		params += " --style raw"
	}
	if data.Quality > 0 {
		params += fmt.Sprintf(" --q %.2f", data.Quality)
	}
	if data.Tile {
		params += " --tile "
	}
	if data.CRef != "" {
		params += fmt.Sprintf(" --cref %s", data.CRef)
		if data.Cw > 0 {
			params += fmt.Sprintf(" --cw %d", data.Cw)
		} else {
			params += " --cw 100"
		}
	}

	if data.SRef != "" {
		params += fmt.Sprintf(" --sref %s", data.SRef)
	}
	if data.Model != "" && !strings.Contains(params, "--v") && !strings.Contains(params, "--niji") {
		params += fmt.Sprintf(" %s", data.Model)
	}

	// 处理融图和换脸的提示词
	if data.TaskType == types.TaskSwapFace.String() || data.TaskType == types.TaskBlend.String() {
		params = fmt.Sprintf("%s:%s", data.TaskType, strings.Join(data.ImgArr, ","))
	}

	// 如果本地图片上传的是相对地址，处理成绝对地址
	for k, v := range data.ImgArr {
		if !strings.HasPrefix(v, "http") {
			data.ImgArr[k] = fmt.Sprintf("http://localhost:5678/%s", strings.TrimLeft(v, "/"))
		}
	}

	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	// generate task id
	taskId, err := h.snowflake.Next(true)
	if err != nil {
		resp.ERROR(c, "error with generate task id: "+err.Error())
		return
	}
	task := types.MjTask{
		ClientId:         data.ClientId,
		TaskId:           taskId,
		Type:             types.TaskType(data.TaskType),
		Prompt:           data.Prompt,
		NegPrompt:        data.NegPrompt,
		Params:           params,
		UserId:           userId,
		ImgArr:           data.ImgArr,
		Mode:             h.App.SysConfig.MjMode,
		TranslateModelId: h.App.SysConfig.TranslateModelId,
	}
	job := model.MidJourneyJob{
		Type:      data.TaskType,
		UserId:    userId,
		TaskId:    taskId,
		TaskInfo:  utils.JsonEncode(task),
		Progress:  0,
		Prompt:    fmt.Sprintf("%s %s", data.Prompt, params),
		Power:     h.App.SysConfig.MjPower,
		CreatedAt: time.Now(),
	}
	opt := "绘图"
	if data.TaskType == types.TaskBlend.String() {
		job.Prompt = "融图：" + strings.Join(data.ImgArr, ",")
		opt = "融图"
	} else if data.TaskType == types.TaskSwapFace.String() {
		job.Prompt = "换脸：" + strings.Join(data.ImgArr, ",")
		opt = "换脸"
	}

	if res := h.DB.Create(&job); res.Error != nil || res.RowsAffected == 0 {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	task.Id = job.Id
	h.mjService.PushTask(task)

	// update user's power
	err = h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "mid-journey",
		Remark: fmt.Sprintf("%s操作，任务ID：%s", opt, job.TaskId),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

type reqVo struct {
	Index       int    `json:"index"`
	ClientId    string `json:"client_id"`
	ChannelId   string `json:"channel_id"`
	MessageId   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

// Upscale send upscale command to MidJourney Bot
func (h *MidJourneyHandler) Upscale(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if !h.preCheck(c) {
		return
	}

	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	taskId, _ := h.snowflake.Next(true)
	task := types.MjTask{
		ClientId:    data.ClientId,
		Type:        types.TaskUpscale,
		UserId:      userId,
		ChannelId:   data.ChannelId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
		Mode:        h.App.SysConfig.MjMode,
	}
	job := model.MidJourneyJob{
		Type:      types.TaskUpscale.String(),
		UserId:    userId,
		TaskId:    taskId,
		TaskInfo:  utils.JsonEncode(task),
		Progress:  0,
		Power:     h.App.SysConfig.MjActionPower,
		CreatedAt: time.Now(),
	}
	if res := h.DB.Create(&job); res.Error != nil || res.RowsAffected == 0 {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	task.Id = job.Id
	h.mjService.PushTask(task)

	// update user's power
	err := h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "mid-journey",
		Remark: fmt.Sprintf("Upscale 操作，任务ID：%s", job.TaskId),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// Variation send variation command to MidJourney Bot
func (h *MidJourneyHandler) Variation(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if !h.preCheck(c) {
		return
	}

	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	taskId, _ := h.snowflake.Next(true)
	task := types.MjTask{
		Type:        types.TaskVariation,
		ClientId:    data.ClientId,
		UserId:      userId,
		Index:       data.Index,
		ChannelId:   data.ChannelId,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
		Mode:        h.App.SysConfig.MjMode,
	}
	job := model.MidJourneyJob{
		Type:      types.TaskVariation.String(),
		ChannelId: data.ChannelId,
		UserId:    userId,
		TaskId:    taskId,
		TaskInfo:  utils.JsonEncode(task),
		Progress:  0,
		Power:     h.App.SysConfig.MjActionPower,
		CreatedAt: time.Now(),
	}
	if res := h.DB.Create(&job); res.Error != nil || res.RowsAffected == 0 {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	task.Id = job.Id
	h.mjService.PushTask(task)

	err := h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "mid-journey",
		Remark: fmt.Sprintf("Variation 操作，任务ID：%s", job.TaskId),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// ImgWall 照片墙
func (h *MidJourneyHandler) ImgWall(c *gin.Context) {
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)
	err, jobs := h.getData(true, 0, page, pageSize, true)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, jobs)
}

// JobList 获取 MJ 任务列表
func (h *MidJourneyHandler) JobList(c *gin.Context) {
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
func (h *MidJourneyHandler) getData(finish bool, userId uint, page int, pageSize int, publish bool) (error, vo.Page) {
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
		session = session.Where("publish = ?", publish)
	}
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}

	// 统计总数
	var total int64
	session.Model(&model.MidJourneyJob{}).Count(&total)

	var items []model.MidJourneyJob
	res := session.Find(&items)
	if res.Error != nil {
		return res.Error, vo.Page{}
	}

	var jobs = make([]vo.MidJourneyJob, 0)
	for _, item := range items {
		var job vo.MidJourneyJob
		err := utils.CopyObject(item, &job)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return nil, vo.NewPage(total, page, pageSize, jobs)
}

// Remove remove task image
func (h *MidJourneyHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetInt(c, "user_id", 0)
	var job model.MidJourneyJob
	if res := h.DB.Where("id = ? AND user_id = ?", id, userId).First(&job); res.Error != nil {
		resp.ERROR(c, "记录不存在")
		return
	}

	// remove job
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

// Publish 发布图片到画廊显示
func (h *MidJourneyHandler) Publish(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetInt(c, "user_id", 0)
	action := h.GetBool(c, "action") // 发布动作，true => 发布，false => 取消分享
	err := h.DB.Model(&model.MidJourneyJob{Id: uint(id), UserId: userId}).UpdateColumn("publish", action).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
