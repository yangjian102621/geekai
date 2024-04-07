package main

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/handler/admin"
	"chatplus/handler/chatimpl"
	logger2 "chatplus/logger"
	"chatplus/service"
	"chatplus/service/mj"
	"chatplus/service/oss"
	"chatplus/service/payment"
	"chatplus/service/sd"
	"chatplus/service/sms"
	"chatplus/service/wx"
	"chatplus/store"
	"context"
	"embed"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
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
	log.Println("AppLifecycle OnStart")
	return nil
}

// OnStop 应用程序停止时执行
func (l *AppLifecycle) OnStop(context.Context) error {
	log.Println("AppLifecycle OnStop")
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
	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	logger.Info("Loading config file: ", configFile)
	if !debug {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic Error:", err)
			}
		}()
	}

	app := fx.New(
		// 初始化配置应用配置
		fx.Provide(func() *types.AppConfig {
			config, err := core.LoadConfig(configFile)
			if err != nil {
				log.Fatal(err)
			}
			config.Path = configFile
			if debug {
				_ = core.SaveConfig(config)
			}
			return config
		}),
		// 创建应用服务
		fx.Provide(core.NewServer),
		// 初始化
		fx.Invoke(func(s *core.AppServer, client *redis.Client) {
			s.Init(debug, client)
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
		fx.Provide(handler.NewChatRoleHandler),
		fx.Provide(handler.NewUserHandler),
		fx.Provide(chatimpl.NewChatHandler),
		fx.Provide(handler.NewUploadHandler),
		fx.Provide(handler.NewSmsHandler),
		fx.Provide(handler.NewRewardHandler),
		fx.Provide(handler.NewCaptchaHandler),
		fx.Provide(handler.NewMidJourneyHandler),
		fx.Provide(handler.NewChatModelHandler),
		fx.Provide(handler.NewSdJobHandler),
		fx.Provide(handler.NewPaymentHandler),
		fx.Provide(handler.NewOrderHandler),
		fx.Provide(handler.NewProductHandler),
		fx.Provide(handler.NewConfigHandler),
		fx.Provide(handler.NewPowerLogHandler),

		fx.Provide(admin.NewConfigHandler),
		fx.Provide(admin.NewAdminHandler),
		fx.Provide(admin.NewApiKeyHandler),
		fx.Provide(admin.NewUserHandler),
		fx.Provide(admin.NewChatRoleHandler),
		fx.Provide(admin.NewRewardHandler),
		fx.Provide(admin.NewDashboardHandler),
		fx.Provide(admin.NewChatModelHandler),
		fx.Provide(admin.NewProductHandler),
		fx.Provide(admin.NewOrderHandler),
		fx.Provide(admin.NewChatHandler),
		fx.Provide(admin.NewPowerLogHandler),

		// 创建服务
		fx.Provide(sms.NewSendServiceManager),
		fx.Provide(func(config *types.AppConfig) *service.CaptchaService {
			return service.NewCaptchaService(config.ApiConfig)
		}),
		fx.Provide(oss.NewUploaderManager),
		fx.Provide(mj.NewService),

		// 邮件服务
		fx.Provide(service.NewSmtpService),

		// 微信机器人服务
		fx.Provide(wx.NewWeChatBot),
		fx.Invoke(func(config *types.AppConfig, bot *wx.Bot) {
			if config.WeChatBot {
				err := bot.Run()
				if err != nil {
					logger.Error("微信登录失败：", err)
				}
			}
		}),

		// MidJourney service pool
		fx.Provide(mj.NewServicePool),
		fx.Invoke(func(pool *mj.ServicePool) {
			if pool.HasAvailableService() {
				pool.DownloadImages()
				pool.CheckTaskNotify()
				pool.SyncTaskProgress()
			}
		}),

		// Stable Diffusion 机器人
		fx.Provide(sd.NewServicePool),
		fx.Invoke(func(pool *sd.ServicePool) {
			if pool.HasAvailableService() {
				pool.CheckTaskNotify()
				pool.CheckTaskStatus()
			}
		}),

		fx.Provide(payment.NewAlipayService),
		fx.Provide(payment.NewHuPiPay),
		fx.Provide(payment.NewPayJS),
		fx.Provide(service.NewSnowflake),
		fx.Provide(service.NewXXLJobExecutor),
		fx.Invoke(func(exec *service.XXLJobExecutor, config *types.AppConfig) {
			if config.XXLConfig.Enabled {
				go func() {
					log.Fatal(exec.Run())
				}()
			}
		}),

		// 注册路由
		fx.Invoke(func(s *core.AppServer, h *handler.ChatRoleHandler) {
			group := s.Engine.Group("/api/role/")
			group.GET("list", h.List)
			group.POST("update", h.UpdateRole)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.UserHandler) {
			group := s.Engine.Group("/api/user/")
			group.POST("register", h.Register)
			group.POST("login", h.Login)
			group.GET("logout", h.Logout)
			group.GET("session", h.Session)
			group.GET("profile", h.Profile)
			group.POST("profile/update", h.ProfileUpdate)
			group.POST("password", h.UpdatePass)
			group.POST("bind/username", h.BindUsername)
			group.POST("resetPass", h.ResetPass)
		}),
		fx.Invoke(func(s *core.AppServer, h *chatimpl.ChatHandler) {
			group := s.Engine.Group("/api/chat/")
			group.Any("new", h.ChatHandle)
			group.GET("list", h.List)
			group.GET("detail", h.Detail)
			group.POST("update", h.Update)
			group.GET("remove", h.Remove)
			group.GET("history", h.History)
			group.GET("clear", h.Clear)
			group.POST("tokens", h.Tokens)
			group.GET("stop", h.StopGenerate)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.UploadHandler) {
			s.Engine.POST("/api/upload", h.Upload)
			s.Engine.GET("/api/upload/list", h.List)
			s.Engine.GET("/api/upload/remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.SmsHandler) {
			group := s.Engine.Group("/api/sms/")
			group.POST("code", h.SendCode)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.CaptchaHandler) {
			group := s.Engine.Group("/api/captcha/")
			group.GET("get", h.Get)
			group.POST("check", h.Check)
			group.GET("slide/get", h.SlideGet)
			group.POST("slide/check", h.SlideCheck)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.RewardHandler) {
			group := s.Engine.Group("/api/reward/")
			group.POST("verify", h.Verify)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.MidJourneyHandler) {
			group := s.Engine.Group("/api/mj/")
			group.Any("client", h.Client)
			group.POST("image", h.Image)
			group.POST("upscale", h.Upscale)
			group.POST("variation", h.Variation)
			group.GET("jobs", h.JobList)
			group.GET("imgWall", h.ImgWall)
			group.POST("remove", h.Remove)
			group.POST("publish", h.Publish)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.SdJobHandler) {
			group := s.Engine.Group("/api/sd")
			group.Any("client", h.Client)
			group.POST("image", h.Image)
			group.GET("jobs", h.JobList)
			group.GET("imgWall", h.ImgWall)
			group.POST("remove", h.Remove)
			group.POST("publish", h.Publish)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ConfigHandler) {
			group := s.Engine.Group("/api/config/")
			group.GET("get", h.Get)
		}),

		// 管理后台控制器
		fx.Invoke(func(s *core.AppServer, h *admin.ConfigHandler) {
			group := s.Engine.Group("/api/admin/config/")
			group.POST("update", h.Update)
			group.GET("get", h.Get)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ManagerHandler) {
			group := s.Engine.Group("/api/admin/")
			group.POST("login", h.Login)
			group.GET("logout", h.Logout)
			group.GET("session", h.Session)
			group.GET("list", h.List)
			group.POST("save", h.Save)
			group.POST("enable", h.Enable)
			group.GET("remove", h.Remove)
			group.POST("resetPass", h.ResetPass)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ApiKeyHandler) {
			group := s.Engine.Group("/api/admin/apikey/")
			group.POST("save", h.Save)
			group.GET("list", h.List)
			group.POST("set", h.Set)
			group.GET("remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.UserHandler) {
			group := s.Engine.Group("/api/admin/user/")
			group.GET("list", h.List)
			group.POST("save", h.Save)
			group.GET("remove", h.Remove)
			group.GET("loginLog", h.LoginLog)
			group.POST("resetPass", h.ResetPass)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ChatRoleHandler) {
			group := s.Engine.Group("/api/admin/role/")
			group.GET("list", h.List)
			group.POST("save", h.Save)
			group.POST("sort", h.Sort)
			group.POST("set", h.Set)
			group.POST("remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.RewardHandler) {
			group := s.Engine.Group("/api/admin/reward/")
			group.GET("list", h.List)
			group.POST("remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.DashboardHandler) {
			group := s.Engine.Group("/api/admin/dashboard/")
			group.GET("stats", h.Stats)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ChatModelHandler) {
			group := s.Engine.Group("/api/model/")
			group.GET("list", h.List)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ChatModelHandler) {
			group := s.Engine.Group("/api/admin/model/")
			group.POST("save", h.Save)
			group.GET("list", h.List)
			group.POST("set", h.Set)
			group.POST("sort", h.Sort)
			group.GET("remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.PaymentHandler) {
			group := s.Engine.Group("/api/payment/")
			group.GET("doPay", h.DoPay)
			group.GET("payWays", h.GetPayWays)
			group.POST("query", h.OrderQuery)
			group.POST("qrcode", h.PayQrcode)
			group.POST("mobile", h.Mobile)
			group.POST("alipay/notify", h.AlipayNotify)
			group.POST("hupipay/notify", h.HuPiPayNotify)
			group.POST("payjs/notify", h.PayJsNotify)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ProductHandler) {
			group := s.Engine.Group("/api/admin/product/")
			group.POST("save", h.Save)
			group.GET("list", h.List)
			group.POST("enable", h.Enable)
			group.POST("sort", h.Sort)
			group.GET("remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.OrderHandler) {
			group := s.Engine.Group("/api/admin/order/")
			group.POST("list", h.List)
			group.GET("remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.OrderHandler) {
			group := s.Engine.Group("/api/order/")
			group.POST("list", h.List)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ProductHandler) {
			group := s.Engine.Group("/api/product/")
			group.GET("list", h.List)
		}),

		fx.Provide(handler.NewInviteHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.InviteHandler) {
			group := s.Engine.Group("/api/invite/")
			group.GET("code", h.Code)
			group.POST("list", h.List)
			group.GET("hits", h.Hits)
		}),

		fx.Provide(admin.NewFunctionHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.FunctionHandler) {
			group := s.Engine.Group("/api/admin/function/")
			group.POST("save", h.Save)
			group.POST("set", h.Set)
			group.GET("list", h.List)
			group.GET("remove", h.Remove)
			group.GET("token", h.GenToken)
		}),

		// 验证码
		fx.Provide(admin.NewCaptchaHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.CaptchaHandler) {
			group := s.Engine.Group("/api/admin/login/")
			group.GET("captcha", h.GetCaptcha)
		}),

		fx.Provide(admin.NewUploadHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.UploadHandler) {
			s.Engine.POST("/api/admin/upload", h.Upload)
		}),

		fx.Provide(handler.NewFunctionHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.FunctionHandler) {
			group := s.Engine.Group("/api/function/")
			group.POST("weibo", h.WeiBo)
			group.POST("zaobao", h.ZaoBao)
			group.POST("dalle3", h.Dall3)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ChatHandler) {
			group := s.Engine.Group("/api/admin/chat/")
			group.POST("list", h.List)
			group.POST("message", h.Messages)
			group.GET("history", h.History)
			group.GET("remove", h.RemoveChat)
			group.GET("message/remove", h.RemoveMessage)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.PowerLogHandler) {
			group := s.Engine.Group("/api/powerLog/")
			group.POST("list", h.List)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.PowerLogHandler) {
			group := s.Engine.Group("/api/admin/powerLog/")
			group.POST("list", h.List)
		}),
		fx.Provide(admin.NewMenuHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.MenuHandler) {
			group := s.Engine.Group("/api/admin/menu/")
			group.POST("save", h.Save)
			group.GET("list", h.List)
			group.POST("enable", h.Enable)
			group.POST("sort", h.Sort)
			group.GET("remove", h.Remove)
		}),
		fx.Provide(handler.NewMenuHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.MenuHandler) {
			group := s.Engine.Group("/api/menu/")
			group.GET("list", h.List)
		}),
		fx.Invoke(func(s *core.AppServer, db *gorm.DB) {
			go func() {
				err := s.Run(db)
				if err != nil {
					log.Fatal(err)
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
