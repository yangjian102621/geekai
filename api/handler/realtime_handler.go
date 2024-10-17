package handler

import (
	"fmt"
	"geekai/core"
	"geekai/store/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// OpenAI Realtime API Relay Server

type RealtimeHandler struct {
	BaseHandler
}

func NewRealtimeHandler(server *core.AppServer, db *gorm.DB) *RealtimeHandler {
	return &RealtimeHandler{BaseHandler{App: server, DB: db}}
}

func (h *RealtimeHandler) Connection(c *gin.Context) {
	// 获取客户端请求中指定的子协议
	clientProtocols := c.GetHeader("Sec-WebSocket-Protocol")
	md := c.Query("model")

	userId := h.GetLoginUserId(c)
	var user model.User
	if err := h.DB.Where("id", userId).First(&user).Error; err != nil {
		c.Abort()
		return
	}

	// 将 HTTP 协议升级为 Websocket 协议
	subProtocols := strings.Split(clientProtocols, ",")
	ws, err := (&websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: subProtocols,
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Abort()
		return
	}
	defer ws.Close()

	// 目前只针对 VIP 用户可以访问
	if !user.Vip {
		sendError(ws, "当前功能只针对 VIP 用户开放")
		c.Abort()
		return
	}

	var apiKey model.ApiKey
	h.DB.Where("type", "realtime").Where("enabled", true).Order("last_used_at ASC").First(&apiKey)
	if apiKey.Id == 0 {
		sendError(ws, "管理员未配置 Realtime API KEY")
		c.Abort()
		return
	}

	apiURL := fmt.Sprintf("%s/v1/realtime?model=%s", apiKey.ApiURL, md)
	// 连接到真实的后端服务器，传入相同的子协议
	headers := http.Header{}
	// 修正子协议内容
	subProtocols[1] = "openai-insecure-api-key." + apiKey.Value
	if clientProtocols != "" {
		headers.Set("Sec-WebSocket-Protocol", strings.Join(subProtocols, ","))
	}
	backendConn, _, err := websocket.DefaultDialer.Dial(apiURL, headers)
	if err != nil {
		sendError(ws, "桥接后端 API 失败："+err.Error())
		c.Abort()
		return
	}
	defer backendConn.Close()

	// 确保协议一致性，如果失败返回
	if ws.Subprotocol() != backendConn.Subprotocol() {
		sendError(ws, "Websocket 子协议不匹配")
		c.Abort()
		return
	}

	// 更新API KEY 最后使用时间
	h.DB.Model(&model.ApiKey{}).Where("id", apiKey.Id).UpdateColumn("last_used_at", time.Now().Unix())

	// 开始双向转发
	errorChan := make(chan error, 2)
	go relay(ws, backendConn, errorChan)
	go relay(backendConn, ws, errorChan)

	// 等待其中一个连接关闭
	err = <-errorChan
	logger.Infof("Relay ended: %v", err)
}

func relay(src, dst *websocket.Conn, errorChan chan error) {
	for {
		messageType, message, err := src.ReadMessage()
		if err != nil {
			errorChan <- err
			return
		}
		err = dst.WriteMessage(messageType, message)
		if err != nil {
			errorChan <- err
			return
		}
	}
}

func sendError(ws *websocket.Conn, message string) {
	err := ws.WriteJSON(map[string]string{"event_id": "event_01", "type": "error", "error": message})
	if err != nil {
		logger.Error(err)
	}
}
