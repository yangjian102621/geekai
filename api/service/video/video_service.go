package video

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/log"
	"geekai/service"
	"geekai/service/oss"
	"geekai/service/video/adapters"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var logger = log.GetLogger()

type Service struct {
	db            *gorm.DB
	uploadManager *oss.UploaderManager
	taskQueue     *store.RedisQueue
	userService   *service.UserService
	adapters      map[string]adapters.VideoAdapter // provider -> adapter
}

func NewService(db *gorm.DB, manager *oss.UploaderManager, redisCli *redis.Client, userService *service.UserService) *Service {
	service := &Service{
		db:            db,
		taskQueue:     store.NewRedisQueue("Video_Task_Queue", redisCli),
		uploadManager: manager,
		userService:   userService,
		adapters:      make(map[string]VideoAdapter),
	}

	// 注册所有适配器
	service.registerAdapters()

	return service
}

// VideoAdapter 类型别名，指向 adapters.VideoAdapter
type VideoAdapter = adapters.VideoAdapter

// registerAdapters 注册所有视频生成适配器
func (s *Service) registerAdapters() {
	// 注册 Veo 适配器
	veoAdapter := adapters.NewVeoAdapter(s.db)
	s.adapters[veoAdapter.GetProvider()] = veoAdapter

	// 注册 Sora 适配器
	soraAdapter := adapters.NewSoraAdapter(s.db)
	s.adapters[soraAdapter.GetProvider()] = soraAdapter

	// 注册 Luma 适配器
	lumaAdapter := adapters.NewLumaAdapter(s.db)
	s.adapters[lumaAdapter.GetProvider()] = lumaAdapter

	// 注册可灵适配器
	kelingAdapter := adapters.NewKelingAdapter(s.db)
	s.adapters[kelingAdapter.GetProvider()] = kelingAdapter

	// 注册 MiniMax 适配器
	minimaxAdapter := adapters.NewMiniMaxAdapter(s.db)
	s.adapters[minimaxAdapter.GetProvider()] = minimaxAdapter

	// 注册 Wan 适配器
	wanAdapter := adapters.NewWanAdapter(s.db)
	s.adapters[wanAdapter.GetProvider()] = wanAdapter

	// 注册 Doubao 适配器
	doubaoAdapter := adapters.NewDoubaoAdapter(s.db)
	s.adapters[doubaoAdapter.GetProvider()] = doubaoAdapter
}

// getAdapter 获取指定 provider 的适配器
func (s *Service) getAdapter(provider string) (adapters.VideoAdapter, error) {
	adapter, ok := s.adapters[provider]
	if !ok {
		return nil, fmt.Errorf("不支持的视频生成服务提供商: %s", provider)
	}
	return adapter, nil
}

// getVideoConfig 获取视频配置
func (s *Service) getVideoConfig() (*types.VideoConfig, error) {
	return GetVideoConfig(s.db)
}

// CreateTask 统一的创建任务方法
func (s *Service) CreateTask(task types.VideoTask) (adapters.CreateTaskResponse, error) {
	// 获取适配器
	adapter, err := s.getAdapter(task.Type)
	if err != nil {
		return adapters.CreateTaskResponse{}, err
	}

	// 获取视频配置
	videoConfig, err := s.getVideoConfig()
	if err != nil {
		return adapters.CreateTaskResponse{}, err
	}

	// 调用适配器创建任务
	return adapter.CreateTask(task, videoConfig)
}

// QueryTask 统一的查询任务方法
func (s *Service) QueryTask(provider string, taskId string, channel string, modelKey string) (adapters.QueryTaskResponse, error) {
	// 获取适配器
	adapter, err := s.getAdapter(provider)
	if err != nil {
		return adapters.QueryTaskResponse{}, err
	}

	// 获取视频配置
	videoConfig, err := s.getVideoConfig()
	if err != nil {
		return adapters.QueryTaskResponse{}, err
	}

	// 调用适配器查询任务
	return adapter.QueryTask(taskId, channel, videoConfig)
}

func (s *Service) PushTask(task types.VideoTask) {
	logger.Infof("[video] push task to queue jobId=%d type=%s", task.Id, task.Type)
	if err := s.taskQueue.RPush(task); err != nil {
		logger.Errorf("[video] push task to queue failed jobId=%d: %v", task.Id, err)
	}
}

func (s *Service) Run() {
	// 将数据库中未提交的任务加载到队列
	var jobs []model.VideoJob
	s.db.Where("task_id", "").Where("progress", 0).Find(&jobs)
	for _, v := range jobs {
		var task types.VideoTask
		err := utils.JsonDecode(v.Params, &task)
		if err != nil {
			logger.Errorf("decode task info with error: %v", err)
			continue
		}
		task.Id = v.Id
		s.PushTask(task)
	}
	logger.Infof("[video] job consumer started, loaded %d pending jobs from DB", len(jobs))
	go func() {
		for {
			var task types.VideoTask
			err := s.taskQueue.LPop(&task)
			if err != nil {
				logger.Errorf("taking task with error: %v", err)
				continue
			}

			logger.Debugf("[video] submitting task jobId=%d type=%s prompt=%q", task.Id, task.Type, task.Prompt)
			r, err := s.CreateTask(task)
			if err != nil {
				logger.Errorf("[video] submit failed jobId=%d type=%s: %v", task.Id, task.Type, err)
				err = s.db.Model(&model.VideoJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
					"err_msg": err.Error(),
					"status":  types.VideoStatusFailed,
				}).Error
				if err != nil {
					logger.Errorf("update task with error: %v", err)
				}
				continue
			}

			logger.Infof("[video] submit success jobId=%d type=%s taskId=%s channel=%s", task.Id, task.Type, r.TaskId, r.Channel)
			err = s.db.Model(&model.VideoJob{Id: task.Id}).UpdateColumns(map[string]interface{}{
				"task_id": r.TaskId,
				"channel": r.Channel,
				"status":  types.VideoStatusPending,
			}).Error
			if err != nil {
				logger.Errorf("update task with error: %v", err)
				s.PushTask(task)
			}

		}
	}()
}

func (s *Service) DownloadFiles() {
	go func() {
		var items []model.VideoJob
		logger.Info("[video] download files started")
		for {
			err := s.db.Where("status", types.VideoStatusDownloading).Find(&items).Error
			if err != nil {
				logger.Errorf("get downloading tasks with error: %v", err)
				continue
			}

			for _, v := range items {
				if v.VideoURL == "" {
					continue
				}

				logger.Infof("try download video: %s", v.VideoURL)
				videoURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(v.VideoURL, ".mp4", true)
				if err != nil {
					logger.Errorf("download video with error: %v", err)
					continue
				}
				logger.Infof("download video success: %s", videoURL)
				s.db.Model(&model.VideoJob{Id: v.Id}).UpdateColumns(map[string]any{
					"video_url": videoURL,
					"status":    types.VideoStatusSuccess,
					"progress":  100,
				})
			}

			time.Sleep(time.Second * 10)
		}
	}()
}

// SyncTaskProgress 异步拉取任务
func (s *Service) SyncTaskProgress() {
	go func() {
		logger.Info("[video] task status poller started")
		var jobs []model.VideoJob
		for {
			res := s.db.Where("status IN ?", []string{types.VideoStatusInProgress, types.VideoStatusPending}).Where("task_id <> ?", "").Find(&jobs)
			if res.Error != nil {
				continue
			}
			if len(jobs) > 0 {
				logger.Infof("[video] polling task status, in_progress count=%d", len(jobs))
			}

			for _, job := range jobs {
				// 检查任务是否超时（超过 2 小时）
				if time.Since(job.CreatedAt) > 2*time.Hour {
					logger.Warnf("[video] task timeout jobId=%d taskId=%s created_at=%s", job.Id, job.TaskId, job.CreatedAt.Format(time.RFC3339))
					err := s.db.Model(&model.VideoJob{Id: job.Id}).UpdateColumns(map[string]any{
						"status":  types.VideoStatusFailed,
						"err_msg": "任务超时",
					}).Error
					if err != nil {
						logger.Errorf("[video] update timeout task failed jobId=%d: %v", job.Id, err)
					}
					continue
				}

				modelKey := ""
				var videoTask types.VideoTask
				if err := json.Unmarshal([]byte(job.Params), &videoTask); err == nil {
					if paramsMap, ok := videoTask.Params.(map[string]any); ok {
						if model, ok := paramsMap["model"].(string); ok {
							modelKey = model
						}
					}
				}

				logger.Debugf("[video] querying task jobId=%d taskId=%s provider=%s", job.Id, job.TaskId, job.Type)
				task, err := s.QueryTask(job.Type, job.TaskId, job.Channel, modelKey)
				if err != nil {
					logger.Errorf("[video] query failed jobId=%d taskId=%s: %v", job.Id, job.TaskId, err)
					// 更新任务信息
					s.db.Model(&model.VideoJob{Id: job.Id}).UpdateColumns(map[string]any{
						"status":  types.VideoStatusFailed,
						"err_msg": err.Error(),
					})
					continue
				}

				logger.Debugf("[video] task status jobId=%d taskId=%s status=%s", job.Id, job.TaskId, task.Status)
				logger.Debugf("[video] output=%s", task.Output)

				if task.Status == types.VideoStatusSuccess {
					data := map[string]any{
						"status":   types.VideoStatusDownloading,
						"progress": 100,
						"output":   task.Output,
					}
					if task.VideoURL != "" {
						data["video_url"] = task.VideoURL
					}
					err = s.db.Model(&model.VideoJob{Id: job.Id}).UpdateColumns(data).Error
					if err != nil {
						logger.Errorf("更新数据库失败：%v", err)
						continue
					}
					logger.Infof("[video] task completed jobId=%d taskId=%s", job.Id, job.TaskId)
				} else if task.Status == "failed" {
					logger.Warnf("[video] task failed jobId=%d taskId=%s err=%s", job.Id, job.TaskId, task.ErrMsg)
					s.db.Model(&model.VideoJob{Id: job.Id}).UpdateColumns(map[string]any{
						"status":  types.VideoStatusFailed,
						"err_msg": task.ErrMsg,
					})
				} else {
					s.db.Model(&model.VideoJob{Id: job.Id}).UpdateColumns(map[string]any{
						"status":   task.Status,
						"progress": task.Progress,
					})
					continue
				}

			}

			// 找出失败的任务，并恢复其扣减算力
			s.db.Select("id", "user_id", "power", "task_id", "err_msg", "type").
				Where("status", types.VideoStatusFailed).Where("power > ?", 0).Find(&jobs)
			for _, job := range jobs {
				err := s.userService.IncreasePower(job.UserId, job.Power, model.PowerLog{
					Type:   types.PowerRefund,
					Model:  job.Type,
					Remark: fmt.Sprintf("%s 任务失败，退回算力。任务ID：%s，Err:%s", job.Type, job.TaskId, job.ErrMsg),
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
