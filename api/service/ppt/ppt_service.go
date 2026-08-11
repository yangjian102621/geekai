package ppt

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// ErrInsufficientPower 用户算力不足以完成本次 PPT 任务（由 handler 映射文案）
var ErrInsufficientPower = errors.New("insufficient power for ppt task")

var (
	// ErrPptTaskNotFound 表示任务不存在
	ErrPptTaskNotFound = errors.New("ppt task not found")
	// ErrPptTaskNotDeletable 表示任务状态不允许删除
	ErrPptTaskNotDeletable = errors.New("ppt task not deletable")
	// ErrPptTaskBusy 任务正在处理中
	ErrPptTaskBusy = errors.New("ppt task is processing")
	// ErrPptTaskNotResumable 无法继续生成（无分镜占位或已完成）
	ErrPptTaskNotResumable = errors.New("ppt task cannot be resumed")
	// ErrPptSlideNotFound 指定 slide_index 不存在
	ErrPptSlideNotFound = errors.New("ppt slide not found")
	// ErrPptSlideNoImage 该页尚无配图
	ErrPptSlideNoImage = errors.New("ppt slide has no image")
	// ErrPptInvalidVersionIndex 历史版本下标无效
	ErrPptInvalidVersionIndex = errors.New("invalid slide version index")
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

// SlideData 单页 PPT 数据
type SlideData struct {
	SlideIndex   int                       `json:"slide_index"`
	Theme        string                    `json:"theme"`
	Title        string                    `json:"title"`
	Points       []string                  `json:"points"`
	ImagePrompt  string                    `json:"image_prompt"`
	ImageURL     string                    `json:"image_url"`
	ImageHistory []vo.PPTSlideImageVersion `json:"image_history,omitempty"`
}

// Task PPT 生成任务（用于业务层与 API 返回，持久化在 DB）
type Task struct {
	TaskID    string      `json:"task_id"`
	UserID    uint        `json:"user_id"`
	Status    TaskStatus  `json:"status"`
	Content   string      `json:"content"`
	Prompt    string      `json:"prompt"`
	Language  string      `json:"language"`
	Mode      string      `json:"mode"`
	Pages     int         `json:"pages"`
	Total     int         `json:"total_slides"`
	Completed int         `json:"completed_slides"`
	Slides    []SlideData `json:"slides"`
	Title     string      `json:"title"`
	Thumb     string      `json:"thumb"`

	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TaskSummaryMap 返回任务摘要 map（公共字段），供 handler 补充独有字段后返回。
func (t *Task) TaskSummaryMap() map[string]any {
	return map[string]any{
		"task_id":          t.TaskID,
		"status":           t.Status,
		"total_slides":     t.Total,
		"completed_slides": t.Completed,
		"created_at":       t.CreatedAt.Unix(),
		"updated_at":       t.UpdatedAt.Unix(),
		"title":            t.Title,
		"thumb":            t.Thumb,
	}
}

// Progress 任务进度信息
type Progress struct {
	Total     int `json:"total_slides"`
	Completed int `json:"completed_slides"`
}

// PptService PPT 任务与生成流程（持久化、LLM 分镜、生图、转存、算力）
type PptService struct {
	db            *gorm.DB
	userService   *service.UserService
	uploadManager *oss.UploaderManager
	llm           *LLMClient
	// slidesLock 全局互斥：串行化 slides JSON 的写库，避免并发覆盖（写库极短，可接受排队）
	slidesLock sync.Mutex
}

// NewPptService 创建 PptService
func NewPptService(db *gorm.DB, userService *service.UserService, uploadManager *oss.UploaderManager) *PptService {
	return &PptService{
		db:            db,
		userService:   userService,
		uploadManager: uploadManager,
		llm:           NewLLMClient(),
	}
}

// GenerateNotebookContent 将原始文档文本提炼成 PPT 可用的大纲 content。
func (s *PptService) GenerateNotebookContent(ctx context.Context, cfg types.PPTConfig, rawDocText, designPrompt, language string) (string, error) {
	if s.llm == nil {
		s.llm = NewLLMClient()
	}
	return s.llm.GenerateNotebookContent(ctx, cfg, rawDocText, designPrompt, language)
}

func slideToVO(s SlideData) vo.PPTSlideData {
	return vo.PPTSlideData{
		SlideIndex:   s.SlideIndex,
		Theme:        s.Theme,
		Title:        s.Title,
		Points:       s.Points,
		ImagePrompt:  s.ImagePrompt,
		ImageURL:     s.ImageURL,
		ImageHistory: vo.PPTSlideImageVersions(s.ImageHistory),
	}
}

func voToSlide(s vo.PPTSlideData) SlideData {
	return SlideData{
		SlideIndex:   s.SlideIndex,
		Theme:        s.Theme,
		Title:        s.Title,
		Points:       s.Points,
		ImagePrompt:  s.ImagePrompt,
		ImageURL:     s.ImageURL,
		ImageHistory: []vo.PPTSlideImageVersion(s.ImageHistory),
	}
}

func voSlidesToBiz(slides vo.PPTSlides) []SlideData {
	out := make([]SlideData, len(slides))
	for i, sv := range slides {
		out[i] = voToSlide(sv)
	}
	return out
}

func voSlidesToBizNormalized(slides vo.PPTSlides) []SlideData {
	out := make([]SlideData, len(slides))
	for i, sv := range slides {
		sd := voToSlide(sv)
		normalizeSlideImageHistory(&sd)
		out[i] = sd
	}
	return out
}

// normalizeSlideImageHistory 旧数据仅有 image_url 时补全 image_history，便于前端展示历史
func normalizeSlideImageHistory(s *SlideData) {
	if strings.TrimSpace(s.ImageURL) != "" && len(s.ImageHistory) == 0 {
		s.ImageHistory = []vo.PPTSlideImageVersion{
			{ImageURL: s.ImageURL, Prompt: strings.TrimSpace(s.ImagePrompt)},
		}
	}
}

// DerivePPTThumbFromSlides 按 slide_index 升序取第一张有图 URL（与 vo.PPTSlides 规则一致）
func DerivePPTThumbFromSlides(slides []SlideData) string {
	if len(slides) == 0 {
		return ""
	}
	cp := make([]SlideData, len(slides))
	copy(cp, slides)
	sort.Slice(cp, func(i, j int) bool {
		return cp[i].SlideIndex < cp[j].SlideIndex
	})
	for _, s := range cp {
		if strings.TrimSpace(s.ImageURL) != "" {
			return s.ImageURL
		}
	}
	return ""
}

func truncateTitleRunes(s string, max int) string {
	if max <= 0 {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

func taskToModel(t *Task) *model.PPTJob {
	now := time.Now()
	job := &model.PPTJob{
		TaskId:          t.TaskID,
		UserId:          t.UserID,
		Status:          string(t.Status),
		ErrMsg:          t.ErrorMessage,
		Prompt:          t.Prompt,
		Title:           t.Title,
		Thumb:           t.Thumb,
		Content:         t.Content,
		Params:          vo.PPTParams{Language: t.Language, Mode: t.Mode, Pages: t.Pages},
		Slides:          nil,
		TotalSlides:     t.Total,
		CompletedSlides: t.Completed,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if t.CreatedAt.IsZero() {
		job.CreatedAt = now
		job.UpdatedAt = now
	} else {
		job.CreatedAt = t.CreatedAt
		job.UpdatedAt = t.UpdatedAt
	}
	return job
}

func modelToTask(j *model.PPTJob) *Task {
	slides := voSlidesToBizNormalized(j.Slides)
	return &Task{
		TaskID:       j.TaskId,
		UserID:       j.UserId,
		Status:       TaskStatus(j.Status),
		Content:      j.Content,
		Prompt:       j.Prompt,
		Title:        j.Title,
		Thumb:        j.Thumb,
		Language:     j.Params.Language,
		Mode:         j.Params.Mode,
		Pages:        j.Params.Pages,
		Total:        j.TotalSlides,
		Completed:    j.CompletedSlides,
		Slides:       slides,
		ErrorMessage: j.ErrMsg,
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
	}
}

// BuildPendingTask 校验算力与页数，组装待写入的 Task（未落库）
func (s *PptService) BuildPendingTask(taskID string, userID uint, userPower int, content, prompt, language, mode string, reqPages int) (*Task, types.PPTConfig, error) {
	cfg, err := s.loadPPTConfig()
	if err != nil {
		return nil, cfg, err
	}

	// 目标页数：用户指定时取 min(请求页数, 服务端上限)；未指定(0)时按服务端上限作为默认生成规模
	effectivePages := cfg.MaxSlidesPerTask
	if reqPages > 0 {
		effectivePages = reqPages
		if effectivePages > cfg.MaxSlidesPerTask {
			effectivePages = cfg.MaxSlidesPerTask
		}
	}

	estimatePower := effectivePages * cfg.PowerCostPerSlide
	if estimatePower > 0 && userPower < estimatePower {
		return nil, cfg, ErrInsufficientPower
	}

	effectiveMode := mode
	if effectiveMode != "detailed" && effectiveMode != "slides" {
		effectiveMode = "slides"
	}

	task := &Task{
		TaskID:   taskID,
		UserID:   userID,
		Status:   TaskStatusPending,
		Content:  content,
		Prompt:   prompt,
		Language: language,
		Pages:    effectivePages,
		Mode:     effectiveMode,
	}
	return task, cfg, nil
}

// CreateTask 创建新任务并写入数据库（调用大模型生成列表标题后落库）
func (s *PptService) CreateTask(ctx context.Context, task *Task, cfg types.PPTConfig) error {
	if s.llm == nil {
		s.llm = NewLLMClient()
	}
	title, err := s.llm.GeneratePPTTitle(ctx, cfg, task.Content, task.Language)
	if err != nil {
		logger.Warnf("GeneratePPTTitle failed task_id=%s: %v", task.TaskID, err)
		title = "未命名演示文稿"
	} else {
		title = strings.TrimSpace(title)
		if title == "" {
			title = "未命名演示文稿"
		}
	}
	task.Title = truncateTitleRunes(title, 255)

	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	task.Status = TaskStatusPending
	job := taskToModel(task)
	return s.db.Create(job).Error
}

// GetTask 从数据库获取任务
func (s *PptService) GetTask(taskID string) (*Task, bool) {
	var job model.PPTJob
	err := s.db.Where("task_id = ?", taskID).First(&job).Error
	if err != nil || job.TaskId == "" {
		return nil, false
	}
	return modelToTask(&job), true
}

// UpdateStatus 更新任务状态
func (s *PptService) UpdateStatus(taskID string, status TaskStatus) {
	s.db.Model(&model.PPTJob{}).Where("task_id = ?", taskID).
		Updates(map[string]interface{}{"status": string(status), "updated_at": time.Now()})
}

func countSlidesWithImage(slides []SlideData) int {
	n := 0
	for _, sl := range slides {
		if strings.TrimSpace(sl.ImageURL) != "" {
			n++
		}
	}
	return n
}

func slidePlansToOutlines(plans []slidePlan) []SlideData {
	out := make([]SlideData, len(plans))
	for i, p := range plans {
		out[i] = SlideData{
			SlideIndex:  p.SlideIndex,
			Theme:       p.Theme,
			Title:       p.Title,
			Points:      p.Points,
			ImagePrompt: p.ImagePrompt,
			ImageURL:    "",
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].SlideIndex < out[j].SlideIndex
	})
	return out
}

// saveSlidesOutline 分镜一出即落库：每页含 theme/title/points/image_prompt，image_url 为空
func (s *PptService) saveSlidesOutline(taskID string, total int, slides []SlideData) error {
	s.slidesLock.Lock()
	defer s.slidesLock.Unlock()

	voSlides := make(vo.PPTSlides, len(slides))
	for i := range slides {
		voSlides[i] = slideToVO(slides[i])
	}
	completed := countSlidesWithImage(slides)
	return s.db.Model(&model.PPTJob{}).Where("task_id = ?", taskID).Updates(map[string]interface{}{
		"slides":           voSlides,
		"total_slides":     total,
		"completed_slides": completed,
		"updated_at":       time.Now(),
	}).Error
}

// ApplySlideImage 按 slide_index 原地写入 image_url，并刷新 completed_slides、thumb
func (s *PptService) ApplySlideImage(taskID string, slide SlideData) error {
	s.slidesLock.Lock()
	defer s.slidesLock.Unlock()

	var job model.PPTJob
	if err := s.db.Where("task_id = ?", taskID).First(&job).Error; err != nil {
		return err
	}
	slides := job.Slides
	found := false
	for i := range slides {
		if slides[i].SlideIndex == slide.SlideIndex {
			slides[i].ImageURL = slide.ImageURL
			if strings.TrimSpace(slide.ImageURL) != "" && len(slides[i].ImageHistory) == 0 {
				slides[i].ImageHistory = vo.PPTSlideImageVersions{
					{ImageURL: slide.ImageURL, Prompt: strings.TrimSpace(slide.ImagePrompt)},
				}
			}
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("slide index %d not found", slide.SlideIndex)
	}
	job.Slides = slides
	biz := voSlidesToBiz(slides)
	return s.refreshJobMeta(&job, biz)
}

// refreshJobMeta 刷新 job 的 CompletedSlides/Thumb/UpdatedAt 并 Save。
// 调用前必须已持有 slidesLock。
func (s *PptService) refreshJobMeta(job *model.PPTJob, biz []SlideData) error {
	job.CompletedSlides = countSlidesWithImage(biz)
	job.Thumb = DerivePPTThumbFromSlides(biz)
	job.UpdatedAt = time.Now()
	return s.db.Save(job).Error
}

func (s *PptService) validateSlideOutline(task *Task) error {
	if task.Total <= 0 {
		return ErrPptTaskNotResumable
	}
	if len(task.Slides) < task.Total {
		return ErrPptTaskNotResumable
	}
	seen := make(map[int]bool, len(task.Slides))
	for _, sl := range task.Slides {
		seen[sl.SlideIndex] = true
	}
	for i := 1; i <= task.Total; i++ {
		if !seen[i] {
			return ErrPptTaskNotResumable
		}
	}
	return nil
}

func slidesNeedingImages(task *Task) []SlideData {
	var need []SlideData
	for _, sl := range task.Slides {
		if strings.TrimSpace(sl.ImageURL) == "" {
			need = append(need, sl)
		}
	}
	sort.Slice(need, func(i, j int) bool {
		return need[i].SlideIndex < need[j].SlideIndex
	})
	return need
}

func (s *PptService) userPower(userID uint) (int, error) {
	var u model.User
	if err := s.db.Where("id = ?", userID).First(&u).Error; err != nil {
		return 0, err
	}
	return u.Power, nil
}

// runSlideImageJobs 为给定幻灯片列表并发生图（每张成功后 ApplySlideImage）
func (s *PptService) runSlideImageJobs(ctx context.Context, task *Task, cfg types.PPTConfig, generator ImageGenerator, slides []SlideData) error {
	if len(slides) == 0 {
		return nil
	}
	if cfg.MaxConcurrentRequests <= 0 {
		cfg.MaxConcurrentRequests = 3
	}
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(cfg.MaxConcurrentRequests)
	for _, item := range slides {
		slide := item
		group.Go(func() error {
			imgURL, err := generator.Generate(ctx, slide.ImagePrompt)
			if err != nil {
				return err
			}
			storedURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(imgURL, ".png", false)
			if err != nil {
				return fmt.Errorf("转存图片失败：%w", err)
			}
			full := slide
			full.ImageURL = storedURL
			if err := s.ApplySlideImage(task.TaskID, full); err != nil {
				return err
			}
			if cfg.PowerCostPerSlide > 0 {
				err = s.userService.DecreasePower(task.UserID, cfg.PowerCostPerSlide, model.PowerLog{
					Type:   types.PowerConsume,
					Model:  generator.Provider(),
					Remark: fmt.Sprintf("PPT 任务 %s 第 %d 页图片生成", task.TaskID, slide.SlideIndex),
				})
				if err != nil {
					return fmt.Errorf("扣减算力失败：%v", err)
				}
			}
			return nil
		})
	}
	return group.Wait()
}

// startSlideImageGenerationAsync 在后台为 missing 页并发生图；若 setProcessingBeforeRun 为 true 则先置为 processing（用户主动 resume）。
func (s *PptService) startSlideImageGenerationAsync(task *Task, missing []SlideData, setProcessingBeforeRun bool) error {
	if len(missing) == 0 {
		return nil
	}
	cfg, err := s.loadPPTConfig()
	if err != nil {
		return err
	}
	cost := len(missing) * cfg.PowerCostPerSlide
	if cost > 0 {
		power, err := s.userPower(task.UserID)
		if err != nil {
			return err
		}
		if power < cost {
			return ErrInsufficientPower
		}
	}

	generator, err := NewImageGenerator(cfg)
	if err != nil {
		return fmt.Errorf("初始化图片生成器失败：%w", err)
	}

	if setProcessingBeforeRun {
		s.UpdateStatus(task.TaskID, TaskStatusProcessing)
	}
	t := task
	go func() {
		bg := context.Background()
		if err := s.runSlideImageJobs(bg, t, cfg, generator, missing); err != nil {
			s.MarkAsFailed(task.TaskID, fmt.Sprintf("图片生成失败：%v", err))
			return
		}
		s.UpdateStatus(task.TaskID, TaskStatusCompleted)
	}()
	return nil
}

// RecoverStaleProcessingTasks 进程启动时扫描 DB 中仍为 processing 且存在缺图页的任务，重新拉起生图协程（用于服务中断后的恢复）。
func (s *PptService) RecoverStaleProcessingTasks() {
	var jobs []model.PPTJob
	if err := s.db.Where("status = ?", string(TaskStatusProcessing)).Find(&jobs).Error; err != nil {
		logger.Warnf("PPT recover: list processing jobs failed: %v", err)
		return
	}
	for i := range jobs {
		task := modelToTask(&jobs[i])
		missing := slidesNeedingImages(task)
		if len(missing) == 0 {
			s.UpdateStatus(task.TaskID, TaskStatusCompleted)
			logger.Infof("PPT recover: task %s was processing but all slides had images, marked completed", task.TaskID)
			continue
		}
		if err := s.validateSlideOutline(task); err != nil {
			logger.Warnf("PPT recover: task %s skip (invalid outline): %v", task.TaskID, err)
			continue
		}
		if err := s.startSlideImageGenerationAsync(task, missing, false); err != nil {
			if errors.Is(err, ErrInsufficientPower) {
				logger.Warnf("PPT recover: task %s skip (insufficient power for %d slides)", task.TaskID, len(missing))
				continue
			}
			logger.Warnf("PPT recover: task %s failed to restart: %v", task.TaskID, err)
			continue
		}
		logger.Infof("PPT recover: restarted image generation for task %s (%d slides)", task.TaskID, len(missing))
	}
}

// ResumeTask 继续为缺图页生图（需完整分镜占位；processing 时返回 ErrPptTaskBusy）
func (s *PptService) ResumeTask(ctx context.Context, taskID string, userID uint) error {
	task, ok := s.GetTask(taskID)
	if !ok {
		return ErrPptTaskNotFound
	}
	if task.UserID != userID {
		return ErrPptTaskNotFound
	}
	if task.Status == TaskStatusProcessing {
		return ErrPptTaskBusy
	}
	if task.Status == TaskStatusCompleted {
		return ErrPptTaskNotResumable
	}
	if err := s.validateSlideOutline(task); err != nil {
		return err
	}
	missing := slidesNeedingImages(task)
	if len(missing) == 0 {
		s.UpdateStatus(taskID, TaskStatusCompleted)
		return nil
	}

	return s.startSlideImageGenerationAsync(task, missing, true)
}

// MarkAsFailed 标记任务失败
func (s *PptService) MarkAsFailed(taskID string, msg string) {
	s.db.Model(&model.PPTJob{}).Where("task_id = ?", taskID).
		Updates(map[string]interface{}{"status": string(TaskStatusFailed), "err_msg": msg, "updated_at": time.Now()})
}

// DeleteTask 删除任务并删除关联生成图片
// 仅允许删除 completed / failed 状态的任务，避免并发任务生成过程被打断。
func (s *PptService) DeleteTask(taskID string, userID uint) error {
	var job model.PPTJob
	if err := s.db.Where("task_id = ? AND user_id = ?", taskID, userID).First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || job.TaskId == "" {
			return ErrPptTaskNotFound
		}
		return err
	}

	if job.Status != string(TaskStatusCompleted) && job.Status != string(TaskStatusFailed) {
		return ErrPptTaskNotDeletable
	}

	// 删除所有幻灯片对应的图片对象。
	uploader := s.uploadManager.GetUploadHandler()
	for _, slide := range job.Slides {
		if slide.ImageURL == "" {
			continue
		}
		if err := uploader.Delete(slide.ImageURL); err != nil {
			// 图片可能已过期/不存在/对象已被清理，此时不应阻断“删除任务记录”的主流程。
			// 这里只记录日志，确保数据库记录被删除后前端认为任务删除成功。
			logger.Warnf("delete ppt image failed (task_id=%s, url=%s): %v", taskID, slide.ImageURL, err)
		}
	}

	// 最后删除任务记录（slides 会随之从数据库消失）。
	return s.db.Where("task_id = ? AND user_id = ?", taskID, userID).Delete(&model.PPTJob{}).Error
}

// List 返回所有任务列表（从数据库按创建时间倒序），供调用方按用户/状态过滤与分页
func (s *PptService) List() []*Task {
	var jobs []model.PPTJob
	s.db.Order("created_at DESC").Find(&jobs)
	tasks := make([]*Task, 0, len(jobs))
	for i := range jobs {
		tasks = append(tasks, modelToTask(&jobs[i]))
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})
	return tasks
}

// ListUserTasks 当前用户的任务分页列表（缺 title/thumb 时补全并写库）
func (s *PptService) ListUserTasks(ctx context.Context, userID uint, page, pageSize int) ([]*Task, int) {
	all := s.List()
	filtered := make([]*Task, 0, len(all))
	for _, t := range all {
		if t.UserID == userID {
			filtered = append(filtered, t)
		}
	}
	total := len(filtered)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	slice := filtered[start:end]
	cfg, cfgErr := s.loadPPTConfig()
	if cfgErr != nil {
		logger.Warnf("ListUserTasks loadPPTConfig: %v", cfgErr)
	}
	if s.llm == nil {
		s.llm = NewLLMClient()
	}
	for _, t := range slice {
		s.ensureTaskMeta(ctx, t, cfg)
	}
	return slice, total
}

// EnsureTaskMeta 对外暴露标题/缩略图补全逻辑，供管理端列表/详情复用。
func (s *PptService) EnsureTaskMeta(ctx context.Context, task *Task) {
	if task == nil {
		return
	}
	cfg, cfgErr := s.loadPPTConfig()
	if cfgErr != nil {
		logger.Warnf("EnsureTaskMeta loadPPTConfig: %v", cfgErr)
	}
	if s.llm == nil {
		s.llm = NewLLMClient()
	}
	s.ensureTaskMeta(ctx, task, cfg)
}

func (s *PptService) ensureTaskMeta(ctx context.Context, task *Task, cfg types.PPTConfig) {
	updates := map[string]interface{}{}
	if strings.TrimSpace(task.Title) == "" && strings.TrimSpace(task.Content) != "" {
		title := "未命名演示文稿"
		if cfg.OutlineLLMApiURL != "" && cfg.OutlineLLMApiKey != "" {
			ti, err := s.llm.GeneratePPTTitle(ctx, cfg, task.Content, task.Language)
			if err != nil {
				logger.Warnf("ensureTaskMeta GeneratePPTTitle task_id=%s: %v", task.TaskID, err)
			} else {
				ti = strings.TrimSpace(ti)
				if ti != "" {
					title = truncateTitleRunes(ti, 255)
				}
			}
		}
		task.Title = title
		updates["title"] = task.Title
	}
	if task.Thumb == "" && len(task.Slides) > 0 {
		thumb := DerivePPTThumbFromSlides(task.Slides)
		if thumb != "" {
			task.Thumb = thumb
			updates["thumb"] = thumb
		}
	}
	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		_ = s.db.Model(&model.PPTJob{}).Where("task_id = ?", task.TaskID).Updates(updates).Error
	}
}

// ListAdminJobs 管理后台任务列表（筛选 + 分页）
func (s *PptService) ListAdminJobs(ctx context.Context, page, pageSize, filterUserID int, status string) ([]*Task, int) {
	items := s.List()
	filtered := make([]*Task, 0, len(items))
	for _, t := range items {
		if filterUserID > 0 && int(t.UserID) != filterUserID {
			continue
		}
		if status != "" && string(t.Status) != status {
			continue
		}
		filtered = append(filtered, t)
	}
	total := len(filtered)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	slice := filtered[start:end]
	for _, t := range slice {
		s.EnsureTaskMeta(ctx, t)
	}
	return slice, total
}

// Stats 任务状态统计（管理后台）
func (s *PptService) Stats() (total, completed, processing, failed, pending int64) {
	for _, t := range s.List() {
		total++
		switch t.Status {
		case TaskStatusCompleted:
			completed++
		case TaskStatusProcessing:
			processing++
		case TaskStatusFailed:
			failed++
		case TaskStatusPending:
			pending++
		}
	}
	return
}

// RunTask 执行 PPT 生成：分镜、并发生图、转存图片、扣算力
func (s *PptService) RunTask(ctx context.Context, task *Task, cfg types.PPTConfig) {
	s.UpdateStatus(task.TaskID, TaskStatusProcessing)

	maxPages := task.Pages
	if maxPages <= 0 {
		maxPages = cfg.MaxSlidesPerTask
	}
	plans, err := s.llm.GenerateSlides(ctx, cfg, task.Content, task.Prompt, task.Language, task.Mode, maxPages)
	if err != nil {
		s.MarkAsFailed(task.TaskID, fmt.Sprintf("生成分镜失败：%v", err))
		return
	}

	if len(plans) == 0 {
		s.MarkAsFailed(task.TaskID, "分镜结果为空")
		return
	}

	total := len(plans)
	if cfg.MaxSlidesPerTask > 0 && total > cfg.MaxSlidesPerTask {
		plans = plans[:cfg.MaxSlidesPerTask]
		total = len(plans)
	}

	outlines := slidePlansToOutlines(plans)
	if err := s.saveSlidesOutline(task.TaskID, total, outlines); err != nil {
		s.MarkAsFailed(task.TaskID, fmt.Sprintf("保存分镜占位失败：%v", err))
		return
	}

	generator, err := NewImageGenerator(cfg)
	if err != nil {
		s.MarkAsFailed(task.TaskID, fmt.Sprintf("初始化图片生成器失败：%v", err))
		return
	}

	if err := s.runSlideImageJobs(ctx, task, cfg, generator, outlines); err != nil {
		s.MarkAsFailed(task.TaskID, fmt.Sprintf("图片生成失败：%v", err))
		return
	}

	s.UpdateStatus(task.TaskID, TaskStatusCompleted)
}

func (s *PptService) loadPPTConfig() (types.PPTConfig, error) {
	var cfgModel model.Config
	var pptCfg types.PPTConfig

	err := s.db.Where("name", types.ConfigKeyPPT).First(&cfgModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			pptCfg.MaxSlidesPerTask = 30
			pptCfg.MaxConcurrentRequests = 3
			pptCfg.QPSLimit = 1
			pptCfg.PowerCostPerSlide = 0
			return pptCfg, nil
		}
		return pptCfg, err
	}

	err = utils.JsonDecode(cfgModel.Value, &pptCfg)
	if err != nil {
		return pptCfg, err
	}

	legacyMax10 := pptCfg.MaxSlidesPerTask == 10
	if pptCfg.MaxSlidesPerTask <= 0 {
		pptCfg.MaxSlidesPerTask = 30
	}
	if legacyMax10 {
		// 与前端 PPT 页数控件 max=30 对齐；历史默认 10 会导致用户选择 12/15 仍被截断为 10
		pptCfg.MaxSlidesPerTask = 30
	}
	if pptCfg.MaxConcurrentRequests <= 0 {
		pptCfg.MaxConcurrentRequests = 3
	}
	if pptCfg.QPSLimit <= 0 {
		pptCfg.QPSLimit = 1
	}

	if legacyMax10 {
		val := utils.JsonEncode(pptCfg)
		_ = s.db.Model(&model.Config{}).Where("name = ?", types.ConfigKeyPPT).Update("value", val)
	}

	return pptCfg, nil
}
