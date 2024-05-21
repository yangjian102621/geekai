package payment

import (
	"chatplus/core/types"
	"chatplus/store/model"
	chatPlusUtils "chatplus/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"log"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type WxpayService struct {
	config             *types.WxpayConfig
	client             *core.Client
	certificateVisitor core.CertificateGetter
}

func NewWxpayService(appConfig *types.AppConfig) *WxpayService {
	config := appConfig.WxpayConfig
	if !config.Enabled {
		logger.Info("Disabled Wxpay service")
		return nil
	}
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(config.PrivateKey)
	if err != nil {
		log.Print("load merchant private key error")
		return nil
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(config.MchId, config.CertificateSerialNo, mchPrivateKey, config.MchKey),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil
	}
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err2 := downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, config.CertificateSerialNo, config.MchId, config.MchKey)
	if err2 != nil {
		logger.Error("支付回调校验失败，请检查应用私钥配置文件")
		return nil
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(config.MchId)
	return &WxpayService{&config, client, certificateVisitor}
}

func (s *WxpayService) Pay(user model.User, order model.Order) (resp *jsapi.PrepayWithRequestPaymentResponse, err error) {
	return s.PayUrlMobile(user, order)
}

func (s *WxpayService) PayUrlMobile(user model.User, order model.Order) (resp *jsapi.PrepayWithRequestPaymentResponse, err error) {
	svc := jsapi.JsapiApiService{Client: s.client}
	var outTradeNo = chatPlusUtils.RandString(16)
	resp, result, err := svc.PrepayWithRequestPayment(context.Background(),
		jsapi.PrepayRequest{
			Appid:         core.String(s.config.AppId),
			Mchid:         core.String(s.config.MchId),
			Description:   core.String(order.Subject),
			OutTradeNo:    core.String(outTradeNo),
			TimeExpire:    core.Time(time.Now()),
			Attach:        core.String(order.OrderNo),
			NotifyUrl:     core.String(s.config.NotifyURL),
			SupportFapiao: core.Bool(false),
			Amount: &jsapi.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(int64(order.Amount * 100)),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(user.OfficialOpenid),
			},
			SceneInfo: &jsapi.SceneInfo{
				PayerClientIp: core.String("127.0.0.1"),
			},
			SettleInfo: &jsapi.SettleInfo{
				ProfitSharing: core.Bool(false),
			},
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call Prepay err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
	return resp, err
}

func (s *WxpayService) TradeVerify(c *gin.Context) (orderNo string, tradeNo string, code int) {
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, _ := notify.NewRSANotifyHandler(s.config.MchKey, verifiers.NewSHA256WithRSAVerifier(s.certificateVisitor))
	// 2. 获取商户号对应的微信支付平台证书访问器
	transaction := new(payments.Transaction)
	_, err3 := handler.ParseNotifyRequest(context.Background(), c.Request, transaction)
	// 如果验签未通过，或者解密失败
	if err3 != nil {
		return "0", "0", 401
	}
	return *transaction.Attach, *transaction.OutTradeNo, 200
}
