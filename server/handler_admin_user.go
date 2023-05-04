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

	var users = make([]UserVo, 0)
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
			users = append(users, user2vo(user))
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
		user.ExpiredTime = utils.Str2stamp(v.(string))
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

type UserVo struct {
	Name           string   `json:"name"`
	MaxCalls       int      `json:"max_calls"`         // 最多调用次数，如果为 0 则表示不限制
	RemainingCalls int      `json:"remaining_calls"`   // 剩余调用次数
	EnableHistory  bool     `json:"enable_history"`    // 是否启用聊天记录
	Status         bool     `json:"status"`            // 当前状态
	Term           int      `json:"term" default:"30"` // 会员有效期，单位：天
	ActiveTime     string   `json:"active_time"`       // 激活时间
	ExpiredTime    string   `json:"expired_time"`      // 到期时间
	ApiKey         string   `json:"api_key"`           // OpenAI  API KEY
	ChatRoles      []string `json:"chat_roles"`        // 当前用户已订阅的聊天角色 map[role_key] => 0/1
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

	users := make([]UserVo, 0)
	for _, u := range GetUsers() {
		users = append(users, user2vo(u))
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: users})
}

// 将 User 实体转为 UserVo 实体
func user2vo(user types.User) UserVo {
	vo := UserVo{
		Name:           user.Name,
		MaxCalls:       user.MaxCalls,
		RemainingCalls: user.RemainingCalls,
		EnableHistory:  user.EnableHistory,
		Status:         user.Status,
		Term:           user.Term,
		ActiveTime:     utils.Stamp2str(user.ActiveTime),
		ExpiredTime:    utils.Stamp2str(user.ExpiredTime),
		ChatRoles:      make([]string, 0),
	}
	for k := range user.ChatRoles {
		vo.ChatRoles = append(vo.ChatRoles, k)
	}
	return vo
}
