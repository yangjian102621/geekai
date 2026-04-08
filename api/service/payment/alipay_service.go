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
	"net/http"
	"os"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

type AlipayService struct {
	client *alipay.Client
	config *types.AlipayConfig
}

var logger = logger2.GetLogger()

func NewAlipayService(sysConfig *types.SystemConfig) (*AlipayService, error) {
	config := sysConfig.Payment.Alipay
	if !config.Enabled {
		logger.Debug("Disabled Alipay service")
	}

	service := &AlipayService{config: &config}
	if config.Enabled {
		err := service.UpdateConfig(&config)
		if err != nil {
			logger.Errorf("支付宝服务初始化失败: %v", err)
		}
	}

	return service, nil
}

func (s *AlipayService) UpdateConfig(config *types.AlipayConfig) error {
	client, err := alipay.NewClient(config.AppId, config.PrivateKey, !config.SandBox)
	if err != nil {
		return fmt.Errorf("error with initialize alipay service: %v", err)
	}

	s.client = client
	s.config = config
	if os.Getenv("GEEKAI_DEBUG") == "true" {
		logger.Info("Alipay Debug mode is enabled")
		client.DebugSwitch = gopay.DebugOn
	}
	return nil
}

func (s *AlipayService) Pay(params PayRequest) (string, error) {
	bm := make(gopay.BodyMap)
	bm.Set("subject", params.Subject)
	bm.Set("out_trade_no", params.OutTradeNo)
	bm.Set("total_amount", params.TotalFee)
	return s.client.TradeWapPay(context.Background(), bm)
}

func (s *AlipayService) Query(outTradeNo string) (OrderInfo, error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo)
	rsp, err := s.client.TradeQuery(context.Background(), bm)
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with trade query: %v", err)
	}

	switch rsp.Response.TradeStatus {
	case "TRADE_SUCCESS":
		logger.Debugf("支付宝查询订单成功：%+v", rsp.Response)
		return OrderInfo{
			OutTradeNo: rsp.Response.OutTradeNo,
			TradeId:    rsp.Response.TradeNo,
			Amount:     rsp.Response.TotalAmount,
			Status:     Success,
			PayTime:    rsp.Response.SendPayDate,
		}, nil
	case "TRADE_CLOSED":
		return OrderInfo{Status: Closed}, nil
	default:
		return OrderInfo{}, fmt.Errorf("error with trade query: %v", rsp.Response.TradeStatus)
	}
}

// TradeVerify 交易验证
func (s *AlipayService) TradeVerify(request *http.Request) (OrderInfo, error) {
	notifyReq, err := alipay.ParseNotifyToBodyMap(request) // c.Request 是 gin 框架的写法
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with parse notify request: %v", err)
	}

	_, err = alipay.VerifySignWithCert(s.config.AlipayPublicKey, notifyReq)
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with verify sign: %v", err)
	}

	return s.Query(request.Form.Get("out_trade_no"))
}

var _ PayService = (*AlipayService)(nil)
