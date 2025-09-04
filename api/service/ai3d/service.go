package ai3d

import (
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"net/url"
	"path/filepath"
	"strings"
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
	userService   *service.UserService
	uploadManager *oss.UploaderManager
}

// NewService 创建3D生成服务
func NewService(db *gorm.DB, redisCli *redis.Client, tencentClient *Tencent3DClient, giteeClient *Gitee3DClient, userService *service.UserService, uploadManager *oss.UploaderManager) *Service {
	return &Service{
		db:            db,
		taskQueue:     store.NewRedisQueue("3D_Task_Queue", redisCli),
		tencentClient: tencentClient,
		giteeClient:   giteeClient,
		userService:   userService,
		uploadManager: uploadManager,
	}
}

// CreateJob 创建3D生成任务
func (s *Service) CreateJob(userId uint, request vo.AI3DJobParams) (*model.AI3DJob, error) {
	switch request.Type {
	case types.AI3DTaskTypeGitee:
		if s.giteeClient == nil {
			return nil, fmt.Errorf("模力方舟 3D 服务未初始化")
		}
		if !s.giteeClient.GetConfig().Enabled {
			return nil, fmt.Errorf("模力方舟 3D 服务未启用")
		}

	case types.AI3DTaskTypeTencent:
		if s.tencentClient == nil {
			return nil, fmt.Errorf("腾讯云 3D 服务未初始化")
		}
		if !s.tencentClient.GetConfig().Enabled {
			return nil, fmt.Errorf("腾讯云 3D 服务未启用")
		}

	default:
		return nil, fmt.Errorf("不支持的 3D 服务类型: %s", request.Type)
	}

	// 创建任务记录
	job := &model.AI3DJob{
		UserId:     userId,
		Type:       request.Type,
		Power:      request.Power,
		Model:      request.Model,
		Status:     types.AI3DJobStatusPending,
		PreviewURL: request.ImageURL,
	}

	job.Params = utils.JsonEncode(request)

	// 保存到数据库
	if err := s.db.Create(job).Error; err != nil {
		return nil, fmt.Errorf("failed to create 3D job: %v", err)
	}

	// 更新用户算力
	err := s.userService.DecreasePower(userId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  job.Model,
		Remark: fmt.Sprintf("创建3D任务，消耗%d算力", job.Power),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user power: %v", err)
	}

	// 将任务添加到队列
	request.JobId = job.Id
	s.PushTask(request)

	return job, nil
}

// PushTask 将任务添加到队列
func (s *Service) PushTask(job vo.AI3DJobParams) {
	logger.Infof("add a new 3D task to the queue: %+v", job)
	if err := s.taskQueue.RPush(job); err != nil {
		logger.Errorf("push 3D task to queue failed: %v", err)
	}
}

// Run 启动任务处理器
func (s *Service) Run() {
	logger.Info("Starting 3D job consumer...")
	go func() {
		for {
			var params vo.AI3DJobParams
			err := s.taskQueue.LPop(&params)
			if err != nil {
				logger.Errorf("taking 3D task with error: %v", err)
				continue
			}
			logger.Infof("handle a new 3D task: %+v", params)
			go func() {
				if err := s.processJob(&params); err != nil {
					logger.Errorf("error processing 3D job: %v", err)
					s.updateJobStatus(params.JobId, types.AI3DJobStatusFailed, err.Error())
				}
			}()
		}
	}()

	go s.pollJobStatus()
}

// processJob 处理3D任务
func (s *Service) processJob(params *vo.AI3DJobParams) error {
	// 更新状态为处理中
	s.updateJobStatus(params.JobId, types.AI3DJobStatusProcessing, "")

	var taskId string
	var err error

	// 根据类型选择客户端
	switch params.Type {
	case types.AI3DTaskTypeTencent:
		if s.tencentClient == nil {
			return fmt.Errorf("tencent 3D client not initialized")
		}
		tencentParams := Tencent3DParams{
			Prompt:       params.Prompt,
			ImageURL:     params.ImageURL,
			ResultFormat: params.FileFormat,
			EnablePBR:    params.EnablePBR,
		}
		taskId, err = s.tencentClient.SubmitJob(tencentParams)
	case types.AI3DTaskTypeGitee:
		if s.giteeClient == nil {
			return fmt.Errorf("gitee 3D client not initialized")
		}
		giteeParams := Gitee3DParams{
			Model:             params.Model,
			Texture:           params.Texture,
			Seed:              params.Seed,
			NumInferenceSteps: params.NumInferenceSteps,
			GuidanceScale:     params.GuidanceScale,
			OctreeResolution:  params.OctreeResolution,
			ImageURL:          params.ImageURL,
		}
		if params.Model == "Hunyuan3D-2" {
			giteeParams.Type = strings.ToLower(params.FileFormat)
		} else {
			giteeParams.FileFormat = strings.ToLower(params.FileFormat)
		}
		taskId, err = s.giteeClient.SubmitJob(giteeParams)
	default:
		return fmt.Errorf("unsupported 3D API type: %s", params.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to submit 3D job: %v", err)
	}

	// 更新任务ID
	s.db.Model(model.AI3DJob{}).Where("id = ?", params.JobId).Update("task_id", taskId)

	return nil
}

// pollJobStatus 轮询任务状态
func (s *Service) pollJobStatus() {
	// 10秒轮询一次
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var jobs []model.AI3DJob
		s.db.Where("status IN (?)", []string{types.AI3DJobStatusProcessing, types.AI3DJobStatusPending}).Find(&jobs)
		if len(jobs) == 0 {
			logger.Debug("no 3D jobs to poll, sleep 10s")
			continue
		}

		for _, job := range jobs {
			// 15 分钟超时
			if job.CreatedAt.Before(time.Now().Add(-20 * time.Minute)) {
				s.updateJobStatus(job.Id, types.AI3DJobStatusFailed, "task timeout")
				continue
			}

			result, err := s.queryJobStatus(&job)
			if err != nil {
				logger.Errorf("failed to query job status: %v", err)
				continue
			}

			updates := map[string]any{
				"status":   result.Status,
				"raw_data": result.RawData,
				"err_msg":  result.ErrorMsg,
			}
			if result.FileURL != "" {
				// 下载文件到本地
				url, err := s.uploadManager.GetUploadHandler().PutUrlFile(result.FileURL, getFileExt(result.FileURL), false)
				if err != nil {
					logger.Errorf("failed to download file: %v", err)
					continue
				}
				updates["file_url"] = url
				logger.Infof("download file: %s", url)
			}
			if result.PreviewURL != "" {
				url, err := s.uploadManager.GetUploadHandler().PutUrlFile(result.PreviewURL, getFileExt(result.PreviewURL), false)
				if err != nil {
					logger.Errorf("failed to download preview image: %v", err)
					continue
				}
				updates["preview_url"] = url
				logger.Infof("download preview image: %s", url)
			}

			s.db.Model(&model.AI3DJob{}).Where("id = ?", job.Id).Updates(updates)

		}
	}
}

// queryJobStatus 查询任务状态
func (s *Service) queryJobStatus(job *model.AI3DJob) (*types.AI3DJobResult, error) {
	switch job.Type {
	case types.AI3DTaskTypeTencent:
		if s.tencentClient == nil {
			return nil, fmt.Errorf("tencent 3D client not initialized")
		}
		return s.tencentClient.QueryJob(job.TaskId)
	case types.AI3DTaskTypeGitee:
		if s.giteeClient == nil {
			return nil, fmt.Errorf("gitee 3D client not initialized")
		}
		return s.giteeClient.QueryJob(job.TaskId)
	default:
		return nil, fmt.Errorf("unsupported 3D API type: %s", job.Type)
	}
}

// updateJobStatus 更新任务状态
func (s *Service) updateJobStatus(jobId uint, status string, errMsg string) error {

	return s.db.Model(model.AI3DJob{}).Where("id = ?", jobId).Updates(map[string]any{
		"status":  status,
		"err_msg": errMsg,
	}).Error
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
			FileURL:    job.FileURL,
			PreviewURL: job.PreviewURL,
			Model:      job.Model,
			Status:     job.Status,
			ErrMsg:     job.ErrMsg,
			CreatedAt:  job.CreatedAt.Unix(),
			UpdatedAt:  job.UpdatedAt.Unix(),
		}
		_ = utils.JsonDecode(job.Params, &jobVO.Params)
		jobList = append(jobList, jobVO)
	}

	return &vo.Page{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Items:    jobList,
	}, nil
}

// DeleteJob 删除任务
func (s *Service) DeleteUserJob(id uint, userId uint) error {
	var job model.AI3DJob
	err := s.db.Where("id = ?", id).Where("user_id = ?", userId).First(&job).Error
	if err != nil {
		return err
	}

	tx := s.db.Begin()
	err = tx.Delete(&job).Error
	if err != nil {
		return err
	}

	// 失败的任务要退回算力
	if job.Status == types.AI3DJobStatusFailed {
		err = s.userService.IncreasePower(userId, job.Power, model.PowerLog{
			Type:   types.PowerRefund,
			Model:  job.Model,
			Remark: fmt.Sprintf("删除任务，退回%d算力", job.Power),
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
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

// getFileExt 获取文件扩展名
func getFileExt(fileURL string) string {
	parse, err := url.Parse(fileURL)
	if err != nil {
		return ""
	}
	ext := filepath.Ext(parse.Path)
	if ext == "" {
		return ".glb"
	}
	return ext
}
