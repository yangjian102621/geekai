package mj

import "chatplus/core/types"

type Client interface {
	Imagine(task types.MjTask) (ImageRes, error)
	Blend(task types.MjTask) (ImageRes, error)
	SwapFace(task types.MjTask) (ImageRes, error)
	Upscale(task types.MjTask) (ImageRes, error)
	Variation(task types.MjTask) (ImageRes, error)
	QueryTask(taskId string) (QueryRes, error)
}

type ImageReq struct {
	BotType       string      `json:"botType,omitempty"`
	Prompt        string      `json:"prompt,omitempty"`
	Dimensions    string      `json:"dimensions,omitempty"`
	Base64Array   []string    `json:"base64Array,omitempty"`
	AccountFilter interface{} `json:"accountFilter,omitempty"`
	NotifyHook    string      `json:"notifyHook,omitempty"`
	State         string      `json:"state,omitempty"`
}

type ImageRes struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Properties  struct {
	} `json:"properties"`
	Result string `json:"result"`
}

type ErrRes struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type QueryRes struct {
	Action  string `json:"action"`
	Buttons []struct {
		CustomId string `json:"customId"`
		Emoji    string `json:"emoji"`
		Label    string `json:"label"`
		Style    int    `json:"style"`
		Type     int    `json:"type"`
	} `json:"buttons"`
	Description string `json:"description"`
	FailReason  string `json:"failReason"`
	FinishTime  int    `json:"finishTime"`
	Id          string `json:"id"`
	ImageUrl    string `json:"imageUrl"`
	Progress    string `json:"progress"`
	Prompt      string `json:"prompt"`
	PromptEn    string `json:"promptEn"`
	Properties  struct {
	} `json:"properties"`
	StartTime  int    `json:"startTime"`
	State      string `json:"state"`
	Status     string `json:"status"`
	SubmitTime int    `json:"submitTime"`
}
