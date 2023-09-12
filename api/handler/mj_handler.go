package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/oss"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
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

	err := h.notifyHandler(c, data)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 解除任务锁定
	if data.Status == Finished || data.Status == Stopped {
		h.redis.Del(c, service.MjRunningJobKey)
	}
	resp.SUCCESS(c)

}

func (h *MidJourneyHandler) notifyHandler(c *gin.Context, data notifyData) error {
	taskString, err := h.redis.Get(c, service.MjRunningJobKey).Result()
	if err != nil { // 过期任务，丢弃
		logger.Warn("任务已过期：", err)
		return nil
	}

	var task service.MjTask
	err = utils.JsonDecode(taskString, &task)
	if err != nil { // 非标准任务，丢弃
		logger.Warn("任务解析失败：", err)
		return nil
	}

	if task.Src == service.TaskSrcImg { // 绘画任务
		logger.Error(err)
		var job model.MidJourneyJob
		res := h.db.First(&job, task.Id)
		if res.Error != nil {
			logger.Warn("非法任务：", err)
			return nil
		}
		job.MessageId = data.MessageId
		job.ReferenceId = data.ReferenceId
		job.Progress = data.Progress
		job.Prompt = data.Prompt

		// download image
		if data.Progress == 100 {
			imgURL, err := h.uploaderManager.GetUploadHandler().PutImg(data.Image.URL)
			if err != nil {
				logger.Error("error with download img: ", err.Error())
				return err
			}
			job.ImgURL = imgURL
		} else {
			// 使用图片代理
			job.ImgURL = fmt.Sprintf("/api/mj/proxy?url=%s", data.Image.URL)
		}
		res = h.db.Updates(&job)
		if res.Error != nil {
			logger.Error("error with update job: ", err.Error())
			return res.Error
		}

	} else if task.Src == service.TaskSrcChat { // 聊天任务
		var job model.MidJourneyJob
		res := h.db.Where("message_id = ?", data.MessageId).First(&job)
		if res.Error == nil {
			logger.Warn("重复消息：", data.MessageId)
			return nil
		}

		wsClient := h.App.MjTaskClients.Get(task.Id)
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
				return err
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
				return res.Error
			}

			// save the job
			job.UserId = task.UserId
			job.MessageId = data.MessageId
			job.ReferenceId = data.ReferenceId
			job.Prompt = data.Prompt
			job.ImgURL = imgURL
			job.Progress = data.Progress
			job.CreatedAt = time.Now()
			res = tx.Create(&job)
			if res.Error != nil {
				tx.Rollback()
				return res.Error
			}
			tx.Commit()
		}

		if wsClient == nil { // 客户端断线，则丢弃
			logger.Errorf("Client is offline: %+v", data)
			return nil
		}

		if data.Status == Finished {
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsEnd})
			// delete client
			h.App.MjTaskClients.Delete(task.Id)
		} else {
			data.Image.URL = fmt.Sprintf("/api/mj/proxy?url=%s", data.Image.URL)
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
		}
	}

	return nil
}

func (h *MidJourneyHandler) Proxy(c *gin.Context) {
	logger.Info(c.Request.Host, c.Request.Proto)
	return
	url := c.Query("url")
	image, err := utils.DownloadImage(url, h.App.Config.ProxyURL)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "data:image/png;base64,"+base64.StdEncoding.EncodeToString(image))
}

type reqVo struct {
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
	wsClient := h.App.ChatClients.Get(data.SessionId)
	if wsClient == nil {
		resp.ERROR(c, "No Websocket client online")
		return
	}
	userId, _ := c.Get(types.LoginUserID)
	h.mjService.PushTask(service.MjTask{
		Id:          data.SessionId,
		Src:         service.TaskSrcChat,
		Type:        service.Upscale,
		Prompt:      data.Prompt,
		UserId:      utils.IntValue(utils.InterfaceToString(userId), 0),
		RoleId:      data.RoleId,
		Icon:        data.Icon,
		ChatId:      data.ChatId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})

	content := fmt.Sprintf("**%s** 已推送 upscale 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
	utils.ReplyMessage(wsClient, content)
	if h.App.MjTaskClients.Get(data.SessionId) == nil {
		h.App.MjTaskClients.Put(data.SessionId, wsClient)
	}
	resp.SUCCESS(c)
}

func (h *MidJourneyHandler) Variation(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil || data.SessionId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	wsClient := h.App.ChatClients.Get(data.SessionId)
	if wsClient == nil {
		resp.ERROR(c, "No Websocket client online")
		return
	}

	userId, _ := c.Get(types.LoginUserID)
	h.mjService.PushTask(service.MjTask{
		Id:          data.SessionId,
		Src:         service.TaskSrcChat,
		Type:        service.Variation,
		Prompt:      data.Prompt,
		UserId:      utils.IntValue(utils.InterfaceToString(userId), 0),
		RoleId:      data.RoleId,
		Icon:        data.Icon,
		ChatId:      data.ChatId,
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})
	content := fmt.Sprintf("**%s** 已推送 variation 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
	utils.ReplyMessage(wsClient, content)
	if h.App.MjTaskClients.Get(data.SessionId) == nil {
		h.App.MjTaskClients.Put(data.SessionId, wsClient)
	}
	resp.SUCCESS(c)
}
