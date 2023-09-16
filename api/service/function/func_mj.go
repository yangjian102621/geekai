package function

import (
	"chatplus/service"
	"chatplus/utils"
	"fmt"
)

// AI 绘画函数

type FuncMidJourney struct {
	name    string
	service *service.MjService
}

func NewMidJourneyFunc(mjService *service.MjService) FuncMidJourney {
	return FuncMidJourney{
		name:    "MidJourney AI 绘画",
		service: mjService}
}

func (f FuncMidJourney) Invoke(params map[string]interface{}) (string, error) {
	logger.Infof("MJ 绘画参数：%+v", params)
	prompt := utils.InterfaceToString(params["prompt"])
	if !utils.IsEmptyValue(params["ar"]) {
		prompt = fmt.Sprintf("%s --ar %s", prompt, params["ar"])
		delete(params, "ar")
	}
	if !utils.IsEmptyValue(params["s"]) {
		prompt = fmt.Sprintf("%s --s %s", prompt, params["s"])
		delete(params, "s")
	}
	if !utils.IsEmptyValue(params["seed"]) {
		prompt = fmt.Sprintf("%s --seed %s", prompt, params["seed"])
		delete(params, "seed")
	}
	if !utils.IsEmptyValue(params["no"]) {
		prompt = fmt.Sprintf("%s --no %s", prompt, params["no"])
		delete(params, "no")
	}
	if !utils.IsEmptyValue(params["niji"]) {
		prompt = fmt.Sprintf("%s --niji %s", prompt, params["niji"])
		delete(params, "niji")
	} else {
		prompt = prompt + " --v 5.2"
	}

	f.service.PushTask(service.MjTask{
		Id:     utils.InterfaceToString(params["session_id"]),
		Src:    service.TaskSrcChat,
		Type:   service.Image,
		Prompt: prompt,
		UserId: utils.IntValue(utils.InterfaceToString(params["user_id"]), 0),
		RoleId: utils.IntValue(utils.InterfaceToString(params["role_id"]), 0),
		Icon:   utils.InterfaceToString(params["icon"]),
		ChatId: utils.InterfaceToString(params["chat_id"]),
	})
	return prompt, nil
}

func (f FuncMidJourney) Name() string {
	return f.name
}

var _ Function = &FuncMidJourney{}
