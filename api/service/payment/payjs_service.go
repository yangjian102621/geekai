package payment

import (
	"chatplus/core/types"
	"chatplus/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type PayJS struct {
	config *types.JPayConfig
}

func NewPayJS(appConfig *types.AppConfig) *PayJS {
	return &PayJS{
		config: &appConfig.JPayConfig,
	}
}

type JPayReq struct {
	TotalFee   int    `json:"total_fee"`
	OutTradeNo string `json:"out_trade_no"`
	Subject    string `json:"body"`
	NotifyURL  string `json:"notify_url"`
}
type JPayReps struct {
	CodeUrl    string `json:"code_url"`
	OutTradeNo string `json:"out_trade_no"`
	OrderId    string `json:"payjs_order_id"`
	Qrcode     string `json:"qrcode"`
	ReturnCode int    `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	Sign       string `json:"sign"`
	TotalFee   string `json:"total_fee"`
}

func (r JPayReps) IsOK() bool {
	return r.ReturnMsg == "SUCCESS"
}

func (js *PayJS) Pay(param JPayReq) JPayReps {
	param.NotifyURL = js.config.NotifyURL
	var p = url.Values{}
	encode := utils.JsonEncode(param)
	m := make(map[string]interface{})
	_ = utils.JsonDecode(encode, &m)
	for k, v := range m {
		p.Add(k, fmt.Sprintf("%v", v))
	}
	p.Add("mchid", js.config.AppId)

	p.Add("sign", sign(p, js.config.PrivateKey))

	cli := http.Client{}
	r, err := cli.PostForm(js.config.ApiURL, p)
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

func sign(params url.Values, priKey string) string {
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
	src += "&key=" + priKey

	md5bs := md5.Sum([]byte(src))
	md5res := hex.EncodeToString(md5bs[:])
	return strings.ToUpper(md5res)
}
