package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/service/function"
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
	Start    = TaskStatus("Started")
	Running  = TaskStatus("Running")
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

func (h *MidJourneyHandler) Notify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != h.App.Config.ExtConfig.Token {
		resp.NotAuth(c)
		return
	}

	var data struct {
		MessageId   string     `json:"message_id"`
		ReferenceId string     `json:"reference_id"`
		Image       Image      `json:"image"`
		Content     string     `json:"content"`
		Prompt      string     `json:"prompt"`
		Status      TaskStatus `json:"status"`
		Progress    int        `json:"progress"`
	}
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	logger.Debugf("收到 MidJourney 回调请求：%+v", data)
	h.lock.Lock()
	defer h.lock.Unlock()

	taskString, err := h.redis.Get(c, service.MjRunningJobKey).Result()
	if err != nil {
		resp.SUCCESS(c) // 过期任务，丢弃
		return
	}

	var task service.MjTask
	err = utils.JsonDecode(taskString, &task)
	if err != nil {
		resp.SUCCESS(c) // 非标准任务，丢弃
		return
	}

	if task.Src == service.TaskSrcImg { // 绘画任务
		logger.Error(err)
		var job model.MidJourneyJob
		res := h.db.First(&job, task.Id)
		if res.Error != nil {
			resp.SUCCESS(c) // 非法任务，丢弃
			return
		}
		job.MessageId = data.MessageId
		job.ReferenceId = data.ReferenceId
		job.Progress = data.Progress

		// download image
		if data.Progress == 100 {
			imgURL, err := h.uploaderManager.GetUploadHandler().PutImg(data.Image.URL)
			if err != nil {
				resp.ERROR(c, "error with download img: "+err.Error())
				return
			}
			job.ImgURL = imgURL
		} else {
			// 使用图片代理
			job.ImgURL = fmt.Sprintf("/api/mj/proxy?url=%s", data.Image.URL)
		}
		res = h.db.Updates(&job)
		if res.Error != nil {
			resp.ERROR(c, "error with update job: "+err.Error())
			return
		}

		resp.SUCCESS(c)

	} else if task.Src == service.TaskSrcChat { // 聊天任务
		var job model.MidJourneyJob
		res := h.db.Where("message_id = ?", data.MessageId).First(&job)
		if res.Error == nil {
			resp.SUCCESS(c)
			return
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
				resp.ERROR(c, err.Error())
				return
			}

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
			res := h.db.Create(&message)
			if res.Error != nil {
				logger.Error("error with save chat history message: ", res.Error)
			}

			// save the job
			job.UserId = task.UserId
			job.MessageId = data.MessageId
			job.ReferenceId = data.ReferenceId
			job.Prompt = data.Prompt
			job.ImgURL = imgURL
			job.Progress = data.Progress
			job.CreatedAt = time.Now()
			res = h.db.Create(&job)
			if res.Error != nil {
				logger.Error("error with save MidJourney Job: ", res.Error)
			}
		}

		if wsClient == nil { // 客户端断线，则丢弃
			logger.Errorf("Client is offline: %+v", data)
			resp.SUCCESS(c, "Client is offline")
			return
		}

		if data.Status == Finished {
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsEnd})
			// delete client
			h.App.MjTaskClients.Delete(task.Id)
		} else {
			//// 使用代理临时转发图片
			//if data.Image.URL != "" {
			//	image, err := utils.DownloadImage(data.Image.URL, h.App.Config.ProxyURL)
			//	if err == nil {
			//		data.Image.URL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
			//	}
			//}
			data.Image.URL = fmt.Sprintf("/api/mj/proxy?url=%s", data.Image.URL)
			utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
		}
		resp.SUCCESS(c, "SUCCESS")
	}

}

func (h *MidJourneyHandler) Proxy(c *gin.Context) {
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
	Key         string `json:"key"`
	Prompt      string `json:"prompt"`
}

// Upscale send upscale command to MidJourney Bot
func (h *MidJourneyHandler) Upscale(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil ||
		data.SessionId == "" ||
		data.Key == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	wsClient := h.App.ChatClients.Get(data.SessionId)
	if wsClient == nil {
		resp.ERROR(c, "No Websocket client online")
		return
	}

	h.mjService.PushTask(service.MjTask{
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})
	err := n.Upscale(function.MjUpscaleReq{
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	content := fmt.Sprintf("**%s** 已推送 upscale 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
	utils.ReplyMessage(wsClient, content)
	if h.App.MjTaskClients.Get(data.Key) == nil {
		h.App.MjTaskClients.Put(data.Key, wsClient)
	}
	resp.SUCCESS(c)
}

func (h *MidJourneyHandler) Variation(c *gin.Context) {
	var data reqVo
	if err := c.ShouldBindJSON(&data); err != nil ||
		data.SessionId == "" ||
		data.Key == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	wsClient := h.App.ChatClients.Get(data.SessionId)
	if wsClient == nil {
		resp.ERROR(c, "No Websocket client online")
		return
	}

	err := h.mjFunc.Variation(function.MjVariationReq{
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	content := fmt.Sprintf("**%s** 已推送 variation 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
	utils.ReplyMessage(wsClient, content)
	if h.App.MjTaskClients.Get(data.Key) == nil {
		h.App.MjTaskClients.Put(data.Key, wsClient)
	}
	resp.SUCCESS(c)
}
