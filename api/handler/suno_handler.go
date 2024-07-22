package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SunoHandler struct {
	BaseHandler
}

func NewSunoHandler(app *core.AppServer, db *gorm.DB) *SunoHandler {
	return &SunoHandler{
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

// Client WebSocket 客户端，用于通知任务状态变更
func (h *SunoHandler) Client(c *gin.Context) {
	//ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	//if err != nil {
	//	logger.Error(err)
	//	c.Abort()
	//	return
	//}
	//
	//userId := h.GetInt(c, "user_id", 0)
	//if userId == 0 {
	//	logger.Info("Invalid user ID")
	//	c.Abort()
	//	return
	//}
	//
	////client := types.NewWsClient(ws)
	//logger.Infof("New websocket connected, IP: %s", c.RemoteIP())
}

func (h *SunoHandler) Create(c *gin.Context) {

}

func (h *SunoHandler) List(c *gin.Context) {

}

func (h *SunoHandler) Remove(c *gin.Context) {

}

func (h *SunoHandler) Publish(c *gin.Context) {

}
