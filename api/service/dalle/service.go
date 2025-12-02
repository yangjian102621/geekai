package dalle

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"
	"io"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

// DALL-E 绘画服务

type Service struct {
	httpClient    *req.Client
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	taskQueue     *store.RedisQueue
	userService   *service.UserService
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, redisCli *redis.Client, userService *service.UserService) *Service {
	return &Service{
		httpClient:    req.C().SetTimeout(time.Minute * 3),
		db:            db,
		taskQueue:     store.NewRedisQueue("DallE_Task_Queue", redisCli),
		uploadManager: manager,
		userService:   userService,
	}
}

// PushTask push a new mj task in to task queue
func (s *Service) PushTask(task types.DallTask) {
	logger.Infof("add a new DALL-E task to the task list: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) Run() {
	// 将数据库中未提交的任务加载到队列
	var jobs []model.DallJob
	s.db.Where("progress", 0).Find(&jobs)
	for _, v := range jobs {
		var task types.DallTask
		err := utils.JsonDecode(v.TaskInfo, &task)
		if err != nil {
			logger.Errorf("decode task info with error: %v", err)
			continue
		}
		task.Id = v.Id
		s.PushTask(task)
	}

	logger.Info("Starting DALL-E job consumer...")
	go func() {
		for {
			var task types.DallTask
			err := s.taskQueue.LPop(&task)
			if err != nil {
				logger.Errorf("taking task with error: %v", err)
				continue
			}
			logger.Infof("handle a new DALL-E task: %+v", task)
			go func() {
				_, err = s.Image(task, false)
				if err != nil {
					logger.Errorf("error with image task: %v", err)
					s.db.Model(&model.DallJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
						"progress": service.FailTaskProgress,
						"err_msg":  err.Error(),
					})
				}
			}()
		}
	}()
}

type imgReq struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	N       int    `json:"n,omitempty"`
	Size    string `json:"size,omitempty"`
	Quality string `json:"quality,omitempty"`
	Style   string `json:"style,omitempty"`
}

type imgRes struct {
	Created int64 `json:"created"`
	Data    []struct {
		RevisedPrompt string `json:"revised_prompt,omitempty"`
		Url           string `json:"url,omitempty"`
		B64Json       string `json:"b64_json,omitempty"`
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

func (s *Service) Image(task types.DallTask, sync bool) (string, error) {
	logger.Debugf("绘画参数：%+v", task)
	prompt := task.Prompt
	// translate prompt
	if utils.HasChinese(prompt) {
		content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.TranslatePromptTemplate, prompt), task.TranslateModelId)
		if err == nil {
			prompt = content
			logger.Debugf("重写后提示词：%s", prompt)
		}
	}

	var chatModel model.ChatModel
	s.db.Where("id = ?", task.ModelId).First(&chatModel)

	// get image generation API KEY
	var apiKey model.ApiKey
	session := s.db.Where("enabled", true)
	if chatModel.KeyId > 0 {
		session = session.Where("id = ?", chatModel.KeyId)
	} else {
		session = session.Where("type = ?", "dalle")
	}
	err := session.Order("last_used_at ASC").First(&apiKey).Error
	if err != nil {
		return "", fmt.Errorf("no available Image Generation api key: %v", err)
	}

	var res imgRes
	var errRes ErrRes
	if len(apiKey.ProxyURL) > 5 {
		s.httpClient.SetProxyURL(apiKey.ProxyURL).R()
	}
	apiURL := fmt.Sprintf("%s/v1/images/generations", apiKey.ApiURL)
	reqBody := imgReq{
		Model:   chatModel.Value,
		Prompt:  prompt,
		N:       1,
		Size:    task.Size,
		Style:   task.Style,
		Quality: task.Quality,
	}
	logger.Infof("Channel:%s, API KEY:%s, BODY: %+v", apiURL, apiKey.Value, reqBody)
	r, err := s.httpClient.R().SetHeader("Body-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(reqBody).
		SetErrorResult(&errRes).
		SetSuccessResult(&res).
		Post(apiURL)
	if err != nil {
		logger.Errorf("error with send request: %v", err)
		return "", fmt.Errorf("error with send request: %v", err)
	}

	if r.IsErrorState() {
		logger.Errorf("error with send request, status: %s, %+v", r.Status, errRes.Error)
		return "", fmt.Errorf("error with send request, status: %s, %+v", r.Status, errRes.Error)
	}

	all, _ := io.ReadAll(r.Body)
	logger.Debugf("response: %+v", string(all))

	// update the api key last use time
	s.db.Model(&apiKey).UpdateColumn("last_used_at", time.Now().Unix())
	var imgURL string
	var data = map[string]interface{}{
		"progress": 100,
		"prompt":   prompt,
	}
	// 如果返回的是base64，则需要上传到oss
	if res.Data[0].B64Json != "" {
		imgURL, err = s.uploadManager.GetUploadHandler().PutBase64(res.Data[0].B64Json)
		if err != nil {
			return "", fmt.Errorf("error with upload image: %v", err)
		}
		logger.Infof("upload image to oss: %s", imgURL)
		data["img_url"] = imgURL
	} else {
		imgURL = res.Data[0].Url
	}
	data["org_url"] = imgURL
	// update task progress
	err = s.db.Model(&model.DallJob{Id: task.Id}).UpdateColumns(data).Error
	if err != nil {
		return "", fmt.Errorf("err with update database: %v", err)
	}

	var content string
	if sync {
		imgURL, err := s.downloadImage(task.Id, res.Data[0].Url)
		if err != nil {
			return "", fmt.Errorf("error with download image: %v", err)
		}
		content = fmt.Sprintf("```\n%s\n```\n下面是我为你创作的图片：\n\n![](%s)\n", prompt, imgURL)
	}

	return content, nil
}

func (s *Service) CheckTaskStatus() {
	go func() {
		logger.Info("Running DALL-E task status checking ...")
		for {
			// 检查未完成任务进度
			var jobs []model.DallJob
			s.db.Where("progress < ?", 100).Find(&jobs)
			for _, job := range jobs {
				// 超时的任务标记为失败
				if time.Since(job.CreatedAt) > time.Minute*10 {
					job.Progress = service.FailTaskProgress
					job.ErrMsg = "任务超时"
					s.db.Updates(&job)
				}
			}

			// 找出失败的任务，并恢复其扣减算力
			s.db.Where("progress", service.FailTaskProgress).Where("power > ?", 0).Find(&jobs)
			for _, job := range jobs {
				var task types.DallTask
				err := utils.JsonDecode(job.TaskInfo, &task)
				if err != nil {
					continue
				}
				err = s.userService.IncreasePower(int(job.UserId), job.Power, model.PowerLog{
					Type:   types.PowerRefund,
					Model:  task.ModelName,
					Remark: fmt.Sprintf("任务失败，退回算力。任务ID：%d，Err: %s", job.Id, job.ErrMsg),
				})
				if err != nil {
					continue
				}
				// 更新任务状态
				s.db.Model(&job).UpdateColumn("power", 0)
			}
			time.Sleep(time.Second * 10)
		}
	}()
}

func (s *Service) DownloadImages() {
	go func() {
		var items []model.DallJob
		for {
			res := s.db.Where("img_url = ? AND progress = ?", "", 100).Find(&items)
			if res.Error != nil {
				continue
			}

			// download images
			for _, v := range items {
				if v.OrgURL == "" {
					continue
				}

				logger.Infof("try to download image: %s", v.OrgURL)
				imgURL, err := s.downloadImage(v.Id, v.OrgURL)
				if err != nil {
					logger.Error("error with download image: %s, error: %v", imgURL, err)
					continue
				} else {
					logger.Infof("download image %s successfully.", v.OrgURL)
				}

			}

			time.Sleep(time.Second * 5)
		}
	}()
}

func (s *Service) downloadImage(jobId uint, orgURL string) (string, error) {
	// sava image
	imgURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(orgURL, false)
	if err != nil {
		return "", err
	}

	// update img_url
	res := s.db.Model(&model.DallJob{Id: jobId}).UpdateColumn("img_url", imgURL)
	if res.Error != nil {
		return "", err
	}
	return imgURL, nil
}
