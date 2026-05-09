package jimeng

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"

	"github.com/go-redis/redis/v8"
)

var logger = logger2.GetLogger()

// Service 即梦服务（合并了消费者功能）
type Service struct {
	db        *gorm.DB
	redis     *redis.Client
	taskQueue *store.RedisQueue
	client    *Client
	ctx       context.Context
	cancel    context.CancelFunc
	running   bool
	uploader  *oss.UploaderManager
}

// NewService 创建即梦服务
func NewService(db *gorm.DB, redisCli *redis.Client, uploader *oss.UploaderManager, client *Client) *Service {
	taskQueue := store.NewRedisQueue("JimengTaskQueue", redisCli)
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		db:        db,
		redis:     redisCli,
		taskQueue: taskQueue,
		client:    client,
		ctx:       ctx,
		cancel:    cancel,
		running:   false,
		uploader:  uploader,
	}
}

// Start 启动服务（包含消费者）
func (s *Service) Start() {
	if s.running {
		return
	}
	logger.Info("Starting Jimeng service and task consumer...")
	s.running = true
	go s.consumeTasks()
	go s.pollTaskStatus()
}

// Stop 停止服务
func (s *Service) Stop() {
	if !s.running {
		return
	}
	logger.Info("Stopping Jimeng service and task consumer...")
	s.running = false
	s.cancel()
}

// consumeTasks 消费任务
func (s *Service) consumeTasks() {
	for {
		select {
		case <-s.ctx.Done():
			logger.Info("Jimeng task consumer stopped")
			return
		default:
			s.processNextTask()
		}
	}
}

// processNextTask 处理下一个任务
func (s *Service) processNextTask() {
	var jobId uint
	if err := s.taskQueue.LPop(&jobId); err != nil {
		// 队列为空，等待1秒后重试
		time.Sleep(time.Second)
		return
	}

	logger.Infof("Processing Jimeng task: job_id=%d", jobId)

	if err := s.ProcessTask(jobId); err != nil {
		logger.Errorf("process jimeng task failed: job_id=%d, error=%v", jobId, err)
		s.UpdateJobStatus(jobId, types.JMTaskStatusFailed, err.Error())
	} else {
		logger.Infof("Jimeng task processed successfully: job_id=%d", jobId)
	}
}

// CreateTask 创建任务
func (s *Service) CreateTask(userId uint, req *types.JimengTaskRequest) (*model.JimengJob, error) {
	// 生成任务ID
	taskId := utils.RandString(20)

	// 创建任务记录
	job := &model.JimengJob{
		UserId:    userId,
		TaskId:    taskId,
		Type:      req.TaskType,
		ReqKey:    req.ReqKey,
		Prompt:    req.Prompt,
		Params:    utils.JsonEncode(req),
		Status:    types.JMTaskStatusInQueue,
		Power:     req.Power,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	if err := s.db.Create(job).Error; err != nil {
		return nil, fmt.Errorf("create jimeng job failed: %w", err)
	}

	// 推送到任务队列
	if err := s.taskQueue.RPush(job.Id); err != nil {
		return nil, fmt.Errorf("push jimeng task to queue failed: %w", err)
	}

	return job, nil
}

// ProcessTask 处理任务
func (s *Service) ProcessTask(jobId uint) error {
	// 获取任务记录
	var job model.JimengJob
	if err := s.db.First(&job, jobId).Error; err != nil {
		return fmt.Errorf("get jimeng job failed: %w", err)
	}

	// 更新任务状态为处理中
	if err := s.UpdateJobStatus(job.Id, types.JMTaskStatusGenerating, ""); err != nil {
		return fmt.Errorf("update job status failed: %w", err)
	}

	// 解析任务参数
	var req types.JimengTaskRequest
	err := utils.JsonDecode(job.Params, &req)
	if err != nil {
		return fmt.Errorf("parse task params failed: %w", err)
	}

	// 构建请求并提交任务
	params, err := s.buildTaskRequest(&req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("build task request failed: %v", err))
	}

	// 数字人任务，先识别主体
	if req.TaskType == types.JMTaskTypeVirtualHuman {
		if err := s.client.AvatarRecognition(req.ImageUrls[0], req.RecognizeKey); err != nil {
			return s.handleTaskError(job.Id, fmt.Sprintf("avatar recognition failed: %v", err))
		}
	}

	// 同步任务 ，后台执行
	if req.ReqKey == DoubaoSeedream40ReqKey {
		go func() {
			resp, err := s.client.SubmitSyncImageTask(req)
			if err != nil {
				_ = s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
				return
			}
			logger.Infof("同步任务提交成功: %+v", resp)
			// 更新原始数据
			rawData, _ := json.Marshal(resp)
			updates := map[string]any{
				"raw_data": string(rawData),
			}
			if resp.Error != nil {
				updates["status"] = types.JMTaskStatusFailed
				updates["err_msg"] = resp.Error.Message
				s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(updates)
				return
			}

			// 更新任务状态
			updates["status"] = types.JMTaskStatusSuccess
			// 下载图片
			imgUrl, err := s.uploader.GetUploadHandler().PutUrlFile(*resp.Data[0].Url, ".png", false)
			if err == nil {
				updates["img_url"] = imgUrl
			}
			s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(updates)
		}()
		return nil
	}

	logger.Debugf("提交即梦任务: %+v", params)
	// 异步任务 ，前台执行
	resp, err := s.client.SubmitTask(params)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
	}

	if resp.Code != CodeSuccess {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %s", resp.Message))
	}

	// 更新任务ID和原始数据
	rawData, _ := json.Marshal(resp)
	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(map[string]any{
		"task_id":    resp.Data.TaskId,
		"raw_data":   string(rawData),
		"updated_at": time.Now(),
	}).Error; err != nil {
		logger.Errorf("update jimeng job task_id failed: %v", err)
	}

	return nil
}

// buildTaskRequest 构建任务请求（统一的参数解析）
func (s *Service) buildTaskRequest(req *types.JimengTaskRequest) (map[string]any, error) {
	var params map[string]any
	err := utils.JsonDecode(utils.JsonEncode(req), &params)
	if err != nil {
		return nil, fmt.Errorf("parse task params failed: %w", err)
	}
	// 把 size 转成 width 和 height
	if size, ok := params["size"]; ok {
		if sizeStr, ok := size.(string); ok {
			if strings.Contains(sizeStr, "x") {
				sizes := strings.Split(sizeStr, "x")
				params["width"] = sizes[0]
				params["height"] = sizes[1]
			}
		}
		delete(params, "size")
	}

	// duration 转成 frames
	if duration, ok := params["duration"]; ok {
		if secs, ok := duration.(int); ok {
			params["frames"] = secs*24 + 1
		}
		delete(params, "duration")
	}

	// 单独处理图片特效任务
	if req.ReqKey == ImageEffectReqKey {
		params["image_input1"] = req.ImageUrls[0]
		delete(params, "image_urls")
	}

	// 动作迁移，数字人任务参数处理
	if req.TaskType == types.JMTaskTypeVirtualHuman || req.TaskType == types.JMTaskTypeActionTransfer {
		params["image_url"] = req.ImageUrls[0]
		delete(params, "image_urls")
	}
	if req.RecognizeKey != "" {
		delete(params, "recognize_key")
	}

	// 删除多余参数，剩下的就是各个任务自己专有参数了
	delete(params, "type")
	delete(params, "power")
	return params, nil
}

// pollTaskStatus 轮询任务状态
func (s *Service) pollTaskStatus() {

	for {
		var jobs []model.JimengJob
		s.db.Where("status IN (?)", []types.JMTaskStatus{types.JMTaskStatusGenerating, types.JMTaskStatusInQueue}).Find(&jobs)
		if len(jobs) == 0 {
			logger.Debugf("no jimeng task to poll, sleep 10s")
			time.Sleep(10 * time.Second)
			continue
		}

		for _, job := range jobs {
			// 任务超时处理
			if job.UpdatedAt.Before(time.Now().Add(-10 * time.Minute)) {
				s.handleTaskError(job.Id, "task timeout")
				continue
			}

			// 豆包生图 4.0 是同步任务，不需要轮询
			if job.ReqKey == DoubaoSeedream40ReqKey {
				continue
			}

			// 查询任务状态
			resp, err := s.client.QueryTask(&QueryTaskRequest{
				ReqKey:  job.ReqKey,
				TaskId:  job.TaskId,
				ReqJson: `{"return_url":true}`,
			}, ASyncActionGetResult)

			if err != nil {
				s.handleTaskError(job.Id, fmt.Sprintf("query task failed: %s", err.Error()))
				continue
			}

			// 更新原始数据
			rawData, _ := json.Marshal(resp)
			s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Update("raw_data", string(rawData))

			if resp.Code != CodeSuccess {
				s.handleTaskError(job.Id, fmt.Sprintf("query task failed: %s", resp.Message))
				continue
			}

			switch resp.Data.Status {
			case types.JMTaskStatusDone:
				// 判断任务是否成功
				if resp.Message != "Success" {
					s.handleTaskError(job.Id, fmt.Sprintf("task failed: %s", resp.Data.AlgorithmBaseResp.StatusMessage))
					continue
				}

				// 任务完成，更新结果
				updates := map[string]any{
					"status":     types.JMTaskStatusSuccess,
					"updated_at": time.Now(),
				}

				// 设置结果URL
				if len(resp.Data.ImageUrls) > 0 {
					imgUrl, err := s.uploader.GetUploadHandler().PutUrlFile(resp.Data.ImageUrls[0], ".png", false)
					if err != nil {
						logger.Errorf("upload image failed: %v", err)
						imgUrl = resp.Data.ImageUrls[0]
					}
					updates["img_url"] = imgUrl
				}
				if resp.Data.VideoUrl != "" {
					videoUrl, err := s.uploader.GetUploadHandler().PutUrlFile(resp.Data.VideoUrl, ".mp4", false)
					if err != nil {
						logger.Errorf("upload video failed: %v", err)
						videoUrl = resp.Data.VideoUrl
					}
					updates["video_url"] = videoUrl
				}

				s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(updates)
			case types.JMTaskStatusInQueue, types.JMTaskStatusGenerating:
				// 任务处理中
				s.UpdateJobStatus(job.Id, types.JMTaskStatusGenerating, "")

			case types.JMTaskStatusNotFound:
				// 任务未找到
				s.handleTaskError(job.Id, "task not found")

			case types.JMTaskStatusExpired:
				continue
			default:
				logger.Warnf("unknown task status: %s", resp.Data.Status)
			}

		}

		time.Sleep(5 * time.Second)

	}

}

// UpdateJobStatus 更新任务状态
func (s *Service) UpdateJobStatus(jobId uint, status types.JMTaskStatus, errMsg string) error {
	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	if errMsg != "" {
		updates["err_msg"] = errMsg
	}
	return s.db.Model(&model.JimengJob{}).Where("id = ?", jobId).Updates(updates).Error
}

// handleTaskError 处理任务错误
func (s *Service) handleTaskError(jobId uint, errMsg string) error {
	logger.Errorf("Jimeng task error (job_id: %d): %s", jobId, errMsg)
	return s.UpdateJobStatus(jobId, types.JMTaskStatusFailed, errMsg)
}

// PushTaskToQueue 推送任务到队列（用于手动重试）
func (s *Service) PushTaskToQueue(jobId uint) error {
	return s.taskQueue.RPush(jobId)
}

// GetTaskStats 获取任务统计信息
func (s *Service) GetTaskStats() (map[string]any, error) {
	type StatResult struct {
		Status types.JMTaskStatus `json:"status"`
		Count  int64              `json:"count"`
	}

	var stats []StatResult
	err := s.db.Model(&model.JimengJob{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&stats).Error
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"total":      int64(0),
		"completed":  int64(0),
		"processing": int64(0),
		"failed":     int64(0),
		"pending":    int64(0),
	}

	for _, stat := range stats {
		result["total"] = result["total"].(int64) + stat.Count
		result[string(stat.Status)] = stat.Count
	}

	return result, nil
}

// GetJob 获取任务
func (s *Service) GetJob(jobId uint) (*model.JimengJob, error) {
	var job model.JimengJob
	if err := s.db.First(&job, jobId).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
