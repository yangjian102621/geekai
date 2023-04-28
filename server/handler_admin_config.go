package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"openai/types"
	"openai/utils"
)

func (s *Server) TestHandle(c *gin.Context) {
	roles := types.GetDefaultChatRole()
	for _, v := range roles {
		_ = PutChatRole(v)
	}
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: GetChatRoles()})

}

// ConfigSetHandle set configs
func (s *Server) ConfigSetHandle(c *gin.Context) {
	var data map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if model, ok := data["model"]; ok {
		s.Config.Chat.Model = model.(string)
	}
	// Temperature
	if temperature, ok := data["temperature"]; ok {
		s.Config.Chat.Temperature = temperature.(float32)
	}
	// max_users
	if maxTokens, ok := data["max_tokens"]; ok {
		s.Config.Chat.MaxTokens = int(maxTokens.(float64))
	}
	// enable Context
	if enableContext, ok := data["enable_context"]; ok {
		s.Config.Chat.EnableContext = enableContext.(bool)
	}
	if expireTime, ok := data["chat_context_expire_time"]; ok {
		s.Config.Chat.ChatContextExpireTime = int(expireTime.(float64))
	}
	// enable auth
	if enableAuth, ok := data["enable_auth"]; ok {
		s.Config.EnableAuth = enableAuth.(bool)
	}

	if debug, ok := data["debug"]; ok {
		s.DebugMode = debug.(bool)
	}
	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config})
}

// AddProxyHandle 添加一个代理
func (s *Server) AddProxyHandle(c *gin.Context) {
	var data struct {
		Proxy string `json:"proxy"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if data.Proxy != "" {
		if !utils.ContainsStr(s.Config.ProxyURL, data.Proxy) {
			s.Config.ProxyURL = append(s.Config.ProxyURL, data.Proxy)
		}
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.ProxyURL})
}

// RemoveProxyHandle 删除一个代理
func (s *Server) RemoveProxyHandle(c *gin.Context) {
	var data struct {
		Proxy string `json:"proxy"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for i, v := range s.Config.ProxyURL {
		if v == data.Proxy {
			s.Config.ProxyURL = append(s.Config.ProxyURL[:i], s.Config.ProxyURL[i+1:]...)
			break
		}
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.ProxyURL})
}
