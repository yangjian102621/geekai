package payment

import (
	"chatplus/core/types"
	"chatplus/utils"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HuPiPayService struct {
	appId     string
	appSecret string
	apiURL    string
}

func NewHuPiPay(config *types.AppConfig) *HuPiPayService {
	return &HuPiPayService{
		appId:     config.HuPiPayConfig.AppId,
		appSecret: config.HuPiPayConfig.AppSecret,
		apiURL:    config.HuPiPayConfig.ApiURL,
	}
}

type HuPiPayReq struct {
	AppId        string `json:"appid"`
	Version      string `json:"version"`
	TradeOrderId string `json:"trade_order_id"`
	TotalFee     string `json:"total_fee"`
	Title        string `json:"title"`
	NotifyURL    string `json:"notify_url"`
	ReturnURL    string `json:"return_url"`
	WapName      string `json:"wap_name"`
	CallbackURL  string `json:"callback_url"`
	Time         string `json:"time"`
	NonceStr     string `json:"nonce_str"`
}

type HuPiResp struct {
	Openid    interface{} `json:"openid"`
	UrlQrcode string      `json:"url_qrcode"`
	URL       string      `json:"url"`
	ErrCode   int         `json:"errcode"`
	ErrMsg    string      `json:"errmsg,omitempty"`
}

// Pay 执行支付请求操作
func (s *HuPiPayService) Pay(params HuPiPayReq) (HuPiResp, error) {
	data := url.Values{}
	simple := strconv.FormatInt(time.Now().Unix(), 10)
	params.AppId = s.appId
	params.Time = simple
	params.NonceStr = simple
	encode := utils.JsonEncode(params)
	m := make(map[string]string)
	_ = utils.JsonDecode(encode, &m)
	for k, v := range m {
		data.Add(k, fmt.Sprintf("%v", v))
	}
	encode = utils.JsonEncode(params)
	m = make(map[string]string)
	_ = utils.JsonDecode(encode, &m)
	data.Add("hash", s.sign(m))
	apiURL := fmt.Sprintf("%s/payment/do.html", s.apiURL)
	logger.Info(apiURL)
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return HuPiResp{}, fmt.Errorf("error with requst api: %v", err)
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return HuPiResp{}, fmt.Errorf("error with reading response: %v", err)
	}

	var res HuPiResp
	err = utils.JsonDecode(string(all), &res)
	if err != nil {
		return HuPiResp{}, fmt.Errorf("error with decode payment result: %v", err)
	}

	if res.ErrCode != 0 {
		return HuPiResp{}, fmt.Errorf("error with generate pay url: %s", res.ErrMsg)
	}

	return res, nil
}

// Sign 签名方法
func (s *HuPiPayService) sign(params map[string]string) string {
	var data string
	keys := make([]string, 0, 0)
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	//拼接
	for _, k := range keys {
		data = fmt.Sprintf("%s%s=%s&", data, k, params[k])
	}
	data = strings.Trim(data, "&")
	data = fmt.Sprintf("%s%s", data, s.appSecret)
	m := md5.New()
	m.Write([]byte(data))
	sign := fmt.Sprintf("%x", m.Sum(nil))
	return sign
}

// Check 校验订单状态
func (s *HuPiPayService) Check(tradeNo string) error {
	data := url.Values{}
	data.Add("appid", s.appId)
	data.Add("open_order_id", tradeNo)
	stamp := strconv.FormatInt(time.Now().Unix(), 10)
	data.Add("time", stamp)
	data.Add("nonce_str", stamp)
	// 生成签名
	encode := utils.JsonEncode(data)
	m := make(map[string]string)
	err := utils.JsonDecode(encode, &m)
	data.Add("sign", s.sign(m))

	apiURL := fmt.Sprintf("%s/payment/query.html", s.apiURL)
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("error with http reqeust: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error with reading response: %v", err)
	}

	var r struct {
		ErrCode int `json:"errcode"`
		Data    struct {
			Status      string `json:"status"`
			OpenOrderId string `json:"open_order_id"`
		} `json:"data,omitempty"`
		ErrMsg string `json:"errmsg"`
		Hash   string `json:"hash"`
	}
	err = utils.JsonDecode(string(body), &r)
	if err != nil {
		return fmt.Errorf("error with decode response: %v", err)
	}

	if r.ErrCode == 0 && r.Data.Status == "OD" {
		return nil
	} else {
		return errors.New("order not paid")
	}
}
