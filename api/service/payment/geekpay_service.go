package payment

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"geekai/core/types"
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
}

// Pay 支付订单
func (s *GeekPayService) Pay(params GeekPayParams) (string, error) {
	if params.Type == "wechat" {
		params.Type = "wxpay"
	}
	p := map[string]string{
		"pid":          s.config.AppId,
		"method":       params.Method,
		"device":       params.Device,
		"type":         params.Type,
		"out_trade_no": params.OutTradeNo,
		"name":         params.Name,
		"money":        params.Money,
		"clientip":     params.ClientIP,
		"sub_openid":   params.SubOpenId,
		"sub_appid":    params.SubAppId,
		"notify_url":   s.config.NotifyURL,
		"return_url":   s.config.ReturnURL,
		"timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
	}
	sign, err := s.Sign(p)
	if err != nil {
		return "", err
	}
	p["sign"] = sign
	p["sign_type"] = "RSA"
	return s.sendRequest(s.config.ApiURL, p)
}

func (s *GeekPayService) Sign(params map[string]string) (string, error) {
	// 按字母顺序排序参数
	var keys []string
	for k := range params {
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
	signString := strings.TrimSuffix(signStr.String(), "&")

	// 使用RSA私钥签名
	block, _ := pem.Decode([]byte(s.config.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to decode private key")
	}

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	hashed := sha256.Sum256([]byte(signString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(fmt.Sprintf("failed to sign: %v", err))
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (s *GeekPayService) sendRequest(apiEndpoint string, params map[string]string) (string, error) {
	form := url.Values{}
	for k, v := range params {
		form.Add(k, v)
	}

	resp, err := http.PostForm(apiEndpoint, form)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
