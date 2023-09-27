package fun

import (
	"chatplus/core/types"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"strings"
	"time"
)

// 每日早报函数实现

type FuncZaoBao struct {
	name   string
	config types.ChatPlusApiConfig
	client *req.Client
}

func NewZaoBao(config types.ChatPlusApiConfig) FuncZaoBao {
	return FuncZaoBao{
		name:   "每日早报",
		config: config,
		client: req.C().SetTimeout(10 * time.Second)}
}

func (f FuncZaoBao) Invoke(map[string]interface{}) (string, error) {
	if f.config.Token == "" {
		return "", errors.New("无效的 API Token")
	}

	url := fmt.Sprintf("%s/api/zaobao/fetch", f.config.ApiURL)
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
	builder = append(builder, fmt.Sprintf("**%s 早报：**", res.Data.UpdatedAt))
	for _, v := range res.Data.Items {
		builder = append(builder, v.Title)
	}
	builder = append(builder, fmt.Sprintf("%s", res.Data.Title))
	return strings.Join(builder, "\n\n"), nil
}

func (f FuncZaoBao) Name() string {
	return f.name
}

var _ Function = &FuncZaoBao{}
