package admin

import (
	"strconv"

	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/service/jimeng"
	"geekai/store/model"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminJimengHandler 管理后台即梦AI处理器
type AdminJimengHandler struct {
	handler.BaseHandler
	jimengService *jimeng.Service
}

// NewAdminJimengHandler 创建管理后台即梦AI处理器
func NewAdminJimengHandler(app *core.AppServer, db *gorm.DB, jimengService *jimeng.Service) *AdminJimengHandler {
	return &AdminJimengHandler{
		BaseHandler:   handler.BaseHandler{App: app, DB: db},
		jimengService: jimengService,
	}
}

// RegisterRoutes 注册即梦AI管理后台路由
func (h *AdminJimengHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/jimeng/")
	rg.GET("/jobs", h.Jobs)
	rg.GET("/jobs/:id", h.JobDetail)
	rg.DELETE("/jobs/:id", h.Remove)
	rg.POST("/jobs/batch-remove", h.BatchRemove)
	rg.GET("/stats", h.Stats)
	rg.GET("/config", h.GetConfig)
	rg.POST("/config/update", h.UpdateConfig)
}

// Jobs 获取任务列表
func (h *AdminJimengHandler) Jobs(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	userId := h.GetInt(c, "user_id", 0)
	taskType := h.GetTrim(c, "type")
	status := h.GetTrim(c, "status")

	var tasks []model.JimengJob
	var total int64

	session := h.DB.Model(&model.JimengJob{})

	// 构建查询条件
	if userId > 0 {
		session = session.Where("user_id = ?", userId)
	}
	if taskType != "" {
		session = session.Where("type = ?", taskType)
	}
	if status != "" {
		session = session.Where("status = ?", status)
	}

	// 获取总数
	err := session.Count(&total).Error
	if err != nil {
		resp.ERROR(c, "获取任务数量失败")
		return
	}

	// 获取数据
	offset := (page - 1) * pageSize
	err = session.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		resp.ERROR(c, "获取任务列表失败")
		return
	}

	resp.SUCCESS(c, gin.H{
		"jobs":      tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// JobDetail 获取任务详情
func (h *AdminJimengHandler) JobDetail(c *gin.Context) {
	idStr := c.Param("id")
	jobId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	var job model.JimengJob
	err = h.DB.Where("id = ?", jobId).First(&job).Error
	if err != nil {
		resp.ERROR(c, "任务不存在")
		return
	}

	resp.SUCCESS(c, job)
}

// Remove 删除任务
func (h *AdminJimengHandler) Remove(c *gin.Context) {
	idStr := c.Param("id")
	jobId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	err = h.DB.Where("id = ?", jobId).Delete(&model.JimengJob{}).Error
	if err != nil {
		resp.ERROR(c, "删除任务失败")
		return
	}

	resp.SUCCESS(c, gin.H{})
}

// BatchRemove 批量删除任务
func (h *AdminJimengHandler) BatchRemove(c *gin.Context) {
	var req struct {
		JobIds []uint `json:"job_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	result := h.DB.Where("id IN ?", req.JobIds).Delete(&model.JimengJob{})
	if result.Error != nil {
		resp.ERROR(c, "批量删除失败")
		return
	}

	resp.SUCCESS(c, gin.H{
		"message":       "批量删除成功",
		"deleted_count": result.RowsAffected,
	})
}

// Stats 获取统计信息
func (h *AdminJimengHandler) Stats(c *gin.Context) {
	type StatResult struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var stats []StatResult
	err := h.DB.Model(&model.JimengJob{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&stats).Error
	if err != nil {
		resp.ERROR(c, "获取统计信息失败")
		return
	}

	// 整理统计数据
	result := gin.H{
		"totalTasks":      int64(0),
		"completedTasks":  int64(0),
		"processingTasks": int64(0),
		"failedTasks":     int64(0),
		"pendingTasks":    int64(0),
	}

	for _, stat := range stats {
		result["totalTasks"] = result["totalTasks"].(int64) + stat.Count
		switch stat.Status {
		case "completed":
			result["completedTasks"] = stat.Count
		case "processing":
			result["processingTasks"] = stat.Count
		case "failed":
			result["failedTasks"] = stat.Count
		case "pending":
			result["pendingTasks"] = stat.Count
		}
	}

	resp.SUCCESS(c, result)
}

// GetConfig 获取即梦AI配置
func (h *AdminJimengHandler) GetConfig(c *gin.Context) {
	var config model.Config
	err := h.DB.Debug().Where("name", "jimeng").First(&config).Error
	if err != nil {
		// 如果配置不存在，返回默认配置
		defaultConfig := types.JimengConfig{
			AccessKey: "",
			SecretKey: "",
			Power: types.JimengPower{
				TextToImage:  10,
				ImageToImage: 15,
				ImageEdit:    20,
				ImageEffects: 25,
				TextToVideo:  30,
				ImageToVideo: 35,
			},
		}
		resp.SUCCESS(c, defaultConfig)
		return
	}

	var jimengConfig types.JimengConfig
	err = utils.JsonDecode(config.Value, &jimengConfig)
	if err != nil {
		resp.ERROR(c, "解析配置失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, jimengConfig)
}

// UpdateConfig 更新即梦AI配置
func (h *AdminJimengHandler) UpdateConfig(c *gin.Context) {
	var req types.JimengConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	// 验证必填字段
	if req.AccessKey == "" {
		resp.ERROR(c, "AccessKey不能为空")
		return
	}
	if req.SecretKey == "" {
		resp.ERROR(c, "SecretKey不能为空")
		return
	}

	// 验证算力配置
	if req.Power.TextToImage <= 0 {
		resp.ERROR(c, "文生图算力必须大于0")
		return
	}
	if req.Power.ImageToImage <= 0 {
		resp.ERROR(c, "图生图算力必须大于0")
		return
	}
	if req.Power.ImageEdit <= 0 {
		resp.ERROR(c, "图片编辑算力必须大于0")
		return
	}
	if req.Power.ImageEffects <= 0 {
		resp.ERROR(c, "图片特效算力必须大于0")
		return
	}
	if req.Power.TextToVideo <= 0 {
		resp.ERROR(c, "文生视频算力必须大于0")
		return
	}
	if req.Power.ImageToVideo <= 0 {
		resp.ERROR(c, "图生视频算力必须大于0")
		return
	}

	// 保存配置
	tx := h.DB.Begin()
	value := utils.JsonEncode(&req)
	config := model.Config{Name: "jimeng", Value: value}

	err := tx.FirstOrCreate(&config, model.Config{Name: "jimeng"}).Error
	if err != nil {
		resp.ERROR(c, "保存配置失败: "+err.Error())
		return
	}

	if config.Id > 0 {
		config.Value = value
		err = tx.Updates(&config).Error
		if err != nil {
			resp.ERROR(c, "更新配置失败: "+err.Error())
			return
		}
	}

	// 更新服务中的客户端配置
	updateErr := h.jimengService.UpdateClientConfig(req.AccessKey, req.SecretKey)
	if updateErr != nil {
		resp.ERROR(c, updateErr.Error())
		tx.Rollback()
		return
	}
	tx.Commit()

	resp.SUCCESS(c, gin.H{"message": "配置更新成功"})
}
