package function

import (
	"chatplus/core/types"
	"chatplus/utils"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"time"
)

// AI 绘画函数

type FuncMidJourney struct {
	name   string
	config types.ChatPlusExtConfig
	client *req.Client
}

func NewMidJourneyFunc(config types.ChatPlusExtConfig) FuncMidJourney {
	return FuncMidJourney{
		name:   "MidJourney AI 绘画",
		config: config,
		client: req.C().SetTimeout(10 * time.Second)}
}

func (f FuncMidJourney) Invoke(params map[string]interface{}) (string, error) {
	if f.config.Token == "" {
		return "", errors.New("无效的 API Token")
	}

	logger.Infof("MJ 绘画参数：%+v", params)
	prompt := utils.InterfaceToString(params["prompt"])
	if !utils.IsEmptyValue(params["ar"]) {
		prompt = prompt + fmt.Sprintf(" --ar %v", params["ar"])
		delete(params, "ar")
	}
	prompt = prompt + " --niji 5"
	var res types.BizVo
	r, err := f.client.R().
		SetHeader("Authorization", f.config.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetSuccessResult(&res).Post(f.config.ApiURL)
	if err != nil || r.IsErrorState() {
		return "", fmt.Errorf("%v%v", r.String(), err)
	}

	if res.Code != types.Success {
		return "", errors.New(res.Message)
	}

	return prompt, nil
}

func (f FuncMidJourney) Name() string {
	return f.name
}

var _ Function = &FuncMidJourney{}
