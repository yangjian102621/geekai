package server

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	logger2 "openai/logger"
	"openai/types"
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
	Config     *types.Config
	ConfigPath string
	History    map[string][]types.Message
}

func NewServer(configPath string) (*Server, error) {
	// load service configs
	config, err := types.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return &Server{
		Config:     config,
		ConfigPath: configPath,
		History:    make(map[string][]types.Message, 16)}, nil
}

func (s *Server) Run(webRoot embed.FS, path string) {
	gin.SetMode(gin.DebugMode)
	engine := gin.Default()
	engine.Use(sessionMiddleware(s.Config))
	engine.Use(corsMiddleware())
	engine.Use(AuthorizeMiddleware())

	engine.GET("/hello", Hello)
	engine.Any("/api/chat", s.ChatHandle)
	engine.POST("/api/config/set", s.ConfigSetHandle)

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
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, Session-Name, Session")
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
func AuthorizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		//if c.Request.URL.Path == "/login" {
		//	c.Next()
		//	return
		//}

		//sessionName := c.GetHeader("Session-Name")
		//session, err := c.Cookie(sessionName)
		//if err == nil {
		//	c.Request.Header.Set(utils.SessionKey, session)
		//	c.Next()
		//} else {
		//	logger.Fatal(err)
		//	c.Abort()
		//	c.JSON(http.StatusUnauthorized, "No session data found")
		//}
	}
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": fmt.Sprintf("HELLO, XWEBSSH !!!")})
}
