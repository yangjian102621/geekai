package mj

import (
	"chatplus/core/types"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

type ApiClient struct {
	client *req.Client
	Config types.MidJourneyApiConfig
}

func NewApiClient(config types.MidJourneyApiConfig) *ApiClient {
	client := req.C().SetTimeout(10 * time.Second)
	return &ApiClient{client: client, Config: config}
}

// refer to https://open.czhanai.com/platform/doc

func (c *ApiClient) Imagine(prompt string) (error, string) {
	interactionsReq := MjAPIRequest{
		Prompt: prompt,
		Type:   "relax",
	}
	var res MjAPIResponse
	logger.Infof("request: %+v", interactionsReq)
	r, err := c.client.R().SetHeader("Auth-Token", c.Config.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(&interactionsReq).
		SetSuccessResult(&res).
		SetErrorResult(&res).
		Post(fmt.Sprintf("%s/aiapi/v1/mj/draw", c.Config.ApiURL))
	if err != nil || r.IsErrorState() {
		// read status code
		logger.Errorf("error with http request: %+v %+v %+v", err, r.StatusCode, res.Msg)
		return fmt.Errorf("error with http request: %+w %+v", err, r.Err), ""
	}
	logger.Infof("response: %+v", res)
	return nil, res.Data.PaintingSign
}

// Upscale 放大指定的图片
func (c *ApiClient) Upscale(index int, messageId string, hash string) error {

	return nil
}

// Variation  以指定的图片的视角进行变换再创作，注意需要在对应的频道中关闭 Remix 变换，否则 Variation 指令将不会生效
func (c *ApiClient) Variation(index int, messageId string, hash string) error {

	return nil
}

func (c *ApiClient) CheckStatus(taskId string) (error, *MjTaskStatusData) {
	interactionsReq := NewCheckStatusRequest(taskId)
	var res MjAPICheckStatuResponse
	r, err := c.client.R().SetHeader("Auth-Token", c.Config.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(interactionsReq).
		SetSuccessResult(&res).
		Post(fmt.Sprintf("%s/aiapi/v1/mj/task/progress", c.Config.ApiURL))
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("error with http request: %w%v", err, r.Err), nil
	}
	return nil, res.Data
}
