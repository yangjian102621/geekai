package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"openai/types"
	"openai/utils"
)

func (s *Server) TestHandle(c *gin.Context) {
	var data map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if v, ok := data["opt"]; ok && v == "init_user" {
		users := GetUsers()
		for _, user := range users {
			user.Status = true
			_ = PutUser(user)
		}

		c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: GetUsers()})
	}

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
	if accessKey, ok := data["access_key"]; ok {
		s.Config.AccessKey = accessKey.(string)
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

	// 保存配置文件
	err = utils.SaveConfig(s.Config, s.ConfigPath)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config file"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: s.Config})
}

// SetDebugHandle 开启/关闭调试模式
func (s *Server) SetDebugHandle(c *gin.Context) {
	var data struct {
		Debug bool `json:"debug"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid args"})
		return
	}

	s.DebugMode = data.Debug
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg})
}

// AddUserHandle 添加 Username
func (s *Server) AddUserHandle(c *gin.Context) {
	var data struct {
		Name          string `json:"name"`
		MaxCalls      int    `json:"max_calls"`
		EnableHistory bool   `json:"enable_history"`
		Term          int    `json:"term"` // 有效期
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid args"})
		return
	}

	// 参数处理
	if data.Name == "" || data.MaxCalls < 0 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid args"})
		return
	}

	// 检查当前要添加的 Username 是否已经存在
	_, err = GetUser(data.Name)
	if err == nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Username " + data.Name + " already exists"})
		return
	}

	user := types.User{
		Name:           data.Name,
		MaxCalls:       data.MaxCalls,
		RemainingCalls: data.MaxCalls,
		EnableHistory:  data.EnableHistory,
		Term:           data.Term,
		Status:         true}
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
		Term          int  `json:"term"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil || data.MaxCalls <= 0 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid args"})
		return
	}

	var users = make([]types.User, 0)
	for i := 0; i < data.Number; i++ {
		name := utils.RandString(12)
		_, err := GetUser(name)
		for err == nil {
			name = utils.RandString(12)
		}
		user := types.User{
			Name:           name,
			MaxCalls:       data.MaxCalls,
			RemainingCalls: data.MaxCalls,
			EnableHistory:  data.EnableHistory,
			Term:           data.Term,
			Status:         true}
		err = PutUser(user)
		if err == nil {
			users = append(users, user)
		}
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: users})
}

func (s *Server) SetUserHandle(c *gin.Context) {
	var data map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid args"})
		return
	}

	var user *types.User
	if name, ok := data["name"]; ok {
		user, err = GetUser(name.(string))
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "User not found"})
			return
		}
	}
	var maxCalls int
	if v, ok := data["max_calls"]; ok {
		maxCalls = int(v.(float64))
	}
	if maxCalls < 0 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid args"})
		return
	} else if maxCalls > 0 {
		user.RemainingCalls += maxCalls - user.MaxCalls
		if user.RemainingCalls < 0 {
			user.RemainingCalls = 0
		}
	}

	if v, ok := data["status"]; ok {
		user.Status = v.(bool)
	}
	if v, ok := data["enable_history"]; ok {
		user.EnableHistory = v.(bool)
	}
	if v, ok := data["remaining_calls"]; ok {
		user.RemainingCalls = int(v.(float64))
	}
	if v, ok := data["expired_time"]; ok {
		user.ExpiredTime = int64(v.(float64))
	}
	if v, ok := data["api_key"]; ok {
		user.ApiKey = v.(string)
	}

	err = PutUser(*user)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save configs"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: user})
}

// RemoveUserHandle 删除 Username
func (s *Server) RemoveUserHandle(c *gin.Context) {
	var data struct {
		Name string `json:"name"`
	}
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
	var data struct {
		ApiKey string `json:"api_key"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if len(data.ApiKey) > 20 {
		s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys, data.ApiKey)
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
		if v == data.ApiKey {
			s.Config.Chat.ApiKeys = append(s.Config.Chat.ApiKeys[:i], s.Config.Chat.ApiKeys[i+1:]...)
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
	var data map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var key string
	if v, ok := data["key"]; !ok {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Must specified the role key"})
		return
	} else {
		key = v.(string)
	}

	role, err := GetChatRole(key)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Role key not exists"})
		return
	}

	if v, ok := data["name"]; ok {
		role.Name = v.(string)
	}
	if v, ok := data["hello_msg"]; ok {
		role.HelloMsg = v.(string)
	}
	if v, ok := data["icon"]; ok {
		role.Icon = v.(string)
	}
	if v, ok := data["enable"]; ok {
		role.Enable = v.(bool)
	}
	if v, ok := data["context"]; ok {
		bytes, _ := json.Marshal(v)
		_ = json.Unmarshal(bytes, &role.Context)
	}

	err = PutChatRole(*role)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: data})
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
