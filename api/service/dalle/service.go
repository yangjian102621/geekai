package dalle

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"
	"github.com/go-redis/redis/v8"
	"time"

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
	notifyQueue   *store.RedisQueue
	Clients       *types.LMap[uint, *types.WsClient] // UserId => Client
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, redisCli *redis.Client) *Service {
	return &Service{
		httpClient:    req.C().SetTimeout(time.Minute * 3),
		db:            db,
		taskQueue:     store.NewRedisQueue("DallE_Task_Queue", redisCli),
		notifyQueue:   store.NewRedisQueue("DallE_Notify_Queue", redisCli),
		Clients:       types.NewLMap[uint, *types.WsClient](),
		uploadManager: manager,
	}
}

// PushTask push a new mj task in to task queue
func (s *Service) PushTask(task types.DallTask) {
	logger.Infof("add a new DALL-E task to the task list: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) Run() {
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
			_, err = s.Image(task, false)
			if err != nil {
				logger.Errorf("error with image task: %v", err)
				s.db.Model(&model.DallJob{Id: task.JobId}).UpdateColumns(map[string]interface{}{
					"progress": service.FailTaskProgress,
					"err_msg":  err.Error(),
				})
				s.notifyQueue.RPush(service.NotifyMessage{UserId: int(task.UserId), JobId: int(task.JobId), Message: service.TaskStatusFailed})
			}
		}
	}()
}

type imgReq struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	N       int    `json:"n"`
	Size    string `json:"size"`
	Quality string `json:"quality"`
	Style   string `json:"style"`
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

func (s *Service) Image(task types.DallTask, sync bool) (string, error) {
	logger.Debugf("绘画参数：%+v", task)
	prompt := task.Prompt
	// translate prompt
	if utils.HasChinese(prompt) {
		content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.RewritePromptTemplate, prompt), "gpt-4o-mini")
		if err == nil {
			prompt = content
			logger.Debugf("重写后提示词：%s", prompt)
		}
	}

	var user model.User
	s.db.Where("id", task.UserId).First(&user)
	if user.Power < task.Power {
		return "", errors.New("insufficient of power")
	}

	// 更新用户算力
	tx := s.db.Model(&model.User{}).Where("id", user.Id).UpdateColumn("power", gorm.Expr("power - ?", task.Power))
	// 记录算力变化日志
	if tx.Error == nil && tx.RowsAffected > 0 {
		var u model.User
		s.db.Where("id", user.Id).First(&u)
		s.db.Create(&model.PowerLog{
			UserId:    user.Id,
			Username:  user.Username,
			Type:      types.PowerConsume,
			Amount:    task.Power,
			Balance:   u.Power,
			Mark:      types.PowerSub,
			Model:     "dall-e-3",
			Remark:    fmt.Sprintf("绘画提示词：%s", utils.CutWords(task.Prompt, 10)),
			CreatedAt: time.Now(),
		})
	}

	// get image generation API KEY
	var apiKey model.ApiKey
	tx = s.db.Where("type", "dalle").
		Where("enabled", true).
		Order("last_used_at ASC").First(&apiKey)
	if tx.Error != nil {
		return "", fmt.Errorf("no available DALL-E api key: %v", tx.Error)
	}

	var res imgRes
	var errRes ErrRes
	if len(apiKey.ProxyURL) > 5 {
		s.httpClient.SetProxyURL(apiKey.ProxyURL).R()
	}
	apiURL := fmt.Sprintf("%s/v1/images/generations", apiKey.ApiURL)
	reqBody := imgReq{
		Model:   "dall-e-3",
		Prompt:  prompt,
		N:       1,
		Size:    task.Size,
		Style:   task.Style,
		Quality: task.Quality,
	}
	logger.Infof("Channel:%s, API KEY:%s, BODY: %+v", apiURL, apiKey.Value, reqBody)
	r, err := s.httpClient.R().SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(reqBody).
		SetErrorResult(&errRes).
		SetSuccessResult(&res).
		Post(apiURL)
	if err != nil {
		return "", fmt.Errorf("error with send request: %v", err)
	}

	if r.IsErrorState() {
		return "", fmt.Errorf("error with send request, status: %s, %+v", r.Status, errRes.Error)
	}
	// update the api key last use time
	s.db.Model(&apiKey).UpdateColumn("last_used_at", time.Now().Unix())
	// update task progress
	tx = s.db.Model(&model.DallJob{Id: task.JobId}).UpdateColumns(map[string]interface{}{
		"progress": 100,
		"org_url":  res.Data[0].Url,
		"prompt":   prompt,
	})
	if tx.Error != nil {
		return "", fmt.Errorf("err with update database: %v", tx.Error)
	}

	s.notifyQueue.RPush(service.NotifyMessage{UserId: int(task.UserId), JobId: int(task.JobId), Message: service.TaskStatusFailed})
	var content string
	if sync {
		imgURL, err := s.downloadImage(task.JobId, int(task.UserId), res.Data[0].Url)
		if err != nil {
			return "", fmt.Errorf("error with download image: %v", err)
		}
		content = fmt.Sprintf("```\n%s\n```\n下面是我为你创作的图片：\n\n![](%s)\n", prompt, imgURL)
	}

	return content, nil
}

func (s *Service) CheckTaskNotify() {
	go func() {
		logger.Info("Running DALL-E task notify checking ...")
		for {
			var message service.NotifyMessage
			err := s.notifyQueue.LPop(&message)
			if err != nil {
				continue
			}
			client := s.Clients.Get(uint(message.UserId))
			if client == nil {
				continue
			}
			err = client.Send([]byte(message.Message))
			if err != nil {
				continue
			}
		}
	}()
}

func (s *Service) CheckTaskStatus() {
	go func() {
		logger.Info("Running DALL-E task status checking ...")
		for {
			var jobs []model.DallJob
			res := s.db.Where("progress < ?", 100).Find(&jobs)
			if res.Error != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, job := range jobs {
				// 超时的任务标记为失败
				if time.Now().Sub(job.CreatedAt) > time.Minute*10 {
					job.Progress = service.FailTaskProgress
					job.ErrMsg = "任务超时"
					s.db.Updates(&job)
				}
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
				imgURL, err := s.downloadImage(v.Id, int(v.UserId), v.OrgURL)
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

func (s *Service) downloadImage(jobId uint, userId int, orgURL string) (string, error) {
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
	s.notifyQueue.RPush(service.NotifyMessage{UserId: userId, JobId: int(jobId), Message: service.TaskStatusFinished})
	return imgURL, nil
}
