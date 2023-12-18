package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/mj"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type MidJourneyHandler struct {
	BaseHandler
	db        *gorm.DB
	pool      *mj.ServicePool
	snowflake *service.Snowflake
}

func NewMidJourneyHandler(app *core.AppServer, db *gorm.DB, snowflake *service.Snowflake, pool *mj.ServicePool) *MidJourneyHandler {
	h := MidJourneyHandler{
		db:        db,
		snowflake: snowflake,
		pool:      pool,
	}
	h.App = app
	return &h
}

func (h *MidJourneyHandler) preCheck(c *gin.Context) bool {
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if user.ImgCalls <= 0 {
		resp.ERROR(c, "您的绘图次数不足，请联系管理员充值！")
		return false
	}

	if !h.pool.HasAvailableService() {
		resp.ERROR(c, "MidJourney 池子中没有没有可用的服务！")
		return false
	}

	return true

}

// Image 创建一个绘画任务
func (h *MidJourneyHandler) Image(c *gin.Context) {
	var data struct {
		SessionId string  `json:"session_id"`
		Prompt    string  `json:"prompt"`
		NegPrompt string  `json:"neg_prompt"`
		Rate      string  `json:"rate"`
		Model     string  `json:"model"`
		Chaos     int     `json:"chaos"`
		Raw       bool    `json:"raw"`
		Seed      int64   `json:"seed"`
		Stylize   int     `json:"stylize"`
		Img       string  `json:"img"`
		Tile      bool    `json:"tile"`
		Quality   float32 `json:"quality"`
		Weight    float32 `json:"weight"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	if !h.preCheck(c) {
		return
	}

	var prompt = data.Prompt
	if data.Rate != "" && !strings.Contains(prompt, "--ar") {
		prompt += " --ar " + data.Rate
	}
	if data.Seed > 0 && !strings.Contains(prompt, "--seed") {
		prompt += fmt.Sprintf(" --seed %d", data.Seed)
	}
	if data.Stylize > 0 && !strings.Contains(prompt, "--s") && !strings.Contains(prompt, "--stylize") {
		prompt += fmt.Sprintf(" --s %d", data.Stylize)
	}
	if data.Chaos > 0 && !strings.Contains(prompt, "--c") && !strings.Contains(prompt, "--chaos") {
		prompt += fmt.Sprintf(" --c %d", data.Chaos)
	}
	if data.Img != "" {
		prompt = fmt.Sprintf("%s %s", data.Img, prompt)
		if data.Weight > 0 {
			prompt += fmt.Sprintf(" --iw %f", data.Weight)
		}
	}
	if data.Raw {
		prompt += " --style raw"
	}
	if data.Quality > 0 {
		prompt += fmt.Sprintf(" --q %.2f", data.Quality)
	}
	if data.NegPrompt != "" {
		prompt += fmt.Sprintf(" --no %s", data.NegPrompt)
	}
	if data.Tile {
		prompt += " --tile "
	}
	if data.Model != "" && !strings.Contains(prompt, "--v") && !strings.Contains(prompt, "--niji") {
		prompt += fmt.Sprintf(" %s", data.Model)
	}

	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	// generate task id
	taskId, err := h.snowflake.Next(true)
	if err != nil {
		resp.ERROR(c, "error with generate task id: "+err.Error())
		return
	}
	job := model.MidJourneyJob{
		Type:      types.TaskImage.String(),
		UserId:    userId,
		TaskId:    taskId,
		Progress:  0,
		Prompt:    prompt,
		CreatedAt: time.Now(),
	}
	if res := h.db.Create(&job); res.Error != nil {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	h.pool.PushTask(types.MjTask{
		Id:        int(job.Id),
		SessionId: data.SessionId,
		Type:      types.TaskImage,
		Prompt:    fmt.Sprintf("%s %s", taskId, prompt),
		UserId:    userId,
	})
	resp.SUCCESS(c)
}

type reqVo struct {
	TaskId      string `json:"task_id"`
	Index       int    `json:"index"`
	ChannelId   string `json:"channel_id"`
	MessageId   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
	SessionId   string `json:"session_id"`
	Prompt      string `json:"prompt"`
	ChatId      string `json:"chat_id"`
	RoleId      int    `json:"role_id"`
	Icon        string `json:"icon"`
}

// Upscale send upscale command to MidJourney Bot
func (h *MidJourneyHandler) Upscale(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil || data.SessionId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if !h.preCheck(c) {
		return
	}

	idValue, _ := c.Get(types.LoginUserID)
	jobId := 0
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	job := model.MidJourneyJob{
		Type:        types.TaskUpscale.String(),
		ReferenceId: data.MessageId,
		UserId:      userId,
		TaskId:      data.TaskId,
		Progress:    0,
		Prompt:      data.Prompt,
		CreatedAt:   time.Now(),
	}
	if res := h.db.Create(&job); res.Error != nil {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	h.pool.PushTask(types.MjTask{
		Id:          jobId,
		SessionId:   data.SessionId,
		Type:        types.TaskUpscale,
		Prompt:      data.Prompt,
		UserId:      userId,
		ChannelId:   data.ChannelId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})
	resp.SUCCESS(c)
}

// Variation send variation command to MidJourney Bot
func (h *MidJourneyHandler) Variation(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil || data.SessionId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if !h.preCheck(c) {
		return
	}

	idValue, _ := c.Get(types.LoginUserID)
	jobId := 0
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)

	job := model.MidJourneyJob{
		Type:        types.TaskVariation.String(),
		ReferenceId: data.MessageId,
		UserId:      userId,
		TaskId:      data.TaskId,
		Progress:    0,
		Prompt:      data.Prompt,
		CreatedAt:   time.Now(),
	}
	if res := h.db.Create(&job); res.Error != nil {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	h.pool.PushTask(types.MjTask{
		Id:          jobId,
		SessionId:   data.SessionId,
		Type:        types.TaskVariation,
		Prompt:      data.Prompt,
		UserId:      userId,
		Index:       data.Index,
		ChannelId:   data.ChannelId,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})
	resp.SUCCESS(c)
}

// JobList 获取 MJ 任务列表
func (h *MidJourneyHandler) JobList(c *gin.Context) {
	status := h.GetInt(c, "status", 0)
	userId := h.GetInt(c, "user_id", 0)
	page := h.GetInt(c, "page", 0)
	pageSize := h.GetInt(c, "page_size", 0)

	session := h.db.Session(&gorm.Session{})
	if status == 1 {
		session = session.Where("progress = ?", 100).Order("id DESC")
	} else {
		session = session.Where("progress < ?", 100).Order("id ASC")
	}
	if userId > 0 {
		session = session.Where("user_id = ?", userId)
	}
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}

	var items []model.MidJourneyJob
	res := session.Find(&items)
	if res.Error != nil {
		resp.ERROR(c, types.NoData)
		return
	}

	var jobs = make([]vo.MidJourneyJob, 0)
	for _, item := range items {
		var job vo.MidJourneyJob
		err := utils.CopyObject(item, &job)
		if err != nil {
			continue
		}

		if job.Progress == -1 {
			h.db.Delete(&model.MidJourneyJob{Id: job.Id})
		}

		if item.Progress < 100 {
			// 10 分钟还没完成的任务直接删除
			if time.Now().Sub(item.CreatedAt) > time.Minute*10 {
				h.db.Delete(&item)
				continue
			}

			// 正在运行中任务使用代理访问图片
			if item.ImgURL == "" && item.OrgURL != "" {
				image, err := utils.DownloadImage(item.OrgURL, h.App.Config.ProxyURL)
				if err == nil {
					job.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				}
			}
		}

		jobs = append(jobs, job)
	}
	resp.SUCCESS(c, jobs)
}
