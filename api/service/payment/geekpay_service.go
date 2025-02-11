package payment

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// GeekPayService Geek 支付服务
type GeekPayService struct {
	config *types.GeekPayConfig
}

func NewJPayService(appConfig *types.AppConfig) *GeekPayService {
	return &GeekPayService{
		config: &appConfig.GeekPayConfig,
	}
}

type GeekPayParams struct {
	Method     string `json:"method"`       // 接口类型
	Device     string `json:"device"`       // 设备类型
	Type       string `json:"type"`         // 支付方式
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	Name       string `json:"name"`         // 商品名称
	Money      string `json:"money"`        // 商品金额
	ClientIP   string `json:"clientip"`     //用户IP地址
	SubOpenId  string `json:"sub_openid"`   // 微信用户 openid，仅小程序支付需要
	SubAppId   string `json:"sub_appid"`    // 小程序 AppId，仅小程序支付需要
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
}

// Pay 支付订单
func (s *GeekPayService) Pay(params GeekPayParams) (*GeekPayResp, error) {
	p := map[string]string{
		"pid": s.config.AppId,
		//"method":       params.Method,
		"device":       params.Device,
		"type":         params.Type,
		"out_trade_no": params.OutTradeNo,
		"name":         params.Name,
		"money":        params.Money,
		"clientip":     params.ClientIP,
		"notify_url":   params.NotifyURL,
		"return_url":   params.ReturnURL,
		"timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
	}
	p["sign"] = s.Sign(p)
	p["sign_type"] = "MD5"
	return s.sendRequest(s.config.ApiURL, p)
}

func (s *GeekPayService) Sign(params map[string]string) string {
	// 按字母顺序排序参数
	var keys []string
	for k := range params {
		if params[k] == "" || k == "sign" || k == "sign_type" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建待签名字符串
	var signStr strings.Builder
	for _, k := range keys {
		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(params[k])
		signStr.WriteString("&")
	}
	signString := strings.TrimSuffix(signStr.String(), "&") + s.config.PrivateKey

	return utils.Md5(signString)
}

type GeekPayResp struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	TradeNo   string `json:"trade_no"`
	PayURL    string `json:"payurl"`
	QrCode    string `json:"qrcode"`
	UrlScheme string `json:"urlscheme"` // 小程序跳转支付链接
}

func (s *GeekPayService) sendRequest(endpoint string, params map[string]string) (*GeekPayResp, error) {
	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}

	apiURL := fmt.Sprintf("%s/mapi.php", endpoint)
	logger.Infof(apiURL)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 取消 SSL 证书验证
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm(apiURL, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	logger.Debugf(string(body))
	if err != nil {
		return nil, err
	}

	var r GeekPayResp
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, errors.New("当前支付渠道暂不支持")
	}
	if r.Code != 1 {
		return nil, errors.New(r.Msg)
	}
	return &r, nil
}
