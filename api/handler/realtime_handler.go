package handler

import (
	"encoding/json"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
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
	userService *service.UserService
}

func NewRealtimeHandler(server *core.AppServer, db *gorm.DB, userService *service.UserService) *RealtimeHandler {
	return &RealtimeHandler{BaseHandler: BaseHandler{App: server, DB: db}, userService: userService}
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

// OpenAI 实时语音对话，一次性对话
func (h *RealtimeHandler) VoiceChat(c *gin.Context) {
	var apiKey model.ApiKey
	err := h.DB.Session(&gorm.Session{}).Where("type", "realtime").Where("enabled", true).First(&apiKey).Error
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("error with fetch OpenAI API KEY：%v", err))
		return
	}

	// 检查用户是否还有算力
	userId := h.GetLoginUserId(c)
	var user model.User
	if err := h.DB.Where("id", userId).First(&user).Error; err != nil {
		resp.ERROR(c, fmt.Sprintf("error with fetch user：%v", err))
		return
	}

	if user.Power < h.App.SysConfig.AdvanceVoicePower {
		resp.ERROR(c, "当前用户算力不足，无法使用该功能")
		return
	}

	var response utils.OpenAIResponse
	client := req.C()
	if len(apiKey.ProxyURL) > 5 {
		client.SetProxyURL(apiKey.ApiURL)
	}
	apiURL := fmt.Sprintf("%s/v1/chat/completions", apiKey.ApiURL)
	logger.Infof("Sending %s request, API KEY:%s, PROXY: %s, Model: %s", apiKey.ApiURL, apiURL, apiKey.ProxyURL, "advanced-voice")
	r, err := client.R().SetHeader("Body-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(types.ApiRequest{
			Model:       "advanced-voice",
			Temperature: 0.9,
			MaxTokens:   1024,
			Stream:      false,
			Messages: []interface{}{types.Message{
				Role:    "user",
				Content: "实时语音通话",
			}},
		}).Post(apiURL)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("请求 OpenAI API失败：%v", err))
		return
	}

	if r.IsErrorState() {
		resp.ERROR(c, fmt.Sprintf("请求 OpenAI API失败：%v", r.Status))
		return
	}

	body, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("解析API数据失败：%v, %s", err, string(body)))
	}

	// 更新 API KEY 的最后使用时间
	h.DB.Model(&apiKey).UpdateColumn("last_used_at", time.Now().Unix())

	// 扣减算力
	err = h.userService.DecreasePower(userId, h.App.SysConfig.AdvanceVoicePower, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "advanced-voice",
		Remark: "实时语音通话",
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	logger.Infof("Response: %v", response.Choices[0].Message.Content)

	// 提取链接
	re := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	links := re.FindAllStringSubmatch(response.Choices[0].Message.Content, -1)
	var url = ""
	if len(links) > 0 {
		url = links[0][2]
	}
	resp.SUCCESS(c, url)
}
