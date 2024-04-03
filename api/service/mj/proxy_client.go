package mj

import (
	"chatplus/core/types"
	"chatplus/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"io"
)

// ProxyClient MidJourney Proxy Client
type ProxyClient struct {
	Config types.MjProxyConfig
	apiURL string
}

func NewProxyClient(config types.MjProxyConfig) *ProxyClient {
	return &ProxyClient{Config: config, apiURL: config.ApiURL}
}

func (c *ProxyClient) Imagine(task types.MjTask) (ImageRes, error) {
	apiURL := fmt.Sprintf("%s/mj/submit/imagine", c.apiURL)
	prompt := fmt.Sprintf("%s %s", task.Prompt, task.Params)
	if task.NegPrompt != "" {
		prompt += fmt.Sprintf(" --no %s", task.NegPrompt)
	}
	body := ImageReq{
		Prompt:      prompt,
		Base64Array: make([]string, 0),
	}
	// 生成图片 Base64 编码
	if len(task.ImgArr) > 0 {
		imageData, err := utils.DownloadImage(task.ImgArr[0], "")
		if err != nil {
			logger.Error("error with download image: ", err)
		} else {
			body.Base64Array = append(body.Base64Array, "data:image/png;base64,"+base64.StdEncoding.EncodeToString(imageData))
		}

	}
	logger.Info("API URL: ", apiURL)
	var res ImageRes
	var errRes ErrRes
	r, err := req.C().R().
		SetHeader("mj-api-secret", c.Config.ApiKey).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&errRes).
		Post(apiURL)
	if err != nil {
		return ImageRes{}, fmt.Errorf("请求 API %s 出错：%v", apiURL, err)
	}

	if r.IsErrorState() {
		errStr, _ := io.ReadAll(r.Body)
		return ImageRes{}, fmt.Errorf("API 返回错误：%s，%v", errRes.Error.Message, string(errStr))
	}

	return res, nil
}

// Blend 融图
func (c *ProxyClient) Blend(task types.MjTask) (ImageRes, error) {
	apiURL := fmt.Sprintf("%s/mj/submit/blend", c.apiURL)
	body := ImageReq{
		Dimensions:  "SQUARE",
		Base64Array: make([]string, 0),
	}
	// 生成图片 Base64 编码
	if len(task.ImgArr) > 0 {
		for _, imgURL := range task.ImgArr {
			imageData, err := utils.DownloadImage(imgURL, "")
			if err != nil {
				logger.Error("error with download image: ", err)
			} else {
				body.Base64Array = append(body.Base64Array, "data:image/png;base64,"+base64.StdEncoding.EncodeToString(imageData))
			}
		}
	}
	var res ImageRes
	var errRes ErrRes
	r, err := req.C().R().
		SetHeader("mj-api-secret", c.Config.ApiKey).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&errRes).
		Post(apiURL)
	if err != nil {
		return ImageRes{}, fmt.Errorf("请求 API %s 出错：%v", apiURL, err)
	}

	if r.IsErrorState() {
		return ImageRes{}, fmt.Errorf("API 返回错误：%s", errRes.Error.Message)
	}

	return res, nil
}

// SwapFace 换脸
func (c *ProxyClient) SwapFace(_ types.MjTask) (ImageRes, error) {
	return ImageRes{}, errors.New("MidJourney-Proxy暂未实现该功能，请使用 MidJourney-Plus")
}

// Upscale 放大指定的图片
func (c *ProxyClient) Upscale(task types.MjTask) (ImageRes, error) {
	body := map[string]interface{}{
		"action": "UPSCALE",
		"index":  task.Index,
		"taskId": task.MessageId,
	}
	apiURL := fmt.Sprintf("%s/mj/submit/change", c.apiURL)
	var res ImageRes
	var errRes ErrRes
	r, err := req.C().R().
		SetHeader("mj-api-secret", c.Config.ApiKey).
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
func (c *ProxyClient) Variation(task types.MjTask) (ImageRes, error) {
	body := map[string]interface{}{
		"action": "VARIATION",
		"index":  task.Index,
		"taskId": task.MessageId,
	}
	apiURL := fmt.Sprintf("%s/mj/submit/change", c.apiURL)
	var res ImageRes
	var errRes ErrRes
	r, err := req.C().R().
		SetHeader("mj-api-secret", c.Config.ApiKey).
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

func (c *ProxyClient) QueryTask(taskId string) (QueryRes, error) {
	apiURL := fmt.Sprintf("%s/mj/task/%s/fetch", c.apiURL, taskId)
	var res QueryRes
	r, err := req.C().R().SetHeader("mj-api-secret", c.Config.ApiKey).
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

var _ Client = &ProxyClient{}
