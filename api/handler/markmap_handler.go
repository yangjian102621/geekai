package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/utils"
	"github.com/gorilla/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MarkMapHandler 生成思维导图
type MarkMapHandler struct {
	BaseHandler
	clients *types.LMap[uint, *types.WsClient]
}

func NewMarkMapHandler(app *core.AppServer, db *gorm.DB) *MarkMapHandler {
	return &MarkMapHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
		clients:     types.NewLMap[uint, *types.WsClient](),
	}
}

func (h *MarkMapHandler) Client(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	modelId := h.GetInt(c, "model_id", 0)
	userId := h.GetLoginUserId(c)
	logger.Info(modelId)
	client := types.NewWsClient(ws)

	// 保存会话连接
	h.clients.Put(userId, client)
	go func() {
		for {
			_, msg, err := client.Receive()
			if err != nil {
				client.Close()
				h.clients.Delete(userId)
				return
			}

			var message types.WsMessage
			err = utils.JsonDecode(string(msg), &message)
			if err != nil {
				continue
			}

			// 心跳消息
			if message.Type == "heartbeat" {
				logger.Debug("收到 Chat 心跳消息：", message.Content)
				continue
			}

			logger.Info("Receive a message: ", message.Content)

		}
	}()
}
