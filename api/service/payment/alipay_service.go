package payment

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"net/http"
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

	client, err := alipay.NewClient(config.AppId, priKey, !config.SandBox)
	if err != nil {
		return nil, fmt.Errorf("error with initialize alipay service: %v", err)
	}

	//client.DebugSwitch = gopay.DebugOn // 开启调试模式
	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8). // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2) // 设置签名类型，不设置默认 RSA2

	if err = client.SetCertSnByPath(config.PublicKey, config.RootCert, config.AlipayPublicKey); err != nil {
		return nil, fmt.Errorf("error with load payment public key: %v", err)
	}

	return &AlipayService{config: &config, client: client}, nil
}

type AlipayParams struct {
	OutTradeNo string `json:"out_trade_no"`
	Subject    string `json:"subject"`
	TotalFee   string `json:"total_fee"`
	ReturnURL  string `json:"return_url"`
	NotifyURL  string `json:"notify_url"`
}

func (s *AlipayService) PayMobile(params AlipayParams) (string, error) {
	bm := make(gopay.BodyMap)
	bm.Set("subject", params.Subject)
	bm.Set("out_trade_no", params.OutTradeNo)
	bm.Set("quit_url", params.ReturnURL)
	bm.Set("total_amount", params.TotalFee)
	bm.Set("product_code", "QUICK_WAP_WAY")
	return s.client.SetNotifyUrl(params.NotifyURL).SetReturnUrl(params.ReturnURL).TradeWapPay(context.Background(), bm)
}

func (s *AlipayService) PayPC(params AlipayParams) (string, error) {
	bm := make(gopay.BodyMap)
	bm.Set("subject", params.Subject)
	bm.Set("out_trade_no", params.OutTradeNo)
	bm.Set("total_amount", params.TotalFee)
	bm.Set("product_code", "FAST_INSTANT_TRADE_PAY")
	return s.client.SetNotifyUrl(params.NotifyURL).SetReturnUrl(params.ReturnURL).TradePagePay(context.Background(), bm)
}

// TradeVerify 交易验证
func (s *AlipayService) TradeVerify(request *http.Request) NotifyVo {
	notifyReq, err := alipay.ParseNotifyToBodyMap(request) // c.Request 是 gin 框架的写法
	if err != nil {
		return NotifyVo{
			Status:  Failure,
			Message: "error with parse notify request: " + err.Error(),
		}
	}

	_, err = alipay.VerifySignWithCert(s.config.AlipayPublicKey, notifyReq)
	if err != nil {
		return NotifyVo{
			Status:  Failure,
			Message: "error with verify sign: " + err.Error(),
		}
	}

	return s.TradeQuery(request.Form.Get("out_trade_no"))
}

func (s *AlipayService) TradeQuery(outTradeNo string) NotifyVo {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo)

	//查询订单
	rsp, err := s.client.TradeQuery(context.Background(), bm)
	if err != nil {
		return NotifyVo{
			Status:  Failure,
			Message: "异步查询验证订单信息发生错误" + outTradeNo + err.Error(),
		}
	}

	if rsp.Response.TradeStatus == "TRADE_SUCCESS" {
		return NotifyVo{
			Status:     Success,
			OutTradeNo: rsp.Response.OutTradeNo,
			TradeId:    rsp.Response.TradeNo,
			Amount:     rsp.Response.TotalAmount,
			Subject:    rsp.Response.Subject,
			Message:    "OK",
		}
	} else {
		return NotifyVo{
			Status:  Failure,
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
