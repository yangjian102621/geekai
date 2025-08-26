package core

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/store/model"
	"geekai/utils"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"github.com/shirou/gopsutil/host"
	"gorm.io/gorm"
)

// AuthConfig 定义授权配置
type AuthConfig struct {
	ExactPaths  map[string]bool // 精确匹配的路径
	PrefixPaths map[string]bool // 前缀匹配的路径
}

var authConfig = &AuthConfig{
	ExactPaths: map[string]bool{
		"/api/user/login":           false,
		"/api/user/logout":          false,
		"/api/user/resetPass":       false,
		"/api/user/register":        false,
		"/api/user/clogin":          false,
		"/api/user/clogin/callback": false,
		"/api/user/signin":          false,
		"/api/admin/login":          false,
		"/api/admin/logout":         false,
		"/api/admin/login/captcha":  false,
		"/api/app/list":             false,
		"/api/app/type/list":        false,
		"/api/app/list/user":        false,
		"/api/model/list":           false,
		"/api/mj/imgWall":           false,
		"/api/mj/notify":            false,
		"/api/invite/hits":          false,
		"/api/sd/imgWall":           false,
		"/api/dall/imgWall":         false,
		"/api/product/list":         false,
		"/api/menu/list":            false,
		"/api/markMap/client":       false,
		"/api/payment/doPay":        false,
		"/api/payment/payWays":      false,
		"/api/download":             false,
		"/api/dall/models":          false,
		"/api/chat/message":         false, // 聊天接口需要特殊处理
		"/api/realtime":             false, // 实时通信接口需要特殊处理
		"/api/realtime/voice":       false, // 语音聊天接口需要特殊处理
	},
	PrefixPaths: map[string]bool{
		"/api/test/":           false,
		"/api/payment/notify/": false,
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
	Redis     *redis.Client
}

func NewServer(appConfig *types.AppConfig, redis *redis.Client, sysConfig *types.SystemConfig) *AppServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	return &AppServer{
		Config:    appConfig,
		Redis:     redis,
		Engine:    gin.Default(),
		SysConfig: sysConfig,
	}
}

func (s *AppServer) Init(client *redis.Client) {
	s.Engine.Use(middleware.ParameterHandlerMiddleware())
	s.Engine.Use(middleware.StaticMiddleware())
	s.Engine.Use(errorHandler)
	// 添加静态资源访问
	s.Engine.Static("/static", s.Config.StaticDir)
}

func (s *AppServer) Run(db *gorm.DB) error {
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

	// 统计安装信息
	go func() {
		info, err := host.Info()
		if err == nil {
			apiURL := fmt.Sprintf("%s/api/installs/push", types.GeekAPIURL)
			timestamp := time.Now().Unix()
			product := "geekai-plus"
			signStr := fmt.Sprintf("%s#%s#%d", product, info.HostID, timestamp)
			sign := utils.Sha256(signStr)
			resp, err := req.C().R().SetBody(map[string]interface{}{"product": product, "device_id": info.HostID, "timestamp": timestamp, "sign": sign}).Post(apiURL)
			if err == nil {
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
