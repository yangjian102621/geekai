package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"openai/types"
	"openai/utils"
	"strings"
)

func (s *Server) TestHandle(c *gin.Context) {
	roles := types.GetDefaultChatRole()
	for _, v := range roles {
		_ = PutChatRole(v)
	}
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: GetChatRoles()})

}

func (s *Server) ConfigGetHandle(c *gin.Context)  {
	data := struct {
		Title string `json:"title"`
		ConsoleTitle string `json:"console_title"`
		ProxyURL string `json:"proxy_url"`
		Model string `json:"model"`
		Temperature float32 `json:"temperature"`
		MaxTokens int `json:"max_tokens"`
		ChatContextExpireTime int `json:"chat_context_expire_time"`
		EnableContext bool `json:"enable_context"`
	}{
		Title: s.Config.Title,
		ConsoleTitle: s.Config.ConsoleTitle,
		ProxyURL: strings.Join(s.Config.ProxyURL, ","),
		Model:  s.Config.Chat.Model,
		Temperature: s.Config.Chat.Temperature,
		MaxTokens: s.Config.Chat.MaxTokens,
		EnableContext: s.Config.Chat.EnableContext,
		ChatContextExpireTime: s.Config.Chat.ChatContextExpireTime,
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: data})
}

// ConfigSetHandle set configs
func (s *Server) ConfigSetHandle(c *gin.Context) {
	var data struct {
		Title string `json:"title"`
		ConsoleTitle string `json:"console_title"`
		ProxyURL string `json:"proxy_url"`
		Model string `json:"model"`
		Temperature float32 `json:"temperature"`
		MaxTokens int `json:"max_tokens"`
		ChatContextExpireTime int `json:"chat_context_expire_time"`
		EnableContext bool `json:"enable_context"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	s.Config.Title = data.Title
	s.Config.ConsoleTitle = data.ConsoleTitle
	s.Config.ProxyURL = strings.Split(data.ProxyURL, ",")
	s.Config.Chat.Model = data.Model
	s.Config.Chat.Temperature = data.Temperature
	s.Config.Chat.MaxTokens = data.MaxTokens
	s.Config.Chat.EnableContext = data.EnableContext
	s.Config.Chat.ChatContextExpireTime = data.ChatContextExpireTime

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
