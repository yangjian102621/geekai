package plus

import (
	"chatplus/core/types"
	logger2 "chatplus/logger"
	"chatplus/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"io"

	"github.com/gin-gonic/gin"
)

var logger = logger2.GetLogger()

// Client MidJourney Plus Client
type Client struct {
	Config types.MidJourneyPlusConfig
	apiURL string
}

func NewClient(config types.MidJourneyPlusConfig) *Client {
	var apiURL string
	if config.CdnURL != "" {
		apiURL = config.CdnURL
	} else {
		apiURL = config.ApiURL
	}
	if config.Mode == "" {
		config.Mode = "fast"
	}
	return &Client{Config: config, apiURL: apiURL}
}

type ImageReq struct {
	BotType       string   `json:"botType"`
	Prompt        string   `json:"prompt,omitempty"`
	Dimensions    string   `json:"dimensions,omitempty"`
	Base64Array   []string `json:"base64Array,omitempty"`
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

func (c *Client) Imagine(task types.MjTask) (ImageRes, error) {
	apiURL := fmt.Sprintf("%s/mj-%s/mj/submit/imagine", c.apiURL, c.Config.Mode)
	body := ImageReq{
		BotType:     "MID_JOURNEY",
		Prompt:      task.Prompt,
		NotifyHook:  c.Config.NotifyURL,
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
		logger.Errorf("API 返回：%s, API URL: %s", string(errStr), apiURL)
		return ImageRes{}, fmt.Errorf("请求 API 出错：%v", err)
	}

	if r.IsErrorState() {
		errStr, _ := io.ReadAll(r.Body)
		return ImageRes{}, fmt.Errorf("API 返回错误：%s，%v", errRes.Error.Message, string(errStr))
	}

	return res, nil
}

// Blend 融图
func (c *Client) Blend(task types.MjTask) (ImageRes, error) {
	apiURL := fmt.Sprintf("%s/mj-%s/mj/submit/blend", c.apiURL, c.Config.Mode)
	body := ImageReq{
		BotType:     "MID_JOURNEY",
		Dimensions:  "SQUARE",
		NotifyHook:  c.Config.NotifyURL,
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

// SwapFace 换脸
func (c *Client) SwapFace(task types.MjTask) (ImageRes, error) {
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
		"notifyHook": c.Config.NotifyURL,
		"state":      "",
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
func (c *Client) Upscale(task types.MjTask) (ImageRes, error) {
	body := map[string]string{
		"customId":   fmt.Sprintf("MJ::JOB::upsample::%d::%s", task.Index, task.MessageHash),
		"taskId":     task.MessageId,
		"notifyHook": c.Config.NotifyURL,
	}
	apiURL := fmt.Sprintf("%s/mj/submit/action", c.apiURL)
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
func (c *Client) Variation(task types.MjTask) (ImageRes, error) {
	body := map[string]string{
		"customId":   fmt.Sprintf("MJ::JOB::variation::%d::%s", task.Index, task.MessageHash),
		"taskId":     task.MessageId,
		"notifyHook": c.Config.NotifyURL,
	}
	apiURL := fmt.Sprintf("%s/mj/submit/action", c.apiURL)
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
	apiURL := fmt.Sprintf("%s/mj/task/%s/fetch", c.apiURL, taskId)
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
