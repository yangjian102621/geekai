package core

import (
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
)

type AppServer struct {
	AppConfig    *types.AppConfig
	Engine       *gin.Engine
	ChatContexts *types.LMap[string, []types.Message] // 聊天上下文 Map [chatId] => []Message
	ChatConfig   *types.ChatConfig                    // 聊天配置

	// 保存 Websocket 会话 UserId, 每个 UserId 只能连接一次
	// 防止第三方直接连接 socket 调用 OpenAI API
	ChatSession   *types.LMap[string, types.ChatSession]  //map[sessionId]UserId
	ChatClients   *types.LMap[string, *types.WsClient]    // Websocket 连接集合
	ReqCancelFunc *types.LMap[string, context.CancelFunc] // HttpClient 请求取消 handle function
}

func NewServer(appConfig *types.AppConfig) *AppServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	return &AppServer{
		AppConfig:     appConfig,
		Engine:        gin.Default(),
		ChatContexts:  types.NewLMap[string, []types.Message](),
		ChatSession:   types.NewLMap[string, types.ChatSession](),
		ChatClients:   types.NewLMap[string, *types.WsClient](),
		ReqCancelFunc: types.NewLMap[string, context.CancelFunc](),
	}
}

func (s *AppServer) Init(debug bool) {
	if debug { // 调试模式允许跨域请求 API
		logger.Info("Enabled debug mode")
		s.Engine.Use(corsMiddleware())
	}
	s.Engine.Use(sessionMiddleware(s.AppConfig))
	s.Engine.Use(authorizeMiddleware(s))
	s.Engine.Use(errorHandler)
	//gob.Register(model.User{})
}

func (s *AppServer) Run(db *gorm.DB) error {
	// load chat config from database
	var config model.Config
	res := db.Where("marker", "chat").First(&config)
	if res.Error != nil {
		return res.Error
	}
	err := utils.JsonDecode(config.Config, &s.ChatConfig)
	if err != nil {
		return err
	}
	logger.Infof("http://%s", s.AppConfig.Listen)
	return s.Engine.Run(s.AppConfig.Listen)
}

// 全局异常处理
func errorHandler(c *gin.Context) {
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

// 会话处理
func sessionMiddleware(config *types.AppConfig) gin.HandlerFunc {
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
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, ChatGPT-TOKEN, ADMIN-SESSION-TOKEN")
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

// 用户授权验证
func authorizeMiddleware(s *AppServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/login" ||
			c.Request.URL.Path == "/api/admin/login" ||
			c.Request.URL.Path == "/api/user/register" ||
			c.Request.URL.Path == "/api/admin/config/get" {
			c.Next()
			return
		}

		// WebSocket 连接请求验证
		if c.Request.URL.Path == "/api/chat" {
			sessionId := c.Query("sessionId")
			session := s.ChatSession.Get(sessionId)
			if session.ClientIP == c.ClientIP() {
				c.Next()
			} else {
				c.Abort()
			}
			return
		}
		session := sessions.Default(c)
		var value interface{}
		if strings.Contains(c.Request.URL.Path, "/api/admin/") {
			value = session.Get(types.SessionAdmin)
		} else {
			value = session.Get(types.SessionUser)
		}
		if value != nil {
			c.Next()
		} else {
			resp.NotAuth(c)
			c.Abort()
		}
	}
}
