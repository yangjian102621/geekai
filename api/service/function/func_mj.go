package function

import (
	"chatplus/service"
	"chatplus/utils"
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
