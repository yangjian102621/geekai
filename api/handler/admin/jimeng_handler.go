package admin

import (
	"strconv"

	"geekai/core"
	"geekai/handler"
	"geekai/service/jimeng"
	"geekai/store/model"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
)

// AdminJimengHandler 管理后台即梦AI处理器
type AdminJimengHandler struct {
	handler.BaseHandler
	jimengService *jimeng.Service
}

// NewAdminJimengHandler 创建管理后台即梦AI处理器
func NewAdminJimengHandler(app *core.AppServer, jimengService *jimeng.Service) *AdminJimengHandler {
	return &AdminJimengHandler{
		BaseHandler:   handler.BaseHandler{App: app},
		jimengService: jimengService,
	}
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