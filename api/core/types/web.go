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

// WsMessage Websocket message
type WsMessage struct {
	Type    WsMsgType   `json:"type"` // 消息类别，start, end, img
	Content interface{} `json:"content"`
}

type WsMsgType string

const (
	WsStart  = WsMsgType("start")
	WsMiddle = WsMsgType("middle")
	WsEnd    = WsMsgType("end")
	WsErr    = WsMsgType("error")
)

type BizCode int

const (
	Success       = BizCode(0)
	Failed        = BizCode(1)
	NotAuthorized = BizCode(401) // 未授权

	OkMsg       = "Success"
	ErrorMsg    = "系统开小差了"
	InvalidArgs = "非法参数或参数解析失败"
)
