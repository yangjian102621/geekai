package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service/function"
	"chatplus/service/oss"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	leveldb         *store.LevelDB
	db              *gorm.DB
	mjFunc          function.FuncMidJourney
	uploaderManager *oss.UploaderManager
	lock            sync.Mutex
}

func NewMidJourneyHandler(
	app *core.AppServer,
	leveldb *store.LevelDB,
	db *gorm.DB,
	manager *oss.UploaderManager,
	functions map[string]function.Function) *MidJourneyHandler {
	h := MidJourneyHandler{
		leveldb:         leveldb,
		db:              db,
		uploaderManager: manager,
		lock:            sync.Mutex{},
		mjFunc:          functions[types.FuncMidJourney].(function.FuncMidJourney)}
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
		Key         string     `json:"key"`
	}
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	logger.Debugf("收到 MidJourney 回调请求：%+v", data)
	h.lock.Lock()
	defer h.lock.Unlock()

	// the job is saved
	var job model.MidJourneyJob
	res := h.db.Where("message_id = ?", data.MessageId).First(&job)
	if res.Error == nil {
		resp.SUCCESS(c)
		return
	}

	data.Key = utils.Sha256(data.Prompt)
	//logger.Info(data.Prompt, ",", key)
	if data.Status == Finished {
		var task types.MjTask
		err := h.leveldb.Get(types.TaskStorePrefix+data.Key, &task)
		if err != nil {
			logger.Error("error with get MidJourney task: ", err)
			resp.SUCCESS(c)
			return
		}
		// download image
		imgURL, err := h.uploaderManager.GetUploadHandler().PutImg(data.Image.URL)
		if err != nil {
			logger.Error("error with download image: ", err)
			resp.SUCCESS(c)
			return
		}

		data.Image.URL = imgURL
		message := model.HistoryMessage{
			UserId:     task.UserId,
			ChatId:     task.ChatId,
			RoleId:     task.RoleId,
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
		job.ChatId = task.ChatId
		job.MessageId = data.MessageId
		job.ReferenceId = data.ReferenceId
		job.Content = data.Content
		job.Prompt = data.Prompt
		job.Image = utils.JsonEncode(data.Image)
		job.Hash = data.Image.Hash
		job.CreatedAt = time.Now()
		res = h.db.Create(&job)
		if res.Error != nil {
			logger.Error("error with save MidJourney Job: ", res.Error)
		}
	}

	// 推送消息到客户端
	wsClient := h.App.MjTaskClients.Get(data.Key)
	if wsClient == nil { // 客户端断线，则丢弃
		logger.Errorf("Client is offline: %+v", data)
		resp.SUCCESS(c, "Client is offline")
		return
	}

	if data.Status == Finished {
		utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
		utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsEnd})
		// delete client
		h.App.MjTaskClients.Delete(data.Key)
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
	resp.SUCCESS(c, "SUCCESS")
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

	err := h.mjFunc.Upscale(function.MjUpscaleReq{
		Index:       data.Index,
		MessageId:   data.MessageId,
		MessageHash: data.MessageHash,
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	content := fmt.Sprintf("**%s** 已推送 Upscale 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
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
	content := fmt.Sprintf("**%s** 已推送 Variation 任务到 MidJourney 机器人，请耐心等待任务执行...", data.Prompt)
	utils.ReplyMessage(wsClient, content)
	if h.App.MjTaskClients.Get(data.Key) == nil {
		h.App.MjTaskClients.Put(data.Key, wsClient)
	}
	resp.SUCCESS(c)
}
