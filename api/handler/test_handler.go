package handler

import (
	"geekai/service"
	"geekai/service/payment"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type TestHandler struct {
	db        *gorm.DB
	snowflake *service.Snowflake
	js        *payment.JPayService
}

func NewTestHandler(db *gorm.DB, snowflake *service.Snowflake, js *payment.JPayService) *TestHandler {
	return &TestHandler{db: db, snowflake: snowflake, js: js}
}

func (h *TestHandler) SseTest(c *gin.Context) {
	//c.Header("Content-Type", "text/event-stream")
	//c.Header("Cache-Control", "no-cache")
	//c.Header("Connection", "keep-alive")
	//
	//
	//// 模拟实时数据更新
	//for i := 0; i < 10; i++ {
	//	// 发送 SSE 数据
	//	_, err := fmt.Fprintf(c.Writer, "data: %v\n\n", data)
	//	if err != nil {
	//		return
	//	}
	//	c.Writer.Flush()            // 确保立即发送数据
	//	time.Sleep(1 * time.Second) // 每秒发送一次数据
	//}
	//c.Abort()
}

func (h *TestHandler) PostTest(c *gin.Context) {
	var data struct {
		Message string `json:"message"`
		UserId  uint   `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将参数存储在上下文中
	c.Set("data", data)
	c.Next()
}
