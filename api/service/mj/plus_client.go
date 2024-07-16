package mj

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/base64"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/service"
	"geekai/utils"
	"github.com/imroc/req/v3"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// PlusClient MidJourney Plus ProxyClient
type PlusClient struct {
	Config         types.MjPlusConfig
	apiURL         string
	client         *req.Client
	licenseService *service.LicenseService
}

func NewPlusClient(config types.MjPlusConfig, licenseService *service.LicenseService) *PlusClient {
	return &PlusClient{
		Config:         config,
		apiURL:         config.ApiURL,
		client:         req.C().SetTimeout(time.Minute).SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
		licenseService: licenseService,
	}
}

func (c *PlusClient) preCheck() error {
	return c.licenseService.IsValidApiURL(c.Config.ApiURL)
}

func (c *PlusClient) Imagine(task types.MjTask) (ImageRes, error) {
	if err := c.preCheck(); err != nil {
		return ImageRes{}, err
	}

	apiURL := fmt.Sprintf("%s/mj-%s/mj/submit/imagine", c.apiURL, c.Config.Mode)
	prompt := fmt.Sprintf("%s %s", task.Prompt, task.Params)
	if task.NegPrompt != "" {
		prompt += fmt.Sprintf(" --no %s", task.NegPrompt)
	}
	body := ImageReq{
		BotType:     "MID_JOURNEY",
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
	return c.doRequest(body, apiURL)
}

// Blend 融图
func (c *PlusClient) Blend(task types.MjTask) (ImageRes, error) {
	if err := c.preCheck(); err != nil {
		return ImageRes{}, err
	}

	apiURL := fmt.Sprintf("%s/mj-%s/mj/submit/blend", c.apiURL, c.Config.Mode)
	logger.Info("API URL: ", apiURL)
	body := ImageReq{
		BotType:     "MID_JOURNEY",
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
	return c.doRequest(body, apiURL)
}

// SwapFace 换脸
func (c *PlusClient) SwapFace(task types.MjTask) (ImageRes, error) {
	if err := c.preCheck(); err != nil {
		return ImageRes{}, err
	}

	apiURL := fmt.Sprintf("%s/mj-%s/mj/insight-face/swap", c.apiURL, c.Config.Mode)
	// 生成图片 Base64 编码
	if len(task.ImgArr) != 2 {
		return ImageRes{}, errors.New("参数错误，必须上传2张图片")
	}
	var sourceBase64 string
	var targetBase64 string
	imageData, err := utils.DownloadImage(task.ImgArr[0], "")
	if err != nil {
		logger.Error("error with download image: ", err)
	} else {
		sourceBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData)
	}
	imageData, err = utils.DownloadImage(task.ImgArr[1], "")
	if err != nil {
		logger.Error("error with download image: ", err)
	} else {
		targetBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData)
	}

	body := gin.H{
		"sourceBase64": sourceBase64,
		"targetBase64": targetBase64,
		"accountFilter": gin.H{
			"instanceId": "",
		},
		"state": "",
	}
	return c.doRequest(body, apiURL)
}

// Upscale 放大指定的图片
func (c *PlusClient) Upscale(task types.MjTask) (ImageRes, error) {
	if err := c.preCheck(); err != nil {
		return ImageRes{}, err
	}

	body := map[string]string{
		"customId": fmt.Sprintf("MJ::JOB::upsample::%d::%s", task.Index, task.MessageHash),
		"taskId":   task.MessageId,
	}
	apiURL := fmt.Sprintf("%s/mj-%s/mj/submit/action", c.apiURL, c.Config.Mode)
	return c.doRequest(body, apiURL)
}

// Variation  以指定的图片的视角进行变换再创作，注意需要在对应的频道中关闭 Remix 变换，否则 Variation 指令将不会生效
func (c *PlusClient) Variation(task types.MjTask) (ImageRes, error) {
	if err := c.preCheck(); err != nil {
		return ImageRes{}, err
	}

	body := map[string]string{
		"customId": fmt.Sprintf("MJ::JOB::variation::%d::%s", task.Index, task.MessageHash),
		"taskId":   task.MessageId,
	}
	apiURL := fmt.Sprintf("%s/mj-%s/mj/submit/action", c.apiURL, c.Config.Mode)

	return c.doRequest(body, apiURL)
}

func (c *PlusClient) doRequest(body interface{}, apiURL string) (ImageRes, error) {
	var res ImageRes
	var errRes ErrRes
	logger.Info("API URL: ", apiURL)
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+c.Config.ApiKey).
		SetBody(body).
		SetSuccessResult(&res).
		SetErrorResult(&errRes).
		Post(apiURL)
	if err != nil {
		errMsg := err.Error()
		if r != nil {
			errStr, _ := io.ReadAll(r.Body)
			logger.Error("请求 API 出错：", string(errStr))
			errMsg = errMsg + " " + string(errStr)
		}
		return ImageRes{}, fmt.Errorf("请求 API 出错：%v", errMsg)
	}

	if r.IsErrorState() {
		return ImageRes{}, fmt.Errorf("API 返回错误：%s", errRes.Error.Message)
	}

	return res, nil
}

func (c *PlusClient) QueryTask(taskId string) (QueryRes, error) {
	apiURL := fmt.Sprintf("%s/mj/task/%s/fetch", c.apiURL, taskId)
	var res QueryRes
	r, err := c.client.R().SetHeader("Authorization", "Bearer "+c.Config.ApiKey).
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

var _ Client = &PlusClient{}
