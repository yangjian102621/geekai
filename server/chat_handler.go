package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) Chat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": fmt.Sprintf("HELLO, ChatGPT !!!")})
}
