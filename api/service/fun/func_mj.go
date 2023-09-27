package fun

import (
	"chatplus/core/types"
	"chatplus/service/mj"
	"chatplus/utils"
)

// AI 绘画函数

type FuncMidJourney struct {
	name    string
	service *mj.Service
}

func NewMidJourneyFunc(mjService *mj.Service) FuncMidJourney {
	return FuncMidJourney{
		name:    "MidJourney AI 绘画",
		service: mjService}
}

func (f FuncMidJourney) Invoke(params map[string]interface{}) (string, error) {
	logger.Infof("MJ 绘画参数：%+v", params)
	prompt := utils.InterfaceToString(params["prompt"])
	f.service.PushTask(types.MjTask{
		SessionId: utils.InterfaceToString(params["session_id"]),
		Src:       types.TaskSrcChat,
		Type:      types.TaskImage,
		Prompt:    prompt,
		UserId:    utils.IntValue(utils.InterfaceToString(params["user_id"]), 0),
		RoleId:    utils.IntValue(utils.InterfaceToString(params["role_id"]), 0),
		Icon:      utils.InterfaceToString(params["icon"]),
		ChatId:    utils.InterfaceToString(params["chat_id"]),
	})
	return prompt, nil
}

func (f FuncMidJourney) Name() string {
	return f.name
}

var _ Function = &FuncMidJourney{}
