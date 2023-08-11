package utils

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"encoding/json"
)

var logger = logger2.GetLogger()

// ReplyChunkMessage 回复客户片段端消息
func ReplyChunkMessage(client types.Client, message types.WsMessage) {
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error for decoding json data: %v", err.Error())
		return
	}
	err = client.(*types.WsClient).Send(msg)
	if err != nil {
		logger.Errorf("Error for reply message: %v", err.Error())
	}
}

// ReplyMessage 回复客户端一条完整的消息
func ReplyMessage(ws types.Client, message string) {
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsStart})
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsMiddle, Content: message})
	ReplyChunkMessage(ws, types.WsMessage{Type: types.WsEnd})
}
