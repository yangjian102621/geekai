package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service/ppt"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PPTHandler 管理后台 PPT 生成配置处理器
type PPTHandler struct {
	handler.BaseHandler
	pptService *ppt.PptService
}

// NewPPTHandler 创建管理后台 PPT 配置处理器
func NewPPTHandler(app *core.AppServer, db *gorm.DB, pptService *ppt.PptService) *PPTHandler {
	return &PPTHandler{
		BaseHandler: handler.BaseHandler{App: app, DB: db},
		pptService:  pptService,
	}
}

// RegisterRoutes 注册 PPT 配置相关路由
func (h *PPTHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/ppt/")
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("config", h.GetConfig)
		rg.POST("config/update", h.UpdateConfig)
		rg.GET("jobs", h.Jobs)
		rg.GET("jobs/:task_id", h.JobDetail)
		rg.GET("jobs/:task_id/export", h.ExportJob)
		rg.GET("stats", h.Stats)
	}
}

// GetConfig 获取 PPT 生成配置
func (h *PPTHandler) GetConfig(c *gin.Context) {
	var cfg model.Config
	err := h.DB.Where("name", types.ConfigKeyPPT).First(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回一个默认空配置
			resp.SUCCESS(c, types.PPTConfig{
				OutlineLLMModel:       "gpt-4o-mini",
				MaxSlidesPerTask:      30,
				PowerCostPerSlide:     0,
				MaxConcurrentRequests: 3,
				QPSLimit:              1,
				NanoBananaModel:       "nano-banana",
				NanoBananaAspectRatio: "16:9",
				SeedreamSize:          "1920x1080",
			})
			return
		}
		resp.ERROR(c, "获取配置失败: "+err.Error())
		return
	}

	var pptConfig types.PPTConfig
	err = utils.JsonDecode(cfg.Value, &pptConfig)
	if err != nil {
		resp.ERROR(c, "解析配置失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, pptConfig)
}

// UpdateConfig 更新 PPT 生成配置
func (h *PPTHandler) UpdateConfig(c *gin.Context) {
	var req types.PPTConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	// 基础校验
	if req.MaxSlidesPerTask <= 0 {
		resp.ERROR(c, "单个任务最多 PPT 页数必须大于 0")
		return
	}

	if req.PowerCostPerSlide < 0 {
		resp.ERROR(c, "每张 PPT 图片消耗算力不能小于 0")
		return
	}

	if req.MaxConcurrentRequests <= 0 {
		req.MaxConcurrentRequests = 3
	}
	if req.QPSLimit <= 0 {
		req.QPSLimit = 1
	}

	// 根据当前生图提供方做必填校验
	switch req.ActiveImageProvider {
	case types.PPTImageProviderNanoBanana:
		if req.NanoBananaApiURL == "" {
			resp.ERROR(c, "Nano Banana API 地址不能为空")
			return
		}
		if req.NanoBananaApiKey == "" {
			resp.ERROR(c, "Nano Banana API Key 不能为空")
			return
		}
	case types.PPTImageProviderSeedream:
		if req.SeedreamBaseURL == "" {
			resp.ERROR(c, "Seedream Base URL 不能为空")
			return
		}
		if req.SeedreamApiKey == "" {
			resp.ERROR(c, "Seedream API Key 不能为空")
			return
		}
		if req.SeedreamModel == "" {
			resp.ERROR(c, "Seedream 模型 ID 不能为空")
			return
		}
	default:
		// 允许为空，未来可以扩展更多 provider
	}

	value := utils.JsonEncode(&req)
	var cfg model.Config
	err := h.DB.Where("name", types.ConfigKeyPPT).First(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cfg.Name = types.ConfigKeyPPT
			cfg.Value = value
			if err = h.DB.Create(&cfg).Error; err != nil {
				resp.ERROR(c, "创建配置失败: "+err.Error())
				return
			}
			resp.SUCCESS(c, gin.H{"message": "配置创建成功"})
			return
		}
		resp.ERROR(c, "获取配置失败: "+err.Error())
		return
	}

	cfg.Value = value
	if err = h.DB.Updates(&cfg).Error; err != nil {
		resp.ERROR(c, "更新配置失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{"message": "配置更新成功"})
}

// Jobs 管理后台查看 PPT 任务列表（内存任务）
func (h *PPTHandler) Jobs(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	filterUserId := h.GetInt(c, "user_id", 0)
	status := h.GetTrim(c, "status")

	filtered, total := h.pptService.ListAdminJobs(c.Request.Context(), page, pageSize, filterUserId, status)

	jobs := make([]gin.H, 0, len(filtered))
	for _, t := range filtered {
		job := t.TaskSummaryMap()
		job["user_id"] = t.UserID
		job["error_message"] = t.ErrorMessage
		jobs = append(jobs, job)
	}

	resp.SUCCESS(c, gin.H{
		"jobs":      jobs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func buildAdminPPTTaskDetail(task *ppt.Task) gin.H {
	percentage := 0
	if task.Total > 0 {
		percentage = int(float64(task.Completed) / float64(task.Total) * 100)
	}

	return gin.H{
		"task_id": task.TaskID,
		"user_id": task.UserID,
		"status":  task.Status,
		"progress": gin.H{
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
	}
}

// JobDetail 管理后台查看指定 PPT 任务详情
func (h *PPTHandler) JobDetail(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	task, ok := h.pptService.GetTask(taskID)
	if !ok {
		resp.ERROR(c, "任务不存在")
		return
	}
	h.pptService.EnsureTaskMeta(c.Request.Context(), task)
	resp.SUCCESS(c, buildAdminPPTTaskDetail(task))
}

// ExportJob 管理后台导出 PPT 任务
func (h *PPTHandler) ExportJob(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	ef, ok := ppt.ParseExportFormat(c.Query("format"))
	if !ok {
		resp.ERROR(c, "format 参数无效，支持 pdf 或 pptx")
		return
	}

	task, exists := h.pptService.GetTask(taskID)
	if !exists {
		resp.ERROR(c, "任务不存在")
		return
	}
	if task.Status != ppt.TaskStatusCompleted {
		resp.ERROR(c, "仅已完成任务可导出")
		return
	}

	h.pptService.EnsureTaskMeta(c.Request.Context(), task)

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
	c.Data(200, ppt.ExportMimeType(ef), data)
}

// Stats PPT 任务统计信息
func (h *PPTHandler) Stats(c *gin.Context) {
	total, completed, processing, failed, pending := h.pptService.Stats()

	resp.SUCCESS(c, gin.H{
		"totalTasks":      total,
		"completedTasks":  completed,
		"processingTasks": processing,
		"failedTasks":     failed,
		"pendingTasks":    pending,
	})
}
