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
	userService   *service.UserService
	wsService     *service.WebsocketService
	clientIds     map[uint]string
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, redisCli *redis.Client, userService *service.UserService, wsService *service.WebsocketService) *Service {
	return &Service{
		httpClient:    req.C().SetTimeout(time.Minute * 3),
		db:            db,
		taskQueue:     store.NewRedisQueue("DallE_Task_Queue", redisCli),
		notifyQueue:   store.NewRedisQueue("DallE_Notify_Queue", redisCli),
		wsService:     wsService,
		uploadManager: manager,
		userService:   userService,
		clientIds:     map[uint]string{},
	}
}

// PushTask push a new mj task in to task queue
func (s *Service) PushTask(task types.DallTask) {
	logger.Infof("add a new DALL-E task to the task list: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) Run() {
	// 将数据库中未提交的人物加载到队列
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
			s.clientIds[task.Id] = task.ClientId
			_, err = s.Image(task, false)
			if err != nil {
				logger.Errorf("error with image task: %v", err)
				s.db.Model(&model.DallJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
					"progress": service.FailTaskProgress,
					"err_msg":  err.Error(),
				})
				s.notifyQueue.RPush(service.NotifyMessage{ClientId: task.ClientId, UserId: int(task.UserId), JobId: int(task.Id), Message: service.TaskStatusFailed})
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
		content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.TranslatePromptTemplate, prompt), task.TranslateModelId)
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

	// 扣减算力
	err := s.userService.DecreasePower(int(user.Id), task.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "dall-e-3",
		Remark: fmt.Sprintf("绘画提示词：%s", utils.CutWords(task.Prompt, 10)),
	})
	if err != nil {
		return "", fmt.Errorf("error with decrease power: %v", err)
	}

	// get image generation API KEY
	var apiKey model.ApiKey
	err = s.db.Where("type", "dalle").
		Where("enabled", true).
		Order("last_used_at ASC").First(&apiKey).Error
	if err != nil {
		return "", fmt.Errorf("no available DALL-E api key: %v", err)
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
	r, err := s.httpClient.R().SetHeader("Body-Type", "application/json").
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
	err = s.db.Model(&model.DallJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
		"progress": 100,
		"org_url":  res.Data[0].Url,
		"prompt":   prompt,
	}).Error
	if err != nil {
		return "", fmt.Errorf("err with update database: %v", err)
	}

	s.notifyQueue.RPush(service.NotifyMessage{ClientId: task.ClientId, UserId: int(task.UserId), JobId: int(task.Id), Message: service.TaskStatusFailed})
	var content string
	if sync {
		imgURL, err := s.downloadImage(task.Id, int(task.UserId), res.Data[0].Url)
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

			logger.Debugf("notify message: %+v", message)
			client := s.wsService.Clients.Get(message.ClientId)
			if client == nil {
				continue
			}
			utils.SendChannelMsg(client, types.ChDall, message.Message)
		}
	}()
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
				if time.Now().Sub(job.CreatedAt) > time.Minute*10 {
					job.Progress = service.FailTaskProgress
					job.ErrMsg = "任务超时"
					s.db.Updates(&job)
				}
			}

			// 找出失败的任务，并恢复其扣减算力
			s.db.Where("progress", service.FailTaskProgress).Where("power > ?", 0).Find(&jobs)
			for _, job := range jobs {
				err := s.userService.IncreasePower(int(job.UserId), job.Power, model.PowerLog{
					Type:   types.PowerRefund,
					Model:  "dall-e-3",
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
	s.notifyQueue.RPush(service.NotifyMessage{ClientId: s.clientIds[jobId], UserId: userId, JobId: int(jobId), Message: service.TaskStatusFinished})
	return imgURL, nil
}
