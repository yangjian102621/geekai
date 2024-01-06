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
	Body       string `json:"body"`
	NotifyURL  string `json:"notify_url"`
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

func (pj *PayJS) Pay(param JPayReq) (string, error) {
	var p = url.Values{}
	encode := utils.JsonEncode(param)
	m := make(map[string]interface{})
	_ = utils.JsonDecode(encode, &m)
	for k, v := range m {
		p.Add(k, fmt.Sprintf("%v", v))
	}
	p.Add("mchid", pj.config.AppId)

	p.Add("sign", sign(p, pj.config.PrivateKey))

	cli := http.Client{}
	r, err := cli.PostForm(pj.config.ApiURL, p)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
