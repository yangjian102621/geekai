package handler

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
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/ppt"
	"geekai/utils"
	"geekai/utils/resp"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PPTTaskHandler 用户侧 PPT 生成任务处理器（薄层：参数、鉴权、调用 PptService）
type PPTTaskHandler struct {
	BaseHandler
	snowflake  *service.Snowflake
	pptService *ppt.PptService
}

// NewPPTTaskHandler 创建 PPT 任务处理器
func NewPPTTaskHandler(app *core.AppServer, db *gorm.DB, snowflake *service.Snowflake, pptService *ppt.PptService) *PPTTaskHandler {
	return &PPTTaskHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
		snowflake:   snowflake,
		pptService:  pptService,
	}
}

// RegisterRoutes 注册 PPT 任务相关路由
func (h *PPTTaskHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/v1/tasks")
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.POST("generate-slides", h.CreateTask)
		group.POST("generate-slides/from-file", h.CreateTaskFromFile)
		group.GET("", h.ListTasks) // 当前用户任务列表，须在 :task_id 前注册
		group.GET(":task_id/export", h.ExportTask)
		group.POST(":task_id/resume", h.ResumeTask)
		group.POST(":task_id/slides/:slide_index/edit-image", h.EditSlideImage)
		group.PATCH(":task_id/slides/:slide_index/active-image", h.SetActiveSlideImage)
		group.GET(":task_id", h.GetTask)
		group.DELETE(":task_id", h.DeleteTask)
	}
}

type createPPTTaskRequest struct {
	Content  string `json:"content"`
	Prompt   string `json:"prompt"`
	Language string `json:"language"` // 如 zh-CN, en，约束分镜输出语言
	Pages    int    `json:"pages"`    // 目标页数，0 表示用配置默认值
	Mode     string `json:"mode"`     // 生成模式：detailed=详细演示文稿，slides=演示用幻灯片，空则默认 slides
}

type createPPTTaskFromFileRequest struct {
	FileURL  string `json:"file_url"`
	Prompt   string `json:"prompt"`
	Language string `json:"language"` // 如 zh-CN, en，约束分镜输出语言
	Pages    int    `json:"pages"`    // 目标页数，0 表示用配置默认值
	Mode     string `json:"mode"`     // 生成模式：detailed=详细演示文稿，slides=演示用幻灯片，空则默认 slides
}

// CreateTask 创建 PPT 生成任务（异步）
func (h *PPTTaskHandler) CreateTask(c *gin.Context) {
	var req createPPTTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	taskID, err := h.snowflake.Next(true)
	if err != nil {
		resp.ERROR(c, "生成任务ID失败："+err.Error())
		return
	}

	task, cfg, err := h.pptService.BuildPendingTask(taskID, user.Id, int(user.Power), req.Content, req.Prompt, req.Language, req.Mode, req.Pages)
	if err != nil {
		if errors.Is(err, ppt.ErrInsufficientPower) {
			resp.ERROR(c, "当前用户算力不足以完成本次 PPT 生成任务！")
			return
		}
		resp.ERROR(c, "加载 PPT 配置或校验失败："+err.Error())
		return
	}

	h.pptService.CreateTask(c.Request.Context(), task, cfg)
	go h.pptService.RunTask(context.Background(), task, cfg)

	resp.SUCCESS(c, map[string]any{
		"task_id": taskID,
		"status":  ppt.TaskStatusPending,
	})
}

// CreateTaskFromFile 创建 PPT 生成任务（基于上传材料文本提炼）
// 支持文件：PDF、Word（doc/docx）、TXT、Markdown（md/markdown）
func (h *PPTTaskHandler) CreateTaskFromFile(c *gin.Context) {
	var req createPPTTaskFromFileRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.FileURL) == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	fileURL := strings.TrimSpace(req.FileURL)

	// 去掉 query，避免 ".pdf?xxx" 这种导致 ext 识别失败
	urlNoQuery := strings.Split(fileURL, "?")[0]
	ext := strings.ToLower(filepath.Ext(urlNoQuery))
	if ext == "" {
		resp.ERROR(c, "不支持的文件格式")
		return
	}

	allowedExts := map[string]bool{
		".pdf":      true,
		".doc":      true,
		".docx":     true,
		".txt":      true,
		".md":       true,
		".markdown": true,
	}
	if !allowedExts[ext] {
		resp.ERROR(c, "不支持的文件格式")
		return
	}

	designPrompt := strings.TrimSpace(req.Prompt)
	language := strings.TrimSpace(req.Language)
	if language == "" {
		language = "zh-CN"
	}
	mode := strings.TrimSpace(req.Mode)

	pages := req.Pages

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	taskID, err := h.snowflake.Next(true)
	if err != nil {
		resp.ERROR(c, "生成任务ID失败："+err.Error())
		return
	}

	// 1) 下载文件并提取原始文本（Notebook 风格提炼前的输入）
	rawText, err := h.extractMaterialTextFromURL(c.Request.Context(), fileURL, ext)
	if err != nil || strings.TrimSpace(rawText) == "" {
		if err == nil {
			err = errors.New("empty extracted text")
		}
		resp.ERROR(c, "读取文件失败："+err.Error())
		return
	}

	// 2) 校验算力与页数，组装待写入的 Task（先在提炼前做 power 检查）
	task, cfg, err := h.pptService.BuildPendingTask(taskID, user.Id, int(user.Power), rawText, designPrompt, language, mode, pages)
	if err != nil {
		if errors.Is(err, ppt.ErrInsufficientPower) {
			resp.ERROR(c, "当前用户算力不足以完成本次 PPT 生成任务！")
			return
		}
		resp.ERROR(c, "加载 PPT 配置或校验失败："+err.Error())
		return
	}

	// 3) NotebookLM 风格提炼：rawText -> PPT 可用的 content（大纲/结构化要点）
	notebookContent, err := h.pptService.GenerateNotebookContent(c.Request.Context(), cfg, rawText, designPrompt, language)
	if err != nil || strings.TrimSpace(notebookContent) == "" {
		if err == nil {
			err = errors.New("empty notebook content")
		}
		resp.ERROR(c, "提炼材料失败："+err.Error())
		return
	}
	task.Content = notebookContent

	// 4) 落库 + 异步生成分镜与图片
	err = h.pptService.CreateTask(c.Request.Context(), task, cfg)
	if err != nil {
		resp.ERROR(c, "创建任务失败："+err.Error())
		return
	}
	go h.pptService.RunTask(context.Background(), task, cfg)

	resp.SUCCESS(c, map[string]any{
		"task_id": taskID,
		"status":  ppt.TaskStatusPending,
	})
}

func (h *PPTTaskHandler) extractMaterialTextFromURL(ctx context.Context, fileURL string, ext string) (string, error) {
	if strings.TrimSpace(fileURL) == "" {
		return "", errors.New("empty file url")
	}

	switch ext {
	case ".txt", ".md", ".markdown":
		b, status, err := utils.FetchURLBytes(ctx, fileURL, "", 30*time.Second, 2, 8<<20)
		if err != nil {
			// status=0 时通常是请求阶段错误（例如 TLS 握手超时）
			return "", fmt.Errorf("download file failed: status=%d: %w", status, err)
		}
		return strings.TrimSpace(string(b)), nil
	default:
		// PDF/Word 等走 Tika
		return utils.ReadFileContent(fileURL, h.App.Config.TikaHost)
	}
}

// ListTasks 当前用户的 PPT 任务列表（分页）
func (h *PPTTaskHandler) ListTasks(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	slice, total := h.pptService.ListUserTasks(c.Request.Context(), user.Id, page, pageSize)

	jobs := make([]map[string]any, 0, len(slice))
	for _, t := range slice {
		job := t.TaskSummaryMap()
		job["prompt"] = t.Prompt
		if t.ErrorMessage != "" {
			job["error_message"] = t.ErrorMessage
		}
		jobs = append(jobs, job)
	}

	resp.SUCCESS(c, map[string]any{
		"jobs":      jobs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ExportTask 导出任务幻灯片为 PDF 或 PPTX（按图片逐页）
func (h *PPTTaskHandler) ExportTask(c *gin.Context) {
	taskID := strings.TrimSpace(c.Param("task_id"))
	if taskID == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	ef, ok := ppt.ParseExportFormat(c.Query("format"))
	if !ok {
		resp.ERROR(c, "format 参数无效，支持 pdf 或 pptx")
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	task, exists := h.pptService.GetTask(taskID)
	if !exists {
		resp.ERROR(c, "任务不存在")
		return
	}
	if user.Id != task.UserID {
		resp.NotAuth(c)
		return
	}
	if task.Status != ppt.TaskStatusCompleted {
		resp.ERROR(c, "仅已完成任务可导出")
		return
	}

	ossCfg := types.OSSConfig{}
	if h.App.SysConfig != nil {
		ossCfg = h.App.SysConfig.OSS
	}
	data, err := ppt.BuildExportBytes(c.Request.Context(), task.Slides, ef, ossCfg, h.App.Config)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	base := ppt.SanitizeExportBaseName(task.Title, task.TaskID)
	filename := base + ppt.ExportFileExt(ef)
	c.Header("Content-Disposition", ppt.ContentDispositionAttachment(filename))
	c.Data(http.StatusOK, ppt.ExportMimeType(ef), data)
}

// ResumeTask 继续生成缺图页（POST /api/v1/tasks/:task_id/resume）
func (h *PPTTaskHandler) ResumeTask(c *gin.Context) {
	taskID := strings.TrimSpace(c.Param("task_id"))
	if taskID == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	err = h.pptService.ResumeTask(c.Request.Context(), taskID, user.Id)
	if err != nil {
		if errors.Is(err, ppt.ErrPptTaskNotFound) {
			resp.ERROR(c, "任务不存在")
			return
		}
		if errors.Is(err, ppt.ErrPptTaskBusy) {
			resp.ERROR(c, "任务正在处理中，请稍后再试")
			return
		}
		if errors.Is(err, ppt.ErrPptTaskNotResumable) {
			resp.ERROR(c, "当前任务无法继续生成（已完成或分镜数据不完整）")
			return
		}
		if errors.Is(err, ppt.ErrInsufficientPower) {
			resp.ERROR(c, "当前用户算力不足以完成剩余图片生成")
			return
		}
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, map[string]any{
		"task_id": taskID,
		"status":  ppt.TaskStatusProcessing,
	})
}

type editSlideImageRequest struct {
	Prompt string `json:"prompt"`
}

type activeSlideImageRequest struct {
	VersionIndex int `json:"version_index"`
}

// EditSlideImage 图生图编辑当前页（基于激活图）
func (h *PPTTaskHandler) EditSlideImage(c *gin.Context) {
	taskID := strings.TrimSpace(c.Param("task_id"))
	slideIndexStr := strings.TrimSpace(c.Param("slide_index"))
	if taskID == "" || slideIndexStr == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	slideIndex, err := strconv.Atoi(slideIndexStr)
	if err != nil || slideIndex < 1 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var req editSlideImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	ossCfg := types.OSSConfig{}
	if h.App.SysConfig != nil {
		ossCfg = h.App.SysConfig.OSS
	}

	slides, err := h.pptService.EditSlideImage(c.Request.Context(), taskID, user.Id, slideIndex, req.Prompt, ossCfg, h.App.Config)
	if err != nil {
		if errors.Is(err, ppt.ErrPptTaskNotFound) {
			resp.ERROR(c, "任务不存在")
			return
		}
		if errors.Is(err, ppt.ErrPptSlideNotFound) {
			resp.ERROR(c, "幻灯片不存在")
			return
		}
		if errors.Is(err, ppt.ErrPptSlideNoImage) {
			resp.ERROR(c, "该页暂无配图，无法编辑")
			return
		}
		if errors.Is(err, ppt.ErrInsufficientPower) {
			resp.ERROR(c, "当前用户算力不足以完成本次编辑")
			return
		}
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, map[string]any{"slides": slides})
}

// SetActiveSlideImage 切换当前页激活的历史版本
func (h *PPTTaskHandler) SetActiveSlideImage(c *gin.Context) {
	taskID := strings.TrimSpace(c.Param("task_id"))
	slideIndexStr := strings.TrimSpace(c.Param("slide_index"))
	if taskID == "" || slideIndexStr == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	slideIndex, err := strconv.Atoi(slideIndexStr)
	if err != nil || slideIndex < 1 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var req activeSlideImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	slides, err := h.pptService.SetActiveSlideVersion(taskID, user.Id, slideIndex, req.VersionIndex)
	if err != nil {
		if errors.Is(err, ppt.ErrPptTaskNotFound) {
			resp.ERROR(c, "任务不存在")
			return
		}
		if errors.Is(err, ppt.ErrPptSlideNotFound) {
			resp.ERROR(c, "幻灯片不存在")
			return
		}
		if errors.Is(err, ppt.ErrPptInvalidVersionIndex) {
			resp.ERROR(c, "无效的历史版本序号")
			return
		}
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, map[string]any{"slides": slides})
}

// GetTask 查询 PPT 任务进度
func (h *PPTTaskHandler) GetTask(c *gin.Context) {
	taskId := c.Param("task_id")
	if taskId == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	task, ok := h.pptService.GetTask(taskId)
	if !ok {
		resp.ERROR(c, "任务不存在")
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil || user.Id != task.UserID {
		resp.NotAuth(c)
		return
	}

	percentage := 0
	if task.Total > 0 {
		percentage = int(float64(task.Completed) / float64(task.Total) * 100)
	}

	resp.SUCCESS(c, map[string]any{
		"task_id": task.TaskID,
		"status":  task.Status,
		"progress": map[string]any{
			"total_slides":     task.Total,
			"completed_slides": task.Completed,
			"percentage":       percentage,
		},
		"slides":        task.Slides,
		"error_message": task.ErrorMessage,
		"content":       task.Content,
		"prompt":        task.Prompt,
		"title":         task.Title,
		"thumb":         task.Thumb,
	})
}

// DeleteTask 删除用户 PPT 任务（仅允许 completed / failed），同时删除该任务生成的图片对象。
func (h *PPTTaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	if err := h.pptService.DeleteTask(taskID, user.Id); err != nil {
		if errors.Is(err, ppt.ErrPptTaskNotFound) {
			resp.ERROR(c, "任务不存在")
			return
		}
		if errors.Is(err, ppt.ErrPptTaskNotDeletable) {
			resp.ERROR(c, "仅允许删除已完成/已失败任务")
			return
		}
		// 前端会统一在文案里拼接“删除失败：”，这里避免重复。
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{"message": "删除成功"})
}
