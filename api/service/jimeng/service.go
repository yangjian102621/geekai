package jimeng

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	logger2 "geekai/logger"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"

	"geekai/core/types"

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
func NewService(db *gorm.DB, redisCli *redis.Client, uploader *oss.UploaderManager) *Service {
	taskQueue := store.NewRedisQueue("JimengTaskQueue", redisCli)
	// 从数据库加载配置
	var config model.Config
	db.Where("name = ?", "Jimeng").First(&config)
	var jimengConfig types.JimengConfig
	if config.Id > 0 {
		_ = utils.JsonDecode(config.Value, &jimengConfig)
	}
	client := NewClient(jimengConfig.AccessKey, jimengConfig.SecretKey)

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
		s.UpdateJobStatus(jobId, model.JMTaskStatusFailed, err.Error())
	} else {
		logger.Infof("Jimeng task processed successfully: job_id=%d", jobId)
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
		Status:     model.JMTaskStatusInQueue,
		Power:      req.Power,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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
	if err := s.UpdateJobStatus(job.Id, model.JMTaskStatusGenerating, ""); err != nil {
		return fmt.Errorf("update job status failed: %w", err)
	}

	// 构建请求并提交任务
	req, err := s.buildTaskRequest(&job)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("build task request failed: %v", err))
	}

	logger.Infof("提交即梦任务: %+v", req)

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
		logger.Errorf("update jimeng job task_id failed: %v", err)
	}

	return nil
}

// buildTaskRequest 构建任务请求（统一的参数解析）
func (s *Service) buildTaskRequest(job *model.JimengJob) (*SubmitTaskRequest, error) {
	// 解析任务参数
	var params map[string]any
	if err := json.Unmarshal([]byte(job.TaskParams), &params); err != nil {
		return nil, fmt.Errorf("parse task params failed: %w", err)
	}

	// 构建基础请求
	req := &SubmitTaskRequest{
		ReqKey: job.ReqKey,
		Prompt: job.Prompt,
	}

	// 根据任务类型设置特定参数
	switch job.Type {
	case model.JMTaskTypeTextToImage:
		s.setTextToImageParams(req, params)
	case model.JMTaskTypeImageToImage:
		s.setImageToImageParams(req, params)
	case model.JMTaskTypeImageEdit:
		s.setImageEditParams(req, params)
	case model.JMTaskTypeImageEffects:
		s.setImageEffectsParams(req, params)
	case model.JMTaskTypeTextToVideo:
		s.setTextToVideoParams(req, params)
	case model.JMTaskTypeImageToVideo:
		s.setImageToVideoParams(req, params)
	default:
		return nil, fmt.Errorf("unsupported task type: %s", job.Type)
	}

	return req, nil
}

// setTextToImageParams 设置文生图参数
func (s *Service) setTextToImageParams(req *SubmitTaskRequest, params map[string]any) {
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
}

// setImageToImageParams 设置图生图参数
func (s *Service) setImageToImageParams(req *SubmitTaskRequest, params map[string]any) {
	if imageInput, ok := params["image_input"].(string); ok {
		req.ImageInput = imageInput
	}
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
	s.setCommonParams(req, params) // 复用通用参数
}

// setImageEditParams 设置图像编辑参数
func (s *Service) setImageEditParams(req *SubmitTaskRequest, params map[string]any) {
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
	if scale, ok := params["scale"]; ok {
		if scaleVal, ok := scale.(float64); ok {
			req.Scale = scaleVal
		}
	}
	s.setCommonParams(req, params)
}

// setImageEffectsParams 设置图像特效参数
func (s *Service) setImageEffectsParams(req *SubmitTaskRequest, params map[string]any) {
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
}

// setTextToVideoParams 设置文生视频参数
func (s *Service) setTextToVideoParams(req *SubmitTaskRequest, params map[string]any) {
	if aspectRatio, ok := params["aspect_ratio"].(string); ok {
		req.AspectRatio = aspectRatio
	}
	s.setCommonParams(req, params)
}

// setImageToVideoParams 设置图生视频参数
func (s *Service) setImageToVideoParams(req *SubmitTaskRequest, params map[string]any) {
	s.setImageEditParams(req, params) // 复用图像编辑的参数设置
	if aspectRatio, ok := params["aspect_ratio"].(string); ok {
		req.AspectRatio = aspectRatio
	}
}

// setCommonParams 设置通用参数（seed, width, height等）
func (s *Service) setCommonParams(req *SubmitTaskRequest, params map[string]any) {
	if seed, ok := params["seed"]; ok {
		if seedVal, err := strconv.ParseInt(fmt.Sprintf("%.0f", seed), 10, 64); err == nil {
			req.Seed = seedVal
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
}

// pollTaskStatus 轮询任务状态
func (s *Service) pollTaskStatus() {

	for {
		var jobs []model.JimengJob
		s.db.Where("status IN (?)", []model.JMTaskStatus{model.JMTaskStatusGenerating, model.JMTaskStatusInQueue}).Find(&jobs)
		if len(jobs) == 0 {
			logger.Debugf("no jimeng task to poll, sleep 10s")
			time.Sleep(10 * time.Second)
			continue
		}

		for _, job := range jobs {
			// 任务超时处理
			if job.UpdatedAt.Before(time.Now().Add(-5 * time.Minute)) {
				s.handleTaskError(job.Id, "task timeout")
				continue
			}

			// 查询任务状态
			resp, err := s.client.QueryTask(&QueryTaskRequest{
				ReqKey:  job.ReqKey,
				TaskId:  job.TaskId,
				ReqJson: `{"return_url":true}`,
			})

			if err != nil {
				logger.Errorf("query jimeng task status failed: %v", err)
				continue
			}

			// 更新原始数据
			rawData, _ := json.Marshal(resp)
			s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Update("raw_data", string(rawData))

			if resp.Code != 10000 {
				s.handleTaskError(job.Id, fmt.Sprintf("query task failed: %s", resp.Message))
				continue
			}

			switch resp.Data.Status {
			case model.JMTaskStatusDone:
				// 判断任务是否成功
				if resp.Message != "Success" {
					s.handleTaskError(job.Id, fmt.Sprintf("task failed: %s", resp.Data.AlgorithmBaseResp.StatusMessage))
					continue
				}

				// 任务完成，更新结果
				updates := map[string]any{
					"status":     model.JMTaskStatusSuccess,
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
			case model.JMTaskStatusInQueue, model.JMTaskStatusGenerating:
				// 任务处理中
				s.UpdateJobStatus(job.Id, model.JMTaskStatusGenerating, "")

			case model.JMTaskStatusNotFound:
				// 任务未找到
				s.handleTaskError(job.Id, "task not found")

			case model.JMTaskStatusExpired:
				// 任务过期
				s.handleTaskError(job.Id, "task expired")

			default:
				logger.Warnf("unknown task status: %s", resp.Data.Status)
			}

		}

		time.Sleep(5 * time.Second)

	}

}

// UpdateJobStatus 更新任务状态
func (s *Service) UpdateJobStatus(jobId uint, status model.JMTaskStatus, errMsg string) error {
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
	return s.UpdateJobStatus(jobId, model.JMTaskStatusFailed, errMsg)
}

// PushTaskToQueue 推送任务到队列（用于手动重试）
func (s *Service) PushTaskToQueue(jobId uint) error {
	return s.taskQueue.RPush(jobId)
}

// GetTaskStats 获取任务统计信息
func (s *Service) GetTaskStats() (map[string]any, error) {
	type StatResult struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
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
		result[stat.Status] = stat.Count
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

// testConnection 测试即梦AI连接
func (s *Service) testConnection(accessKey, secretKey string) error {
	testClient := NewClient(accessKey, secretKey)

	// 使用一个简单的查询任务来测试连接
	testReq := &QueryTaskRequest{
		ReqKey: "test_connection",
		TaskId: "test_task_id_12345",
	}

	_, err := testClient.QueryTask(testReq)
	// 即使任务不存在，只要不是认证错误就说明连接正常
	if err != nil {
		// 检查是否是认证错误
		if strings.Contains(err.Error(), "InvalidAccessKey") {
			return fmt.Errorf("认证失败，请检查AccessKey和SecretKey是否正确")
		}
		// 其他错误（如任务不存在）说明连接正常
		return nil
	}
	return nil
}

// UpdateClientConfig 更新客户端配置
func (s *Service) UpdateClientConfig(accessKey, secretKey string) error {
	// 创建新的客户端
	newClient := NewClient(accessKey, secretKey)

	// 测试新客户端是否可用
	err := s.testConnection(accessKey, secretKey)
	if err != nil {
		return err
	}

	// 更新客户端
	s.client = newClient
	return nil
}

var defaultPower = types.JimengPower{
	TextToImage:  20,
	ImageToImage: 20,
	ImageEdit:    20,
	ImageEffects: 20,
	TextToVideo:  300,
	ImageToVideo: 300,
}

// GetConfig 获取即梦AI配置
func (s *Service) GetConfig() *types.JimengConfig {
	var config model.Config
	err := s.db.Where("name", "jimeng").First(&config).Error
	if err != nil {
		// 如果配置不存在，返回默认配置
		return &types.JimengConfig{
			AccessKey: "",
			SecretKey: "",
			Power:     defaultPower,
		}
	}

	var jimengConfig types.JimengConfig
	err = utils.JsonDecode(config.Value, &jimengConfig)
	if err != nil {
		return &types.JimengConfig{
			AccessKey: "",
			SecretKey: "",
			Power:     defaultPower,
		}
	}

	return &jimengConfig
}
