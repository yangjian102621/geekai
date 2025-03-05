package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Websocket 连接处理 handler

type WebsocketHandler struct {
	BaseHandler
	wsService   *service.WebsocketService
	chatHandler *ChatHandler
}

func NewWebsocketHandler(app *core.AppServer, s *service.WebsocketService, db *gorm.DB, chatHandler *ChatHandler) *WebsocketHandler {
	return &WebsocketHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
		chatHandler: chatHandler,
		wsService:   s,
	}
}

func (h *WebsocketHandler) Client(c *gin.Context) {
	clientProtocols := c.GetHeader("Sec-WebSocket-Protocol")
	ws, err := (&websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: strings.Split(clientProtocols, ","),
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Abort()
		return
	}

	clientId := c.Query("client_id")
	client := types.NewWsClient(ws, clientId)
	userId := h.GetLoginUserId(c)
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
				h.wsService.Clients.Delete(clientId)
				break
			}

			var message types.InputMessage
			err = utils.JsonDecode(string(msg), &message)
			if err != nil {
				continue
			}

			logger.Debugf("Receive a message:%+v", message)
			if message.Type == types.MsgTypePing {
				utils.SendChannelMsg(client, types.ChPing, "pong")
				continue
			}

			// 当前只处理聊天消息，其他消息全部丢弃
			var chatMessage types.ChatMessage
			err = utils.JsonDecode(utils.JsonEncode(message.Body), &chatMessage)
			if err != nil || message.Channel != types.ChChat {
				logger.Warnf("invalid message body:%+v", message.Body)
				continue
			}
			var chatRole model.ChatRole
			err = h.DB.First(&chatRole, chatMessage.RoleId).Error
			if err != nil || !chatRole.Enable {
				utils.SendAndFlush(client, "当前聊天角色不存在或者未启用，请更换角色之后再发起对话！！！")
				continue
			}
			// if the role bind a model_id, use role's bind model_id
			if chatRole.ModelId > 0 {
				chatMessage.RoleId = chatRole.ModelId
			}
			// get model info
			var chatModel model.ChatModel
			err = h.DB.Where("id", chatMessage.ModelId).First(&chatModel).Error
			if err != nil || chatModel.Enabled == false {
				utils.SendAndFlush(client, "当前AI模型暂未启用，请更换模型后再发起对话！！！")
				continue
			}

			session := &types.ChatSession{
				ClientIP: c.ClientIP(),
				UserId:   userId,
			}

			// use old chat data override the chat model and role ID
			var chat model.ChatItem
			h.DB.Where("chat_id", chatMessage.ChatId).First(&chat)
			if chat.Id > 0 {
				chatModel.Id = chat.ModelId
				chatMessage.RoleId = int(chat.RoleId)
			}

			session.ChatId = chatMessage.ChatId
			session.Tools = chatMessage.Tools
			session.Stream = chatMessage.Stream
			// 复制模型数据
			err = utils.CopyObject(chatModel, &session.Model)
			if err != nil {
				logger.Error(err, chatModel)
			}
			ctx, cancel := context.WithCancel(context.Background())
			h.chatHandler.ReqCancelFunc.Put(clientId, cancel)
			err = h.chatHandler.sendMessage(ctx, session, chatRole, chatMessage.Content, client)
			if err != nil {
				logger.Error(err)
				utils.SendAndFlush(client, err.Error())
			} else {
				utils.SendMsg(client, types.ReplyMessage{Channel: types.ChChat, Type: types.MsgTypeEnd})
				logger.Infof("回答完毕: %v", message.Body)
			}

		}
	}()
}
