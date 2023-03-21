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
	token := c.Query("token")
	if token != "RockYang" {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
		return
	}

	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	// API key
	if key, ok := data["api_key"]; ok && len(key) > 20 {
		s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys, key)
	}

	// proxy URL
	if proxy, ok := data["proxy"]; ok {
		s.Config.ProxyURL = proxy
	}

	// Model
	if model, ok := data["model"]; ok {
		s.Config.Chat.Model = model
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
		s.Config.Chat.Temperature = float32(v)
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
		s.Config.Chat.MaxTokens = v
	}

	// enable Context
	if enableContext, ok := data["enable_context"]; ok {
		v, err := strconv.ParseBool(enableContext)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "enable_context must be a bool parameter",
			})
			return
		}
		s.Config.Chat.EnableContext = v
	}

	// enable auth
	if enableAuth, ok := data["enable_auth"]; ok {
		v, err := strconv.ParseBool(enableAuth)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "enable_auth must be a bool parameter",
			})
			return
		}
		s.Config.EnableAuth = v
	}

	if token, ok := data["token"]; ok {
		s.Config.Tokens = append(s.Config.Tokens, token)
	}

	// 保存配置文件
	logger.Infof("Config: %+v", s.Config)
	err = types.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}
