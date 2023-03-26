package server

import (
	"embed"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	logger2 "openai/logger"
	"openai/types"
	"openai/utils"
	"os"
	"path/filepath"
	"strings"
)

var logger = logger2.GetLogger()

type StaticFile struct {
	embedFS embed.FS
	path    string
}

func (s StaticFile) Open(name string) (fs.File, error) {
	filename := filepath.Join(s.path, strings.TrimLeft(name, "/"))
	file, err := s.embedFS.Open(filename)
	return file, err
}

type Server struct {
	Config      *types.Config
	ConfigPath  string
	ChatContext map[string][]types.Message // 聊天上下文 [SessionID] => []Messages

	// 保存 Websocket 会话 Token, 每个 Token 只能连接一次
	// 防止第三方直接连接 socket 调用 OpenAI API
	WsSession        map[string]string
	ApiKeyAccessStat map[string]int64 // 记录每个 API Key 的最后访问之间，保持在 15/min 之内
}

func NewServer(configPath string) (*Server, error) {
	// load service configs
	config, err := types.LoadConfig(configPath)
	if config.ChatRoles == nil {
		config.ChatRoles = types.GetDefaultChatRole()
	}
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:           config,
		ConfigPath:       configPath,
		ChatContext:      make(map[string][]types.Message, 16),
		WsSession:        make(map[string]string),
		ApiKeyAccessStat: make(map[string]int64),
	}, nil
}

func (s *Server) Run(webRoot embed.FS, path string, debug bool) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	if debug {
		engine.Use(corsMiddleware())
	}
	engine.Use(sessionMiddleware(s.Config))
	engine.Use(AuthorizeMiddleware(s))

	engine.GET("/hello", Hello)
	engine.GET("/api/session/get", s.GetSessionHandle)
	engine.POST("/api/login", s.LoginHandle)
	engine.Any("/api/chat", s.ChatHandle)
	engine.POST("/api/config/set", s.ConfigSetHandle)
	engine.GET("/api/config/chat-roles/get", s.GetChatRoles)
	engine.POST("api/config/token/add", s.AddToken)
	engine.POST("api/config/token/remove", s.RemoveToken)
	engine.POST("api/config/apikey/add", s.AddApiKey)
	engine.POST("api/config/apikey/remove", s.RemoveApiKey)
	engine.POST("api/config/apikey/list", s.ListApiKeys)
	engine.POST("api/config/role/set", s.UpdateChatRole)

	engine.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/favicon.ico" {
			c.Redirect(http.StatusMovedPermanently, "/chat/"+c.Request.URL.Path)
		}
	})

	// process front-end web static files
	engine.StaticFS("/chat", http.FS(StaticFile{
		embedFS: webRoot,
		path:    path,
	}))

	logger.Infof("http://%s", s.Config.Listen)
	err := engine.Run(s.Config.Listen)

	if err != nil {
		logger.Error("Fail to start server:", err.Error())
		os.Exit(1)
	}

}

func sessionMiddleware(config *types.Config) gin.HandlerFunc {
	// encrypt the cookie
	store := cookie.NewStore([]byte(config.Session.SecretKey))
	store.Options(sessions.Options{
		Path:     config.Session.Path,
		Domain:   config.Session.Domain,
		MaxAge:   config.Session.MaxAge,
		Secure:   config.Session.Secure,
		HttpOnly: config.Session.HttpOnly,
		SameSite: config.Session.SameSite,
	})
	return sessions.Sessions(config.Session.Name, store)
}

// 跨域中间件设置
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 设置允许的请求源
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, ChatGPT-Token")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == http.MethodOptions {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

// AuthorizeMiddleware 用户授权验证
func AuthorizeMiddleware(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.Config.EnableAuth ||
			c.Request.URL.Path == "/api/login" ||
			c.Request.URL.Path == "/api/config/chat-roles/get" ||
			!strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next()
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/api/config") {
			accessKey := c.Query("access_key")
			if accessKey != "RockYang" {
				c.Abort()
				c.JSON(http.StatusOK, types.BizVo{Code: types.NotAuthorized, Message: "No Permissions"})
			} else {
				c.Next()
			}
			return
		}

		// WebSocket 连接请求验证
		if c.Request.URL.Path == "/api/chat" {
			tokenName := c.Query("token")
			if addr, ok := s.WsSession[tokenName]; ok && addr == c.ClientIP() {
				// 每个令牌只能连接一次
				//delete(s.WsSession, tokenName)
				c.Next()
			} else {
				c.Abort()
			}
			return
		}

		tokenName := c.GetHeader(types.TokenName)
		session := sessions.Default(c)
		userInfo := session.Get(tokenName)
		if userInfo != nil {
			c.Set(types.SessionKey, userInfo)
			c.Next()
		} else {
			c.Abort()
			c.JSON(http.StatusOK, types.BizVo{
				Code:    types.NotAuthorized,
				Message: "Not Authorized",
			})
		}
	}
}

func (s *Server) GetSessionHandle(c *gin.Context) {
	tokenName := c.GetHeader(types.TokenName)
	if addr, ok := s.WsSession[tokenName]; ok && addr == c.ClientIP() {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: addr})
	} else {
		c.JSON(http.StatusOK, types.BizVo{
			Code:    types.NotAuthorized,
			Message: "Not Authorized",
		})
	}

}

func (s *Server) LoginHandle(c *gin.Context) {
	var data map[string]string
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
		return
	}
	token := data["token"]
	if !utils.ContainsItem(s.Config.Tokens, token) {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Invalid token"})
		return
	}

	sessionId := utils.RandString(42)
	session := sessions.Default(c)
	session.Set(sessionId, token)
	err = session.Save()
	if err != nil {
		logger.Error("Error for save session: ", err)
	}
	// 记录客户端 IP 地址
	s.WsSession[sessionId] = c.ClientIP()
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Data: sessionId})
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, types.BizVo{Code: types.Success, Message: "HELLO, ChatGPT !!!"})
}
