package jimeng

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"geekai/logger"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"

	"github.com/go-redis/redis/v8"
)

var serviceLogger = logger.GetLogger()

// Service 即梦服务
type Service struct {
	db        *gorm.DB
	redis     *redis.Client
	taskQueue *store.RedisQueue
	client    *Client
}

// NewService 创建即梦服务
func NewService(db *gorm.DB, redisCli *redis.Client, client *Client) *Service {
	taskQueue := store.NewRedisQueue("JimengTaskQueue", redisCli)
	return &Service{
		db:        db,
		redis:     redisCli,
		taskQueue: taskQueue,
		client:    client,
	}
}

// CreateTask 创建任务
func (s *Service) CreateTask(userId uint, req *CreateTaskRequest) (*model.JimengJob, error) {
	// 生成任务ID
	taskId := utils.RandString(20)

	// 序列化任务参数
	paramsJson, err := json.Marshal(req.Params)
	if err != nil {
		return nil, fmt.Errorf("marshal task params failed: %w", err)
	}

	// 创建任务记录
	job := &model.JimengJob{
		UserId:     userId,
		TaskId:     taskId,
		Type:       req.Type,
		ReqKey:     req.ReqKey,
		Prompt:     req.Prompt,
		TaskParams: string(paramsJson),
		Status:     model.JimengJobStatusPending,
		Power:      req.Power,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存到数据库
	if err := s.db.Create(job).Error; err != nil {
		return nil, fmt.Errorf("create jimeng job failed: %w", err)
	}

	// 推送到任务队列
	task := map[string]any{
		"job_id": job.Id,
		"type":   job.Type,
	}
	if err := s.taskQueue.RPush(task); err != nil {
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
	if err := s.UpdateJobStatus(job.Id, model.JimengJobStatusProcessing, ""); err != nil {
		return fmt.Errorf("update job status failed: %w", err)
	}

	// 根据任务类型处理
	switch job.Type {
	case model.JimengJobTypeTextToImage:
		return s.processTextToImage(&job)
	case model.JimengJobTypeImageToImagePortrait:
		return s.processImageToImagePortrait(&job)
	case model.JimengJobTypeImageEdit:
		return s.processImageEdit(&job)
	case model.JimengJobTypeImageEffects:
		return s.processImageEffects(&job)
	case model.JimengJobTypeTextToVideo:
		return s.processTextToVideo(&job)
	case model.JimengJobTypeImageToVideo:
		return s.processImageToVideo(&job)
	default:
		return fmt.Errorf("unsupported task type: %s", job.Type)
	}
}

// processTextToImage 处理文生图任务
func (s *Service) processTextToImage(job *model.JimengJob) error {
	// 解析任务参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.TaskParams), &params); err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("parse task params failed: %v", err))
	}

	// 构建请求
	req := &SubmitTaskRequest{
		ReqKey: job.ReqKey,
		Prompt: job.Prompt,
	}

	// 设置参数
	if seed, ok := params["seed"]; ok {
		if seedVal, err := strconv.ParseInt(fmt.Sprintf("%.0f", seed), 10, 64); err == nil {
			req.Seed = seedVal
		}
	}
	if scale, ok := params["scale"]; ok {
		if scaleVal, ok := scale.(float64); ok {
			req.Scale = scaleVal
		}
	}
	if width, ok := params["width"]; ok {
		if widthVal, ok := width.(float64); ok {
			req.Width = int(widthVal)
		}
	}
	if height, ok := params["height"]; ok {
		if heightVal, ok := height.(float64); ok {
			req.Height = int(heightVal)
		}
	}
	if usePreLlm, ok := params["use_pre_llm"]; ok {
		if usePreLlmVal, ok := usePreLlm.(bool); ok {
			req.UsePreLLM = usePreLlmVal
		}
	}

	// 提交异步任务
	resp, err := s.client.SubmitTask(req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
	}

	if resp.Code != 10000 {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %s", resp.Message))
	}

	// 更新任务ID和原始数据
	rawData, _ := json.Marshal(resp)
	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(map[string]any{
		"task_id":    resp.Data.TaskId,
		"raw_data":   string(rawData),
		"updated_at": time.Now(),
	}).Error; err != nil {
		serviceLogger.Errorf("update jimeng job task_id failed: %v", err)
	}

	// 开始轮询任务状态
	return s.pollTaskStatus(job.Id, resp.Data.TaskId, job.ReqKey)
}

// processImageToImagePortrait 处理图生图人像写真任务
func (s *Service) processImageToImagePortrait(job *model.JimengJob) error {
	// 解析任务参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.TaskParams), &params); err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("parse task params failed: %v", err))
	}

	// 构建请求
	req := &SubmitTaskRequest{
		ReqKey: job.ReqKey,
		Prompt: job.Prompt,
	}

	// 设置图像输入
	if imageInput, ok := params["image_input"].(string); ok {
		req.ImageInput = imageInput
	}

	// 设置其他参数
	if gpen, ok := params["gpen"]; ok {
		if gpenVal, ok := gpen.(float64); ok {
			req.Gpen = gpenVal
		}
	}
	if skin, ok := params["skin"]; ok {
		if skinVal, ok := skin.(float64); ok {
			req.Skin = skinVal
		}
	}
	if skinUnifi, ok := params["skin_unifi"]; ok {
		if skinUnifiVal, ok := skinUnifi.(float64); ok {
			req.SkinUnifi = skinUnifiVal
		}
	}
	if genMode, ok := params["gen_mode"].(string); ok {
		req.GenMode = genMode
	}
	if width, ok := params["width"]; ok {
		if widthVal, ok := width.(float64); ok {
			req.Width = int(widthVal)
		}
	}
	if height, ok := params["height"]; ok {
		if heightVal, ok := height.(float64); ok {
			req.Height = int(heightVal)
		}
	}
	if seed, ok := params["seed"]; ok {
		if seedVal, err := strconv.ParseInt(fmt.Sprintf("%.0f", seed), 10, 64); err == nil {
			req.Seed = seedVal
		}
	}

	// 提交异步任务
	resp, err := s.client.SubmitTask(req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
	}

	if resp.Code != 10000 {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %s", resp.Message))
	}

	// 更新任务ID和原始数据
	rawData, _ := json.Marshal(resp)
	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(map[string]any{
		"task_id":    resp.Data.TaskId,
		"raw_data":   string(rawData),
		"updated_at": time.Now(),
	}).Error; err != nil {
		serviceLogger.Errorf("update jimeng job task_id failed: %v", err)
	}

	// 开始轮询任务状态
	return s.pollTaskStatus(job.Id, resp.Data.TaskId, job.ReqKey)
}

// processImageEdit 处理图像编辑任务
func (s *Service) processImageEdit(job *model.JimengJob) error {
	// 解析任务参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.TaskParams), &params); err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("parse task params failed: %v", err))
	}

	// 构建请求
	req := &SubmitTaskRequest{
		ReqKey: job.ReqKey,
		Prompt: job.Prompt,
	}

	// 设置图像输入
	if imageUrls, ok := params["image_urls"].([]any); ok {
		for _, url := range imageUrls {
			if urlStr, ok := url.(string); ok {
				req.ImageUrls = append(req.ImageUrls, urlStr)
			}
		}
	}
	if binaryData, ok := params["binary_data_base64"].([]any); ok {
		for _, data := range binaryData {
			if dataStr, ok := data.(string); ok {
				req.BinaryDataBase64 = append(req.BinaryDataBase64, dataStr)
			}
		}
	}

	// 设置其他参数
	if seed, ok := params["seed"]; ok {
		if seedVal, err := strconv.ParseInt(fmt.Sprintf("%.0f", seed), 10, 64); err == nil {
			req.Seed = seedVal
		}
	}
	if scale, ok := params["scale"]; ok {
		if scaleVal, ok := scale.(float64); ok {
			req.Scale = scaleVal
		}
	}

	// 提交异步任务
	resp, err := s.client.SubmitTask(req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
	}

	if resp.Code != 10000 {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %s", resp.Message))
	}

	// 更新任务ID和原始数据
	rawData, _ := json.Marshal(resp)
	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(map[string]any{
		"task_id":    resp.Data.TaskId,
		"raw_data":   string(rawData),
		"updated_at": time.Now(),
	}).Error; err != nil {
		serviceLogger.Errorf("update jimeng job task_id failed: %v", err)
	}

	// 开始轮询任务状态
	return s.pollTaskStatus(job.Id, resp.Data.TaskId, job.ReqKey)
}

// processImageEffects 处理图像特效任务
func (s *Service) processImageEffects(job *model.JimengJob) error {
	// 解析任务参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.TaskParams), &params); err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("parse task params failed: %v", err))
	}

	// 构建请求
	req := &SubmitTaskRequest{
		ReqKey: job.ReqKey,
	}

	// 设置图像输入
	if imageInput1, ok := params["image_input1"].(string); ok {
		req.ImageInput1 = imageInput1
	}
	if templateId, ok := params["template_id"].(string); ok {
		req.TemplateId = templateId
	}
	if width, ok := params["width"]; ok {
		if widthVal, ok := width.(float64); ok {
			req.Width = int(widthVal)
		}
	}
	if height, ok := params["height"]; ok {
		if heightVal, ok := height.(float64); ok {
			req.Height = int(heightVal)
		}
	}

	// 提交异步任务
	resp, err := s.client.SubmitTask(req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
	}

	if resp.Code != 10000 {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %s", resp.Message))
	}

	// 更新任务ID和原始数据
	rawData, _ := json.Marshal(resp)
	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(map[string]any{
		"task_id":    resp.Data.TaskId,
		"raw_data":   string(rawData),
		"updated_at": time.Now(),
	}).Error; err != nil {
		serviceLogger.Errorf("update jimeng job task_id failed: %v", err)
	}

	// 开始轮询任务状态
	return s.pollTaskStatus(job.Id, resp.Data.TaskId, job.ReqKey)
}

// processTextToVideo 处理文生视频任务
func (s *Service) processTextToVideo(job *model.JimengJob) error {
	// 解析任务参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.TaskParams), &params); err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("parse task params failed: %v", err))
	}

	// 构建请求
	req := &SubmitTaskRequest{
		ReqKey: job.ReqKey,
		Prompt: job.Prompt,
	}

	// 设置参数
	if seed, ok := params["seed"]; ok {
		if seedVal, err := strconv.ParseInt(fmt.Sprintf("%.0f", seed), 10, 64); err == nil {
			req.Seed = seedVal
		}
	}
	if aspectRatio, ok := params["aspect_ratio"].(string); ok {
		req.AspectRatio = aspectRatio
	}

	// 提交异步任务
	resp, err := s.client.SubmitTask(req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
	}

	if resp.Code != 10000 {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %s", resp.Message))
	}

	// 更新任务ID和原始数据
	rawData, _ := json.Marshal(resp)
	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(map[string]any{
		"task_id":    resp.Data.TaskId,
		"raw_data":   string(rawData),
		"updated_at": time.Now(),
	}).Error; err != nil {
		serviceLogger.Errorf("update jimeng job task_id failed: %v", err)
	}

	// 开始轮询任务状态
	return s.pollTaskStatus(job.Id, resp.Data.TaskId, job.ReqKey)
}

// processImageToVideo 处理图生视频任务
func (s *Service) processImageToVideo(job *model.JimengJob) error {
	// 解析任务参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.TaskParams), &params); err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("parse task params failed: %v", err))
	}

	// 构建请求
	req := &SubmitTaskRequest{
		ReqKey: job.ReqKey,
		Prompt: job.Prompt,
	}

	// 设置图像输入
	if imageUrls, ok := params["image_urls"].([]any); ok {
		for _, url := range imageUrls {
			if urlStr, ok := url.(string); ok {
				req.ImageUrls = append(req.ImageUrls, urlStr)
			}
		}
	}
	if binaryData, ok := params["binary_data_base64"].([]any); ok {
		for _, data := range binaryData {
			if dataStr, ok := data.(string); ok {
				req.BinaryDataBase64 = append(req.BinaryDataBase64, dataStr)
			}
		}
	}

	// 设置其他参数
	if seed, ok := params["seed"]; ok {
		if seedVal, err := strconv.ParseInt(fmt.Sprintf("%.0f", seed), 10, 64); err == nil {
			req.Seed = seedVal
		}
	}
	if aspectRatio, ok := params["aspect_ratio"].(string); ok {
		req.AspectRatio = aspectRatio
	}

	// 提交异步任务
	resp, err := s.client.SubmitTask(req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %v", err))
	}

	if resp.Code != 10000 {
		return s.handleTaskError(job.Id, fmt.Sprintf("submit task failed: %s", resp.Message))
	}

	// 更新任务ID和原始数据
	rawData, _ := json.Marshal(resp)
	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(map[string]any{
		"task_id":    resp.Data.TaskId,
		"raw_data":   string(rawData),
		"updated_at": time.Now(),
	}).Error; err != nil {
		serviceLogger.Errorf("update jimeng job task_id failed: %v", err)
	}

	// 开始轮询任务状态
	return s.pollTaskStatus(job.Id, resp.Data.TaskId, job.ReqKey)
}

// pollTaskStatus 轮询任务状态
func (s *Service) pollTaskStatus(jobId uint, taskId, reqKey string) error {
	maxRetries := 60 // 最大重试次数，60次 * 5秒 = 5分钟
	retryCount := 0

	for retryCount < maxRetries {
		time.Sleep(5 * time.Second) // 等待5秒

		// 查询任务状态
		resp, err := s.client.QueryTask(&QueryTaskRequest{
			ReqKey:  reqKey,
			TaskId:  taskId,
			ReqJson: `{"return_url":true}`,
		})

		if err != nil {
			serviceLogger.Errorf("query jimeng task status failed: %v", err)
			retryCount++
			continue
		}

		// 更新原始数据
		rawData, _ := json.Marshal(resp)
		s.db.Model(&model.JimengJob{}).Where("id = ?", jobId).Update("raw_data", string(rawData))

		if resp.Code != 10000 {
			return s.handleTaskError(jobId, fmt.Sprintf("query task failed: %s", resp.Message))
		}

		switch resp.Data.Status {
		case TaskStatusDone:
			// 任务完成，更新结果
			updates := map[string]any{
				"status":     model.JimengJobStatusCompleted,
				"progress":   100,
				"updated_at": time.Now(),
			}

			// 设置结果URL
			if len(resp.Data.ImageUrls) > 0 {
				updates["img_url"] = resp.Data.ImageUrls[0]
			}
			if resp.Data.VideoUrl != "" {
				updates["video_url"] = resp.Data.VideoUrl
			}

			return s.db.Model(&model.JimengJob{}).Where("id = ?", jobId).Updates(updates).Error

		case TaskStatusInQueue:
			// 任务在队列中
			s.UpdateJobProgress(jobId, 10)

		case TaskStatusGenerating:
			// 任务处理中
			s.UpdateJobProgress(jobId, 50)

		case TaskStatusNotFound, TaskStatusExpired:
			// 任务未找到或已过期
			return s.handleTaskError(jobId, fmt.Sprintf("task not found or expired: %s", resp.Data.Status))

		default:
			serviceLogger.Warnf("unknown task status: %s", resp.Data.Status)
		}

		retryCount++
	}

	// 超时处理
	return s.handleTaskError(jobId, "task timeout")
}

// UpdateJobStatus 更新任务状态
func (s *Service) UpdateJobStatus(jobId uint, status, errMsg string) error {
	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	if errMsg != "" {
		updates["err_msg"] = errMsg
	}
	return s.db.Model(&model.JimengJob{}).Where("id = ?", jobId).Updates(updates).Error
}

// UpdateJobProgress 更新任务进度
func (s *Service) UpdateJobProgress(jobId uint, progress int) error {
	return s.db.Model(&model.JimengJob{}).Where("id = ?", jobId).Updates(map[string]any{
		"progress":   progress,
		"updated_at": time.Now(),
	}).Error
}

// handleTaskError 处理任务错误
func (s *Service) handleTaskError(jobId uint, errMsg string) error {
	serviceLogger.Errorf("Jimeng task error (job_id: %d): %s", jobId, errMsg)
	return s.UpdateJobStatus(jobId, model.JimengJobStatusFailed, errMsg)
}

// GetJob 获取任务
func (s *Service) GetJob(jobId uint) (*model.JimengJob, error) {
	var job model.JimengJob
	if err := s.db.First(&job, jobId).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// GetUserJobs 获取用户任务列表
func (s *Service) GetUserJobs(userId uint, page, pageSize int) ([]*model.JimengJob, int64, error) {
	var jobs []*model.JimengJob
	var total int64

	query := s.db.Model(&model.JimengJob{}).Where("user_id = ?", userId)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// GetPendingTaskCount 获取用户未完成任务数量
func (s *Service) GetPendingTaskCount(userId uint) (int64, error) {
	var count int64
	err := s.db.Model(&model.JimengJob{}).Where("user_id = ? AND status IN (?)", userId,
		[]string{model.JimengJobStatusPending, model.JimengJobStatusProcessing}).Count(&count).Error
	return count, err
}

// DeleteJob 删除任务
func (s *Service) DeleteJob(jobId uint, userId uint) error {
	return s.db.Where("id = ? AND user_id = ?", jobId, userId).Delete(&model.JimengJob{}).Error
}

// PushTaskToQueue 推送任务到队列
func (s *Service) PushTaskToQueue(task map[string]interface{}) error {
	return s.taskQueue.RPush(task)
}
