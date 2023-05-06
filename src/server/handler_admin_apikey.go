package server

import (
	"chatplus/types"
	"chatplus/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddApiKeyHandle 添加一个 API key
func (s *Server) AddApiKeyHandle(c *gin.Context) {
	var data struct {
		ApiKey string `json:"api_key"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	// 过滤已存在的 Key
	for _, key := range s.Config.Chat.ApiKeys {
		if key.Value == data.ApiKey {
			c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "API KEY 已存在"})
			return
		}
	}
	if len(data.ApiKey) > 20 {
		s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys, types.APIKey{Value: data.ApiKey})
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.Chat.ApiKeys})
}

// RemoveApiKeyHandle 移除一个 API key
func (s *Server) RemoveApiKeyHandle(c *gin.Context) {
	var data struct {
		ApiKey string `json:"api_key"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	for i, v := range s.Config.Chat.ApiKeys {
		if v.Value == data.ApiKey {
			s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys[:i], s.Config.Chat.ApiKeys[i+1:]...)
			break
		}
	}

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.Chat.ApiKeys})
}

// ListApiKeysHandle 获取 API key 列表
func (s *Server) ListApiKeysHandle(c *gin.Context) {
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config.Chat.ApiKeys})
}
