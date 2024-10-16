package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// 实时 API 中继器

type RealtimeHandler struct {
	BaseHandler
}

func NewRealtimeHandler() *RealtimeHandler {
	return &RealtimeHandler{}
}

func (h *RealtimeHandler) Connection(c *gin.Context) {
	// 获取客户端请求中指定的子协议
	clientProtocols := c.GetHeader("Sec-WebSocket-Protocol")
	logger.Info(clientProtocols)

	// 升级HTTP连接为WebSocket，并传入客户端请求的子协议
	upgrader := websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{clientProtocols},
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Abort()
		return
	}
	defer ws.Close()

	// 连接到真实的后端服务器，传入相同的子协议
	headers := http.Header{}
	if clientProtocols != "" {
		headers.Set("Sec-WebSocket-Protocol", clientProtocols)
	}
	for key, values := range headers {
		for _, value := range values {
			logger.Infof("%s: %s", key, value)
		}
	}
	backendConn, _, err := websocket.DefaultDialer.Dial("wss://api.geekai.pro/v1/realtime?model=gpt-4o-realtime-preview-2024-10-01", headers)
	if err != nil {
		log.Printf("Failed to connect to backend: %v", err)
		return
	}
	defer backendConn.Close()

	//logger.Info(ws.Subprotocol(), ",", backendConn.Subprotocol())
	//// 确保协议一致性，如果失败返回
	//if ws.Subprotocol() != backendConn.Subprotocol() {
	//	log.Println("Subprotocol mismatch")
	//	return
	//}

	// 开始双向转发
	errorChan := make(chan error, 2)
	go relay(ws, backendConn, errorChan)
	go relay(backendConn, ws, errorChan)

	// 等待其中一个连接关闭
	<-errorChan
	log.Println("Relay ended")
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

//func (h *RealtimeHandler) handleMessage(client *RealtimeClient, message []byte) {
//	var event Event
//	err := json.Unmarshal(message, &event)
//	if err != nil {
//		logger.Infof("Error parsing event from client: %s", message)
//		return
//	}
//	logger.Infof("Relaying %q to OpenAI", event.Type)
//	client.Send(event)
//}
//
//func relay(src, dst *websocket.Conn, errorChan chan error) {
//	for {
//		messageType, message, err := src.ReadMessage()
//		if err != nil {
//			errorChan <- err
//			return
//		}
//		err = dst.WriteMessage(messageType, message)
//		if err != nil {
//			errorChan <- err
//			return
//		}
//	}
//}
//
//func NewRealtimeClient(apiKey string) *RealtimeClient {
//	return &RealtimeClient{
//		APIKey: apiKey,
//		send:   make(chan Event, 100),
//	}
//}
//
//func (rc *RealtimeClient) Connect() error {
//	u := url.URL{Scheme: "wss", Host: "api.geekai.pro", Path: "v1/realtime", RawQuery: "model=gpt-4o-realtime-preview-2024-10-01"}
//	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//	if err != nil {
//		return err
//	}
//	rc.conn = c
//
//	go rc.readPump()
//	go rc.writePump()
//
//	return nil
//}
//
//func (rc *RealtimeClient) readPump() {
//	defer rc.conn.Close()
//	for {
//		_, message, err := rc.conn.ReadMessage()
//		if err != nil {
//			log.Println("read error:", err)
//			return
//		}
//		var event Event
//		err = json.Unmarshal(message, &event)
//		if err != nil {
//			log.Println("parse error:", err)
//			continue
//		}
//		rc.send <- event
//	}
//}
//
//func (rc *RealtimeClient) writePump() {
//	defer rc.conn.Close()
//	for event := range rc.send {
//		err := rc.conn.WriteJSON(event)
//		if err != nil {
//			log.Println("write error:", err)
//			return
//		}
//	}
//}
//
//func (rc *RealtimeClient) Send(event Event) {
//	rc.send <- event
//}
