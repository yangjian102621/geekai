package server

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

var ErrConClosed = errors.New("connection closed")

type Client interface {
	Close()
}

// WsClient websocket client
type WsClient struct {
	Conn   *websocket.Conn
	lock   sync.Mutex
	mt     int
	closed bool
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	return &WsClient{
		Conn:   conn,
		lock:   sync.Mutex{},
		mt:     2, // fixed bug for 'Invalid UTF-8 in text frame'
		closed: false,
	}
}

func (wc *WsClient) Send(message []byte) error {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	if wc.closed {
		return ErrConClosed
	}

	return wc.Conn.WriteMessage(wc.mt, message)
}

func (wc *WsClient) Receive() (int, []byte, error) {
	if wc.closed {
		return 0, nil, ErrConClosed
	}

	return wc.Conn.ReadMessage()
}

func (wc *WsClient) Close() {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	if wc.closed {
		return
	}

	_ = wc.Conn.Close()
	wc.closed = true
}
