package core

import (
	"bytes"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type AppServer struct {
	Debug        bool
	Config       *types.AppConfig
	Engine       *gin.Engine
	ChatContexts *types.LMap[string, []types.Message] // 聊天上下文 Map [chatId] => []Message

	SysConfig *types.SystemConfig // system config cache

	// 保存 Websocket 会话 UserId, 每个 UserId 只能连接一次
	// 防止第三方直接连接 socket 调用 OpenAI API
	ChatSession   *types.LMap[string, *types.ChatSession] //map[sessionId]UserId
	ChatClients   *types.LMap[string, *types.WsClient]    // map[sessionId]Websocket 连接集合
	ReqCancelFunc *types.LMap[string, context.CancelFunc] // HttpClient 请求取消 handle function
}

func NewServer(appConfig *types.AppConfig) *AppServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	return &AppServer{
		Debug:         false,
		Config:        appConfig,
		Engine:        gin.Default(),
		ChatContexts:  types.NewLMap[string, []types.Message](),
		ChatSession:   types.NewLMap[string, *types.ChatSession](),
		ChatClients:   types.NewLMap[string, *types.WsClient](),
		ReqCancelFunc: types.NewLMap[string, context.CancelFunc](),
	}
}

func (s *AppServer) Init(debug bool, client *redis.Client) {
	if debug { // 调试模式允许跨域请求 API
		s.Debug = debug
		logger.Info("Enabled debug mode")
	}
	s.Engine.Use(corsMiddleware())
	s.Engine.Use(staticResourceMiddleware())
	s.Engine.Use(authorizeMiddleware(s, client))
	s.Engine.Use(parameterHandlerMiddleware())
	s.Engine.Use(errorHandler)
	// 添加静态资源访问
	s.Engine.Static("/static", s.Config.StaticDir)
}

func (s *AppServer) Run(db *gorm.DB) error {
	// load system configs
	var sysConfig model.Config
	res := db.Where("marker", "system").First(&sysConfig)
	if res.Error != nil {
		return res.Error
	}
	err := utils.JsonDecode(sysConfig.Config, &s.SysConfig)
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
		var tokenString string
		isAdminApi := strings.Contains(c.Request.URL.Path, "/api/admin/")
		if isAdminApi { // 后台管理 API
			tokenString = c.GetHeader(types.AdminAuthHeader)
		} else if c.Request.URL.Path == "/api/chat/new" {
			tokenString = c.Query("token")
		} else {
			tokenString = c.GetHeader(types.UserAuthHeader)
		}

		if tokenString == "" {
			if needLogin(c) {
				resp.ERROR(c, "You should put Authorization in request headers")
				c.Abort()
				return
			} else { // 直接放行
				c.Next()
				return
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok && needLogin(c) {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			if isAdminApi {
				return []byte(s.Config.AdminSession.SecretKey), nil
			} else {
				return []byte(s.Config.Session.SecretKey), nil
			}

		})

		if err != nil && needLogin(c) {
			resp.NotAuth(c, fmt.Sprintf("Error with parse auth token: %v", err))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid && needLogin(c) {
			resp.NotAuth(c, "Token is invalid")
			c.Abort()
			return
		}

		expr := utils.IntValue(utils.InterfaceToString(claims["expired"]), 0)
		if expr > 0 && int64(expr) < time.Now().Unix() && needLogin(c) {
			resp.NotAuth(c, "Token is expired")
			c.Abort()
			return
		}

		key := fmt.Sprintf("users/%v", claims["user_id"])
		if isAdminApi {
			key = fmt.Sprintf("admin/%v", claims["user_id"])
		}
		if _, err := client.Get(context.Background(), key).Result(); err != nil && needLogin(c) {
			resp.NotAuth(c, "Token is not found in redis")
			c.Abort()
			return
		}
		c.Set(types.LoginUserID, claims["user_id"])
	}
}

func needLogin(c *gin.Context) bool {
	if c.Request.URL.Path == "/api/user/login" ||
		c.Request.URL.Path == "/api/user/resetPass" ||
		c.Request.URL.Path == "/api/admin/login" ||
		c.Request.URL.Path == "/api/admin/login/captcha" ||
		c.Request.URL.Path == "/api/user/register" ||
		c.Request.URL.Path == "/api/chat/history" ||
		c.Request.URL.Path == "/api/chat/detail" ||
		c.Request.URL.Path == "/api/chat/list" ||
		c.Request.URL.Path == "/api/role/list" ||
		c.Request.URL.Path == "/api/model/list" ||
		c.Request.URL.Path == "/api/mj/imgWall" ||
		c.Request.URL.Path == "/api/mj/client" ||
		c.Request.URL.Path == "/api/mj/notify" ||
		c.Request.URL.Path == "/api/invite/hits" ||
		c.Request.URL.Path == "/api/sd/imgWall" ||
		c.Request.URL.Path == "/api/sd/client" ||
		c.Request.URL.Path == "/api/config/get" ||
		c.Request.URL.Path == "/api/product/list" ||
		c.Request.URL.Path == "/api/menu/list" ||
		strings.HasPrefix(c.Request.URL.Path, "/api/test") ||
		strings.HasPrefix(c.Request.URL.Path, "/api/function/") ||
		strings.HasPrefix(c.Request.URL.Path, "/api/sms/") ||
		strings.HasPrefix(c.Request.URL.Path, "/api/captcha/") ||
		strings.HasPrefix(c.Request.URL.Path, "/api/payment/") ||
		strings.HasPrefix(c.Request.URL.Path, "/static/") {
		return false
	}
	return true
}

// 统一参数处理
func parameterHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET 参数处理
		params := c.Request.URL.Query()
		for key, values := range params {
			for i, value := range values {
				params[key][i] = strings.TrimSpace(value)
			}
		}
		// update get parameters
		c.Request.URL.RawQuery = params.Encode()
		// skip file upload requests
		contentType := c.Request.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			c.Next()
			return
		}

		if strings.Contains(contentType, "application/json") {
			// process POST JSON request body
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.Next()
				return
			}

			// 还原请求体
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// 将请求体解析为 JSON
			var jsonData map[string]interface{}
			if err := c.ShouldBindJSON(&jsonData); err != nil {
				c.Next()
				return
			}

			// 对 JSON 数据中的字符串值去除两端空格
			trimJSONStrings(jsonData)
			// 更新请求体
			c.Request.Body = io.NopCloser(bytes.NewBufferString(utils.JsonEncode(jsonData)))
		}

		c.Next()
	}
}

// 递归对 JSON 数据中的字符串值去除两端空格
func trimJSONStrings(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			switch valueType := value.(type) {
			case string:
				v[key] = strings.TrimSpace(valueType)
			case map[string]interface{}, []interface{}:
				trimJSONStrings(value)
			}
		}
	case []interface{}:
		for i, value := range v {
			switch valueType := value.(type) {
			case string:
				v[i] = strings.TrimSpace(valueType)
			case map[string]interface{}, []interface{}:
				trimJSONStrings(value)
			}
		}
	}
}

// 静态资源中间件
func staticResourceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		url := c.Request.URL.String()
		// 拦截生成缩略图请求
		if strings.HasPrefix(url, "/static/") && strings.Contains(url, "?imageView2") {
			r := strings.SplitAfter(url, "imageView2")
			size := strings.Split(r[1], "/")
			if len(size) != 8 {
				c.String(http.StatusNotFound, "invalid thumb args")
				return
			}
			with := utils.IntValue(size[3], 0)
			height := utils.IntValue(size[5], 0)
			quality := utils.IntValue(size[7], 75)

			// 打开图片文件
			filePath := strings.TrimLeft(c.Request.URL.Path, "/")
			file, err := os.Open(filePath)
			if err != nil {
				c.String(http.StatusNotFound, "Image not found")
				return
			}
			defer file.Close()

			// 解码图片
			img, _, err := image.Decode(file)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error decoding image")
				return
			}

			var newImg image.Image
			if height == 0 || with == 0 {
				// 固定宽度，高度自适应
				newImg = resize.Resize(uint(with), uint(height), img, resize.Lanczos3)
			} else {
				// 生成缩略图
				newImg = resize.Thumbnail(uint(with), uint(height), img, resize.Lanczos3)
			}
			var buffer bytes.Buffer
			err = jpeg.Encode(&buffer, newImg, &jpeg.Options{Quality: quality})
			if err != nil {
				log.Fatal(err)
			}

			// 设置图片缓存有效期为一年 (365天)
			c.Header("Cache-Control", "max-age=31536000, public")
			// 直接输出图像数据流
			c.Data(http.StatusOK, "image/jpeg", buffer.Bytes())
			c.Abort() // 中断请求
		}
		c.Next()
	}
}
