package wexin

import (
	"chatplus/store/model"
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

// MessageHandler 消息处理
func MessageHandler(msg *openwechat.Message, db *gorm.DB) {
	sender, err := msg.Sender()
	if err != nil {
		return
	}
	
	// 只处理微信支付的推送消息
	if sender.NickName == "微信支付" ||
		msg.MsgType == openwechat.MsgTypeApp ||
		msg.AppMsgType == openwechat.AppMsgTypeUrl {
		// 解析支付金额
		message, err := parseTransactionMessage(msg.Content)
		if err == nil {
			transaction := extractTransaction(message)
			logger.Infof("解析到收款信息：%+v", transaction)
			if transaction.Amount <= 0 {
				return
			}
			var item model.Reward
			res := db.Where("tx_id = ?", transaction.TransId).First(&item)
			if res.Error == nil {
				logger.Infof("当前交易 ID %s 己经存在！", transaction.TransId)
				return
			}

			res = db.Create(&model.Reward{
				TxId:   transaction.TransId,
				Amount: transaction.Amount,
				Remark: transaction.Remark,
				Status: false,
			})
			if res.Error != nil {
				logger.Errorf("交易保存失败，ID: %s", transaction.TransId)
			}
		}
	}
}

// QrCodeCallBack 登录扫码回调，
func QrCodeCallBack(uuid string) {
	logger.Info("请使用微信扫描下面二维码登录")
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Medium)
	logger.Info(q.ToString(true))
}
