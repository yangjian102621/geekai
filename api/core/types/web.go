package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// BizVo 业务返回 VO
type BizVo struct {
	Code     BizCode     `json:"code"`
	Page     int         `json:"page,omitempty"`
	PageSize int         `json:"page_size,omitempty"`
	Total    int         `json:"total,omitempty"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// ReplyMessage 对话回复消息结构
type ReplyMessage struct {
	Channel  WsChannel   `json:"channel"`  // 消息频道，目前只有 chat
	ClientId string      `json:"clientId"` // 客户端ID
	Type     WsMsgType   `json:"type"`     // 消息类别
	Body     interface{} `json:"body"`
}

type WsMsgType string
type WsChannel string

const (
	MsgTypeText = WsMsgType("text") // 输出内容
	MsgTypeEnd  = WsMsgType("end")
	MsgTypeErr  = WsMsgType("error")
	MsgTypePing = WsMsgType("ping") // 心跳消息

	ChPing = WsChannel("ping")
	ChChat = WsChannel("chat")
	ChMj   = WsChannel("mj")
	ChSd   = WsChannel("sd")
	ChDall = WsChannel("dall")
	ChSuno = WsChannel("suno")
	ChLuma = WsChannel("luma")
)

// InputMessage 对话输入消息结构
type InputMessage struct {
	Channel WsChannel   `json:"channel"` // 消息频道
	Type    WsMsgType   `json:"type"`    // 消息类别
	Body    interface{} `json:"body"`
}

type ChatMessage struct {
	Tools   []int  `json:"tools,omitempty"`  // 允许调用工具列表
	Stream  bool   `json:"stream,omitempty"` // 是否采用流式输出
	RoleId  int    `json:"role_id"`
	ModelId int    `json:"model_id"`
	ChatId  string `json:"chat_id"`
	Content string `json:"content"`
}

type BizCode int

const (
	Success       = BizCode(0)
	Failed        = BizCode(1)
	NotAuthorized = BizCode(401) // 未授权

	OkMsg       = "Success"
	ErrorMsg    = "系统开小差了"
	InvalidArgs = "非法参数或参数解析失败"
)
