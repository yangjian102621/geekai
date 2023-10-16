package main

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/handler/admin"
	"chatplus/handler/chatimpl"
	logger2 "chatplus/logger"
	"chatplus/service"
	"chatplus/service/fun"
	"chatplus/service/mj"
	"chatplus/service/oss"
	"chatplus/service/sd"
	"chatplus/service/wx"
	"chatplus/store"
	"context"
	"embed"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

//go:embed res/ip2region.xdb
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

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.toml"
	}
	var debug bool
	debugEnv := os.Getenv("DEBUG")
	if debugEnv == "" {
		debug = true
	} else {
		debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	}
	logger.Info("Loading config file: ", configFile)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Panic Error:", err)
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
		fx.Provide(store.NewLevelDB),
		fx.Provide(store.NewRedisClient),

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

		// 创建函数
		fx.Provide(fun.NewFunctions),

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

		fx.Provide(admin.NewConfigHandler),
		fx.Provide(admin.NewAdminHandler),
		fx.Provide(admin.NewApiKeyHandler),
		fx.Provide(admin.NewUserHandler),
		fx.Provide(admin.NewChatRoleHandler),
		fx.Provide(admin.NewRewardHandler),
		fx.Provide(admin.NewDashboardHandler),
		fx.Provide(admin.NewChatModelHandler),

		// 创建服务
		fx.Provide(service.NewAliYunSmsService),
		fx.Provide(func(config *types.AppConfig) *service.CaptchaService {
			return service.NewCaptchaService(config.ApiConfig)
		}),
		fx.Provide(oss.NewUploaderManager),
		fx.Provide(mj.NewService),

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

		// MidJourney 机器人
		fx.Provide(mj.NewBot),
		fx.Provide(mj.NewClient),
		fx.Invoke(func(config *types.AppConfig, bot *mj.Bot) {
			if config.MjConfig.Enabled {
				err := bot.Run()
				if err != nil {
					log.Fatal("MidJourney 服务启动失败：", err)
				}
			}
		}),
		fx.Invoke(func(config *types.AppConfig, mjService *mj.Service) {
			if config.MjConfig.Enabled {
				go func() {
					mjService.Run()
				}()
			}
		}),

		// Stable Diffusion 机器人
		fx.Provide(sd.NewService),
		fx.Invoke(func(config *types.AppConfig, service *sd.Service) {
			if config.SdConfig.Enabled {
				go func() {
					service.Run()
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
			group.POST("password", h.Password)
			group.POST("bind/mobile", h.BindMobile)
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
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.SmsHandler) {
			group := s.Engine.Group("/api/sms/")
			group.GET("status", h.Status)
			group.POST("code", h.SendCode)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.CaptchaHandler) {
			group := s.Engine.Group("/api/captcha/")
			group.GET("get", h.Get)
			group.POST("check", h.Check)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.RewardHandler) {
			group := s.Engine.Group("/api/reward/")
			group.POST("verify", h.Verify)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.MidJourneyHandler) {
			group := s.Engine.Group("/api/mj/")
			group.POST("image", h.Image)
			group.POST("upscale", h.Upscale)
			group.POST("variation", h.Variation)
			group.GET("jobs", h.JobList)
			group.Any("client", h.Client)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.SdJobHandler) {
			group := s.Engine.Group("/api/sd")
			group.POST("image", h.Image)
			group.GET("jobs", h.JobList)
			group.Any("client", h.Client)
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
			group.GET("migrate", h.Migrate)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ApiKeyHandler) {
			group := s.Engine.Group("/api/admin/apikey/")
			group.POST("save", h.Save)
			group.GET("list", h.List)
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
			group.GET("remove", h.Remove)
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.RewardHandler) {
			group := s.Engine.Group("/api/admin/reward/")
			group.GET("list", h.List)
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
			group.POST("enable", h.Enable)
			group.POST("sort", h.Sort)
			group.GET("remove", h.Remove)
		}),

		fx.Invoke(func(s *core.AppServer, db *gorm.DB) {
			err := s.Run(db)
			if err != nil {
				log.Fatal(err)
			}
		}),

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
