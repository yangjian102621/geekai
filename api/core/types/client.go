package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

var ErrConClosed = errors.New("connection Closed")

// WsClient websocket client
type WsClient struct {
	Conn   *websocket.Conn
	lock   sync.Mutex
	mt     int
	Closed bool
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	return &WsClient{
		Conn:   conn,
		lock:   sync.Mutex{},
		mt:     2, // fixed bug for 'Invalid UTF-8 in text frame'
		Closed: false,
	}
}

func (wc *WsClient) Send(message []byte) error {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	if wc.Closed {
		return ErrConClosed
	}

	return wc.Conn.WriteMessage(wc.mt, message)
}

func (wc *WsClient) SendJson(value interface{}) error {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	if wc.Closed {
		return ErrConClosed
	}
	return wc.Conn.WriteJSON(value)
}

func (wc *WsClient) Receive() (int, []byte, error) {
	if wc.Closed {
		return 0, nil, ErrConClosed
	}

	return wc.Conn.ReadMessage()
}

func (wc *WsClient) Close() {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	if wc.Closed {
		return
	}

	_ = wc.Conn.Close()
	wc.Closed = true
}
