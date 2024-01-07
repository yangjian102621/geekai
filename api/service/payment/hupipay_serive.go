package payment

import (
	"chatplus/core/types"
	"chatplus/utils"
	"crypto/md5"
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
	host      string
}

func NewHuPiPay(config *types.AppConfig) *HuPiPayService {
	return &HuPiPayService{
		appId:     config.HuPiPayConfig.AppId,
		appSecret: config.HuPiPayConfig.AppSecret,
		host:      config.HuPiPayConfig.PayURL,
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

// Pay 执行支付请求操作
func (s *HuPiPayService) Pay(params HuPiPayReq) (string, error) {
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
	data.Add("hash", s.Sign(m))
	resp, err := http.PostForm(s.host, data)
	if err != nil {
		return "error", err
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return "error", err
	}
	return string(all), err
}

// Sign 签名方法
func (s *HuPiPayService) Sign(params map[string]string) string {
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
