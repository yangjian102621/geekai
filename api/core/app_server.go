package core

import (
	"chatplus/core/types"
	"chatplus/service/fun"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type AppServer struct {
	Debug        bool
	Config       *types.AppConfig
	Engine       *gin.Engine
	ChatContexts *types.LMap[string, []interface{}] // 聊天上下文 Map [chatId] => []Message

	ChatConfig *types.ChatConfig   // chat config cache
	SysConfig  *types.SystemConfig // system config cache

	// 保存 Websocket 会话 UserId, 每个 UserId 只能连接一次
	// 防止第三方直接连接 socket 调用 OpenAI API
	ChatSession   *types.LMap[string, *types.ChatSession] //map[sessionId]UserId
	ChatClients   *types.LMap[string, *types.WsClient]    // map[sessionId]Websocket 连接集合
	ReqCancelFunc *types.LMap[string, context.CancelFunc] // HttpClient 请求取消 handle function
	Functions     map[string]fun.Function
}

func NewServer(appConfig *types.AppConfig, functions map[string]fun.Function) *AppServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	return &AppServer{
		Debug:         false,
		Config:        appConfig,
		Engine:        gin.Default(),
		ChatContexts:  types.NewLMap[string, []interface{}](),
		ChatSession:   types.NewLMap[string, *types.ChatSession](),
		ChatClients:   types.NewLMap[string, *types.WsClient](),
		ReqCancelFunc: types.NewLMap[string, context.CancelFunc](),
		Functions:     functions,
	}
}

func (s *AppServer) Init(debug bool, client *redis.Client) {
	if debug { // 调试模式允许跨域请求 API
		s.Debug = debug
		logger.Info("Enabled debug mode")
	}
	s.Engine.Use(corsMiddleware())
	s.Engine.Use(authorizeMiddleware(s, client))
	s.Engine.Use(errorHandler)
	// 添加静态资源访问
	s.Engine.Static("/static", s.Config.StaticDir)
}

func (s *AppServer) Run(db *gorm.DB) error {
	// load chat config from database
	var chatConfig model.Config
	res := db.Where("marker", "chat").First(&chatConfig)
	if res.Error != nil {
		return res.Error
	}
	err := utils.JsonDecode(chatConfig.Config, &s.ChatConfig)
	if err != nil {
		return err
	}
	// load system configs
	var sysConfig model.Config
	res = db.Where("marker", "system").First(&sysConfig)
	if res.Error != nil {
		return res.Error
	}
	err = utils.JsonDecode(sysConfig.Config, &s.SysConfig)
	if err != nil {
		return err
	}
	logger.Infof("http://%s", s.Config.Listen)
	return s.Engine.Run(s.Config.Listen)
}

// 全局异常处理
func errorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Handler Panic: %v", r)
			debug.PrintStack()
			c.JSON(http.StatusOK, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
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
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, Chat-Token, Admin-Authorization")
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
func authorizeMiddleware(s *AppServer, client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/login" ||
			c.Request.URL.Path == "/api/admin/login" ||
			c.Request.URL.Path == "/api/user/register" ||
			c.Request.URL.Path == "/api/reward/notify" ||
			c.Request.URL.Path == "/api/mj/notify" ||
			c.Request.URL.Path == "/api/chat/history" ||
			c.Request.URL.Path == "/api/chat/detail" ||
			c.Request.URL.Path == "/api/role/list" ||
			c.Request.URL.Path == "/api/mj/jobs" ||
			c.Request.URL.Path == "/api/mj/proxy" ||
			c.Request.URL.Path == "/api/sd/jobs" ||
			strings.HasPrefix(c.Request.URL.Path, "/api/sms/") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/captcha/") ||
			strings.HasPrefix(c.Request.URL.Path, "/static/") ||
			c.Request.URL.Path == "/api/admin/config/get" {
			c.Next()
			return
		}

		var tokenString string
		if strings.Contains(c.Request.URL.Path, "/api/admin/") { // 后台管理 API
			tokenString = c.GetHeader(types.AdminAuthHeader)
		} else if c.Request.URL.Path == "/api/chat/new" ||
			c.Request.URL.Path == "/api/mj/client" ||
			c.Request.URL.Path == "/api/sd/client" {
			tokenString = c.Query("token")
		} else {
			tokenString = c.GetHeader(types.UserAuthHeader)
		}
		if tokenString == "" {
			resp.ERROR(c, "You should put Authorization in request headers")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(s.Config.Session.SecretKey), nil
		})

		if err != nil {
			resp.NotAuth(c, fmt.Sprintf("Error with parse auth token: %v", err))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			resp.NotAuth(c, "Token is invalid")
			c.Abort()
			return
		}

		expr := utils.IntValue(utils.InterfaceToString(claims["expired"]), 0)
		if expr > 0 && int64(expr) < time.Now().Unix() {
			resp.NotAuth(c, "Token is expired")
			c.Abort()
			return
		}

		key := fmt.Sprintf("users/%v", claims["user_id"])
		if _, err := client.Get(context.Background(), key).Result(); err != nil {
			resp.NotAuth(c, "Token is not found in redis")
			c.Abort()
			return
		}
		c.Set(types.LoginUserID, claims["user_id"])
	}
}
