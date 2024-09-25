package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// Websocket 连接处理 handler

type WebsocketHandler struct {
	BaseHandler
	wsService *service.WebsocketService
}

func NewWebsocketHandler(app *core.AppServer, s *service.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{
		BaseHandler: BaseHandler{App: app},
		wsService:   s,
	}
}

func (h *WebsocketHandler) Client(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Abort()
		return
	}

	userId := h.GetInt(c, "user_id", 0)
	clientId := c.Query("client")
	client := types.NewWsClient(ws)
	if userId == 0 {
		_ = client.Send([]byte("Invalid user_id"))
		c.Abort()
		return
	}
	var user model.User
	if err := h.DB.Where("id", userId).First(&user).Error; err != nil {
		_ = client.Send([]byte("Invalid user_id"))
		c.Abort()
		return
	}

	h.wsService.Clients.Put(clientId, client)
	logger.Infof("New websocket connected, IP: %s", c.RemoteIP())
	go func() {
		for {
			_, msg, err := client.Receive()
			if err != nil {
				logger.Debugf("close connection: %s", client.Conn.RemoteAddr())
				client.Close()
			}

			var message types.InputMessage
			err = utils.JsonDecode(string(msg), &message)
			if err != nil {
				continue
			}

			logger.Infof("Receive a message:%+v", message)
			if message.Type == types.WsMsgTypePing {
				_ = client.Send([]byte(`{"type":"pong"}`))
			}

			switch message.Channel {

			}

		}
	}()
}
