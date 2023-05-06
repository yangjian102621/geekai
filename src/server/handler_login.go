package server

import (
	"chatplus/types"
	"chatplus/utils"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func (s *Server) LoginHandle(c *gin.Context) {
	var data struct {
		Token string `json:"token"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
		return
	}
	username := strings.TrimSpace(data.Token)
	user, err := GetUser(username)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid user"})
		return
	}

	sessionId := utils.RandString(42)
	session := sessions.Default(c)
	session.Set(sessionId, username)
	err = session.Save()
	if err != nil {
		logger.Error("Error for save session: ", err)
	}
	// 记录客户端 IP 地址
	s.ChatSession[sessionId] = types.ChatSession{ClientIP: c.ClientIP(), Username: username, SessionId: sessionId}
	// 更新用户激活时间
	user.ActiveTime = time.Now().Unix()
	if user.ExpiredTime == 0 {
		activeTime := time.Unix(user.ActiveTime, 0)
		if user.Term == 0 {
			user.Term = 30 // 默认 30 天到期
		}
		user.ExpiredTime = activeTime.Add(time.Hour * 24 * time.Duration(user.Term)).Unix()
	}
	err = PutUser(*user)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Save user info failed"})
		return
	}

	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: struct {
		User      types.User `json:"user"`
		SessionId string     `json:"session_id"`
	}{User: *user, SessionId: sessionId}})
}

// ManagerLoginHandle 管理员登录
func (s *Server) ManagerLoginHandle(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
		return
	}
	username := strings.TrimSpace(data.Username)
	password := strings.TrimSpace(data.Password)
	if username == s.Config.Manager.Username && password == s.Config.Manager.Password {
		sessionId := utils.RandString(42)
		session := sessions.Default(c)
		session.Set(sessionId, username)
		err = session.Save()
		// 记录登录信息
		s.ChatSession[sessionId] = types.ChatSession{ClientIP: c.ClientIP(), Username: username, SessionId: sessionId}
		c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: struct {
			User      types.Manager `json:"user"`
			SessionId string        `json:"session_id"`
		}{User: data, SessionId: sessionId}})
	} else {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "用户名或者密码错误"})
	}
}

// LogoutHandle 注销
func (s *Server) LogoutHandle(c *gin.Context) {
	sessionId := c.GetHeader(types.TokenName)
	session := sessions.Default(c)
	session.Delete(sessionId)
	err := session.Save()
	if err != nil {
		logger.Error("Error for save session: ", err)
	}
	// 删除 websocket 会话列表
	delete(s.ChatSession, sessionId)
	// 关闭 socket 连接
	if client, ok := s.ChatClients[sessionId]; ok {
		client.Close()
	}
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success})
}

func (s *Server) GetSessionHandle(c *gin.Context) {
	sessionId := c.GetHeader(types.TokenName)
	if session, ok := s.ChatSession[sessionId]; ok && session.ClientIP == c.ClientIP() {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Success})
	} else {
		c.JSON(http.StatusOK, types.BizVo{
			Code:    types.NotAuthorized,
			Message: "Not Authorized",
		})
	}

}
