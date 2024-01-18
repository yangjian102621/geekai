package wx

import (
	logger2 "chatplus/logger"
	"chatplus/store/model"
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
	"os"
	"strconv"
)

// 微信收款机器人
var logger = logger2.GetLogger()

type Bot struct {
	bot   *openwechat.Bot
	token string
	db    *gorm.DB
}

func NewWeChatBot(db *gorm.DB) *Bot {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	return &Bot{
		bot: bot,
		db:  db,
	}
}

func (b *Bot) Run() error {
	logger.Info("Starting WeChat Bot...")

	// set message handler
	b.bot.MessageHandler = func(msg *openwechat.Message) {
		b.messageHandler(msg)
	}
	// scan code login callback
	b.bot.UUIDCallback = b.qrCodeCallBack
	debug, err := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if debug {
		reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
		err = b.bot.HotLogin(reloadStorage, true)
	} else {
		err = b.bot.Login()
	}
	if err != nil {
		return err
	}

	logger.Info("微信登录成功！")
	return nil
}

// message handler
func (b *Bot) messageHandler(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		return
	}

	// 只处理微信支付的推送消息
	if sender.NickName == "微信支付" ||
		msg.MsgType == openwechat.MsgTypeApp ||
		msg.AppMsgType == openwechat.AppMsgTypeUrl {
		// 解析支付金额
		message := parseTransactionMessage(msg.Content)
		transaction := extractTransaction(message)
		logger.Infof("解析到收款信息：%+v", transaction)
		if transaction.TransId != "" {
			var item model.Reward
			res := b.db.Where("tx_id = ?", transaction.TransId).First(&item)
			if item.Id > 0 {
				logger.Error("当前交易 ID 己经存在！")
				return
			}

			res = b.db.Create(&model.Reward{
				TxId:   transaction.TransId,
				Amount: transaction.Amount,
				Remark: transaction.Remark,
				Status: false,
			})
			if res.Error != nil {
				logger.Errorf("交易保存失败: %v", res.Error)
			}
		}
	}
}

func (b *Bot) qrCodeCallBack(uuid string) {
	logger.Info("请使用微信扫描下面二维码登录")
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Medium)
	logger.Info(q.ToString(true))
}
