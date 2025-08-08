package handler

import (
	"fmt"

	"geekai/core"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/jimeng"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JimengHandler 即梦AI处理器
type JimengHandler struct {
	BaseHandler
	jimengService *jimeng.Service
	userService   *service.UserService
}

// NewJimengHandler 创建即梦AI处理器
func NewJimengHandler(app *core.AppServer, jimengService *jimeng.Service, db *gorm.DB, userService *service.UserService) *JimengHandler {
	return &JimengHandler{
		BaseHandler:   BaseHandler{App: app, DB: db},
		jimengService: jimengService,
		userService:   userService,
	}
}

// RegisterRoutes 注册路由，新增统一任务接口
func (h *JimengHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/jimeng")
	rg.POST("task", h.CreateTask)            // 只保留统一任务接口
	rg.GET("power-config", h.GetPowerConfig) // 新增算力配置接口
	rg.POST("jobs", h.Jobs)
	rg.GET("remove", h.Remove)
	rg.GET("retry", h.Retry)
}

// JimengTaskRequest 统一任务请求结构体
// 支持所有生图和生成视频类型
type JimengTaskRequest struct {
	TaskType         string   `json:"task_type" binding:"required"`
	Prompt           string   `json:"prompt"`
	ImageInput       string   `json:"image_input"`
	ImageUrls        []string `json:"image_urls"`
	BinaryDataBase64 []string `json:"binary_data_base64"`
	Scale            float64  `json:"scale"`
	Width            int      `json:"width"`
	Height           int      `json:"height"`
	Gpen             float64  `json:"gpen"`
	Skin             float64  `json:"skin"`
	SkinUnifi        float64  `json:"skin_unifi"`
	GenMode          string   `json:"gen_mode"`
	Seed             int64    `json:"seed"`
	UsePreLLM        bool     `json:"use_pre_llm"`
	TemplateId       string   `json:"template_id"`
	AspectRatio      string   `json:"aspect_ratio"`
}

// CreateTask 统一任务创建接口
func (h *JimengHandler) CreateTask(c *gin.Context) {
	var req JimengTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 新增：除图像特效外，其他任务类型必须有提示词
	if req.TaskType != "image_effects" && req.Prompt == "" {
		resp.ERROR(c, "提示词不能为空")
		return
	}
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	if req.Width == 0 {
		req.Width = 1328
	}
	if req.Height == 0 {
		req.Height = 1328
	}
	if req.Seed == 0 {
		req.Seed = -1
	}

	var powerCost int
	var taskType model.JMTaskType
	var params map[string]any
	var reqKey string
	var modelName string

	switch req.TaskType {
	case "text_to_image":
		powerCost = h.getPowerFromConfig(model.JMTaskTypeTextToImage)
		taskType = model.JMTaskTypeTextToImage
		reqKey = jimeng.ReqKeyTextToImage
		modelName = "即梦文生图"
		if req.Scale == 0 {
			req.Scale = 2.5
		}
		params = map[string]any{
			"seed":        req.Seed,
			"scale":       req.Scale,
			"width":       req.Width,
			"height":      req.Height,
			"use_pre_llm": req.UsePreLLM,
		}
	case "image_to_image":
		powerCost = h.getPowerFromConfig(model.JMTaskTypeImageToImage)
		taskType = model.JMTaskTypeImageToImage
		reqKey = jimeng.ReqKeyImageToImagePortrait
		modelName = "即梦图生图"
		if req.Gpen == 0 {
			req.Gpen = 0.4
		}
		if req.Skin == 0 {
			req.Skin = 0.3
		}
		if req.GenMode == "" {
			if req.Prompt != "" {
				req.GenMode = jimeng.GenModeCreative
			} else {
				req.GenMode = jimeng.GenModeReference
			}
		}
		params = map[string]any{
			"image_input": req.ImageInput,
			"width":       req.Width,
			"height":      req.Height,
			"gpen":        req.Gpen,
			"skin":        req.Skin,
			"skin_unifi":  req.SkinUnifi,
			"gen_mode":    req.GenMode,
			"seed":        req.Seed,
		}
	case "image_edit":
		powerCost = h.getPowerFromConfig(model.JMTaskTypeImageEdit)
		taskType = model.JMTaskTypeImageEdit
		reqKey = jimeng.ReqKeyImageEdit
		modelName = "即梦图像编辑"
		if req.Scale == 0 {
			req.Scale = 0.5
		}
		params = map[string]any{
			"seed":  req.Seed,
			"scale": req.Scale,
		}
		if len(req.ImageUrls) > 0 {
			params["image_urls"] = req.ImageUrls
		}
		if len(req.BinaryDataBase64) > 0 {
			params["binary_data_base64"] = req.BinaryDataBase64
		}
	case "image_effects":
		powerCost = h.getPowerFromConfig(model.JMTaskTypeImageEffects)
		taskType = model.JMTaskTypeImageEffects
		reqKey = jimeng.ReqKeyImageEffects
		modelName = "即梦图像特效"
		if req.Width == 0 {
			req.Width = 1328
		}
		if req.Height == 0 {
			req.Height = 1328
		}
		params = map[string]any{
			"image_input1": req.ImageInput,
			"template_id":  req.TemplateId,
			"width":        req.Width,
			"height":       req.Height,
		}
	case "text_to_video":
		powerCost = h.getPowerFromConfig(model.JMTaskTypeTextToVideo)
		taskType = model.JMTaskTypeTextToVideo
		reqKey = jimeng.ReqKeyTextToVideo
		modelName = "即梦文生视频"
		if req.Seed == 0 {
			req.Seed = -1
		}
		if req.AspectRatio == "" {
			req.AspectRatio = jimeng.AspectRatio16_9
		}
		params = map[string]any{
			"seed":         req.Seed,
			"aspect_ratio": req.AspectRatio,
		}
	case "image_to_video":
		powerCost = h.getPowerFromConfig(model.JMTaskTypeImageToVideo)
		taskType = model.JMTaskTypeImageToVideo
		reqKey = jimeng.ReqKeyImageToVideo
		modelName = "即梦图生视频"
		if req.Seed == 0 {
			req.Seed = -1
		}
		params = map[string]any{
			"seed":         req.Seed,
			"aspect_ratio": req.AspectRatio,
		}
		if len(req.ImageUrls) > 0 {
			params["image_urls"] = req.ImageUrls
		}
		if len(req.BinaryDataBase64) > 0 {
			params["binary_data_base64"] = req.BinaryDataBase64
		}
	default:
		resp.ERROR(c, "不支持的任务类型")
		return
	}

	if user.Power < powerCost {
		resp.ERROR(c, fmt.Sprintf("算力不足，需要%d算力", powerCost))
		return
	}

	taskReq := &jimeng.CreateTaskRequest{
		Type:   taskType,
		Prompt: req.Prompt,
		Params: params,
		ReqKey: reqKey,
		Power:  powerCost,
	}

	job, err := h.jimengService.CreateTask(user.Id, taskReq)
	if err != nil {
		logger.Errorf("create jimeng task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	h.userService.DecreasePower(user.Id, powerCost, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "jimeng",
		Remark: fmt.Sprintf("%s，任务ID：%d", modelName, job.Id),
	})

	resp.SUCCESS(c, job)
}

// Jobs 获取任务列表
func (h *JimengHandler) Jobs(c *gin.Context) {
	userId := h.GetLoginUserId(c)

	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Filter   string `json:"filter"`
		Ids      []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var jobs []model.JimengJob
	var total int64
	query := h.DB.Model(&model.JimengJob{}).Where("user_id = ?", userId)

	switch req.Filter {
	case "image":
		query = query.Where("type IN (?)", []model.JMTaskType{
			model.JMTaskTypeTextToImage,
			model.JMTaskTypeImageToImage,
			model.JMTaskTypeImageEdit,
			model.JMTaskTypeImageEffects,
		})
	case "video":
		query = query.Where("type IN (?)", []model.JMTaskType{
			model.JMTaskTypeTextToVideo,
			model.JMTaskTypeImageToVideo,
		})
	}

	if len(req.Ids) > 0 {
		query = query.Where("id IN (?)", req.Ids)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("updated_at DESC").Offset(offset).Limit(req.PageSize).Find(&jobs).Error; err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 填充 VO
	var jobVos []vo.JimengJob
	for _, job := range jobs {
		var jobVo vo.JimengJob
		err := utils.CopyObject(job, &jobVo)
		if err != nil {
			continue
		}
		jobVo.CreatedAt = job.CreatedAt.Unix()
		jobVos = append(jobVos, jobVo)
	}
	resp.SUCCESS(c, vo.NewPage(total, req.Page, req.PageSize, jobVos))
}

// Remove 删除任务
func (h *JimengHandler) Remove(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	jobId := h.GetInt(c, "id", 0)
	if jobId == 0 {
		resp.ERROR(c, "参数错误")
		return
	}

	// 获取任务，判断状态
	job, err := h.jimengService.GetJob(uint(jobId))
	if err != nil {
		resp.ERROR(c, "任务不存在")
		return
	}
	if job.UserId != user.Id {
		resp.ERROR(c, "无权限操作")
		return
	}

	// 正在运行中的任务不能删除
	if job.Status == model.JMTaskStatusGenerating || job.Status == model.JMTaskStatusInQueue {
		resp.ERROR(c, "正在运行中的任务不能删除，否则无法退回算力")
		return
	}

	tx := h.DB.Begin()
	if err := tx.Where("id = ? AND user_id = ?", jobId, user.Id).Delete(&model.JimengJob{}).Error; err != nil {
		logger.Errorf("delete jimeng job failed: %v", err)
		resp.ERROR(c, "删除任务失败")
		return
	}

	// 失败任务删除后退回算力
	if job.Status != model.JMTaskStatusFailed {
		err = h.userService.IncreasePower(user.Id, job.Power, model.PowerLog{
			Type:   types.PowerRefund,
			Model:  "jimeng",
			Remark: fmt.Sprintf("删除任务，退回%d算力", job.Power),
		})
		if err != nil {
			resp.ERROR(c, "退回算力失败")
			tx.Rollback()
			return
		}
	}

	tx.Commit()

	resp.SUCCESS(c, gin.H{})
}

// Retry 重试任务
func (h *JimengHandler) Retry(c *gin.Context) {
	userId := h.GetLoginUserId(c)

	jobId := h.GetInt(c, "id", 0)
	if jobId == 0 {
		resp.ERROR(c, "参数错误")
		return
	}

	// 检查任务是否存在且属于当前用户
	job, err := h.jimengService.GetJob(uint(jobId))
	if err != nil {
		resp.ERROR(c, "任务不存在")
		return
	}

	if job.UserId != userId {
		resp.ERROR(c, "无权限操作")
		return
	}

	// 只有失败的任务才能重试
	if job.Status != model.JMTaskStatusFailed {
		resp.ERROR(c, "只有失败的任务才能重试")
		return
	}

	// 重置任务状态
	if err := h.jimengService.UpdateJobStatus(uint(jobId), model.JMTaskStatusInQueue, ""); err != nil {
		logger.Errorf("reset job status failed: %v", err)
		resp.ERROR(c, "重置任务状态失败")
		return
	}

	// 重新推送到队列
	if err := h.jimengService.PushTaskToQueue(uint(jobId)); err != nil {
		logger.Errorf("push retry task to queue failed: %v", err)
		resp.ERROR(c, "推送重试任务失败")
		return
	}

	resp.SUCCESS(c, gin.H{"message": "重试任务已提交"})
}

// getPowerFromConfig 从配置中获取指定类型的算力消耗
func (h *JimengHandler) getPowerFromConfig(taskType model.JMTaskType) int {
	config := h.jimengService.GetConfig()

	switch taskType {
	case model.JMTaskTypeTextToImage:
		return config.Power.TextToImage
	case model.JMTaskTypeImageToImage:
		return config.Power.ImageToImage
	case model.JMTaskTypeImageEdit:
		return config.Power.ImageEdit
	case model.JMTaskTypeImageEffects:
		return config.Power.ImageEffects
	case model.JMTaskTypeTextToVideo:
		return config.Power.TextToVideo
	case model.JMTaskTypeImageToVideo:
		return config.Power.ImageToVideo
	default:
		return 10
	}
}

// GetPowerConfig 获取即梦各任务类型算力消耗配置
func (h *JimengHandler) GetPowerConfig(c *gin.Context) {
	config := h.jimengService.GetConfig()
	resp.SUCCESS(c, gin.H{
		"text_to_image":  config.Power.TextToImage,
		"image_to_image": config.Power.ImageToImage,
		"image_edit":     config.Power.ImageEdit,
		"image_effects":  config.Power.ImageEffects,
		"text_to_video":  config.Power.TextToVideo,
		"image_to_video": config.Power.ImageToVideo,
	})
}
