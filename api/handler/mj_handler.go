package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/oss"
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
	"sync"
	"time"
)

type TaskStatus string

const (
	Stopped  = TaskStatus("Stopped")
	Finished = TaskStatus("Finished")
)

type Image struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Hash     string `json:"hash"`
}

type MidJourneyHandler struct {
	BaseHandler
	redis           *redis.Client
	db              *gorm.DB
	mjService       *service.MjService
	uploaderManager *oss.UploaderManager
	lock            sync.Mutex
	clients         *types.LMap[string, *types.WsClient]
}

func NewMidJourneyHandler(
	app *core.AppServer,
	client *redis.Client,
	db *gorm.DB,
	manager *oss.UploaderManager,
	mjService *service.MjService) *MidJourneyHandler {
	h := MidJourneyHandler{
		redis:           client,
		db:              db,
		uploaderManager: manager,
		lock:            sync.Mutex{},
		mjService:       mjService,
		clients:         types.NewLMap[string, *types.WsClient](),
	}
	h.App = app
	return &h
}

type notifyData struct {
	MessageId   string     `json:"message_id"`
	ReferenceId string     `json:"reference_id"`
	Image       Image      `json:"image"`
	Content     string     `json:"content"`
	Prompt      string     `json:"prompt"`
	Status      TaskStatus `json:"status"`
	Progress    int        `json:"progress"`
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
	// 删除旧的连接
	h.clients.Delete(sessionId)
	h.clients.Put(sessionId, client)
	logger.Infof("New websocket connected, IP: %s", c.ClientIP())
}

func (h *MidJourneyHandler) Notify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != h.App.Config.ExtConfig.Token {
		resp.NotAuth(c)
		return
	}
	var data notifyData
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	logger.Debugf("收到 MidJourney 回调请求：%+v", data)

	h.lock.Lock()
	defer h.lock.Unlock()

	err, finished := h.notifyHandler(c, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 解除任务锁定
	if finished && (data.Status == Finished || data.Status == Stopped) {
		h.redis.Del(c, service.MjRunningJobKey)
	}
	resp.SUCCESS(c)

}

func (h *MidJourneyHandler) notifyHandler(c *gin.Context, data notifyData) (error, bool) {
	taskString, err := h.redis.Get(c, service.MjRunningJobKey).Result()
	if err != nil { // 过期任务，丢弃
		logger.Warn("任务已过期：", err)
		return nil, true
	}

	var task service.MjTask
	err = utils.JsonDecode(taskString, &task)
	if err != nil { // 非标准任务，丢弃
		logger.Warn("任务解析失败：", err)
		return nil, false
	}

	var job model.MidJourneyJob
	res := h.db.Where("message_id = ?", data.MessageId).First(&job)
	if res.Error == nil && data.Status == Finished {
		logger.Warn("重复消息：", data.MessageId)
		return nil, false
	}

	if task.Src == service.TaskSrcImg { // 绘画任务
		var job model.MidJourneyJob
		res := h.db.Where("id = ?", task.Id).First(&job)
		if res.Error != nil {
			logger.Warn("非法任务：", res.Error)
			return nil, false
		}
		job.MessageId = data.MessageId
		job.ReferenceId = data.ReferenceId
		job.Progress = data.Progress
		job.Prompt = data.Prompt
		job.Hash = data.Image.Hash

		// 任务完成，将最终的图片下载下来
		if data.Progress == 100 {
			imgURL, err := h.uploaderManager.GetUploadHandler().PutImg(data.Image.URL)
			if err != nil {
				logger.Error("error with download img: ", err.Error())
				return err, false
			}
			job.ImgURL = imgURL
		} else {
			// 临时图片直接保存，访问的时候使用代理进行转发
			job.ImgURL = data.Image.URL
		}
		res = h.db.Updates(&job)
		if res.Error != nil {
			logger.Error("error with update job: ", res.Error)
			return res.Error, false
		}

		var jobVo vo.MidJourneyJob
		err := utils.CopyObject(job, &jobVo)
		if err == nil {
			if data.Progress < 100 {
				image, err := utils.DownloadImage(jobVo.ImgURL, h.App.Config.ProxyURL)
				if err == nil {
					jobVo.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				}
			}

			// 推送任务到前端
			client := h.clients.Get(task.SessionId)
			if client != nil {
				utils.ReplyChunkMessage(client, jobVo)
			}
		}

	} else if task.Src == service.TaskSrcChat { // 聊天任务
		wsClient := h.App.MjTaskClients.Get(task.SessionId)
		if data.Status == Finished {
			if wsClient != nil && data.ReferenceId != "" {
				content := fmt.Sprintf("**%s** 任务执行成功，正在从 MidJourney 服务器下载图片，请稍后...", data.Prompt)
				utils.ReplyMessage(wsClient, content)
			}
			// download image
			imgURL, err := h.uploaderManager.GetUploadHandler().PutImg(data.Image.URL)
			if err != nil {
				logger.Error("error with download image: ", err)
				if wsClient != nil && data.ReferenceId != "" {
					content := fmt.Sprintf("**%s** 图片下载失败：%s", data.Prompt, err.Error())
					utils.ReplyMessage(wsClient, content)
				}
				return err, false
			}

			tx := h.db.Begin()
			data.Image.URL = imgURL
			message := model.HistoryMessage{
				UserId:     uint(task.UserId),
				ChatId:     task.ChatId,
				RoleId:     uint(task.RoleId),
				Type:       types.MjMsg,
				Icon:       task.Icon,
				Content:    utils.JsonEncode(data),
				Tokens:     0,
				UseContext: false,
			}
			res = tx.Create(&message)
			if res.Error != nil {
				return res.Error, false
			}

			// save the job
			job.UserId = task.UserId
			job.Type = task.Type.String()
			job.MessageId = data.MessageId
			job.ReferenceId = data.ReferenceId
			job.Prompt = data.Prompt
			job.ImgURL = imgURL
			job.Progress = data.Progress
			job.Hash = data.Image.Hash
			job.CreatedAt = time.Now()
			res = tx.Create(&job)
			if res.Error != nil {
				tx.Rollback()
				return res.Error, false
			}
			tx.Commit()
		}

		if wsClient == nil { // 客户端断线，则丢弃
			logger.Errorf("Client is offline: %+v", data)
			return nil, true
		}

		if data.Status == Finished {
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsEnd})
			// 本次绘画完毕，移除客户端
			h.App.MjTaskClients.Delete(task.SessionId)
		} else {
			// 使用代理临时转发图片
			if data.Image.URL != "" {
				image, err := utils.DownloadImage(data.Image.URL, h.App.Config.ProxyURL)
				if err == nil {
					data.Image.URL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
				}
			}
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
		}
	}

	// 更新用户剩余绘图次数
	// TODO: 放大图片是否需要消耗绘图次数？
	if data.Status == Finished {
		h.db.Model(&model.User{}).Where("id = ?", task.UserId).UpdateColumn("img_calls", gorm.Expr("img_calls - ?", 1))
	}

	return nil, true
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
	var data struct {
		SessionId string  `json:"session_id"`
		Prompt    string  `json:"prompt"`
		Rate      string  `json:"rate"`
		Model     string  `json:"model"`
		Chaos     int     `json:"chaos"`
		Raw       bool    `json:"raw"`
		Seed      int64   `json:"seed"`
		Stylize   int     `json:"stylize"`
		Img       string  `json:"img"`
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
	if data.Model != "" && !strings.Contains(prompt, "--v") && !strings.Contains(prompt, "--niji") {
		prompt += data.Model
	}

	idValue, _ := c.Get(types.LoginUserID)
	userId := utils.IntValue(utils.InterfaceToString(idValue), 0)
	job := model.MidJourneyJob{
		Type:      service.Image.String(),
		UserId:    userId,
		Progress:  0,
		Prompt:    prompt,
		CreatedAt: time.Now(),
	}
	if res := h.db.Create(&job); res.Error != nil {
		resp.ERROR(c, "添加任务失败："+res.Error.Error())
		return
	}

	h.mjService.PushTask(service.MjTask{
		Id:        int(job.Id),
		SessionId: data.SessionId,
		Src:       service.TaskSrcImg,
		Type:      service.Image,
		Prompt:    prompt,
		UserId:    userId,
	})

	var jobVo vo.MidJourneyJob
	err := utils.CopyObject(job, &jobVo)
	if err == nil {
		// 推送任务到前端
		client := h.clients.Get(data.SessionId)
		if client != nil {
			utils.ReplyChunkMessage(client, jobVo)
		}
	}
	resp.SUCCESS(c)
}

type reqVo struct {
	Src         string `json:"src"`
	Index       int32  `json:"index"`
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
	src := service.TaskSrc(data.Src)
	if src == service.TaskSrcImg {
		job := model.MidJourneyJob{
			Type:      service.Upscale.String(),
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
			client := h.clients.Get(data.SessionId)
			if client != nil {
				utils.ReplyChunkMessage(client, jobVo)
			}
		}
	}
	h.mjService.PushTask(service.MjTask{
		Id:          jobId,
		SessionId:   data.SessionId,
		Src:         src,
		Type:        service.Upscale,
		Prompt:      data.Prompt,
		UserId:      userId,
		RoleId:      data.RoleId,
		Icon:        data.Icon,
		ChatId:      data.ChatId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})

	wsClient := h.App.ChatClients.Get(data.SessionId)
	if wsClient != nil {
		content := fmt.Sprintf("**%s** 已推送 upscale 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
		utils.ReplyMessage(wsClient, content)
		if h.App.MjTaskClients.Get(data.SessionId) == nil {
			h.App.MjTaskClients.Put(data.SessionId, wsClient)
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
	src := service.TaskSrc(data.Src)
	if src == service.TaskSrcImg {
		job := model.MidJourneyJob{
			Type:      service.Variation.String(),
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
			client := h.clients.Get(data.SessionId)
			if client != nil {
				utils.ReplyChunkMessage(client, jobVo)
			}
		}
	}
	h.mjService.PushTask(service.MjTask{
		Id:          jobId,
		SessionId:   data.SessionId,
		Src:         src,
		Type:        service.Variation,
		Prompt:      data.Prompt,
		UserId:      userId,
		RoleId:      data.RoleId,
		Icon:        data.Icon,
		ChatId:      data.ChatId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})

	// 从聊天窗口发送的请求，记录客户端信息
	wsClient := h.App.ChatClients.Get(data.SessionId)
	if wsClient != nil {
		content := fmt.Sprintf("**%s** 已推送 variation 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
		utils.ReplyMessage(wsClient, content)
		if h.App.MjTaskClients.Get(data.SessionId) == nil {
			h.App.MjTaskClients.Put(data.SessionId, wsClient)
		}
	}
	resp.SUCCESS(c)
}

// JobList 获取 MJ 任务列表
func (h *MidJourneyHandler) JobList(c *gin.Context) {
	status := h.GetInt(c, "status", 0)
	var items []model.MidJourneyJob
	var res *gorm.DB
	userId, _ := c.Get(types.LoginUserID)
	if status == 1 {
		res = h.db.Where("user_id = ? AND progress = 100", userId).Order("id DESC").Find(&items)
	} else {
		res = h.db.Where("user_id = ? AND progress < 100", userId).Order("id ASC").Find(&items)
	}
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
		if item.Progress < 100 && item.ImgURL != "" { // 正在运行中任务使用代理访问图片
			image, err := utils.DownloadImage(item.ImgURL, h.App.Config.ProxyURL)
			if err == nil {
				job.ImgURL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
			}
		}
		jobs = append(jobs, job)
	}
	resp.SUCCESS(c, jobs)
}
