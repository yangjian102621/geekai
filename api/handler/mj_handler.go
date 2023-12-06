package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service/mj"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type MidJourneyHandler struct {
	BaseHandler
	redis     *redis.Client
	db        *gorm.DB
	mjService *mj.Service
}

func NewMidJourneyHandler(
	app *core.AppServer,
	client *redis.Client,
	db *gorm.DB,
	mjService *mj.Service) *MidJourneyHandler {
	h := MidJourneyHandler{
		redis:     client,
		db:        db,
		mjService: mjService,
	}
	h.App = app
	return &h
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *MidJourneyHandler) Client(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	sessionId := c.Query("session_id")
	client := types.NewWsClient(ws)
	h.mjService.Clients.Put(sessionId, client)
	logger.Infof("New websocket connected, IP: %s", c.ClientIP())
}

func (h *MidJourneyHandler) checkLimits(c *gin.Context) bool {
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return false
	}

	if user.ImgCalls <= 0 {
		resp.ERROR(c, "您的绘图次数不足，请联系管理员充值！")
		return false
	}

	return true

}

// Image 创建一个绘画任务
func (h *MidJourneyHandler) Image(c *gin.Context) {
	if !h.App.Config.MjConfig.Enabled {
		resp.ERROR(c, "MidJourney service is disabled")
		return
	}

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
	if !h.checkLimits(c) {
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
	job := model.MidJourneyJob{
		Type:      types.TaskImage.String(),
		UserId:    userId,
		Progress:  0,
		Prompt:    prompt,
		CreatedAt: time.Now(),
	}
	if res := h.db.Create(&job); res.Error != nil {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	h.mjService.PushTask(types.MjTask{
		Id:        int(job.Id),
		SessionId: data.SessionId,
		Src:       types.TaskSrcImg,
		Type:      types.TaskImage,
		Prompt:    prompt,
		UserId:    userId,
	})

	var jobVo vo.MidJourneyJob
	err := utils.CopyObject(job, &jobVo)
	if err == nil {
		// 推送任务到前端
		client := h.mjService.Clients.Get(data.SessionId)
		if client != nil {
			utils.ReplyChunkMessage(client, jobVo)
		}
	}
	resp.SUCCESS(c)
}

type reqVo struct {
	Src         string `json:"src"`
	Index       int    `json:"index"`
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

	if !h.checkLimits(c) {
		return
	}

	idValue, _ := c.Get(types.LoginUserID)
	jobId := 0
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	src := types.TaskSrc(data.Src)
	if src == types.TaskSrcImg {
		job := model.MidJourneyJob{
			Type:      types.TaskUpscale.String(),
			UserId:    userId,
			Hash:      data.MessageHash,
			Progress:  0,
			Prompt:    data.Prompt,
			CreatedAt: time.Now(),
		}
		if res := h.db.Create(&job); res.Error == nil {
			jobId = int(job.Id)
		} else {
			resp.ERROR(c, "添加任务失败："+res.Error.Error())
			return
		}

		var jobVo vo.MidJourneyJob
		err := utils.CopyObject(job, &jobVo)
		if err == nil {
			// 推送任务到前端
			client := h.mjService.Clients.Get(data.SessionId)
			if client != nil {
				utils.ReplyChunkMessage(client, jobVo)
			}
		}
	}
	h.mjService.PushTask(types.MjTask{
		Id:          jobId,
		SessionId:   data.SessionId,
		Src:         src,
		Type:        types.TaskUpscale,
		Prompt:      data.Prompt,
		UserId:      userId,
		RoleId:      data.RoleId,
		Icon:        data.Icon,
		ChatId:      data.ChatId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})

	if src == types.TaskSrcChat {
		wsClient := h.App.ChatClients.Get(data.SessionId)
		if wsClient != nil {
			content := fmt.Sprintf("**%s** 已推送 upscale 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
			utils.ReplyMessage(wsClient, content)
			if h.mjService.ChatClients.Get(data.SessionId) == nil {
				h.mjService.ChatClients.Put(data.SessionId, wsClient)
			}
		}
	}
	resp.SUCCESS(c)
}

// Variation send variation command to MidJourney Bot
func (h *MidJourneyHandler) Variation(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil || data.SessionId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if !h.checkLimits(c) {
		return
	}

	idValue, _ := c.Get(types.LoginUserID)
	jobId := 0
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	src := types.TaskSrc(data.Src)
	if src == types.TaskSrcImg {
		job := model.MidJourneyJob{
			Type:      types.TaskVariation.String(),
			UserId:    userId,
			ImgURL:    "",
			Hash:      data.MessageHash,
			Progress:  0,
			Prompt:    data.Prompt,
			CreatedAt: time.Now(),
		}
		if res := h.db.Create(&job); res.Error == nil {
			jobId = int(job.Id)
		} else {
			resp.ERROR(c, "添加任务失败："+res.Error.Error())
			return
		}

		var jobVo vo.MidJourneyJob
		err := utils.CopyObject(job, &jobVo)
		if err == nil {
			// 推送任务到前端
			client := h.mjService.Clients.Get(data.SessionId)
			if client != nil {
				utils.ReplyChunkMessage(client, jobVo)
			}
		}
	}
	h.mjService.PushTask(types.MjTask{
		Id:          jobId,
		SessionId:   data.SessionId,
		Src:         src,
		Type:        types.TaskVariation,
		Prompt:      data.Prompt,
		UserId:      userId,
		RoleId:      data.RoleId,
		Icon:        data.Icon,
		ChatId:      data.ChatId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})

	if src == types.TaskSrcChat {
		// 从聊天窗口发送的请求，记录客户端信息
		wsClient := h.mjService.ChatClients.Get(data.SessionId)
		if wsClient != nil {
			content := fmt.Sprintf("**%s** 已推送 variation 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
			utils.ReplyMessage(wsClient, content)
			if h.mjService.Clients.Get(data.SessionId) == nil {
				h.mjService.Clients.Put(data.SessionId, wsClient)
			}
		}
	}
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
		if item.Progress < 100 {
			// 10 分钟还没完成的任务直接删除
			if time.Now().Sub(item.CreatedAt) > time.Minute*10 {
				h.db.Delete(&item)
				continue
			}
			if item.ImgURL != "" { // 正在运行中任务使用代理访问图片
				image, err := utils.DownloadImage(item.ImgURL, h.App.Config.ProxyURL)
				if err == nil {
					job.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				}
			}
		}
		jobs = append(jobs, job)
	}
	resp.SUCCESS(c, jobs)
}
