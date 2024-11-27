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
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// Client MidJourney client
type Client struct {
	client         *req.Client
	licenseService *service.LicenseService
	db             *gorm.DB
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
	Result  string `json:"result"`
	Channel string `json:"channel,omitempty"`
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

var logger = logger2.GetLogger()

func NewClient(licenseService *service.LicenseService, db *gorm.DB) *Client {
	return &Client{
		client:         req.C().SetTimeout(time.Minute).SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"),
		licenseService: licenseService,
		db:             db,
	}
}

func (c *Client) Imagine(task types.MjTask) (ImageRes, error) {
	apiPath := fmt.Sprintf("mj-%s/mj/submit/imagine", task.Mode)
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
	return c.doRequest(body, apiPath, task.ChannelId)
}

// Blend 融图
func (c *Client) Blend(task types.MjTask) (ImageRes, error) {
	apiPath := fmt.Sprintf("mj-%s/mj/submit/blend", task.Mode)
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
	return c.doRequest(body, apiPath, task.ChannelId)
}

// SwapFace 换脸
func (c *Client) SwapFace(task types.MjTask) (ImageRes, error) {
	apiPath := fmt.Sprintf("mj-%s/mj/insight-face/swap", task.Mode)
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
	return c.doRequest(body, apiPath, task.ChannelId)
}

// Upscale 放大指定的图片
func (c *Client) Upscale(task types.MjTask) (ImageRes, error) {
	body := map[string]string{
		"customId": fmt.Sprintf("MJ::JOB::upsample::%d::%s", task.Index, task.MessageHash),
		"taskId":   task.MessageId,
	}
	apiPath := fmt.Sprintf("mj-%s/mj/submit/action", task.Mode)
	return c.doRequest(body, apiPath, task.ChannelId)
}

// Variation  以指定的图片的视角进行变换再创作，注意需要在对应的频道中关闭 Remix 变换，否则 Variation 指令将不会生效
func (c *Client) Variation(task types.MjTask) (ImageRes, error) {
	body := map[string]string{
		"customId": fmt.Sprintf("MJ::JOB::variation::%d::%s", task.Index, task.MessageHash),
		"taskId":   task.MessageId,
	}
	apiPath := fmt.Sprintf("mj-%s/mj/submit/action", task.Mode)

	return c.doRequest(body, apiPath, task.ChannelId)
}

func (c *Client) doRequest(body interface{}, apiPath string, channel string) (ImageRes, error) {
	var res ImageRes
	session := c.db.Session(&gorm.Session{}).Where("type", "mj").Where("enabled", true)
	if channel != "" {
		session = session.Where("api_url", channel)
	}

	var apiKey model.ApiKey
	err := session.Order("last_used_at ASC").First(&apiKey).Error
	if err != nil {
		return ImageRes{}, fmt.Errorf("no available MidJourney api key: %v", err)
	}

	if err = c.licenseService.IsValidApiURL(apiKey.ApiURL); err != nil {
		return ImageRes{}, err
	}

	apiURL := fmt.Sprintf("%s/%s", apiKey.ApiURL, apiPath)
	logger.Info("API URL: ", apiURL)
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(body).
		SetSuccessResult(&res).
		Post(apiURL)
	if err != nil {
		return ImageRes{}, fmt.Errorf("请求 API 出错：%v", err)
	}

	if r.IsErrorState() {
		errMsg, _ := io.ReadAll(r.Body)
		return ImageRes{}, fmt.Errorf("API 返回错误：%s", string(errMsg))
	}

	// update the api key last used time
	if err = c.db.Model(&apiKey).Update("last_used_at", time.Now().Unix()).Error; err != nil {
		logger.Error("update api key last used time error: ", err)
	}
	res.Channel = apiKey.ApiURL
	return res, nil
}

func (c *Client) QueryTask(taskId string, channel string) (QueryRes, error) {
	var apiKey model.ApiKey
	err := c.db.Where("type", "mj").Where("enabled", true).Where("api_url", channel).First(&apiKey).Error
	if err != nil {
		return QueryRes{}, fmt.Errorf("no available MidJourney api key: %v", err)
	}
	apiURL := fmt.Sprintf("%s/mj/task/%s/fetch", apiKey.ApiURL, taskId)
	var res QueryRes
	r, err := c.client.R().SetHeader("Authorization", "Bearer "+apiKey.Value).
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
