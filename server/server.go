package server

import (
	"context"
	"embed"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
	logger2 "openai/logger"
	"openai/types"
	"openai/utils"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
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
	Config       *types.Config
	ConfigPath   string
	ChatContexts map[string]types.ChatContext // 聊天上下文 [SessionID+ChatRole] => ChatContext

	// 保存 Websocket 会话 Username, 每个 Username 只能连接一次
	// 防止第三方直接连接 socket 调用 OpenAI API
	ChatSession      map[string]types.ChatSession  //map[sessionId]User
	ApiKeyAccessStat map[string]int64              // 记录每个 API Key 的最后访问之间，保持在 15/min 之内
	ChatClients      map[string]*WsClient          // Websocket 连接集合
	ReqCancelFunc    map[string]context.CancelFunc // HttpClient 请求取消 handle function
	DebugMode        bool                          // 是否开启调试模式
}

func NewServer(configPath string) (*Server, error) {
	// load service configs
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	roles := GetChatRoles()
	if len(roles) == 0 { // 初始化默认聊天角色到 leveldb
		roles = types.GetDefaultChatRole()
		for _, v := range roles {
			err := PutChatRole(v)
			if err != nil {
				return nil, err
			}
		}
	}
	return &Server{
		Config:           config,
		ConfigPath:       configPath,
		ChatContexts:     make(map[string]types.ChatContext, 16),
		ChatSession:      make(map[string]types.ChatSession),
		ChatClients:      make(map[string]*WsClient),
		ReqCancelFunc:    make(map[string]context.CancelFunc),
		ApiKeyAccessStat: make(map[string]int64),
	}, nil
}

func (s *Server) Run(webRoot embed.FS, path string, debug bool) {
	s.DebugMode = debug
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	engine := gin.Default()
	if debug {
		engine.Use(corsMiddleware())
	}
	engine.Use(sessionMiddleware(s.Config))
	engine.Use(AuthorizeMiddleware(s))
	engine.Use(Recover)

	engine.POST("api/test", s.TestHandle)
	engine.GET("api/session/get", s.GetSessionHandle)
	engine.POST("api/login", s.LoginHandle)
	engine.POST("api/logout", s.LogoutHandle)
	engine.Any("api/chat", s.ChatHandle)
	engine.POST("api/chat/stop", s.StopGenerateHandle)
	engine.POST("api/chat/history", s.GetChatHistoryHandle)
	engine.POST("api/chat/history/clear", s.ClearHistoryHandle)
	engine.GET("api/role/hello", s.GetHelloMsgHandle)
	engine.POST("api/img/get", s.GetImgURLHandle)
	engine.POST("api/img/set", s.SetImgURLHandle)

	engine.POST("api/admin/config", s.ConfigSetHandle)
	engine.GET("api/admin/chat-roles/get", s.GetChatRoleListHandle)
	engine.GET("api/admin/chat-roles/add", s.AddChatRoleHandle)
	engine.POST("api/admin/user/add", s.AddUserHandle)
	engine.POST("api/admin/user/batch-add", s.BatchAddUserHandle)
	engine.POST("api/admin/user/set", s.SetUserHandle)
	engine.POST("api/admin/user/list", s.GetUserListHandle)
	engine.POST("api/admin/user/remove", s.RemoveUserHandle)
	engine.POST("api/admin/apikey/add", s.AddApiKeyHandle)
	engine.POST("api/admin/apikey/remove", s.RemoveApiKeyHandle)
	engine.POST("api/admin/apikey/list", s.ListApiKeysHandle)
	engine.POST("api/admin/role/set", s.SetChatRoleHandle)
	engine.POST("api/admin/role/get", s.GetChatRoleHandle)
	engine.POST("api/admin/proxy/add", s.AddProxyHandle)
	engine.POST("api/admin/proxy/remove", s.RemoveProxyHandle)

	engine.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path == "/favicon.ico" {
			c.Redirect(http.StatusMovedPermanently, "/chat/"+c.Request.URL.Path)
		}
		if c.Request.URL.Path == "/" {
			c.Redirect(http.StatusMovedPermanently, "/chat")
		}
	})

	// process front-end web static files
	engine.StaticFS("/chat", http.FS(StaticFile{
		embedFS: webRoot,
		path:    path,
	}))

	// 定时清理过期的会话
	go func() {
		for {
			for key, ctx := range s.ChatContexts {
				// 清理超过 60min 没有更新，则表示为过期会话
				if time.Now().Unix()-ctx.LastAccessTime >= int64(s.Config.Chat.ChatContextExpireTime) {
					logger.Infof("清理会话上下文: %s", key)
					delete(s.ChatContexts, key)
				}
			}
			time.Sleep(time.Second * 5) // 每隔 5 秒钟清理一次
		}
	}()

	logger.Infof("http://%s", s.Config.Listen)
	err := engine.Run(s.Config.Listen)

	if err != nil {
		logger.Error("Fail to start server:", err.Error())
		os.Exit(1)
	}

}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic: %v\n", r)
			debug.PrintStack()
			c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
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
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, ChatGPT-TOKEN, ACCESS-KEY")
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
				logger.Info("Panic info is: %v", err)
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

		//if strings.HasPrefix(c.Request.URL.Path, "/api/config") {
		//	accessKey := c.GetHeader("ACCESS-KEY")
		//	if accessKey != strings.TrimSpace(s.Config.AccessKey) {
		//		c.Abort()
		//		c.JSON(http.StatusOK, types.BizVo{Code: types.NotAuthorized, Message: "No Permissions"})
		//	} else {
		//		c.Next()
		//	}
		//	return
		//}

		// WebSocket 连接请求验证
		if c.Request.URL.Path == "/api/chat" {
			sessionId := c.Query("sessionId")
			if session, ok := s.ChatSession[sessionId]; ok && session.ClientIP == c.ClientIP() {
				c.Next()
			} else {
				c.Abort()
			}
			return
		}

		sessionId := c.GetHeader(types.TokenName)
		session := sessions.Default(c)
		userInfo := session.Get(sessionId)
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

// LogoutHandle 注销
func (s *Server) LogoutHandle(c *gin.Context) {
	var data struct {
		Opt string `json:"opt"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
		return
	}

	if data.Opt == "logout" {
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
	} else {
		c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: "Hack attempt!"})
	}
}
