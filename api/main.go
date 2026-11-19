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
	"geekai/log"
	"geekai/service"
	"geekai/service/image"
	"geekai/service/jimeng"
	"geekai/service/library"
	"geekai/service/mj"
	"geekai/service/moderation"
	"geekai/service/oss"
	"geekai/service/payment"
	"geekai/service/sms"
	"geekai/service/sora"
	"geekai/service/suno"
	"geekai/service/video"
	"geekai/store"
	"io"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var logger = log.GetLogger()

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

// Handlers 汇总全部 52 个路由注册器。
//
// 为什么要显式列出而不是靠 fx 的 group / As：
//  1. fx.As 会替换构造函数的返回类型，*handler.PaymentHandler 这类"既要注册路由、又要单独
//     启动同步任务"的具体类型会失踪；
//  2. fx 的 group 不保证顺序，gin 的路由注册顺序会变成不确定的。
//
// 显式字段 + 显式切片：DI 语义不变、顺序确定、注册表一目了然。
type Handlers struct {
	fx.In

	// 用户侧（30 个）
	ChatApp     *handler.ChatAppHandler
	User        *handler.UserHandler
	Chat        *handler.ChatHandler
	Net         *handler.NetHandler
	Download    *handler.DownloadHandler
	Sms         *handler.SmsHandler
	Redeem      *handler.RedeemHandler
	Captcha     *handler.CaptchaHandler
	Mj          *handler.MidJourneyHandler
	ChatModel   *handler.ChatModelHandler
	ModelAPI    *handler.ModelAPIHandler
	Payment     *handler.PaymentHandler
	Order       *handler.OrderHandler
	Product     *handler.ProductHandler
	Config      *handler.ConfigHandler
	PowerLog    *handler.PowerLogHandler
	Jimeng      *handler.JimengHandler
	WxGzh       *handler.WxGzhHandler
	Invite      *handler.InviteHandler
	Function    *handler.FunctionHandler
	Menu        *handler.MenuHandler
	MarkMap     *handler.MarkMapHandler
	ImageJob    *handler.ImageJobHandler
	Suno        *handler.SunoHandler
	Video       *handler.VideoHandler
	ChatAppType *handler.ChatAppTypeHandler
	Test        *handler.TestHandler
	Prompt      *handler.PromptHandler
	Realtime    *handler.RealtimeHandler
	Library     *handler.LibraryHandler

	// 管理后台（22 个）
	AdminConfig      *admin.ConfigHandler
	AdminManager     *admin.ManagerHandler
	AdminApiKey      *admin.ApiKeyHandler
	AdminUser        *admin.UserHandler
	AdminChatApp     *admin.ChatAppHandler
	AdminRedeem      *admin.RedeemHandler
	AdminDashboard   *admin.DashboardHandler
	AdminChatModel   *admin.ChatModelHandler
	AdminProduct     *admin.ProductHandler
	AdminOrder       *admin.OrderHandler
	AdminPowerLog    *admin.PowerLogHandler
	AdminJimeng      *admin.AdminJimengHandler
	AdminVideo       *admin.VideoHandler
	AdminSuno        *admin.SunoHandler
	AdminChatAppType *admin.ChatAppTypeHandler
	AdminFunction    *admin.FunctionHandler
	AdminUpload      *admin.UploadHandler
	AdminChat        *admin.ChatHandler
	AdminMenu        *admin.MenuHandler
	AdminImage       *admin.ImageHandler
	AdminModeration  *admin.ModerationHandler
	AdminLibrary     *admin.LibraryHandler
}

// wantRegistrars 是"必须被注册的路由注册器总数"，同时等价于 Handlers 的字段数。
// 任何"加了字段忘了加切片"或"加了 handler 忘了加字段"都会在启动期炸掉，而不是静默不生效。
const wantRegistrars = 52

// registerRoutes 是 main.go 里唯一的路由注册点。
// 新增 handler 必须同时在 Handlers 结构体里加字段、在下面的切片里加一项。
func registerRoutes(r *core.Router, hs Handlers) {
	registrars := []core.RouteRegistrar{
		// 用户侧
		hs.User, hs.ChatApp, hs.Chat, hs.Net, hs.Download, hs.Sms, hs.Redeem,
		hs.Captcha, hs.Mj, hs.ChatModel, hs.ModelAPI, hs.Payment, hs.Order,
		hs.Product, hs.Config, hs.PowerLog, hs.Jimeng, hs.WxGzh, hs.Invite,
		hs.Function, hs.Menu, hs.MarkMap, hs.ImageJob, hs.Suno, hs.Video,
		hs.ChatAppType, hs.Test, hs.Prompt, hs.Realtime, hs.Library,
		// 管理后台
		hs.AdminConfig, hs.AdminManager, hs.AdminApiKey, hs.AdminUser,
		hs.AdminChatApp, hs.AdminRedeem, hs.AdminDashboard, hs.AdminChatModel,
		hs.AdminProduct, hs.AdminOrder, hs.AdminPowerLog, hs.AdminJimeng,
		hs.AdminVideo, hs.AdminSuno, hs.AdminChatAppType, hs.AdminFunction,
		hs.AdminUpload, hs.AdminChat, hs.AdminMenu, hs.AdminImage,
		hs.AdminModeration, hs.AdminLibrary,
	}
	if got := reflect.ValueOf(hs).NumField(); got != wantRegistrars {
		logger.Fatalf("Handlers 字段数 %d != %d，有 handler 忘记加进 Handlers 结构体", got, wantRegistrars)
	}
	if len(registrars) != wantRegistrars {
		logger.Fatalf("路由注册器数量 %d != %d，有 handler 忘记加进 registerRoutes 切片", len(registrars), wantRegistrars)
	}
	for _, reg := range registrars {
		reg.RegisterRoutes(r)
	}
	logger.Infof("registered %d route registrars", len(registrars))
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
				logger.Fatal(err)
			}
			config.Path = configFile
			return config
		}),
		// 创建应用服务（构造即完备：Engine 中间件、静态资源与 Router 都在 NewServer 内部装配完毕）
		fx.Provide(core.NewServer),
		// 对外唯一的路由注册入口：handler 只能拿到 *core.Router，拿不到 *gin.Engine
		fx.Provide(func(s *core.AppServer) *core.Router {
			return s.Router
		}),
		// 唯一的路由注册点：52 个 handler 一处登记、一处装配，顺序确定
		fx.Invoke(registerRoutes),
		fx.Provide(func(db *gorm.DB) *types.SystemConfig {
			return core.LoadSystemConfig(db)
		}),

		// 初始化数据库
		fx.Provide(store.NewGormConfig),
		fx.Provide(store.NewMysql),
		fx.Provide(store.NewRedisClient),
		fx.Provide(store.NewLevelDB),
		fx.Provide(func(db *gorm.DB) *library.Service {
			return library.NewService(library.NewRepository(db))
		}),
		fx.Provide(handler.NewLibraryHandler),
		fx.Provide(admin.NewLibraryHandler),

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
		fx.Provide(handler.NewDownloadHandler),
		fx.Provide(handler.NewSmsHandler),
		fx.Provide(handler.NewRedeemHandler),
		fx.Provide(handler.NewCaptchaHandler),
		fx.Provide(handler.NewMidJourneyHandler),
		fx.Provide(handler.NewChatModelHandler),
		fx.Provide(handler.NewModelAPIHandler),
		fx.Provide(handler.NewPaymentHandler),
		fx.Provide(handler.NewOrderHandler),
		fx.Provide(handler.NewProductHandler),
		fx.Provide(handler.NewConfigHandler),
		fx.Provide(handler.NewPowerLogHandler),
		fx.Provide(handler.NewJimengHandler),
		fx.Provide(service.NewWxGzhService),
		fx.Provide(handler.NewWxGzhHandler),
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
		fx.Provide(admin.NewVideoHandler),
		fx.Provide(admin.NewSunoHandler),

		// 邮件服务
		fx.Provide(service.NewSmtpService),
		// Image 服务
		fx.Provide(image.NewService),
		fx.Invoke(func(s *image.Service) {
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

		// Sora service
		fx.Provide(sora.NewSoraService),

		fx.Provide(suno.NewService),
		fx.Invoke(func(s *suno.Service) {
			s.Run()
			// s.SyncTaskProgress()
			// s.DownloadFiles()
		}),
		fx.Provide(video.NewService),
		fx.Invoke(func(s *video.Service) {
			s.Run()
			s.SyncTaskProgress()
			s.DownloadFiles()
		}),

		// 即梦AI 服务
		fx.Provide(jimeng.NewClient),
		fx.Provide(jimeng.NewService),
		fx.Invoke(func(service *jimeng.Service) {
			service.Start()
		}),

		fx.Provide(service.NewSnowflake),

		// 创建短信服务
		fx.Provide(sms.NewAliYunSmsService),
		fx.Provide(sms.NewBaoSmsService),
		fx.Provide(sms.NewTencentSmsService),
		fx.Provide(sms.NewSmsManager),
		fx.Provide(service.NewCaptchaService),
		fx.Provide(service.NewWxLoginService),

		// 支付服务
		fx.Provide(payment.NewAlipayService),
		fx.Provide(payment.NewEPayService),
		fx.Provide(payment.NewWxpayService),
		fx.Provide(payment.NewStripeService),

		// 文件上传服务
		fx.Provide(oss.NewLocalStorage),
		fx.Provide(oss.NewMiniOss),
		fx.Provide(oss.NewQiNiuOss),
		fx.Provide(oss.NewAliYunOss),
		fx.Provide(oss.NewTencentOss),
		fx.Provide(oss.NewUploaderManager),

		// 用户服务
		fx.Provide(service.NewUserService),
		// 积分账本服务：所有改变余额的操作都应经过它（事务 + 幂等键）
		fx.Provide(service.NewLedgerService),

		// 文本审查服务
		fx.Provide(moderation.NewGiteeAIModeration),
		fx.Provide(moderation.NewBaiduAIModeration),
		fx.Provide(moderation.NewTencentAIModeration),
		fx.Provide(moderation.NewServiceManager),
		fx.Provide(admin.NewModerationHandler),

		// 注册路由

		// 管理后台路由注册
		fx.Invoke(func(h *handler.PaymentHandler) {
			h.StartSyncOrders()
		}),

		fx.Provide(handler.NewInviteHandler),

		fx.Provide(admin.NewFunctionHandler),

		fx.Provide(admin.NewUploadHandler),

		fx.Provide(handler.NewFunctionHandler),
		fx.Provide(admin.NewChatHandler),
		fx.Provide(admin.NewMenuHandler),
		fx.Provide(handler.NewMenuHandler),
		fx.Provide(handler.NewMarkMapHandler),
		fx.Provide(handler.NewImageJobHandler),
		fx.Provide(handler.NewSunoHandler),
		fx.Provide(handler.NewVideoHandler),

		// 即梦AI 路由
		fx.Provide(admin.NewChatAppTypeHandler),
		fx.Provide(handler.NewChatAppTypeHandler),
		fx.Provide(handler.NewTestHandler),
		fx.Provide(handler.NewPromptHandler),

		// 微信公众号路由
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
		fx.Provide(handler.NewRealtimeHandler),
	)
	// 启动应用程序
	go func() {
		if err := app.Start(context.Background()); err != nil {
			logger.Fatal(err)
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
		logger.Fatal(err)
	}

}
