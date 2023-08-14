package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	leveldb *store.LevelDB
	db      *gorm.DB
}

func NewMidJourneyHandler(app *core.AppServer, leveldb *store.LevelDB, db *gorm.DB) *MidJourneyHandler {
	h := MidJourneyHandler{leveldb: leveldb, db: db}
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
		Type      string     `json:"type"`
		MessageId string     `json:"message_id"`
		Image     Image      `json:"image"`
		Content   string     `json:"content"`
		Prompt    string     `json:"prompt"`
		Status    TaskStatus `json:"status"`
		Key       string     `json:"key"`
	}
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	key := utils.Sha256(data.Prompt)
	data.Key = key
	//logger.Info(data.Prompt, ",", key)
	if data.Status == Finished {
		var task types.MjTask
		err := h.leveldb.Get(types.TaskStorePrefix+key, &task)
		if err != nil {
			logger.Error("error with get MidJourney task: ", err)
			resp.ERROR(c)
			return
		}

		// TODO: 是否需要把图片下载到本地服务器？

		historyUserMsg := model.HistoryMessage{
			UserId:     task.UserId,
			ChatId:     task.ChatId,
			RoleId:     task.RoleId,
			Type:       types.MjMsg,
			Icon:       task.Icon,
			Content:    utils.JsonEncode(data),
			Tokens:     0,
			UseContext: false,
		}
		res := h.db.Save(&historyUserMsg)
		if res.Error != nil {
			logger.Error("error with save MidJourney message: ", res.Error)
		}

		// delete task from leveldb
		_ = h.leveldb.Delete(types.TaskStorePrefix + key)
	}

	// 推送消息到客户端
	wsClient := h.App.MjTaskClients.Get(key)
	if wsClient == nil { // 客户端断线，则丢弃
		resp.SUCCESS(c, "Client is offline")
		return
	}

	if data.Status == Finished {
		utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
		utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsEnd})
		// delete client
		h.App.MjTaskClients.Delete(key)
	} else {
		utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsMjImg, Content: data})
	}
	resp.SUCCESS(c, "SUCCESS")
}
