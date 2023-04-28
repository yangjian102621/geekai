package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"openai/types"
	"openai/utils"
	"strings"
)

// AddUserHandle 添加 Username
func (s *Server) AddUserHandle(c *gin.Context) {
	var data struct {
		Name          string   `json:"name"`
		MaxCalls      int      `json:"max_calls"`
		EnableHistory bool     `json:"enable_history"`
		Term          int      `json:"term"`       // 有效期
		ChatRoles     []string `json:"chat_roles"` // 订阅角色
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

	var chatRoles = make(map[string]int)
	if len(data.ChatRoles) > 0 {
		if data.ChatRoles[0] == "all" { // 所有的角色
			roles := GetChatRoles()
			for key := range roles {
				chatRoles[key] = 1
			}
		} else {
			for _, key := range data.ChatRoles {
				chatRoles[key] = 1
			}
		}

	}

	user := types.User{
		Name:           data.Name,
		MaxCalls:       data.MaxCalls,
		RemainingCalls: data.MaxCalls,
		EnableHistory:  data.EnableHistory,
		Term:           data.Term,
		ChatRoles:      chatRoles,
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
		Number        int      `json:"number"`
		MaxCalls      int      `json:"max_calls"`
		EnableHistory bool     `json:"enable_history"`
		Term          int      `json:"term"`
		ChatRoles     []string `json:"chat_roles"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil || data.MaxCalls <= 0 {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid args"})
		return
	}

	var chatRoles = make(map[string]int)
	if len(data.ChatRoles) > 0 {
		if data.ChatRoles[0] == "all" { // 所有的角色
			roles := GetChatRoles()
			for key := range roles {
				chatRoles[key] = 1
			}
		} else {
			for _, key := range data.ChatRoles {
				chatRoles[key] = 1
			}
		}

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
			ChatRoles:      chatRoles,
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
		user.MaxCalls = maxCalls
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
	if v, ok := data["chat_roles"]; ok {
		if roles, ok := v.([]interface{}); ok {
			chatRoles := make(map[string]int)
			if roles[0] == "all" {
				roles := GetChatRoles()
				for key := range roles {
					chatRoles[key] = 1
				}
			} else {
				for _, key := range roles {
					key := strings.TrimSpace(fmt.Sprintf("%v", key))
					chatRoles[key] = 1
				}
			}
			user.ChatRoles = chatRoles
		}

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
	username := c.PostForm("username")
	if username != "" {
		user, err := GetUser(username)
		if err != nil {
			c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "User not exists"})
		} else {
			c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: user})
		}
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: GetUsers()})
}
