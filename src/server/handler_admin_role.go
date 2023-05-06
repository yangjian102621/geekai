package server

import (
	"chatplus/types"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllChatRolesHandle 获取所有聊天角色列表
func (s *Server) GetAllChatRolesHandle(c *gin.Context) {
	var rolesOrder = []string{"gpt", "teacher", "translator", "english_trainer", "weekly_report", "girl_friend",
		"kong_zi", "lu_xun", "steve_jobs", "elon_musk", "red_book", "dou_yin", "programmer",
		"seller", "good_comment", "psychiatrist", "artist"}
	var res = make([]interface{}, 0)
	var roles = GetChatRoles()
	for _, k := range rolesOrder {
		if v, ok := roles[k]; ok {
			res = append(res, v)
		}
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: res})
}

// GetChatRoleListHandle 获取当前登录用户的角色列表
func (s *Server) GetChatRoleListHandle(c *gin.Context) {
	sessionId := c.GetHeader(types.TokenName)
	session := s.ChatSession[sessionId]
	user, err := GetUser(session.Username)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Hacker Access!!!"})
		return
	}
	var rolesOrder = []string{"gpt", "teacher", "translator", "english_trainer", "weekly_report", "girl_friend",
		"kong_zi", "lu_xun", "steve_jobs", "elon_musk", "red_book", "dou_yin", "programmer",
		"seller", "good_comment", "psychiatrist", "artist"}
	var res = make([]interface{}, 0)
	var roles = GetChatRoles()
	for _, k := range rolesOrder {
		// 确认当前用户是否订阅了当前角色
		if v, ok := user.ChatRoles[k]; !ok || v != 1 {
			continue
		}

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

// AddChatRoleHandle 添加一个聊天角色
func (s *Server) AddChatRoleHandle(c *gin.Context) {
	var data types.ChatRole
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		logger.Errorf("Error decode json data: %s", err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if data.Key == "" || data.Name == "" || data.Icon == "" {
		c.JSON(http.StatusOK, types.BizVo{Code: types.InvalidParams, Message: "Invalid parameters"})
		return
	}

	err = PutChatRole(data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save levelDB"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: data})
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
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Failed to save levelDB"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: types.OkMsg, Data: data})
}
