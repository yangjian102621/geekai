package payment

import (
	"chatplus/core/types"
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

// Pay 执行支付请求操作
func (s *HuPiPayService) Pay(params map[string]string) (string, error) {
	data := url.Values{}
	simple := strconv.FormatInt(time.Now().Unix(), 10)
	params["appid"] = s.appId
	params["time"] = simple
	params["nonce_str"] = simple
	for k, v := range params {
		data.Add(k, v)
	}
	data.Add("hash", s.Sign(params))
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
	params["appid"] = s.appId
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
