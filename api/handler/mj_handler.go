package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
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
}

func NewMidJourneyHandler(app *core.AppServer) *MidJourneyHandler {
	h := MidJourneyHandler{}
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
	// TODO: 如果绘画任务完成了则将该消息保存到当前会话的聊天历史记录

	wsClient := h.App.MjTaskClients.Get(key)
	if wsClient == nil { // 客户端断线，则丢弃
		resp.SUCCESS(c)
		return
	}

	// 推送消息到客户端
	// TODO: 增加绘画消息类型
	utils.ReplyChunkMessage(wsClient, types.WsMessage{Type: types.WsImg, Content: data})
	resp.ERROR(c, "Error with CallBack")
}
