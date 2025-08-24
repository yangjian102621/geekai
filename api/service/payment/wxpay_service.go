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
	"geekai/utils"
	"net/http"
	"os"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

type WxPayService struct {
	config *types.WxPayConfig
	client *wechat.ClientV3
}

func NewWxpayService(sysConfig *types.SystemConfig) (*WxPayService, error) {
	config := sysConfig.Payment.WxPay
	if !config.Enabled {
		logger.Debug("Disabled WechatPay service")
	}

	service := &WxPayService{config: &config}
	if config.Enabled {
		err := service.UpdateConfig(&config)
		if err != nil {
			logger.Errorf("微信支付服务初始化失败: %v", err)
		}
	}

	return service, nil
}

func (s *WxPayService) UpdateConfig(config *types.WxPayConfig) error {
	client, err := wechat.NewClientV3(config.MchId, config.SerialNo, config.ApiV3Key, config.PrivateKey)
	if err != nil {
		return fmt.Errorf("error with initialize WechatPay service: %v", err)
	}
	err = client.AutoVerifySign()
	if err != nil {
		return fmt.Errorf("error with autoVerifySign: %v", err)
	}
	s.client = client
	if os.Getenv("GEEKAI_DEBUG") == "true" {
		logger.Info("WechatPay Debug mode is enabled")
		client.DebugSwitch = gopay.DebugOn
	}
	s.config = config
	return nil
}

func (s *WxPayService) Pay(params PayRequest) (string, error) {
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("appid", s.config.AppId).
		Set("mchid", s.config.MchId).
		Set("description", params.Subject).
		Set("out_trade_no", params.OutTradeNo).
		Set("time_expire", expire).
		Set("notify_url", params.NotifyURL).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", utils.IntValue(params.TotalFee, 0)).
				Set("currency", "CNY")
		})
	if params.Device == "mobile" {
		bm.SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.Set("payer_client_ip", params.ClientIP)
		}).SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", params.OpenID)
		})
		wxRsp, err := s.client.V3TransactionJsapi(context.Background(), bm)
		if err != nil {
			return "", fmt.Errorf("error with client v3 transaction Jsapi: %v", err)
		}
		if wxRsp.Code != wechat.Success {
			return "", fmt.Errorf("error status with generating pay url: %v", wxRsp.Error)
		}
		return wxRsp.Response.PrepayId, nil
	} else if params.Device == "pc" {
		wxRsp, err := s.client.V3TransactionNative(context.Background(), bm)
		if err != nil {
			return "", fmt.Errorf("error with client v3 transaction Native: %v", err)
		}
		if wxRsp.Code != wechat.Success {
			return "", fmt.Errorf("error status with generating pay url: %v", wxRsp.Error)
		}
		return wxRsp.Response.CodeUrl, nil

	}
	return "", nil
}

func (s *WxPayService) Query(outTradeNo string) (OrderInfo, error) {
	wxRsp, err := s.client.V3TransactionQueryOrder(context.Background(), wechat.OutTradeNo, outTradeNo)
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with client v3 transaction query: %v", err)
	}

	if wxRsp.Code != wechat.Success {
		return OrderInfo{}, fmt.Errorf("error status with querying order: %v", wxRsp.Error)
	}

	if wxRsp.Response.TradeState == "CLOSED" {
		return OrderInfo{Status: Closed}, nil
	}

	orderInfo := OrderInfo{
		OutTradeNo: wxRsp.Response.OutTradeNo,
		TradeId:    wxRsp.Response.TransactionId,
		Amount:     fmt.Sprintf("%d", wxRsp.Response.Amount.Total/100),
		PayTime:    wxRsp.Response.SuccessTime,
	}
	if wxRsp.Response.TradeState == "SUCCESS" {
		orderInfo.Status = Success
	} else {
		orderInfo.Status = Failure
	}
	return orderInfo, nil
}

// TradeVerify 交易验证
func (s *WxPayService) TradeVerify(request *http.Request) (OrderInfo, error) {
	notifyReq, err := wechat.V3ParseNotify(request)
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with client v3 parse notify: %v", err)
	}

	// 解密支付密文，验证订单信息
	result, err := notifyReq.DecryptPayCipherText(s.config.ApiV3Key)
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with client v3 decrypt: %v", err)
	}

	return OrderInfo{
		Status:     Success,
		OutTradeNo: result.OutTradeNo,
		TradeId:    result.TransactionId,
		Amount:     fmt.Sprintf("%.2f", float64(result.Amount.Total)/100),
		PayTime:    result.SuccessTime,
	}, nil
}

// func (s *WechatPayService) PayUrlNative(params WechatPayParams) (string, error) {
// 	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
// 	// 初始化 BodyMap
// 	bm := make(gopay.BodyMap)
// 	bm.Set("appid", s.config.AppId).
// 		Set("mchid", s.config.MchId).
// 		Set("description", params.Subject).
// 		Set("out_trade_no", params.OutTradeNo).
// 		Set("time_expire", expire).
// 		Set("notify_url", params.NotifyURL).
// 		SetBodyMap("amount", func(bm gopay.BodyMap) {
// 			bm.Set("total", params.TotalFee).
// 				Set("currency", "CNY")
// 		})

// 	wxRsp, err := s.client.V3TransactionNative(context.Background(), bm)
// 	if err != nil {
// 		return "", fmt.Errorf("error with client v3 transaction Native: %v", err)
// 	}
// 	if wxRsp.Code != wechat.Success {
// 		return "", fmt.Errorf("error status with generating pay url: %v", wxRsp.Error)
// 	}
// 	return wxRsp.Response.CodeUrl, nil
// }

// func (s *WechatPayService) PayUrlH5(params WechatPayParams) (string, error) {
// 	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
// 	// 初始化 BodyMap
// 	bm := make(gopay.BodyMap)
// 	bm.Set("appid", s.config.AppId).
// 		Set("mchid", s.config.MchId).
// 		Set("description", params.Subject).
// 		Set("out_trade_no", params.OutTradeNo).
// 		Set("time_expire", expire).
// 		Set("notify_url", params.NotifyURL).
// 		SetBodyMap("amount", func(bm gopay.BodyMap) {
// 			bm.Set("total", params.TotalFee).
// 				Set("currency", "CNY")
// 		}).
// 		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
// 			bm.Set("payer_client_ip", params.ClientIP).
// 				SetBodyMap("h5_info", func(bm gopay.BodyMap) {
// 					bm.Set("type", "Wap")
// 				})
// 		})

// 	wxRsp, err := s.client.V3TransactionH5(context.Background(), bm)
// 	if err != nil {
// 		return "", fmt.Errorf("error with client v3 transaction H5: %v", err)
// 	}
// 	if wxRsp.Code != wechat.Success {
// 		return "", fmt.Errorf("error with generating pay url: %v", wxRsp.Error)
// 	}
// 	return wxRsp.Response.H5Url, nil
// }

// type NotifyResponse struct {
// 	Code    string `json:"code"`
// 	Message string `xml:"message"`
// }

var _ PayService = (*WxPayService)(nil)
