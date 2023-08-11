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
		Image   Image      `json:"image"`
		Content string     `json:"content"`
		Status  TaskStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	sessionId := "u7blnft9zqisyrwidjb22j6b78iqc30lv9jtud3k9o"
	wsClient := h.App.ChatClients.Get(sessionId)
	utils.ReplyMessage(wsClient, "![](https://cdn.discordapp.com/attachments/1138713254718361633/1139482452579070053/lal603743923_A_Chinese_girl_walking_barefoot_on_the_beach_weari_df8b6dc0-3b13-478c-8dbb-983015d21661.png)")
	logger.Infof("Data: %+v", data)
	resp.ERROR(c, "Error with CallBack")
}
