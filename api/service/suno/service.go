package suno

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
	Clients       *types.LMap[uint, *types.WsClient] // UserId => Client
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, redisCli *redis.Client) *Service {
	return &Service{
		httpClient:    req.C().SetTimeout(time.Minute * 3),
		db:            db,
		taskQueue:     store.NewRedisQueue("Suno_Task_Queue", redisCli),
		notifyQueue:   store.NewRedisQueue("Suno_Notify_Queue", redisCli),
		Clients:       types.NewLMap[uint, *types.WsClient](),
		uploadManager: manager,
	}
}

func (s *Service) PushTask(task types.SunoTask) {
	logger.Infof("add a new Suno task to the task list: %+v", task)
	s.taskQueue.RPush(task)
}

func (s *Service) Run() {
	// 将数据库中未提交的人物加载到队列
	var jobs []model.SunoJob
	s.db.Where("task_id", "").Find(&jobs)
	for _, v := range jobs {
		s.PushTask(types.SunoTask{
			Id:           v.Id,
			Channel:      v.Channel,
			UserId:       v.UserId,
			Type:         v.Type,
			Title:        v.Title,
			RefTaskId:    v.RefTaskId,
			RefSongId:    v.RefSongId,
			Prompt:       v.Prompt,
			Tags:         v.Tags,
			Model:        v.ModelName,
			Instrumental: v.Instrumental,
			ExtendSecs:   v.ExtendSecs,
		})
	}
	logger.Info("Starting Suno job consumer...")
	go func() {
		for {
			var task types.SunoTask
			err := s.taskQueue.LPop(&task)
			if err != nil {
				logger.Errorf("taking task with error: %v", err)
				continue
			}

			r, err := s.Create(task)
			if err != nil {
				logger.Errorf("create task with error: %v", err)
				s.db.Model(&model.SunoJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
					"err_msg":  err.Error(),
					"progress": service.FailTaskProgress,
				})
				continue
			}

			// 更新任务信息
			s.db.Model(&model.SunoJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
				"task_id": r.Data,
				"channel": r.Channel,
			})
		}
	}()
}

type RespVo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
	Channel string `json:"channel,omitempty"`
}

func (s *Service) Create(task types.SunoTask) (RespVo, error) {
	// 读取 API KEY
	var apiKey model.ApiKey
	session := s.db.Session(&gorm.Session{}).Where("type", "suno").Where("enabled", true)
	if task.Channel != "" {
		session = session.Where("api_url", task.Channel)
	}
	tx := session.Order("last_used_at DESC").First(&apiKey)
	if tx.Error != nil {
		return RespVo{}, errors.New("no available API KEY for Suno")
	}

	reqBody := map[string]interface{}{
		"task_id":           task.RefTaskId,
		"continue_clip_id":  task.RefSongId,
		"continue_at":       task.ExtendSecs,
		"make_instrumental": task.Instrumental,
	}
	// 灵感模式
	if task.Type == 1 {
		reqBody["gpt_description_prompt"] = task.Prompt
	} else { // 自定义模式
		reqBody["prompt"] = task.Prompt
		reqBody["tags"] = task.Tags
		reqBody["mv"] = task.Model
		reqBody["title"] = task.Title
	}

	var res RespVo
	apiURL := fmt.Sprintf("%s/task/suno/v1/submit/music", apiKey.ApiURL)
	logger.Debugf("API URL: %s, request body: %+v", apiURL, reqBody)
	r, err := req.C().R().
		SetHeader("Authorization", "Bearer "+apiKey.Value).
		SetBody(reqBody).
		Post(apiURL)
	if err != nil {
		return RespVo{}, fmt.Errorf("请求 API 出错：%v", err)
	}

	body, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return RespVo{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	if res.Code != "success" {
		return RespVo{}, fmt.Errorf("API 返回失败：%s", res.Message)
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

func (s *Service) DownloadImages() {
	go func() {
		var items []model.SunoJob
		for {
			res := s.db.Where("progress", 102).Find(&items)
			if res.Error != nil {
				continue
			}

			for _, v := range items {
				// 下载图片和音频
				logger.Infof("try download cover image: %s", v.CoverURL)
				coverURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(v.CoverURL, true)
				if err != nil {
					logger.Errorf("download image with error: %v", err)
					continue
				}

				logger.Infof("try download audio: %s", v.AudioURL)
				audioURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(v.AudioURL, true)
				if err != nil {
					logger.Errorf("download audio with error: %v", err)
					continue
				}
				v.CoverURL = coverURL
				v.AudioURL = audioURL
				v.Progress = 100
				s.db.Updates(&v)
				s.notifyQueue.RPush(service.NotifyMessage{UserId: v.UserId, JobId: int(v.Id), Message: service.TaskStatusFinished})
			}

			time.Sleep(time.Second * 10)
		}
	}()
}

// SyncTaskProgress 异步拉取任务
func (s *Service) SyncTaskProgress() {
	go func() {
		var jobs []model.SunoJob
		for {
			res := s.db.Where("progress < ?", 100).Where("task_id <> ?", "").Find(&jobs)
			if res.Error != nil {
				continue
			}

			for _, job := range jobs {
				task, err := s.QueryTask(job.TaskId, job.Channel)
				if err != nil {
					logger.Errorf("query task with error: %v", err)
					continue
				}

				if task.Code != "success" {
					logger.Errorf("query task with error: %v", task.Message)
					continue
				}

				logger.Debugf("task: %+v", task.Data.Status)
				// 任务完成，删除旧任务插入两条新任务
				if task.Data.Status == "SUCCESS" {
					var jobId = job.Id
					var flag = false
					tx := s.db.Begin()
					for _, v := range task.Data.Data {
						job.Id = 0
						job.Progress = 102 // 102 表示资源未下载完成
						job.Title = v.Title
						job.SongId = v.Id
						job.Duration = int(v.Metadata.Duration)
						job.Prompt = v.Metadata.Prompt
						job.Tags = v.Metadata.Tags
						job.ModelName = v.ModelName
						job.RawData = utils.JsonEncode(v)
						job.CoverURL = v.ImageLargeUrl
						job.AudioURL = v.AudioUrl

						if err = tx.Create(&job).Error; err != nil {
							logger.Error("create job with error: %v", err)
							tx.Rollback()
							break
						}
						flag = true
					}

					// 删除旧任务
					if flag {
						if err = tx.Delete(&model.SunoJob{}, "id = ?", jobId).Error; err != nil {
							logger.Error("create job with error: %v", err)
							tx.Rollback()
							continue
						}
					}
					tx.Commit()

				} else if task.Data.FailReason != "" {
					job.Progress = service.FailTaskProgress
					job.ErrMsg = task.Data.FailReason
					s.db.Updates(&job)
					s.notifyQueue.RPush(service.NotifyMessage{UserId: job.UserId, JobId: int(job.Id), Message: service.TaskStatusFailed})
				}
			}

			time.Sleep(time.Second * 10)
		}
	}()
}

type QueryRespVo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskId     string `json:"task_id"`
		Action     string `json:"action"`
		Status     string `json:"status"`
		FailReason string `json:"fail_reason"`
		SubmitTime int    `json:"submit_time"`
		StartTime  int    `json:"start_time"`
		FinishTime int    `json:"finish_time"`
		Progress   string `json:"progress"`
		Data       []struct {
			Id       string `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Metadata struct {
				Tags         string      `json:"tags"`
				Type         string      `json:"type"`
				Prompt       string      `json:"prompt"`
				Stream       bool        `json:"stream"`
				Duration     float64     `json:"duration"`
				ErrorMessage interface{} `json:"error_message"`
			} `json:"metadata"`
			AudioUrl          string `json:"audio_url"`
			ImageUrl          string `json:"image_url"`
			VideoUrl          string `json:"video_url"`
			ModelName         string `json:"model_name"`
			DisplayName       string `json:"display_name"`
			ImageLargeUrl     string `json:"image_large_url"`
			MajorModelVersion string `json:"major_model_version"`
		} `json:"data"`
	} `json:"data"`
}

func (s *Service) QueryTask(taskId string, channel string) (QueryRespVo, error) {
	// 读取 API KEY
	var apiKey model.ApiKey
	tx := s.db.Session(&gorm.Session{}).Where("type", "suno").
		Where("api_url", channel).
		Where("enabled", true).
		Order("last_used_at DESC").First(&apiKey)
	if tx.Error != nil {
		return QueryRespVo{}, errors.New("no available API KEY for Suno")
	}

	apiURL := fmt.Sprintf("%s/task/suno/v1/fetch/%s", apiKey.ApiURL, taskId)
	var res QueryRespVo
	r, err := req.C().R().SetHeader("Authorization", "Bearer "+apiKey.Value).Get(apiURL)

	if err != nil {
		return QueryRespVo{}, fmt.Errorf("请求 API 失败：%v", err)
	}

	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return QueryRespVo{}, fmt.Errorf("解析API数据失败：%v, %s", err, string(body))
	}

	return res, nil
}
