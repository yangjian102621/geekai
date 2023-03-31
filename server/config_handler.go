package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"openai/types"
	"openai/utils"
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

	// Model
	if model, ok := data["model"]; ok {
		s.Config.Chat.Model = model
	}

	if accessKey, ok := data["access_key"]; ok {
		s.Config.AccessKey = accessKey
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

	// max_users
	if maxTokens, ok := data["max_tokens"]; ok {
		v, err := strconv.Atoi(maxTokens)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "max_users must be a int parameter",
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

	if expireTime, ok := data["chat_context_expire_time"]; ok {
		v, err := strconv.Atoi(expireTime)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.InvalidParams,
				Message: "chat_context_expire_time must be a integer parameter",
			})
			return
		}
		s.Config.Chat.ChatContextExpireTime = v
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

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}

// SetDebugHandle 开启/关闭调试模式
func (s *Server) SetDebugHandle(c *gin.Context) {
	var data struct {
		Debug bool `json:"debug"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	s.DebugMode = data.Debug
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}

// AddUserHandle 添加 Username
func (s *Server) AddUserHandle(c *gin.Context) {
	var data types.User
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	// 参数处理
	if data.Name == "" || data.MaxCalls < 0 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	// 检查当前要添加的 Username 是否已经存在
	_, err = GetUser(data.Name)
	if err == nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Username " + data.Name + " already exists"})
		return
	}

	user := types.User{Name: data.Name, MaxCalls: data.MaxCalls, RemainingCalls: data.MaxCalls, EnableHistory: data.EnableHistory}
	err = PutUser(user)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save configs"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: user})
}

// BatchAddUserHandle 批量生成 Username
func (s *Server) BatchAddUserHandle(c *gin.Context) {
	var data struct {
		Number        int  `json:"number"`
		MaxCalls      int  `json:"max_calls"`
		EnableHistory bool `json:"enable_history"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil || data.MaxCalls <= 0 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	var users = make([]string, 0)
	for i := 0; i < data.Number; i++ {
		name := utils.RandString(12)
		_, err := GetUser(name)
		for err == nil {
			name = utils.RandString(12)
		}
		err = PutUser(types.User{Name: name, MaxCalls: data.MaxCalls, RemainingCalls: data.MaxCalls, EnableHistory: data.EnableHistory})
		if err == nil {
			users = append(users, name)
		}
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: users})
}

func (s *Server) SetUserHandle(c *gin.Context) {
	var data types.User
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	// 参数处理
	if data.Name == "" || data.MaxCalls < 0 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	user, err := GetUser(data.Name)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Username not found"})
		return
	}

	if data.MaxCalls > 0 {
		user.RemainingCalls += data.MaxCalls - user.MaxCalls
		if user.RemainingCalls < 0 {
			user.RemainingCalls = 0
		}
	}
	user.MaxCalls = data.MaxCalls
	user.EnableHistory = data.EnableHistory

	err = PutUser(*user)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save configs"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: user})
}

// RemoveUserHandle 删除 Username
func (s *Server) RemoveUserHandle(c *gin.Context) {
	var data types.User
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	}

	err = RemoveUser(data.Name)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save configs"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}

// GetUserListHandle 获取用户列表
func (s *Server) GetUserListHandle(c *gin.Context) {
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: GetUsers()})
}

// AddApiKeyHandle 添加一个 API key
func (s *Server) AddApiKeyHandle(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if key, ok := data["api_key"]; ok && len(key) > 20 {
		s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys, key)
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
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if key, ok := data["api_key"]; ok {
		for i, v := range s.Config.Chat.ApiKeys {
			if v == key {
				s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys[:i], s.Config.Chat.ApiKeys[i+1:]...)
			}
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

// GetChatRoleListHandle 获取聊天角色列表
func (s *Server) GetChatRoleListHandle(c *gin.Context) {
	var rolesOrder = []string{"gpt", "programmer", "teacher", "artist", "philosopher", "lu-xun", "english_trainer", "seller"}
	var res = make([]interface{}, 0)
	var roles = GetChatRoles()
	for _, k := range rolesOrder {
		if v, ok := roles[k]; ok && v.Enable {
			res = append(res, struct {
				Key  string `json:"key"`
				Name string `json:"name"`
				Icon string `json:"icon"`
			}{
				Key:  v.Key,
				Name: v.Name,
				Icon: v.Icon,
			})
		}
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: res})
}

// GetChatRoleHandle 获取指定的角色
func (s *Server) GetChatRoleHandle(c *gin.Context) {
	var data struct {
		Key string `json:"key"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	role, err := GetChatRole(data.Key)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "No role found"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: role})
	return
}

// SetChatRoleHandle 更新某个聊天角色信息，这里只允许更改名称以及启用和禁用角色操作
func (s *Server) SetChatRoleHandle(c *gin.Context) {
	var data types.ChatRole
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if data.Key == "" {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Must specified the role key"})
		return
	}

	_, err = GetChatRole(data.Key)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Role key not exists"})
		return
	}

	err = PutChatRole(data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: data})
}

// AddProxyHandle 添加一个代理
func (s *Server) AddProxyHandle(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if proxy, ok := data["proxy"]; ok {
		if !utils.ContainsStr(s.Config.ProxyURL, proxy) {
			s.Config.ProxyURL = append(s.Config.ProxyURL, proxy)
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
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if proxy, ok := data["proxy"]; ok {
		for i, v := range s.Config.ProxyURL {
			if v == proxy {
				s.Config.ProxyURL = append(s.Config.ProxyURL[:i], s.Config.ProxyURL[i+1:]...)
				break
			}
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
