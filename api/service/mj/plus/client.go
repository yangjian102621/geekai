package plus

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"errors"
	"fmt"
	"io"

	"github.com/imroc/req/v3"
)

var logger = logger2.GetLogger()

// Client MidJourney Plus Client
type Client struct {
	Config types.MidJourneyPlusConfig
}

func NewClient(config types.MidJourneyPlusConfig) *Client {
	return &Client{Config: config}
}

type ImageReq struct {
	BotType       string        `json:"botType"`
	Prompt        string        `json:"prompt"`
	Base64Array   []interface{} `json:"base64Array,omitempty"`
	AccountFilter struct {
		InstanceId          string        `json:"instanceId"`
		Modes               []interface{} `json:"modes"`
		Remix               bool          `json:"remix"`
		RemixAutoConsidered bool          `json:"remixAutoConsidered"`
	} `json:"accountFilter,omitempty"`
	NotifyHook string `json:"notifyHook"`
	State      string `json:"state,omitempty"`
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

func (c *Client) Imagine(prompt string) (ImageRes, error) {
	apiURL := fmt.Sprintf("%s/mj-fast/mj/submit/imagine", c.Config.ApiURL)
	body := ImageReq{
		BotType:    "MID_JOURNEY",
		Prompt:     prompt,
		NotifyHook: c.Config.NotifyURL,
	}
	var res ImageRes
	var errRes ErrRes
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+c.Config.ApiKey).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&errRes).
		Post(apiURL)
	if err != nil {
		errStr, _ := io.ReadAll(r.Body)
		return ImageRes{}, fmt.Errorf("请求 API 出错：%v，%v", err, string(errStr))
	}

	if r.IsErrorState() {
		return ImageRes{}, fmt.Errorf("API 返回错误：%s", errRes.Error.Message)
	}

	return res, nil
}

// Upscale 放大指定的图片
func (c *Client) Upscale(index int, messageId string, hash string) (ImageRes, error) {
	body := map[string]string{
		"customId":   fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, hash),
		"taskId":     messageId,
		"notifyHook": c.Config.NotifyURL,
	}
	apiURL := fmt.Sprintf("%s/mj/submit/action", c.Config.ApiURL)
	var res ImageRes
	var errRes ErrRes
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+c.Config.ApiKey).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&errRes).
		Post(apiURL)
	if err != nil {
		return ImageRes{}, fmt.Errorf("请求 API 出错：%v", err)
	}

	if r.IsErrorState() {
		return ImageRes{}, fmt.Errorf("API 返回错误：%s", errRes.Error.Message)
	}

	return res, nil
}

// Variation  以指定的图片的视角进行变换再创作，注意需要在对应的频道中关闭 Remix 变换，否则 Variation 指令将不会生效
func (c *Client) Variation(index int, messageId string, hash string) (ImageRes, error) {
	body := map[string]string{
		"customId":   fmt.Sprintf("MJ::JOB::variation::%d::%s", index, hash),
		"taskId":     messageId,
		"notifyHook": c.Config.NotifyURL,
	}
	apiURL := fmt.Sprintf("%s/mj/submit/action", c.Config.ApiURL)
	var res ImageRes
	var errRes ErrRes
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+c.Config.ApiKey).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&errRes).
		Post(apiURL)
	if err != nil {
		return ImageRes{}, fmt.Errorf("请求 API 出错：%v", err)
	}

	if r.IsErrorState() {
		return ImageRes{}, fmt.Errorf("API 返回错误：%s", errRes.Error.Message)
	}

	return res, nil
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

func (c *Client) QueryTask(taskId string) (QueryRes, error) {
	apiURL := fmt.Sprintf("%s/mj/task/%s/fetch", c.Config.ApiURL, taskId)
	var res QueryRes
	r, err := req.C().R().SetHeader("Authorization", "Bearer "+c.Config.ApiKey).
		SetSuccessResult(&res).
		Get(apiURL)

	if err != nil {
		return QueryRes{}, err
	}

	if r.IsErrorState() {
		return QueryRes{}, errors.New("error status:" + r.Status)
	}

	return res, nil
}
