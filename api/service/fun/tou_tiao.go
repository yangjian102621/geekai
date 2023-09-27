package fun

import (
	"chatplus/core/types"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"strings"
	"time"
)

// 今日头条函数实现

type FuncHeadlines struct {
	name   string
	config types.ChatPlusApiConfig
	client *req.Client
}

func NewHeadLines(config types.ChatPlusApiConfig) FuncHeadlines {
	return FuncHeadlines{
		name:   "今日头条",
		config: config,
		client: req.C().SetTimeout(10 * time.Second)}
}

func (f FuncHeadlines) Invoke(map[string]interface{}) (string, error) {
	if f.config.Token == "" {
		return "", errors.New("无效的 API Token")
	}

	url := fmt.Sprintf("%s/api/headline/fetch", f.config.ApiURL)
	var res resVo
	r, err := f.client.R().
		SetHeader("AppId", f.config.AppId).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", f.config.Token)).
		SetSuccessResult(&res).Get(url)
	if err != nil || r.IsErrorState() {
		return "", fmt.Errorf("%v%v", err, r.Err)
	}

	if res.Code != types.Success {
		return "", errors.New(res.Message)
	}

	builder := make([]string, 0)
	builder = append(builder, fmt.Sprintf("**%s**，最新更新：%s", res.Data.Title, res.Data.UpdatedAt))
	for i, v := range res.Data.Items {
		builder = append(builder, fmt.Sprintf("%d、 [%s](%s) [%s]", i+1, v.Title, v.Url, v.Remark))
	}
	return strings.Join(builder, "\n\n"), nil
}

func (f FuncHeadlines) Name() string {
	return f.name
}

var _ Function = &FuncHeadlines{}
