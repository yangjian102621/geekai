package function

import (
	"chatplus/utils"
	"errors"
	"fmt"
	"strings"
)

// 今日头条函数实现

type FuncHeadlines struct {
	name   string
	apiURL string
	token  string
}

func NewHeadLines(token string) FuncHeadlines {
	return FuncHeadlines{name: "今日头条", apiURL: "https://v2.alapi.cn/api/tophub/get", token: token}
}

type HeadLineVo struct {
	resVo
	Data struct {
		Name       string `json:"name"`
		LastUpdate string `json:"last_update"`
		List       []struct {
			Title string `json:"title"`
			Link  string `json:"link"`
			Other string `json:"other"`
		} `json:"list"`
	} `json:"data"`
}

func (f FuncHeadlines) Invoke(...interface{}) (string, error) {
	if f.token == "" {
		return "", errors.New("无效的 API Token")
	}

	url := fmt.Sprintf("%s?type=toutiao&token=%s", f.apiURL, f.token)
	bytes, err := utils.HttpGet(url, "")
	if err != nil {
		return "", err
	}
	var res HeadLineVo
	err = utils.JsonDecode(string(bytes), &res)
	if err != nil {
		return "", err
	}

	if res.Code != 200 {
		return "", fmt.Errorf("call api fail: %s", res.Msg)
	}
	builder := make([]string, 0)
	builder = append(builder, fmt.Sprintf("**%s**，最新更新：%s", res.Data.Name, res.Data.LastUpdate))
	for i, v := range res.Data.List {
		builder = append(builder, fmt.Sprintf("%d、 [%s](%s) [%s]", i+1, v.Title, v.Link, v.Other))
	}
	return strings.Join(builder, "\n\n"), nil
}

func (f FuncHeadlines) Name() string {
	return f.name
}
