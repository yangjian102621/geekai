package ai3d

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/store"
	"geekai/store/model"
	"geekai/store/vo"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

// Service 3D生成服务
type Service struct {
	db            *gorm.DB
	taskQueue     *store.RedisQueue
	tencentClient *Tencent3DClient
	giteeClient   *Gitee3DClient
}

// NewService 创建3D生成服务
func NewService(db *gorm.DB, redisCli *redis.Client, tencentClient *Tencent3DClient, giteeClient *Gitee3DClient) *Service {
	return &Service{
		db:            db,
		taskQueue:     store.NewRedisQueue("3D_Task_Queue", redisCli),
		tencentClient: tencentClient,
		giteeClient:   giteeClient,
	}
}

// CreateJob 创建3D生成任务
func (s *Service) CreateJob(userId uint, request vo.AI3DJobCreate) (*model.AI3DJob, error) {
	// 创建任务记录
	job := &model.AI3DJob{
		UserId: userId,
		Type:   request.Type,
		Power:  request.Power,
		Model:  request.Model,
		Status: types.AI3DJobStatusPending,
	}

	// 序列化参数
	params := map[string]any{
		"prompt":    request.Prompt,
		"image_url": request.ImageURL,
		"model":     request.Model,
		"power":     request.Power,
	}
	paramsJSON, _ := json.Marshal(params)
	job.Params = string(paramsJSON)

	// 保存到数据库
	if err := s.db.Create(job).Error; err != nil {
		return nil, fmt.Errorf("failed to create 3D job: %v", err)
	}

	// 将任务添加到队列
	s.PushTask(job)

	return job, nil
}

// PushTask 将任务添加到队列
func (s *Service) PushTask(job *model.AI3DJob) {
	logger.Infof("add a new 3D task to the queue: %+v", job)
	if err := s.taskQueue.RPush(job); err != nil {
		logger.Errorf("push 3D task to queue failed: %v", err)
	}
}

// Run 启动任务处理器
func (s *Service) Run() {
	// 将数据库中未完成的任务加载到队列
	var jobs []model.AI3DJob
	s.db.Where("status IN ?", []string{types.AI3DJobStatusPending, types.AI3DJobStatusProcessing}).Find(&jobs)
	for _, job := range jobs {
		s.PushTask(&job)
	}

	logger.Info("Starting 3D job consumer...")
	go func() {
		for {
			var job model.AI3DJob
			err := s.taskQueue.LPop(&job)
			if err != nil {
				logger.Errorf("taking 3D task with error: %v", err)
				continue
			}
			logger.Infof("handle a new 3D task: %+v", job)
			go func() {
				if err := s.processJob(&job); err != nil {
					logger.Errorf("error processing 3D job: %v", err)
					s.updateJobStatus(&job, types.AI3DJobStatusFailed, 0, err.Error())
				}
			}()
		}
	}()
}

// processJob 处理3D任务
func (s *Service) processJob(job *model.AI3DJob) error {
	// 更新状态为处理中
	s.updateJobStatus(job, types.AI3DJobStatusProcessing, 10, "")

	// 解析参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.Params), &params); err != nil {
		return fmt.Errorf("failed to parse job params: %v", err)
	}

	var taskId string
	var err error

	// 根据类型选择客户端
	switch job.Type {
	case "tencent":
		if s.tencentClient == nil {
			return fmt.Errorf("tencent 3D client not initialized")
		}
		tencentParams := Tencent3DParams{
			Prompt:       s.getString(params, "prompt"),
			ImageURL:     s.getString(params, "image_url"),
			ResultFormat: job.Model,
			EnablePBR:    false,
		}
		taskId, err = s.tencentClient.SubmitJob(tencentParams)
	case "gitee":
		if s.giteeClient == nil {
			return fmt.Errorf("gitee 3D client not initialized")
		}
		giteeParams := Gitee3DParams{
			Prompt:       s.getString(params, "prompt"),
			ImageURL:     s.getString(params, "image_url"),
			ResultFormat: job.Model,
		}
		taskId, err = s.giteeClient.SubmitJob(giteeParams)
	default:
		return fmt.Errorf("unsupported 3D API type: %s", job.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to submit 3D job: %v", err)
	}

	// 更新任务ID
	job.TaskId = taskId
	s.db.Model(job).Update("task_id", taskId)

	// 开始轮询任务状态
	go s.pollJobStatus(job)

	return nil
}

// pollJobStatus 轮询任务状态
func (s *Service) pollJobStatus(job *model.AI3DJob) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := s.queryJobStatus(job)
			if err != nil {
				logger.Errorf("failed to query job status: %v", err)
				continue
			}

			// 更新进度
			s.updateJobStatus(job, result.Status, result.Progress, result.ErrorMsg)

			// 如果任务完成或失败，停止轮询
			if result.Status == types.AI3DJobStatusCompleted || result.Status == types.AI3DJobStatusFailed {
				if result.Status == types.AI3DJobStatusCompleted {
					// 更新结果文件URL
					s.db.Model(job).Updates(map[string]interface{}{
						"img_url":     result.FileURL,
						"preview_url": result.PreviewURL,
					})
				}
				return
			}
		}
	}
}

// queryJobStatus 查询任务状态
func (s *Service) queryJobStatus(job *model.AI3DJob) (*types.AI3DJobResult, error) {
	switch job.Type {
	case "tencent":
		if s.tencentClient == nil {
			return nil, fmt.Errorf("tencent 3D client not initialized")
		}
		return s.tencentClient.QueryJob(job.TaskId)
	case "gitee":
		if s.giteeClient == nil {
			return nil, fmt.Errorf("gitee 3D client not initialized")
		}
		return s.giteeClient.QueryJob(job.TaskId)
	default:
		return nil, fmt.Errorf("unsupported 3D API type: %s", job.Type)
	}
}

// updateJobStatus 更新任务状态
func (s *Service) updateJobStatus(job *model.AI3DJob, status string, progress int, errMsg string) {
	updates := map[string]interface{}{
		"status":     status,
		"progress":   progress,
		"updated_at": time.Now(),
	}
	if errMsg != "" {
		updates["err_msg"] = errMsg
	}

	if err := s.db.Model(job).Updates(updates).Error; err != nil {
		logger.Errorf("failed to update job status: %v", err)
	}
}

// GetJobList 获取任务列表
func (s *Service) GetJobList(userId uint, page, pageSize int) (*vo.Page, error) {
	var total int64
	var jobs []model.AI3DJob

	// 查询总数
	if err := s.db.Model(&model.AI3DJob{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询任务列表
	offset := (page - 1) * pageSize
	if err := s.db.Where("user_id = ?", userId).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&jobs).Error; err != nil {
		return nil, err
	}

	// 转换为VO
	var jobList []vo.AI3DJob
	for _, job := range jobs {
		jobVO := vo.AI3DJob{
			Id:         job.Id,
			UserId:     job.UserId,
			Type:       job.Type,
			Power:      job.Power,
			TaskId:     job.TaskId,
			ImgURL:     job.FileURL,
			PreviewURL: job.PreviewURL,
			Model:      job.Model,
			Status:     job.Status,
			ErrMsg:     job.ErrMsg,
			Params:     job.Params,
			CreatedAt:  job.CreatedAt.Unix(),
			UpdatedAt:  job.UpdatedAt.Unix(),
		}
		jobList = append(jobList, jobVO)
	}

	return &vo.Page{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Items:    jobList,
	}, nil
}

// GetJobById 根据ID获取任务
func (s *Service) GetJobById(id uint) (*model.AI3DJob, error) {
	var job model.AI3DJob
	if err := s.db.Where("id = ?", id).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// DeleteJob 删除任务
func (s *Service) DeleteJob(id uint, userId uint) error {
	var job model.AI3DJob
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&job).Error; err != nil {
		return err
	}

	// 如果任务已完成，退还算力
	if job.Status == types.AI3DJobStatusCompleted {
		// TODO: 实现算力退还逻辑
		logger2.GetLogger().Infof("should refund power %d for user %d", job.Power, userId)
	}

	return s.db.Delete(&job).Error
}

// GetSupportedModels 获取支持的模型列表
func (s *Service) GetSupportedModels() map[string][]types.AI3DModel {

	models := make(map[string][]types.AI3DModel)
	if s.tencentClient != nil {
		models["tencent"] = s.tencentClient.GetSupportedModels()
	}
	if s.giteeClient != nil {
		models["gitee"] = s.giteeClient.GetSupportedModels()
	}
	return models
}

func (s *Service) UpdateConfig(config types.AI3DConfig) {
	if s.tencentClient != nil {
		s.tencentClient.UpdateConfig(config.Tencent)
	}
	if s.giteeClient != nil {
		s.giteeClient.UpdateConfig(config.Gitee)
	}
}

// getString 从map中获取字符串值
func (s *Service) getString(params map[string]interface{}, key string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
