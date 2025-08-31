package main

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"embed"
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/handler/admin"
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/service/dalle"
	"geekai/service/jimeng"
	"geekai/service/mj"
	"geekai/service/oss"
	"geekai/service/payment"
	"geekai/service/sd"
	"geekai/service/sms"
	"geekai/service/suno"
	"geekai/service/video"
	"geekai/store"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

//go:embed res
var xdbFS embed.FS

// AppLifecycle 应用程序生命周期
type AppLifecycle struct {
}

// OnStart 应用程序启动时执行
func (l *AppLifecycle) OnStart(context.Context) error {
	logger.Info("AppLifecycle OnStart")
	return nil
}

// OnStop 应用程序停止时执行
func (l *AppLifecycle) OnStop(context.Context) error {
	logger.Info("AppLifecycle OnStop")
	return nil
}

func NewAppLifeCycle() *AppLifecycle {
	return &AppLifecycle{}
}

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.toml"
	}
	logger.Info("Loading config file: ", configFile)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Panic Error:", err)
			// 打印堆栈信息
			if os.Getenv("GEEKAI_DEBUG") == "true" {
				debug.PrintStack()
			}
		}
	}()

	app := fx.New(
		// 初始化配置应用配置
		fx.Provide(func() *types.AppConfig {
			config, err := core.LoadConfig(configFile)
			if err != nil {
				log.Fatal(err)
			}
			config.Path = configFile
			return config
		}),
		// 创建应用服务
		fx.Provide(core.NewServer),
		// 初始化
		fx.Invoke(func(s *core.AppServer, client *redis.Client) {
			s.Init(client)
		}),
		fx.Provide(func(db *gorm.DB) *types.SystemConfig {
			return core.LoadSystemConfig(db)
		}),

		// 初始化数据库
		fx.Provide(store.NewGormConfig),
		fx.Provide(store.NewMysql),
		fx.Provide(store.NewRedisClient),
		fx.Provide(store.NewLevelDB),

		fx.Provide(func() embed.FS {
			return xdbFS
		}),

		// 创建 Ip2Region 查询对象
		fx.Provide(func() (*xdb.Searcher, error) {
			file, err := xdbFS.Open("res/ip2region.xdb")
			if err != nil {
				return nil, err
			}
			cBuff, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			return xdb.NewWithBuffer(cBuff)
		}),

		// 创建控制器
		fx.Provide(handler.NewChatAppHandler),
		fx.Provide(handler.NewUserHandler),
		fx.Provide(handler.NewChatHandler),
		fx.Provide(handler.NewNetHandler),
		fx.Provide(handler.NewSmsHandler),
		fx.Provide(handler.NewRedeemHandler),
		fx.Provide(handler.NewCaptchaHandler),
		fx.Provide(handler.NewMidJourneyHandler),
		fx.Provide(handler.NewChatModelHandler),
		fx.Provide(handler.NewSdJobHandler),
		fx.Provide(handler.NewPaymentHandler),
		fx.Provide(handler.NewOrderHandler),
		fx.Provide(handler.NewProductHandler),
		fx.Provide(handler.NewConfigHandler),
		fx.Provide(handler.NewPowerLogHandler),
		fx.Provide(handler.NewJimengHandler),

		fx.Provide(service.NewMigrationService),
		fx.Invoke(func(migrationService *service.MigrationService) {
			migrationService.StartMigrate()
		}),

		// 管理后台控制器
		fx.Provide(admin.NewConfigHandler),
		fx.Provide(admin.NewAdminHandler),
		fx.Provide(admin.NewApiKeyHandler),
		fx.Provide(admin.NewUserHandler),
		fx.Provide(admin.NewChatAppHandler),
		fx.Provide(admin.NewRedeemHandler),
		fx.Provide(admin.NewDashboardHandler),
		fx.Provide(admin.NewChatModelHandler),
		fx.Provide(admin.NewProductHandler),
		fx.Provide(admin.NewOrderHandler),
		fx.Provide(admin.NewPowerLogHandler),
		fx.Provide(admin.NewAdminJimengHandler),

		// 邮件服务
		fx.Provide(service.NewSmtpService),
		// License 服务
		fx.Provide(service.NewLicenseService),
		fx.Invoke(func(licenseService *service.LicenseService) {
			licenseService.SyncLicense()
		}),

		// Dalle 服务
		fx.Provide(dalle.NewService),
		fx.Invoke(func(s *dalle.Service) {
			s.Run()
			s.DownloadImages()
			s.CheckTaskStatus()
		}),

		// MidJourney service pool
		fx.Provide(mj.NewService),
		fx.Provide(mj.NewClient),
		fx.Invoke(func(s *mj.Service) {
			s.Run()
			s.SyncTaskProgress()
			s.DownloadImages()
		}),

		// Stable Diffusion 机器人
		fx.Provide(sd.NewService),
		fx.Invoke(func(s *sd.Service, config *types.AppConfig) {
			s.Run()
			s.CheckTaskStatus()
		}),

		fx.Provide(suno.NewService),
		fx.Invoke(func(s *suno.Service) {
			s.Run()
			s.SyncTaskProgress()
			s.DownloadFiles()
		}),
		fx.Provide(video.NewService),
		fx.Invoke(func(s *video.Service) {
			s.Run()
			s.SyncTaskProgress()
			s.DownloadFiles()
		}),

		// 即梦AI 服务
		fx.Provide(jimeng.NewService),
		fx.Invoke(func(service *jimeng.Service) {
			service.Start()
		}),
		fx.Provide(service.NewSnowflake),

		// 创建短信服务
		fx.Provide(sms.NewAliYunSmsService),
		fx.Provide(sms.NewBaoSmsService),
		fx.Provide(sms.NewSmsManager),
		fx.Provide(func(config *types.SystemConfig) *service.CaptchaService {
			return service.NewCaptchaService(config.Captcha)
		}),
		fx.Provide(func(config *types.SystemConfig, client *redis.Client) *service.WxLoginService {
			return service.NewWxLoginService(config.WxLogin, client)
		}),

		// 支付服务
		fx.Provide(payment.NewAlipayService),
		fx.Provide(payment.NewEPayService),
		fx.Provide(payment.NewWxpayService),

		// 文件上传服务
		fx.Provide(oss.NewLocalStorage),
		fx.Provide(oss.NewMiniOss),
		fx.Provide(oss.NewQiNiuOss),
		fx.Provide(oss.NewAliYunOss),
		fx.Provide(oss.NewUploaderManager),

		// 用户服务
		fx.Provide(service.NewUserService),

		// 注册路由
		fx.Invoke(func(s *core.AppServer, h *handler.ChatAppHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.UserHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ChatHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.NetHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.SmsHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.CaptchaHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.RedeemHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.MidJourneyHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.SdJobHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ConfigHandler) {
			h.RegisterRoutes()
		}),

		// 管理后台路由注册
		fx.Invoke(func(s *core.AppServer, h *admin.ConfigHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ManagerHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ApiKeyHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.UserHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ChatAppHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.RedeemHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.DashboardHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ChatModelHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ChatModelHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.PaymentHandler) {
			h.RegisterRoutes()
			h.StartSyncOrders()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ProductHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.OrderHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.OrderHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ProductHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewInviteHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.InviteHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewFunctionHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.FunctionHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewUploadHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.UploadHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewFunctionHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.FunctionHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewChatHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.ChatHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.PowerLogHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.PowerLogHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewMenuHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.MenuHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewMenuHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.MenuHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewMarkMapHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.MarkMapHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewDallJobHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.DallJobHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewSunoHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.SunoHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewVideoHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.VideoHandler) {
			h.RegisterRoutes()
		}),

		// 即梦AI 路由
		fx.Invoke(func(s *core.AppServer, h *handler.JimengHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.AdminJimengHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewChatAppTypeHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.ChatAppTypeHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewChatAppTypeHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.ChatAppTypeHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewTestHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.TestHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewPromptHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.PromptHandler) {
			h.RegisterRoutes()
		}),
		fx.Invoke(func(s *core.AppServer, db *gorm.DB) {
			go func() {
				err := s.Run(db)
				if err != nil {
					logger.Error(err)
					os.Exit(0)
				}
			}()
		}),
		fx.Provide(NewAppLifeCycle),
		// 注册生命周期回调函数
		fx.Invoke(func(lifecycle fx.Lifecycle, lc *AppLifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return lc.OnStart(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return lc.OnStop(ctx)
				},
			})
		}),
		fx.Provide(admin.NewImageHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.ImageHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewMediaHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.MediaHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewRealtimeHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.RealtimeHandler) {
			h.RegisterRoutes()
		}),
	)
	// 启动应用程序
	go func() {
		if err := app.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭应用程序
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Fatal(err)
	}

}
