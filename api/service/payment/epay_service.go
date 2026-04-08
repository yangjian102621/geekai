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

// EPayService 易支付服务
type EPayService struct {
	config *types.EpayConfig
}

func NewEPayService(sysConfig *types.SystemConfig) *EPayService {
	return &EPayService{
		config: &sysConfig.Payment.Epay,
	}
}

func (s *EPayService) UpdateConfig(config *types.EpayConfig) {
	s.config = config
}

// Pay 支付订单
func (s *EPayService) Pay(params PayRequest) (string, error) {
	p := map[string]string{
		"pid":          s.config.AppId,
		"device":       params.Device,
		"type":         params.PayWay,
		"out_trade_no": params.OutTradeNo,
		"name":         params.Subject,
		"money":        params.TotalFee,
		"clientip":     params.ClientIP,
		"notify_url":   params.NotifyURL,
		"return_url":   params.ReturnURL,
		"timestamp":    fmt.Sprintf("%d", time.Now().Unix()),
	}
	p["sign"] = s.Sign(p)
	p["sign_type"] = "MD5"
	resp, err := s.sendRequest(s.config.ApiURL, p)
	if err != nil {
		return "", err
	}
	if resp.Code != 1 {
		return "", errors.New(resp.Msg)
	}
	if resp.PayURL != "" {
		return resp.PayURL, nil
	} else {
		return resp.QrCode, nil
	}
}

func (s *EPayService) Sign(params map[string]string) string {
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

func (s *EPayService) sendRequest(endpoint string, params map[string]string) (*GeekPayResp, error) {
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

func (s *EPayService) Query(outTradeNo string) (OrderInfo, error) {

	params := url.Values{}
	params.Set("act", "order")
	params.Set("pid", s.config.AppId)
	params.Set("key", s.config.PrivateKey)
	params.Set("out_trade_no", outTradeNo)

	apiURL := fmt.Sprintf("%s/api.php?%s", s.config.ApiURL, params.Encode())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(apiURL)
	if err != nil {
		return OrderInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OrderInfo{}, err
	}
	logger.Debugf(string(body))

	var result struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		Status  string `json:"status"`
		Name    string `json:"name"`
		Money   string `json:"money"`
		EndTime string `json:"endtime"`
		TradeNo string `json:"trade_no"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return OrderInfo{}, errors.New("订单查询响应解析失败")
	}
	if result.Code != 1 {
		return OrderInfo{}, errors.New(result.Msg)
	}
	logger.Debugf("订单信息：%+v", result)
	orderInfo := OrderInfo{
		OutTradeNo: outTradeNo,
		TradeId:    result.TradeNo,
		Amount:     result.Money,
		PayTime:    result.EndTime,
	}
	if result.Status == "1" {
		orderInfo.Status = Success
	} else {
		orderInfo.Status = Failure
	}
	return orderInfo, nil
}

var _ PayService = (*EPayService)(nil)
