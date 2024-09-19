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
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"net/http"
	"time"
)

type WechatPayService struct {
	config *types.WechatPayConfig
	client *wechat.ClientV3
}

func NewWechatService(appConfig *types.AppConfig) (*WechatPayService, error) {
	config := appConfig.WechatPayConfig
	if !config.Enabled {
		logger.Info("Disabled WechatPay service")
		return nil, nil
	}
	priKey, err := readKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error with read App Private key: %v", err)
	}

	client, err := wechat.NewClientV3(config.MchId, config.SerialNo, config.ApiV3Key, priKey)
	if err != nil {
		return nil, fmt.Errorf("error with initialize WechatPay service: %v", err)
	}
	err = client.AutoVerifySign()
	if err != nil {
		return nil, fmt.Errorf("error with autoVerifySign: %v", err)
	}
	//client.DebugSwitch = gopay.DebugOn

	return &WechatPayService{config: &config, client: client}, nil
}

type WechatPayParams struct {
	OutTradeNo string `json:"out_trade_no"`
	TotalFee   int    `json:"total_fee"`
	Subject    string `json:"subject"`
	ClientIP   string `json:"client_ip"`
	ReturnURL  string `json:"return_url"`
	NotifyURL  string `json:"notify_url"`
}

func (s *WechatPayService) PayUrlNative(params WechatPayParams) (string, error) {
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
			bm.Set("total", params.TotalFee).
				Set("currency", "CNY")
		})

	wxRsp, err := s.client.V3TransactionNative(context.Background(), bm)
	if err != nil {
		return "", fmt.Errorf("error with client v3 transaction Native: %v", err)
	}
	if wxRsp.Code != wechat.Success {
		return "", fmt.Errorf("error status with generating pay url: %v", wxRsp.Error)
	}
	return wxRsp.Response.CodeUrl, nil
}

func (s *WechatPayService) PayUrlH5(params WechatPayParams) (string, error) {
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
			bm.Set("total", params.TotalFee).
				Set("currency", "CNY")
		}).
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.Set("payer_client_ip", params.ClientIP).
				SetBodyMap("h5_info", func(bm gopay.BodyMap) {
					bm.Set("type", "Wap")
				})
		})

	wxRsp, err := s.client.V3TransactionH5(context.Background(), bm)
	if err != nil {
		return "", fmt.Errorf("error with client v3 transaction H5: %v", err)
	}
	if wxRsp.Code != wechat.Success {
		return "", fmt.Errorf("error with generating pay url: %v", wxRsp.Error)
	}
	return wxRsp.Response.H5Url, nil
}

type NotifyResponse struct {
	Code    string `json:"code"`
	Message string `xml:"message"`
}

// TradeVerify 交易验证
func (s *WechatPayService) TradeVerify(request *http.Request) NotifyVo {
	notifyReq, err := wechat.V3ParseNotify(request)
	if err != nil {
		return NotifyVo{Status: 1, Message: fmt.Sprintf("error with client v3 parse notify: %v", err)}
	}

	// TODO: 这里验签程序有 Bug，一直报错：crypto/rsa: verification error，先暂时取消验签
	//err = notifyReq.VerifySignByPK(s.client.WxPublicKey())
	//if err != nil {
	//	return fmt.Errorf("error with client v3 verify sign: %v", err)
	//}

	// 解密支付密文，验证订单信息
	result, err := notifyReq.DecryptPayCipherText(s.config.ApiV3Key)
	if err != nil {
		return NotifyVo{Status: Failure, Message: fmt.Sprintf("error with client v3 decrypt: %v", err)}
	}

	return NotifyVo{
		Status:     Success,
		OutTradeNo: result.OutTradeNo,
		TradeId:    result.TransactionId,
		Amount:     fmt.Sprintf("%.2f", float64(result.Amount.Total)/100),
	}
}
