package wexin

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"github.com/eatmoreapple/openwechat"
	"gorm.io/gorm"
)

// 微信收款机器人服务
var logger = logger2.GetLogger()

type WeChatBot struct {
	bot       *openwechat.Bot
	db        *gorm.DB
	appConfig *types.AppConfig
}

func NewWeChatBot(db *gorm.DB, config *types.AppConfig) *WeChatBot {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		MessageHandler(msg, db)
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = QrCodeCallBack
	return &WeChatBot{
		bot:       bot,
		db:        db,
		appConfig: config,
	}
}

func (b *WeChatBot) Login() error {
	if !b.appConfig.StartWechatBot {
		return nil
	}

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	// 执行热登录
	err := b.bot.HotLogin(reloadStorage)
	if err != nil {
		logger.Error("login error: %v", err)
		return b.bot.Login()
	}
	logger.Info("微信登录成功！")
	return nil
}
