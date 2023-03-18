package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"openai/types"
	"strconv"
)

// ConfigSetHandle set configs
func (s *Server) ConfigSetHandle(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	// API key
	if key, ok := data["api_key"]; ok && len(key) > 20 {
		s.Config.OpenAi.ApiKeys = append(s.Config.OpenAi.ApiKeys, key)
	}

	// proxy URL
	if proxy, ok := data["proxy"]; ok {
		s.Config.ProxyURL = proxy
	}

	// Model
	if model, ok := data["model"]; ok {
		s.Config.OpenAi.Model = model
	}

	// Temperature
	if temperature, ok := data["temperature"]; ok {
		v, err := strconv.ParseFloat(temperature, 32)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "temperature must be a float parameter",
			})
			return
		}
		s.Config.OpenAi.Temperature = float32(v)
	}

	// max_tokens
	if maxTokens, ok := data["max_tokens"]; ok {
		v, err := strconv.Atoi(maxTokens)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "max_tokens must be a int parameter",
			})
			return
		}
		s.Config.OpenAi.MaxTokens = v

	}

	// 保存配置文件
	logger.Infof("Config: %+v", s.Config)
	types.SaveConfig(s.Config, s.ConfigPath)
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}
