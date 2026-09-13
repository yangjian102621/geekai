package jimeng

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"

	"geekai/core/types"
	"geekai/log"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store"
	"geekai/store/model"
	"geekai/utils"

	"github.com/go-redis/redis/v8"
)

var logger = log.GetLogger()

const seedanceOfficialBaseURL = "https://ark.cn-beijing.volces.com/api/v3"

const (
	jimengMediaRepairInterval = 60 * time.Second
	jimengMediaRepairBatch    = 50
)

// Service 即梦服务（合并了消费者功能）
type Service struct {
	db          *gorm.DB
	redis       *redis.Client
	taskQueue   *store.RedisQueue
	client      *Client
	ctx         context.Context
	cancel      context.CancelFunc
	running     bool
	uploader    *oss.UploaderManager
	userService *service.UserService
}

// NewService 创建即梦服务
func NewService(db *gorm.DB, redisCli *redis.Client, uploader *oss.UploaderManager, client *Client, userService *service.UserService) *Service {
	taskQueue := store.NewRedisQueue("JimengTaskQueue", redisCli)
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		db:          db,
		redis:       redisCli,
		taskQueue:   taskQueue,
		client:      client,
		ctx:         ctx,
		cancel:      cancel,
		running:     false,
		uploader:    uploader,
		userService: userService,
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
	go s.runJimengSuccessMediaRepairLoop()
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
		Status:    types.JMTaskStatusSubmited,
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

	// 解析任务参数
	var req types.JimengTaskRequest
	err := utils.JsonDecode(job.Params, &req)
	if err != nil {
		return fmt.Errorf("parse task params failed: %w", err)
	}

	// 数字人任务，先识别主体
	if req.TaskType == types.JMTaskTypeVirtualHuman {
		if err := s.client.AvatarRecognition(req.ImageUrls[0], req.RecognizeKey); err != nil {
			return s.handleTaskError(job.Id, fmt.Sprintf("avatar recognition failed: %v", err))
		}
	}

	// Seedream 同步生图（Ark）
	if IsSeedreamReqKey(req.ReqKey) {
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

			if len(resp.Data) == 0 || resp.Data[0] == nil || resp.Data[0].Url == nil || strings.TrimSpace(*resp.Data[0].Url) == "" {
				_ = s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Update("raw_data", string(rawData)).Error
				_ = s.handleTaskError(job.Id, "seedream response has no image url")
				return
			}

			remoteURL := strings.TrimSpace(*resp.Data[0].Url)
			ext := filepath.Ext(strings.Split(remoteURL, "?")[0])
			if ext == "" {
				ext = ".png"
			}

			// 更新任务状态
			updates["status"] = types.JMTaskStatusSuccess
			// 转存到本地/OSS（失败时回退为官方临时 URL，与即梦异步任务一致）
			imgURL, err := s.uploader.GetUploadHandler().PutUrlFile(remoteURL, ext, false)
			if err != nil {
				logger.Errorf("jimeng seedream upload image failed, job_id=%d: %v", job.Id, err)
				imgURL = remoteURL
			}
			updates["img_url"] = imgURL
			updates["progress"] = 100
			s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(updates)
		}()
		return nil
	}

	// Seedance 视频任务（DoubaoAdapter）
	if IsSeedanceReqKey(req.ReqKey) {
		return s.submitSeedanceTask(job.Id, &req)
	}

	// 其他请求走即梦 Visual 异步任务
	params, err := s.buildTaskRequest(&req)
	if err != nil {
		return s.handleTaskError(job.Id, fmt.Sprintf("build task request failed: %v", err))
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
		"status":     types.JMTaskStatusInQueue, // 把任务状态改成排队中，以便开启轮询
		"updated_at": time.Now(),
	}).Error; err != nil {
		logger.Errorf("update jimeng job task_id failed: %v", err)
	}

	return nil
}

func (s *Service) submitSeedanceTask(jobId uint, req *types.JimengTaskRequest) error {
	jimengConfig, err := s.getJimengConfig()
	if err != nil {
		return s.handleTaskError(jobId, fmt.Sprintf("load jimeng config failed: %v", err))
	}

	content := s.buildSeedanceContent(req)
	if len(content) == 0 {
		return s.handleTaskError(jobId, "seedance content 不能为空")
	}

	payload := map[string]any{
		"model":      req.ReqKey,
		"content":    content,
		"duration":   req.Duration,
		"ratio":      req.AspectRatio,
		"resolution": req.Resolution,
	}
	// 兼容旧参数：0 让官方走默认值
	if req.Duration == 0 {
		delete(payload, "duration")
	}
	if req.AspectRatio == "" {
		delete(payload, "ratio")
	}
	if req.Resolution == "" {
		delete(payload, "resolution")
	}
	if req.ReturnLastFrame {
		payload["return_last_frame"] = req.ReturnLastFrame
	}
	if req.Watermark != nil {
		payload["watermark"] = *req.Watermark
	}
	if req.GenerateAudio != nil {
		payload["generate_audio"] = *req.GenerateAudio
	}

	resp, rawData, err := s.callSeedanceCreate(payload, jimengConfig)
	if err != nil {
		return s.handleTaskError(jobId, fmt.Sprintf("submit seedance task failed: %v", err))
	}

	logger.Debugf("seedance create response: %+v", resp)

	if err := s.db.Model(&model.JimengJob{}).Where("id = ?", jobId).Updates(map[string]any{
		"task_id":    resp.TaskID,
		"raw_data":   rawData,
		"status":     types.JMTaskStatusInQueue,
		"updated_at": time.Now(),
	}).Error; err != nil {
		logger.Errorf("update seedance task_id failed: %v", err)
	}
	return nil
}

func (s *Service) buildSeedanceContent(req *types.JimengTaskRequest) []types.JMContentItem {
	if len(req.Content) > 0 {
		return req.Content
	}

	content := make([]types.JMContentItem, 0, 4)
	if req.Prompt != "" {
		content = append(content, types.JMContentItem{
			Type: "text",
			Text: req.Prompt,
		})
	}

	if len(req.ImageUrls) > 0 {
		for index, imageURL := range req.ImageUrls {
			if imageURL == "" {
				continue
			}
			role := "reference_image"
			if len(req.ImageUrls) == 1 {
				role = "first_frame"
			} else if len(req.ImageUrls) == 2 {
				if index == 0 {
					role = "first_frame"
				} else {
					role = "last_frame"
				}
			}
			content = append(content, types.JMContentItem{
				Type: "image_url",
				ImageURL: &types.JMAssetRef{
					URL: imageURL,
				},
				Role: role,
			})
		}
	}

	if req.VideoURL != "" {
		content = append(content, types.JMContentItem{
			Type: "video_url",
			VideoURL: &types.JMAssetRef{
				URL: req.VideoURL,
			},
			Role: "reference_video",
		})
	}

	if req.AudioURL != "" {
		content = append(content, types.JMContentItem{
			Type: "audio_url",
			AudioURL: &types.JMAssetRef{
				URL: req.AudioURL,
			},
			Role: "reference_audio",
		})
	}

	return content
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
		if v, ok := duration.(int); ok {
			params["frames"] = v*24 + 1
		} else if v, ok := duration.(float64); ok {
			params["frames"] = int(v*24) + 1
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
		// 找出排队中和处理中的任务进行轮询
		s.db.Where("status IN (?)", []types.JMTaskStatus{
			types.JMTaskStatusGenerating,
			types.JMTaskStatusInQueue}).Find(&jobs)

		for _, job := range jobs {
			// 任务超时处理
			if job.UpdatedAt.Before(time.Now().Add(-10 * time.Minute)) {
				s.handleTaskError(job.Id, "task timeout")
				continue
			}

			// Seedream 为同步任务，不需要轮询
			if IsSeedreamReqKey(job.ReqKey) {
				continue
			}

			if IsSeedanceReqKey(job.ReqKey) {
				if err := s.pollSeedanceTask(&job); err != nil {
					s.handleTaskError(job.Id, err.Error())
				}
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
					"progress":   100,
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

		// 找出失败的任务，并恢复其扣减积分
		s.db.Where("status = ?", types.JMTaskStatusFailed).Where("power > ?", 0).Find(&jobs)
		for _, job := range jobs {
			err := s.userService.IncreasePower(job.UserId, job.Power, model.PowerLog{
				Type:   types.PowerRefund,
				Model:  job.ReqKey,
				Remark: fmt.Sprintf("任务失败，退回积分。任务ID：%d", job.Id),
			})
			if err != nil {
				continue
			}
			// 更新任务状态
			s.db.Model(&job).UpdateColumn("power", 0)
		}

		time.Sleep(5 * time.Second)

	}

}

func (s *Service) pollSeedanceTask(job *model.JimengJob) error {
	jimengConfig, err := s.getJimengConfig()
	if err != nil {
		return fmt.Errorf("load jimeng config failed: %w", err)
	}

	resp, rawData, err := s.callSeedanceQuery(job.TaskId, jimengConfig)
	if err != nil {
		return fmt.Errorf("query seedance task failed: %w", err)
	}

	logger.Debugf("seedance query response: %+v", resp)

	s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Update("raw_data", rawData)

	switch resp.Status {
	case "succeeded":
		updates := map[string]any{
			"status":     types.JMTaskStatusSuccess,
			"updated_at": time.Now(),
			"progress":   100,
		}
		if resp.Content.VideoURL != "" {
			videoURL, upErr := s.uploader.GetUploadHandler().PutUrlFile(resp.Content.VideoURL, ".mp4", false)
			if upErr != nil {
				logger.Errorf("upload seedance video failed: %v", upErr)
				videoURL = resp.Content.VideoURL
			}
			updates["video_url"] = videoURL
		}
		return s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(updates).Error
	case "queued", "running":
		return s.UpdateJobStatus(job.Id, types.JMTaskStatusGenerating, "")
	case "failed", "cancelled":
		errMsg := resp.Error
		if errMsg == "" {
			errMsg = "seedance task failed"
		}
		return fmt.Errorf("%s", errMsg)
	default:
		return nil
	}
}

type seedanceCreateResponse struct {
	ID         string `json:"id"`
	PlatformID string `json:"platform_id"`
}

type seedanceCreateResult struct {
	TaskID string
}

type seedanceQueryResponse struct {
	ID         string `json:"id"`
	PlatformID string `json:"platform_id"`
	Status     string `json:"status"`
	Error      string `json:"error"`
	Content    struct {
		VideoURL string `json:"video_url"`
	} `json:"content"`
}

func (s *Service) callSeedanceCreate(payload map[string]any, jimengConfig *types.JimengConfig) (*seedanceCreateResult, string, error) {
	if jimengConfig == nil || strings.TrimSpace(jimengConfig.ApiKey) == "" {
		return nil, "", fmt.Errorf("jimeng api key 未配置")
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}
	url := fmt.Sprintf("%s/contents/generations/tasks", seedanceOfficialBaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(jimengConfig.ApiKey))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, string(raw), fmt.Errorf("status=%d body=%s", resp.StatusCode, string(raw))
	}
	var parsed seedanceCreateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, string(raw), err
	}
	taskID := parsed.PlatformID
	if taskID == "" {
		taskID = parsed.ID
	}
	if taskID == "" {
		return nil, string(raw), fmt.Errorf("seedance create 响应缺少 task id: %s", string(raw))
	}
	return &seedanceCreateResult{TaskID: taskID}, string(raw), nil
}

func (s *Service) callSeedanceQuery(taskID string, jimengConfig *types.JimengConfig) (*seedanceQueryResponse, string, error) {
	if jimengConfig == nil || strings.TrimSpace(jimengConfig.ApiKey) == "" {
		return nil, "", fmt.Errorf("jimeng api key 未配置")
	}
	url := fmt.Sprintf("%s/contents/generations/tasks/%s", seedanceOfficialBaseURL, taskID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(jimengConfig.ApiKey))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, string(raw), fmt.Errorf("status=%d body=%s", resp.StatusCode, string(raw))
	}
	var parsed seedanceQueryResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, string(raw), err
	}
	return &parsed, string(raw), nil
}

func (s *Service) getJimengConfig() (*types.JimengConfig, error) {
	var configRow model.Config
	if err := s.db.Where("name = ?", types.ConfigKeyJimeng).First(&configRow).Error; err != nil {
		return nil, err
	}
	var jimengConfig types.JimengConfig
	if err := utils.JsonDecode(configRow.Value, &jimengConfig); err != nil {
		return nil, err
	}
	if strings.TrimSpace(jimengConfig.ApiKey) == "" {
		return nil, fmt.Errorf("jimeng api key 未配置")
	}
	return &jimengConfig, nil
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

// runJimengSuccessMediaRepairLoop 定时修复：状态已是 success 但进度未满且媒体地址均为空的任务，从 raw_data 重新解析并转存。
func (s *Service) runJimengSuccessMediaRepairLoop() {
	ticker := time.NewTicker(jimengMediaRepairInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			logger.Info("Jimeng success-media repair loop stopped")
			return
		case <-ticker.C:
			s.repairJimengSuccessJobsMediaOnce()
		}
	}
}

func (s *Service) repairJimengSuccessJobsMediaOnce() {
	var jobs []model.JimengJob
	err := s.db.Where("status = ?", types.JMTaskStatusSuccess).
		Where("progress <> ?", 100).
		Where("(COALESCE(img_url, '') = ? AND COALESCE(video_url, '') = ?)", "", "").
		Where("raw_data IS NOT NULL AND raw_data <> ?", "").
		Order("id ASC").
		Limit(jimengMediaRepairBatch).
		Find(&jobs).Error
	if err != nil {
		logger.Errorf("jimeng media repair query failed: %v", err)
		return
	}
	if len(jobs) == 0 {
		return
	}
	for i := range jobs {
		job := jobs[i]
		remoteImg, remoteVid := parseJimengRawMediaURLs(job.RawData)
		if strings.TrimSpace(remoteImg) == "" && strings.TrimSpace(remoteVid) == "" {
			logger.Warnf("jimeng media repair: job_id=%d no media url in raw_data", job.Id)
			continue
		}
		updates := map[string]any{
			"updated_at": time.Now(),
			"progress":   100,
		}
		if strings.TrimSpace(remoteImg) != "" {
			updates["img_url"] = s.putJimengRemoteMedia(strings.TrimSpace(remoteImg), ".png", false)
		}
		if strings.TrimSpace(remoteVid) != "" {
			updates["video_url"] = s.putJimengRemoteMedia(strings.TrimSpace(remoteVid), ".mp4", true)
		}
		if err := s.db.Model(&model.JimengJob{}).Where("id = ?", job.Id).Updates(updates).Error; err != nil {
			logger.Errorf("jimeng media repair update failed job_id=%d: %v", job.Id, err)
		} else {
			logger.Infof("jimeng media repair ok job_id=%d", job.Id)
		}
	}
}

// parseJimengRawMediaURLs 从 raw_data 解析远程图片/视频地址（即梦异步、Ark Seedream、Seedance）。
func parseJimengRawMediaURLs(raw string) (remoteImg, remoteVideo string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ""
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &root); err != nil {
		return "", ""
	}

	if dataRaw, ok := root["data"]; ok {
		dataBytes := []byte(dataRaw)
		trimmed := bytes.TrimSpace(dataBytes)
		if len(trimmed) > 0 && trimmed[0] == '[' {
			var items []struct {
				Url *string `json:"url"`
			}
			if json.Unmarshal(dataBytes, &items) == nil {
				for _, it := range items {
					if it.Url != nil && strings.TrimSpace(*it.Url) != "" {
						return strings.TrimSpace(*it.Url), ""
					}
				}
			}
		} else {
			var qd struct {
				ImageUrls []string `json:"image_urls"`
				VideoUrl  string   `json:"video_url"`
			}
			if json.Unmarshal(dataBytes, &qd) == nil {
				img := ""
				if len(qd.ImageUrls) > 0 {
					img = strings.TrimSpace(qd.ImageUrls[0])
				}
				vid := strings.TrimSpace(qd.VideoUrl)
				if img != "" || vid != "" {
					return img, vid
				}
			}
		}
	}

	if contentRaw, ok := root["content"]; ok {
		var c struct {
			VideoURL string `json:"video_url"`
		}
		if json.Unmarshal(contentRaw, &c) == nil && strings.TrimSpace(c.VideoURL) != "" {
			return "", strings.TrimSpace(c.VideoURL)
		}
	}

	return "", ""
}

func (s *Service) putJimengRemoteMedia(remote, fallbackExt string, isVideo bool) string {
	remote = strings.TrimSpace(remote)
	if remote == "" {
		return ""
	}
	ext := fallbackExt
	if !isVideo {
		u := remote
		if i := strings.Index(u, "?"); i >= 0 {
			u = u[:i]
		}
		if e := filepath.Ext(u); e != "" {
			ext = e
		}
	} else {
		ext = ".mp4"
	}
	out, err := s.uploader.GetUploadHandler().PutUrlFile(remote, ext, false)
	if err != nil {
		logger.Errorf("jimeng putJimengRemoteMedia failed: %v", err)
		return remote
	}
	return out
}

// GetJob 获取任务
func (s *Service) GetJob(jobId uint) (*model.JimengJob, error) {
	var job model.JimengJob
	if err := s.db.First(&job, jobId).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
