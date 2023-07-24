package function

import (
	"chatplus/utils"
	"errors"
	"fmt"
	"strings"
)

// 微博热搜函数实现

type FuncWeiboHot struct {
	name   string
	apiURL string
	token  string
}

func NewWeiboHot(token string) FuncWeiboHot {
	return FuncWeiboHot{name: "微博热搜", apiURL: "https://v2.alapi.cn/api/new/wbtop", token: token}
}

type WeiBoVo struct {
	resVo
	Data []struct {
		HotWord    string `json:"hot_word"`
		HotWordNum int    `json:"hot_word_num"`
		Url        string `json:"url"`
	} `json:"data"`
}

func (f FuncWeiboHot) Invoke(...interface{}) (string, error) {
	if f.token == "" {
		return "", errors.New("无效的 API Token")
	}

	url := fmt.Sprintf("%s?num=10&token=%s", f.apiURL, f.token)
	bytes, err := utils.HttpGet(url, "")
	if err != nil {
		return "", err
	}
	var res WeiBoVo
	err = utils.JsonDecode(string(bytes), &res)
	if err != nil {
		return "", err
	}

	if res.Code != 200 {
		return "", fmt.Errorf("call api fail: %s", res.Msg)
	}
	builder := make([]string, 0)
	builder = append(builder, "**新浪微博今日热搜：**")
	for i, v := range res.Data {
		builder = append(builder, fmt.Sprintf("%d、 [%s](%s) [热度：%d]", i+1, v.HotWord, v.Url, v.HotWordNum))
	}
	return strings.Join(builder, "\n\n"), nil
}

func (f FuncWeiboHot) Name() string {
	return f.name
}
