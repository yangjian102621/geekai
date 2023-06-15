package param

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func GetTrim(c *gin.Context, key string) string {
	return strings.TrimSpace(c.Query(key))
}

func GetInt(c *gin.Context, key string, defaultValue int) int {
	return intValue(c.Query(key), defaultValue)
}

func PostInt(c *gin.Context, key string, defaultValue int) int {
	return intValue(c.PostForm(key), defaultValue)
}

func intValue(str string, defaultValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetFloat(c *gin.Context, key string) float64 {
	return floatValue(c.Query(key))
}
func PostFloat(c *gin.Context, key string) float64 {
	return floatValue(c.PostForm(key))
}

func floatValue(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return value
}

func GetBool(c *gin.Context, key string) bool {
	return boolValue(c.Query(key))
}
func PostBool(c *gin.Context, key string) bool {
	return boolValue(c.PostForm(key))
}

func boolValue(str string) bool {
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return value
}
