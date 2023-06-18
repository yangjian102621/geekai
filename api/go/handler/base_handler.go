package handler

import (
	"chatplus/core"
	logger2 "chatplus/logger"
	"strconv"
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
	return intValue(c.PostForm(key), defaultValue)
}

func (h *BaseHandler) GetInt(c *gin.Context, key string, defaultValue int) int {
	return intValue(c.Query(key), defaultValue)
}

func intValue(str string, defaultValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}

func (h *BaseHandler) GetFloat(c *gin.Context, key string) float64 {
	return floatValue(c.Query(key))
}
func (h *BaseHandler) PostFloat(c *gin.Context, key string) float64 {
	return floatValue(c.PostForm(key))
}

func floatValue(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return value
}

func (h *BaseHandler) GetBool(c *gin.Context, key string) bool {
	return boolValue(c.Query(key))
}
func (h *BaseHandler) PostBool(c *gin.Context, key string) bool {
	return boolValue(c.PostForm(key))
}

func boolValue(str string) bool {
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return value
}
