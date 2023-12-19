package fun

import (
	"chatplus/core/types"
	"chatplus/service/oss"
	"chatplus/store/model"
	"chatplus/utils"
	"fmt"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

// AI 绘画函数

type FuncImage struct {
	name          string
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	proxyURL      string
}

func NewImageFunc(db *gorm.DB, manager *oss.UploaderManager, config *types.AppConfig) FuncImage {
	return FuncImage{
		db:            db,
		name:          "DALL-E3 绘画",
		uploadManager: manager,
		proxyURL:      config.ProxyURL,
	}
}

type imgReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type imgRes struct {
	Created int64 `json:"created"`
	Data    []struct {
		RevisedPrompt string `json:"revised_prompt"`
		Url           string `json:"url"`
	} `json:"data"`
}

type ErrRes struct {
	Error struct {
		Code    interface{} `json:"code"`
		Message string      `json:"message"`
		Param   interface{} `json:"param"`
		Type    string      `json:"type"`
	} `json:"error"`
}

func (f FuncImage) Invoke(params map[string]interface{}) (string, error) {
	logger.Infof("绘画参数：%+v", params)
	prompt := utils.InterfaceToString(params["prompt"])
	// get image generation API KEY
	var apiKey model.ApiKey
	tx := f.db.Where("platform = ? AND type = ?", types.OpenAI, "img").Order("last_used_at ASC").First(&apiKey)
	if tx.Error != nil {
		return "", fmt.Errorf("error with get generation API KEY: %v", tx.Error)
	}

	// get image generation api URL
	var conf model.Config
	var chatConfig types.ChatConfig
	tx = f.db.Where("marker", "chat").First(&conf)
	if tx.Error != nil {
		return "", fmt.Errorf("error with get chat configs: %v", tx.Error)
	}

	err := utils.JsonDecode(conf.Config, &chatConfig)
	if err != nil {
		return "", fmt.Errorf("error with decode chat config: %v", err)
	}

	apiURL := chatConfig.DallApiURL
	if utils.IsEmptyValue(apiURL) {
		apiURL = "https://api.openai.com/v1/images/generations"
	}
	imgNum := chatConfig.DallImgNum
	if imgNum <= 0 {
		imgNum = 1
	}
	var res imgRes
	var errRes ErrRes
	r, err := req.C().SetProxyURL(f.proxyURL).R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(imgReq{
			Model:  "dall-e-3",
			Prompt: prompt,
			N:      imgNum,
			Size:   "1024x1024",
		}).
		SetErrorResult(&errRes).
		SetSuccessResult(&res).Post(apiURL)
	if err != nil || r.IsErrorState() {
		return "", fmt.Errorf("error with http request: %v%v%s", err, r.Err, errRes.Error.Message)
	}
	// 存储图片
	imgURL, err := f.uploadManager.GetUploadHandler().PutImg(res.Data[0].Url, false)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %s", err.Error())
	}

	//logger.Info(imgURL)
	return fmt.Sprintf("\n\n![](%s)\n", imgURL), nil
}

func (f FuncImage) Name() string {
	return f.name
}

var _ Function = &FuncImage{}
