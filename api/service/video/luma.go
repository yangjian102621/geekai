package video

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
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
	"io"
	"time"

	"github.com/imroc/req/v3"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

type Service struct {
	httpClient    *req.Client
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	taskQueue     *store.RedisQueue
	notifyQueue   *store.RedisQueue
	wsService     *service.WebsocketService
	clientIds     map[uint]string
	userService   *service.UserService
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, redisCli *redis.Client, wsService *service.WebsocketService, userService *service.UserService) *Service {
	return &Service{
		httpClient:    req.C().SetTimeout(time.Minute * 3),
		db:            db,
		taskQueue:     store.NewRedisQueue("Video_Task_Queue", redisCli),
		notifyQueue:   store.NewRedisQueue("Video_Notify_Queue", redisCli),
		wsService:     wsService,
		uploadManager: manager,
		clientIds:     map[uint]string{},
		userService:   userService,
	}
}

func (s *Service) PushTask(task types.VideoTask) {
	logger.Infof("add a new Video task to the task list: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) Run() {
	// 将数据库中未提交的人物加载到队列
	var jobs []model.VideoJob
	s.db.Where("task_id", "").Where("progress", 0).Find(&jobs)
	for _, v := range jobs {
		var task types.VideoTask
		err := utils.JsonDecode(v.TaskInfo, &task)
		if err != nil {
			logger.Errorf("decode task info with error: %v", err)
			continue
		}
		task.Id = v.Id
		s.PushTask(task)
		s.clientIds[v.Id] = task.ClientId
	}
	logger.Info("Starting Video job consumer...")
	go func() {
		for {
			var task types.VideoTask
			err := s.taskQueue.LPop(&task)
			if err != nil {
				logger.Errorf("taking task with error: %v", err)
				continue
			}

			// translate prompt
			if utils.HasChinese(task.Prompt) {
				content, err := utils.OpenAIRequest(s.db, fmt.Sprintf(service.TranslatePromptTemplate, task.Prompt), task.TranslateModelId)
				if err == nil {
					task.Prompt = content
				} else {
					logger.Warnf("error with translate prompt: %v", err)
				}
			}

			if task.ClientId != "" {
				s.clientIds[task.Id] = task.ClientId
			}

			var r LumaRespVo
			r, err = s.LumaCreate(task)
			if err != nil {
				logger.Errorf("create task with error: %v", err)
				err = s.db.Model(&model.VideoJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
					"err_msg":   err.Error(),
					"progress":  service.FailTaskProgress,
					"cover_url": "/images/failed.jpg",
				}).Error
				if err != nil {
					logger.Errorf("update task with error: %v", err)
				}
				s.notifyQueue.RPush(service.NotifyMessage{ClientId: task.ClientId, UserId: task.UserId, JobId: int(task.Id), Message: service.TaskStatusFailed})
				continue
			}

			// 更新任务信息
			err = s.db.Model(&model.VideoJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
				"task_id":    r.Id,
				"channel":    r.Channel,
				"prompt_ext": r.Prompt,
			}).Error
			if err != nil {
				logger.Errorf("update task with error: %v", err)
				s.PushTask(task)
			}
		}
	}()
}

type LumaRespVo struct {
	Id                  string      `json:"id"`
	Prompt              string      `json:"prompt"`
	State               string      `json:"state"`
	CreatedAt           time.Time   `json:"created_at"`
	Video               interface{} `json:"video"`
	Liked               interface{} `json:"liked"`
	EstimateWaitSeconds interface{} `json:"estimate_wait_seconds"`
	Channel             string      `json:"channel,omitempty"`
}

func (s *Service) LumaCreate(task types.VideoTask) (LumaRespVo, error) {
	// 读取 API KEY
	var apiKey model.ApiKey
	session := s.db.Session(&gorm.Session{}).Where("type", "luma").Where("enabled", true)
	if task.Channel != "" {
		session = session.Where("api_url", task.Channel)
	}
	tx := session.Order("last_used_at DESC").First(&apiKey)
	if tx.Error != nil {
		return LumaRespVo{}, errors.New("no available API KEY for Luma")
	}

	reqBody := map[string]interface{}{
		"user_prompt":   task.Prompt,
		"expand_prompt": task.Params.PromptOptimize,
		"loop":          task.Params.Loop,
		"image_url":     task.Params.StartImgURL,
		"image_end_url": task.Params.EndImgURL,
	}
	var res LumaRespVo
	apiURL := fmt.Sprintf("%s/luma/generations", apiKey.ApiURL)
	logger.Debugf("API URL: %s, request body: %+v", apiURL, reqBody)
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(reqBody).
		Post(apiURL)
	if err != nil {
		return LumaRespVo{}, fmt.Errorf("请求 API 出错：%v", err)
	}

	if r.StatusCode != 200 && r.StatusCode != 201 {
		return LumaRespVo{}, fmt.Errorf("请求 API 出错：%d, %s", r.StatusCode, r.String())
	}

	body, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return LumaRespVo{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	// update the last_use_at for api key
	apiKey.LastUsedAt = time.Now().Unix()
	session.Updates(&apiKey)
	res.Channel = apiKey.ApiURL
	return res, nil
}

func (s *Service) CheckTaskNotify() {
	go func() {
		logger.Info("Running Suno task notify checking ...")
		for {
			var message service.NotifyMessage
			err := s.notifyQueue.LPop(&message)
			if err != nil {
				continue
			}
			logger.Debugf("Receive notify message: %+v", message)
			client := s.wsService.Clients.Get(message.ClientId)
			if client == nil {
				continue
			}
			utils.SendChannelMsg(client, types.ChLuma, message.Message)
		}
	}()
}

func (s *Service) DownloadFiles() {
	go func() {
		var items []model.VideoJob
		for {
			res := s.db.Where("progress", 102).Find(&items)
			if res.Error != nil {
				continue
			}

			for _, v := range items {
				if v.WaterURL == "" {
					continue
				}

				logger.Infof("try download video: %s", v.WaterURL)
				videoURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(v.WaterURL, true)
				if err != nil {
					logger.Errorf("download video with error: %v", err)
					continue
				}
				logger.Infof("download video success: %s", videoURL)
				v.WaterURL = videoURL

				if v.VideoURL != "" {
					logger.Infof("try download no water video: %s", v.VideoURL)
					videoURL, err = s.uploadManager.GetUploadHandler().PutUrlFile(v.VideoURL, true)
					if err != nil {
						logger.Errorf("download video with error: %v", err)
						continue
					}
				}
				logger.Info("download no water video success: %s", videoURL)
				v.VideoURL = videoURL
				v.Progress = 100
				s.db.Updates(&v)
				s.notifyQueue.RPush(service.NotifyMessage{ClientId: s.clientIds[v.Id], UserId: v.UserId, JobId: int(v.Id), Message: service.TaskStatusFinished})
			}

			time.Sleep(time.Second * 10)
		}
	}()
}

// SyncTaskProgress 异步拉取任务
func (s *Service) SyncTaskProgress() {
	go func() {
		var jobs []model.VideoJob
		for {
			res := s.db.Where("progress < ?", 100).Where("task_id <> ?", "").Find(&jobs)
			if res.Error != nil {
				continue
			}

			for _, job := range jobs {
				task, err := s.QueryLumaTask(job.TaskId, job.Channel)
				if err != nil {
					logger.Errorf("query task with error: %v", err)
					// 更新任务信息
					s.db.Model(&model.VideoJob{Id: job.Id}).UpdateColumns(map[string]interface{}{
						"progress": service.FailTaskProgress, // 102 表示资源未下载完成,
						"err_msg":  err.Error(),
					})
					continue
				}

				logger.Debugf("task: %+v", task)
				if task.State == "completed" { // 更新任务信息
					data := map[string]interface{}{
						"progress":   102, // 102 表示资源未下载完成,
						"water_url":  task.Video.Url,
						"raw_data":   utils.JsonEncode(task),
						"prompt_ext": task.Prompt,
					}
					if task.Video.DownloadUrl != "" {
						data["video_url"] = task.Video.DownloadUrl
					}
					err = s.db.Model(&model.VideoJob{Id: job.Id}).UpdateColumns(data).Error
					if err != nil {
						logger.Errorf("更新数据库失败：%v", err)
						continue
					}
				}

			}

			// 找出失败的任务，并恢复其扣减算力
			s.db.Where("progress", service.FailTaskProgress).Where("power > ?", 0).Find(&jobs)
			for _, job := range jobs {
				err := s.userService.IncreasePower(job.UserId, job.Power, model.PowerLog{
					Type:   types.PowerRefund,
					Model:  "luma",
					Remark: fmt.Sprintf("Luma 任务失败，退回算力。任务ID：%s，Err:%s", job.TaskId, job.ErrMsg),
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

type LumaTaskVo struct {
	Id    string      `json:"id"`
	Liked interface{} `json:"liked"`
	State string      `json:"state"`
	Video struct {
		Url         string `json:"url"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		DownloadUrl string `json:"download_url"`
	} `json:"video"`
	Prompt              string      `json:"prompt"`
	CreatedAt           time.Time   `json:"created_at"`
	EstimateWaitSeconds interface{} `json:"estimate_wait_seconds"`
}

func (s *Service) QueryLumaTask(taskId string, channel string) (LumaTaskVo, error) {
	// 读取 API KEY
	var apiKey model.ApiKey
	err := s.db.Session(&gorm.Session{}).Where("type", "luma").
		Where("api_url", channel).
		Where("enabled", true).
		Order("last_used_at DESC").First(&apiKey).Error
	if err != nil {
		return LumaTaskVo{}, errors.New("no available API KEY for Luma")
	}

	apiURL := fmt.Sprintf("%s/luma/generations/%s", apiKey.ApiURL, taskId)
	var res LumaTaskVo
	r, err := req.C().R().SetHeader("Authorization", "Bearer "+apiKey.Value).Get(apiURL)

	if err != nil {
		return LumaTaskVo{}, fmt.Errorf("请求 API 失败：%v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return LumaTaskVo{}, fmt.Errorf("API 返回失败：%v", r.String())
	}

	body, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return LumaTaskVo{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	return res, nil
}
