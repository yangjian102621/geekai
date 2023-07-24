package function

import (
	"chatplus/utils"
	"errors"
	"fmt"
	"strings"
)

// 每日早报函数实现

type FuncZaoBao struct {
	name   string
	apiURL string
	token  string
}

func NewZaoBao(token string) FuncZaoBao {
	return FuncZaoBao{name: "每日早报", apiURL: "https://v2.alapi.cn/api/zaobao", token: token}
}

type ZaoBaoVo struct {
	resVo
	Data struct {
		Date  string   `json:"date"`
		News  []string `json:"news"`
		WeiYu string   `json:"weiyu"`
	} `json:"data"`
}

func (f FuncZaoBao) Invoke(...interface{}) (string, error) {
	if f.token == "" {
		return "", errors.New("无效的 API Token")
	}

	url := fmt.Sprintf("%s?format=json&token=%s", f.apiURL, f.token)
	bytes, err := utils.HttpGet(url, "")
	if err != nil {
		return "", err
	}
	var res ZaoBaoVo
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

func (f FuncZaoBao) Name() string {
	return f.name
}
