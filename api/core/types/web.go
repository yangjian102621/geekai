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
	Channel WsChannel   `json:"channel"` // 消息频道，目前只有 chat
	Type    WsMsgType   `json:"type"`    // 消息类别
	Content interface{} `json:"content"`
}

type WsMsgType string
type WsChannel string

const (
	WsMsgTypeContent = WsMsgType("content") // 输出内容
	WsMsgTypeEnd     = WsMsgType("end")
	WsMsgTypeErr     = WsMsgType("error")
	WsMsgTypePing    = WsMsgType("ping") // 心跳消息

	WsChat = WsChannel("chat")
	WsMj   = WsChannel("mj")
	WsSd   = WsChannel("sd")
	WsDall = WsChannel("dall")
	WsSuno = WsChannel("suno")
	WsLuma = WsChannel("luma")
)

// InputMessage 对话输入消息结构
type InputMessage struct {
	Channel WsChannel `json:"channel"` // 消息频道
	Type    WsMsgType `json:"type"`    // 消息类别
	Content string    `json:"content"`
	Tools   []int     `json:"tools"`  // 允许调用工具列表
	Stream  bool      `json:"stream"` // 是否采用流式输出
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
