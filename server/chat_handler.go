package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func (s *Server) Chat(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Fatal(err)
		return
	}
	logger.Infof("New websocket connected, IP: %s", c.Request.RemoteAddr)
	client := NewWsClient(ws)
	go func() {
		for {
			_, message, err := client.Receive()
			if err != nil {
				logger.Error(err)
				client.Close()
				return
			}

			// TODO: 接受消息，调用 ChatGPT 返回消息
			logger.Info(string(message))
			err = client.Send(message)
			if err != nil {
				logger.Error(err)
			}
		}
	}()
}
