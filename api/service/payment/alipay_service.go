package payment

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"log"
	"net/url"
	"os"
)

type AlipayService struct {
	config *types.AlipayConfig
	client *alipay.Client
}

var logger = logger2.GetLogger()

func NewAlipayService(appConfig *types.AppConfig) (*AlipayService, error) {
	config := appConfig.AlipayConfig
	if !config.Enabled {
		logger.Info("Disabled Alipay service")
		return nil, nil
	}
	priKey, err := readKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error with read App Private key: %v", err)
	}

	xClient, err := alipay.New(config.AppId, priKey, !config.SandBox)
	if err != nil {
		return nil, fmt.Errorf("error with initialize alipay service: %v", err)
	}

	if err = xClient.LoadAppCertPublicKeyFromFile(config.PublicKey); err != nil {
		return nil, fmt.Errorf("error with loading App PublicKey: %v", err)
	}
	if err = xClient.LoadAliPayRootCertFromFile(config.RootCert); err != nil {
		return nil, fmt.Errorf("error with loading alipay RootCert: %v", err)
	}
	if err = xClient.LoadAlipayCertPublicKeyFromFile(config.AlipayPublicKey); err != nil {
		return nil, fmt.Errorf("error with loading Alipay PublicKey: %v", err)
	}

	return &AlipayService{config: &config, client: xClient}, nil
}

func (s *AlipayService) PayUrlMobile(outTradeNo string, notifyURL string, returnURL string, Amount string, subject string) (string, error) {
	var p = alipay.TradeWapPay{}
	p.NotifyURL = notifyURL
	p.ReturnURL = returnURL
	p.Subject = subject
	p.OutTradeNo = outTradeNo
	p.TotalAmount = Amount
	p.ProductCode = "QUICK_WAP_WAY"
	res, err := s.client.TradeWapPay(p)
	if err != nil {
		return "", err
	}

	return res.String(), err
}

func (s *AlipayService) PayUrlPc(outTradeNo string, notifyURL string, returnURL string, amount string, subject string) (string, error) {
	var p = alipay.TradePagePay{}
	p.NotifyURL = notifyURL
	p.ReturnURL = returnURL
	p.Subject = subject
	p.OutTradeNo = outTradeNo
	p.TotalAmount = amount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	res, err := s.client.TradePagePay(p)
	if err != nil {
		return "", nil
	}

	return res.String(), err
}

// TradeVerify 交易验证
func (s *AlipayService) TradeVerify(reqForm url.Values) NotifyVo {
	err := s.client.VerifySign(reqForm)
	if err != nil {
		log.Println("异步通知验证签名发生错误", err)
		return NotifyVo{
			Status:  0,
			Message: "异步通知验证签名发生错误",
		}
	}

	return s.TradeQuery(reqForm.Get("out_trade_no"))
}

func (s *AlipayService) TradeQuery(outTradeNo string) NotifyVo {
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	rsp, err := s.client.TradeQuery(p)
	if err != nil {
		return NotifyVo{
			Status:  0,
			Message: "异步查询验证订单信息发生错误" + outTradeNo + err.Error(),
		}
	}

	if rsp.IsSuccess() == true && rsp.TradeStatus == "TRADE_SUCCESS" {
		return NotifyVo{
			Status:     1,
			OutTradeNo: rsp.OutTradeNo,
			TradeNo:    rsp.TradeNo,
			Amount:     rsp.TotalAmount,
			Subject:    rsp.Subject,
			Message:    "OK",
		}
	} else {
		return NotifyVo{
			Status:  0,
			Message: "异步查询验证订单信息发生错误" + outTradeNo,
		}
	}
}

func readKey(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type NotifyVo struct {
	Status     int
	OutTradeNo string
	TradeNo    string
	Amount     string
	Message    string
	Subject    string
}

func (v NotifyVo) Success() bool {
	return v.Status == 1
}
