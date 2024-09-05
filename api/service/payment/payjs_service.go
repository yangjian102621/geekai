package payment

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type JPayService struct {
	config *types.JPayConfig
}

func NewJPayService(appConfig *types.AppConfig) *JPayService {
	return &JPayService{
		config: &appConfig.JPayConfig,
	}
}

type JPayReq struct {
	TotalFee   int    `json:"total_fee"`
	OutTradeNo string `json:"out_trade_no"`
	Subject    string `json:"body"`
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"callback_url"`
}
type JPayReps struct {
	OutTradeNo string `json:"out_trade_no"`
	OrderId    string `json:"payjs_order_id"`
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	Sign       string `json:"Sign"`
	TotalFee   string `json:"total_fee"`
	CodeUrl    string `json:"code_url,omitempty"`
	Qrcode     string `json:"qrcode,omitempty"`
}

func (r JPayReps) IsOK() bool {
	return r.ReturnMsg == "SUCCESS"
}

func (js *JPayService) Pay(param JPayReq) JPayReps {
	param.NotifyURL = js.config.NotifyURL
	var p = url.Values{}
	encode := utils.JsonEncode(param)
	m := make(map[string]interface{})
	_ = utils.JsonDecode(encode, &m)
	for k, v := range m {
		p.Add(k, fmt.Sprintf("%v", v))
	}
	p.Add("mchid", js.config.AppId)

	p.Add("sign", js.sign(p))

	cli := http.Client{}
	apiURL := fmt.Sprintf("%s/api/native", js.config.ApiURL)
	r, err := cli.PostForm(apiURL, p)
	if err != nil {
		return JPayReps{ReturnMsg: err.Error()}
	}
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		return JPayReps{ReturnMsg: err.Error()}
	}

	var data JPayReps
	err = utils.JsonDecode(string(bs), &data)
	if err != nil {
		return JPayReps{ReturnMsg: err.Error()}
	}
	return data
}

func (js *JPayService) PayH5(p url.Values) string {
	p.Add("mchid", js.config.AppId)
	p.Add("sign", js.sign(p))
	return fmt.Sprintf("%s/api/cashier?%s", js.config.ApiURL, p.Encode())
}

func (js *JPayService) sign(params url.Values) string {
	params.Del(`sign`)
	var keys = make([]string, 0, 0)
	for key := range params {
		if params.Get(key) != `` {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(params.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	src += "&key=" + js.config.PrivateKey

	md5bs := md5.Sum([]byte(src))
	md5res := hex.EncodeToString(md5bs[:])
	return strings.ToUpper(md5res)
}

// TradeVerify 查询订单支付状态
// @param tradeNo 支付平台交易 ID
func (js *JPayService) TradeVerify(tradeNo string) error {
	apiURL := fmt.Sprintf("%s/api/check", js.config.ApiURL)
	params := url.Values{}
	params.Add("payjs_order_id", tradeNo)
	params.Add("sign", js.sign(params))
	data := strings.NewReader(params.Encode())
	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", data)
	if err != nil {
		return fmt.Errorf("error with http reqeust: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error with reading response: %v", err)
	}

	var r struct {
		ReturnCode int `json:"return_code"`
		Status     int `json:"status"`
	}
	err = utils.JsonDecode(string(body), &r)
	if err != nil {
		return fmt.Errorf("error with decode response: %v", err)
	}

	if r.ReturnCode == 1 && r.Status == 1 {
		return nil
	} else {
		logger.Errorf("PayJs 支付验证响应：%s", string(body))
		return errors.New("order not paid")
	}
}
