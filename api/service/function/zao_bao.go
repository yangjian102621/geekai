package function

import (
	"chatplus/utils"
	"fmt"
	"strings"
)

// 每日早报函数实现

type FuncZaoBao struct {
	apiURL string
	token  string
}

func NewZaoBao(token string) *FuncZaoBao {
	return &FuncZaoBao{apiURL: "https://v2.alapi.cn/api/zaobao", token: token}
}

type resVo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Date  string   `json:"date"`
		News  []string `json:"news"`
		WeiYu string   `json:"weiyu"`
	} `json:"data"`
}

func (f *FuncZaoBao) Fetch() (string, error) {

	url := fmt.Sprintf("%s?format=json&token=%s", f.apiURL, f.token)
	bytes, err := utils.HttpGet(url, "")
	if err != nil {
		return "", err
	}
	var res resVo
	err = utils.JsonDecode(string(bytes), &res)
	if err != nil {
		return "", err
	}

	if res.Code != 200 {
		return "", fmt.Errorf("call api fail: %s", res.Msg)
	}
	builder := make([]string, 0)
	builder = append(builder, fmt.Sprintf("**%s 早报：**", res.Data.Date))
	builder = append(builder, res.Data.News...)
	builder = append(builder, fmt.Sprintf("%s", res.Data.WeiYu))
	return strings.Join(builder, "\n\n"), nil
}
