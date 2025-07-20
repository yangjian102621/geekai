package core

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"context"
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/imroc/req/v3"
	"github.com/nfnt/resize"
	"github.com/shirou/gopsutil/host"
	"golang.org/x/image/webp"
	"gorm.io/gorm"
)

// AuthConfig 定义授权配置
type AuthConfig struct {
	ExactPaths  map[string]bool // 精确匹配的路径
	PrefixPaths map[string]bool // 前缀匹配的路径
}

var authConfig = &AuthConfig{
	ExactPaths: map[string]bool{
		"/api/user/login":          false,
		"/api/user/logout":         false,
		"/api/user/resetPass":      false,
		"/api/user/register":       false,
		"/api/admin/login":         false,
		"/api/admin/logout":        false,
		"/api/admin/login/captcha": false,
		"/api/chat/history":        false,
		"/api/chat/detail":         false,
		"/api/chat/list":           false,
		"/api/app/list":            false,
		"/api/app/type/list":       false,
		"/api/app/list/user":       false,
		"/api/model/list":          false,
		"/api/mj/imgWall":          false,
		"/api/mj/notify":           false,
		"/api/invite/hits":         false,
		"/api/sd/imgWall":          false,
		"/api/dall/imgWall":        false,
		"/api/product/list":        false,
		"/api/menu/list":           false,
		"/api/markMap/client":      false,
		"/api/payment/doPay":       false,
		"/api/payment/payWays":     false,
		"/api/suno/detail":         false,
		"/api/suno/play":           false,
		"/api/download":            false,
		"/api/dall/models":         false,
	},
	PrefixPaths: map[string]bool{
		"/api/test/":           false,
		"/api/payment/notify/": false,
		"/api/user/clogin":     false,
		"/api/config/":         false,
		"/api/function/":       false,
		"/api/sms/":            false,
		"/api/captcha/":        false,
		"/static/":             false,
	},
}

type AppServer struct {
	Config    *types.AppConfig
	Engine    *gin.Engine
	SysConfig *types.SystemConfig // system config cache
}

func NewServer(appConfig *types.AppConfig) *AppServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	return &AppServer{
		Config: appConfig,
		Engine: gin.Default(),
	}
}

func (s *AppServer) Init(debug bool, client *redis.Client) {
	s.Engine.Use(corsMiddleware())
	s.Engine.Use(staticResourceMiddleware())
	s.Engine.Use(authorizeMiddleware(s, client))
	s.Engine.Use(parameterHandlerMiddleware())
	s.Engine.Use(errorHandler)
	// 添加静态资源访问
	s.Engine.Static("/static", s.Config.StaticDir)
}

func (s *AppServer) Run(db *gorm.DB) error {

	// 重命名 config 表字段
	if db.Migrator().HasColumn(&model.Config{}, "config_json") {
		db.Migrator().RenameColumn(&model.Config{}, "config_json", "value")
	}
	if db.Migrator().HasColumn(&model.Config{}, "marker") {
		db.Migrator().RenameColumn(&model.Config{}, "marker", "name")
	}
	if db.Migrator().HasIndex(&model.Config{}, "idx_chatgpt_configs_key") {
		db.Migrator().DropIndex(&model.Config{}, "idx_chatgpt_configs_key")
	}
	if db.Migrator().HasIndex(&model.Config{}, "marker") {
		db.Migrator().DropIndex(&model.Config{}, "marker")
	}

	// load system configs
	var sysConfig model.Config
	err := db.Where("name", "system").First(&sysConfig).Error
	if err != nil {
		return fmt.Errorf("failed to load system config: %v", err)
	}
	err = utils.JsonDecode(sysConfig.Value, &s.SysConfig)
	if err != nil {
		return fmt.Errorf("failed to decode system config: %v", err)
	}

	// 迁移数据表
	logger.Info("Migrating database tables...")
	db.AutoMigrate(
		&model.ChatItem{},
		&model.ChatMessage{},
		&model.ChatRole{},
		&model.ChatModel{},
		&model.InviteCode{},
		&model.InviteLog{},
		&model.Menu{},
		&model.Order{},
		&model.Product{},
		&model.User{},
		&model.Function{},
		&model.File{},
		&model.Redeem{},
		&model.Config{},
		&model.ApiKey{},
		&model.AdminUser{},
		&model.AppType{},
		&model.SdJob{},
		&model.SunoJob{},
		&model.PowerLog{},
		&model.VideoJob{},
		&model.MidJourneyJob{},
		&model.UserLoginLog{},
		&model.DallJob{},
		&model.JimengJob{},
	)
	// 手动删除字段
	if db.Migrator().HasColumn(&model.Order{}, "deleted_at") {
		db.Migrator().DropColumn(&model.Order{}, "deleted_at")
	}
	if db.Migrator().HasColumn(&model.ChatItem{}, "deleted_at") {
		db.Migrator().DropColumn(&model.ChatItem{}, "deleted_at")
	}
	if db.Migrator().HasColumn(&model.ChatMessage{}, "deleted_at") {
		db.Migrator().DropColumn(&model.ChatMessage{}, "deleted_at")
	}
	if db.Migrator().HasColumn(&model.User{}, "chat_config") {
		db.Migrator().DropColumn(&model.User{}, "chat_config")
	}
	if db.Migrator().HasColumn(&model.ChatModel{}, "category") {
		db.Migrator().DropColumn(&model.ChatModel{}, "category")
	}
	if db.Migrator().HasColumn(&model.ChatModel{}, "description") {
		db.Migrator().DropColumn(&model.ChatModel{}, "description")
	}

	logger.Info("Database tables migrated successfully")

	// 统计安装信息
	go func() {
		info, err := host.Info()
		if err == nil {
			apiURL := fmt.Sprintf("%s/%s", s.Config.ApiConfig.ApiURL, "api/installs/push")
			timestamp := time.Now().Unix()
			product := "geekai-plus"
			signStr := fmt.Sprintf("%s#%s#%d", product, info.HostID, timestamp)
			sign := utils.Sha256(signStr)
			resp, err := req.C().R().SetBody(map[string]interface{}{"product": product, "device_id": info.HostID, "timestamp": timestamp, "sign": sign}).Post(apiURL)
			if err != nil {
				logger.Errorf("register install info failed: %v", err)
			} else {
				logger.Debugf("register install info success: %v", resp.String())
			}
		}
	}()
	logger.Infof("http://%s", s.Config.Listen)
	return s.Engine.Run(s.Config.Listen)
}

// 全局异常处理
func errorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Handler Panic: %v", r)
			debug.PrintStack()
			c.JSON(http.StatusBadRequest, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
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

		// 设置允许的请求源
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "Authorization, Body-Length, Body-Type, Admin-Authorization,content-type")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Body-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")

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
		if !needLogin(c) {
			c.Next()
			return
		}

		clientProtocols := c.GetHeader("Sec-WebSocket-Protocol")
		var tokenString string
		isAdminApi := strings.Contains(c.Request.URL.Path, "/api/admin/")
		if isAdminApi { // 后台管理 API
			tokenString = c.GetHeader(types.AdminAuthHeader)
		} else if clientProtocols != "" { // Websocket 连接
			// 解析子协议内容
			protocols := strings.Split(clientProtocols, ",")
			if protocols[0] == "realtime" {
				tokenString = strings.TrimSpace(protocols[1][25:])
			} else if protocols[0] == "token" {
				tokenString = strings.TrimSpace(protocols[1])
			}
		} else {
			tokenString = c.GetHeader(types.UserAuthHeader)
		}

		if tokenString == "" {
			resp.NotAuth(c, "You should put Authorization in request headers")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			if isAdminApi {
				return []byte(s.Config.AdminSession.SecretKey), nil
			} else {
				return []byte(s.Config.Session.SecretKey), nil
			}

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
		if isAdminApi {
			key = fmt.Sprintf("admin/%v", claims["user_id"])
		}
		if _, err := client.Get(context.Background(), key).Result(); err != nil {
			resp.NotAuth(c, "Token is not found in redis")
			c.Abort()
			return
		}
		c.Set(types.LoginUserID, claims["user_id"])
		c.Next()
	}
}

func needLogin(c *gin.Context) bool {
	path := c.Request.URL.Path

	// 如果不是 API 路径，不需要登录
	if !strings.HasPrefix(path, "/api") {
		return false
	}

	// 检查精确匹配的路径
	if skip, exists := authConfig.ExactPaths[path]; exists {
		return skip
	}

	// 检查前缀匹配的路径
	for prefix, skip := range authConfig.PrefixPaths {
		if strings.HasPrefix(path, prefix) {
			return skip
		}
	}

	return true
}

// 跳过授权
func (s *AppServer) SkipAuth(url string, prefix bool) {
	if prefix {
		authConfig.PrefixPaths[url] = false
	} else {
		authConfig.ExactPaths[url] = false
	}
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
			// for .webp image
			if err != nil {
				img, err = webp.Decode(file)
			}
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
				logger.Error(err)
				c.String(http.StatusInternalServerError, err.Error())
				return
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
