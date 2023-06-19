package main

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/handler/admin"
	logger2 "chatplus/logger"
	"chatplus/service"
	"chatplus/store"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

var configFile string
var debugMode bool

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
	logger.Info("Loading config file: ", configFile)
	app := fx.New(
		// 初始化配置应用配置
		fx.Provide(func() *types.AppConfig {
			config, err := core.LoadConfig(configFile)
			if err != nil {
				log.Fatal(err)
			}
			return config
		}),
		// 创建应用服务
		fx.Provide(core.NewServer),
		// 初始化
		fx.Invoke(func(s *core.AppServer) {
			s.Init(debugMode)
		}),

		// 初始化数据库
		fx.Provide(store.NewGormConfig),
		fx.Provide(store.NewMysql),

		// 创建 Ip2Region 查询对象
		fx.Provide(func() (*xdb.Searcher, error) {
			dbPath := "res/ip2region.xdb"
			cBuff, err := xdb.LoadContentFromFile(dbPath)
			if err != nil {
				return nil, err
			}

			return xdb.NewWithBuffer(cBuff)
		}),

		// 初始化服务
		fx.Provide(store.NewLevelDB),
		fx.Provide(service.NewChatRoleService),
		fx.Invoke(core.InitChatRoles),

		// 创建控制器
		fx.Provide(handler.NewChatRoleHandler),
		fx.Provide(handler.NewUserHandler),
		fx.Provide(handler.NewChatHandler),
		fx.Provide(admin.NewConfigHandler),

		fx.Provide(admin.NewAdminHandler),
		fx.Provide(admin.NewApiKeyHandler),

		// 注册路由
		fx.Invoke(func(s *core.AppServer, h *handler.ChatRoleHandler) {
			group := s.Engine.Group("/api/chat/role/")
			group.GET("list", h.List)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.UserHandler) {
			group := s.Engine.Group("/api/user/")
			group.POST("register", h.Register)
			group.GET("list", h.List)
			group.POST("login", h.Login)
			group.GET("logout", h.Logout)
			group.GET("session", h.Session)
			group.GET("profile", h.Profile)
			group.POST("profile/update", h.ProfileUpdate)
			group.POST("password", h.Password)
		}),
		fx.Invoke(func(s *core.AppServer, h *handler.ChatHandler) {
			group := s.Engine.Group("/api/chat/")
			group.Any("new", h.ChatHandle)
			group.GET("list", h.List)
			group.POST("update", h.Update)
			group.GET("remove", h.Remove)
			group.GET("history", h.History)
			group.GET("clear", h.Clear)
			group.GET("tokens", h.Tokens)
			group.GET("stop", h.StopGenerate)
		}),

		//
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
		}),
		fx.Invoke(func(s *core.AppServer, h *admin.ApiKeyHandler) {
			group := s.Engine.Group("/api/admin/apikey/")
			group.POST("add", h.Add)
			group.GET("list", h.List)
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

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "AppConfig file path (default: config.toml)")
	flag.BoolVar(&debugMode, "debug", true, "Enable debug mode (default: true, recommend to set false in production env)")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Printf(`ChatGPT-Web-Plus, Version: 2.0.0
USAGE: 
  %s [command options]
OPTIONS:
`, os.Args[0])

	flagSet := flag.CommandLine
	order := []string{"config", "debug"}
	for _, name := range order {
		f := flagSet.Lookup(name)
		fmt.Printf("  --%s => %s\n", f.Name, f.Usage)
	}
}
