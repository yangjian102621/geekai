package handler

import (
	"chatplus/core"
	logger2 "chatplus/logger"
	"chatplus/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

var logger = logger2.GetLogger()

type BaseHandler struct {
	App *core.AppServer
}

func (h *BaseHandler) GetTrim(c *gin.Context, key string) string {
	return strings.TrimSpace(c.Query(key))
}

func (h *BaseHandler) PostInt(c *gin.Context, key string, defaultValue int) int {
	return utils.IntValue(c.PostForm(key), defaultValue)
}

func (h *BaseHandler) GetInt(c *gin.Context, key string, defaultValue int) int {
	return utils.IntValue(c.Query(key), defaultValue)
}

func (h *BaseHandler) GetFloat(c *gin.Context, key string) float64 {
	return utils.FloatValue(c.Query(key))
}
func (h *BaseHandler) PostFloat(c *gin.Context, key string) float64 {
	return utils.FloatValue(c.PostForm(key))
}

func (h *BaseHandler) GetBool(c *gin.Context, key string) bool {
	return utils.BoolValue(c.Query(key))
}
func (h *BaseHandler) PostBool(c *gin.Context, key string) bool {
	return utils.BoolValue(c.PostForm(key))
}
