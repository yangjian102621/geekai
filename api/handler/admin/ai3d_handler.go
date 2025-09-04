package admin

import (
	"strconv"

	"geekai/core"
	"geekai/core/types"
	"geekai/service/ai3d"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AI3DHandler 3D管理处理器
type AI3DHandler struct {
	app     *core.AppServer
	db      *gorm.DB
	service *ai3d.Service
}

// NewAI3DHandler 创建3D管理处理器
func NewAI3DHandler(app *core.AppServer, db *gorm.DB, service *ai3d.Service) *AI3DHandler {
	return &AI3DHandler{
		app:     app,
		db:      db,
		service: service,
	}
}

// RegisterRoutes 注册路由
func (h *AI3DHandler) RegisterRoutes() {
	admin := h.app.Engine.Group("/api/admin/ai3d")
	{
		admin.GET("/jobs", h.GetJobList)
		admin.GET("/jobs/:id", h.GetJobDetail)
		admin.DELETE("/jobs/:id", h.DeleteJob)
		admin.GET("/stats", h.GetStats)
		admin.GET("/models", h.GetModels)
		admin.POST("/config", h.SaveConfig)
	}
}

// GetJobList 获取任务列表
func (h *AI3DHandler) GetJobList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	jobType := c.Query("type")
	userIdStr := c.Query("user_id")

	var userId uint
	if userIdStr != "" {
		if id, err := strconv.ParseUint(userIdStr, 10, 32); err == nil {
			userId = uint(id)
		}
	}

	// 构建查询条件
	query := h.db.Model(&model.AI3DJob{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if jobType != "" {
		query = query.Where("type = ?", jobType)
	}

	if userId > 0 {
		query = query.Where("user_id = ?", userId)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var jobs []model.AI3DJob
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&jobs).Error

	if err != nil {
		resp.ERROR(c, "获取任务列表失败")
		return
	}

	// 转换为VO
	var jobList []vo.AI3DJob
	for _, job := range jobs {
		var jobVo vo.AI3DJob
		err = utils.CopyObject(job, &jobVo)
		if err != nil {
			continue
		}
		utils.JsonDecode(job.Params, &jobVo.Params)
		jobVo.CreatedAt = job.CreatedAt.Unix()
		jobVo.UpdatedAt = job.UpdatedAt.Unix()
		jobList = append(jobList, jobVo)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, jobList))
}

// GetJobDetail 获取任务详情
func (h *AI3DHandler) GetJobDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "无效的任务ID")
		return
	}

	var job model.AI3DJob
	err = h.db.First(&job, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.ERROR(c, "任务不存在")
		} else {
			resp.ERROR(c, "获取任务详情失败")
		}
		return
	}

	var jobVo vo.AI3DJob
	err = utils.CopyObject(job, &jobVo)
	if err != nil {
		resp.ERROR(c, "获取任务详情失败")
		return
	}
	utils.JsonDecode(job.Params, &jobVo.Params)
	jobVo.CreatedAt = job.CreatedAt.Unix()
	jobVo.UpdatedAt = job.UpdatedAt.Unix()
	resp.SUCCESS(c, jobVo)
}

// DeleteJob 删除任务
func (h *AI3DHandler) DeleteJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "无效的任务ID")
		return
	}

	// 检查任务是否存在
	var job model.AI3DJob
	err = h.db.First(&job, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.ERROR(c, "任务不存在")
		} else {
			resp.ERROR(c, "获取任务失败")
		}
		return
	}

	// 删除任务
	err = h.db.Delete(&job).Error
	if err != nil {
		resp.ERROR(c, "删除任务失败")
		return
	}

	resp.SUCCESS(c, "删除成功")
}

// GetStats 获取统计数据
func (h *AI3DHandler) GetStats(c *gin.Context) {
	var stats struct {
		Pending    int64 `json:"pending"`
		Processing int64 `json:"processing"`
		Success    int64 `json:"success"`
		Failed     int64 `json:"failed"`
	}

	// 统计各状态的任务数量
	h.db.Model(&model.AI3DJob{}).Where("status = ?", "pending").Count(&stats.Pending)
	h.db.Model(&model.AI3DJob{}).Where("status = ?", "processing").Count(&stats.Processing)
	h.db.Model(&model.AI3DJob{}).Where("status = ?", "success").Count(&stats.Success)
	h.db.Model(&model.AI3DJob{}).Where("status = ?", "failed").Count(&stats.Failed)

	resp.SUCCESS(c, stats)
}

// GetModels 获取配置
func (h *AI3DHandler) GetModels(c *gin.Context) {
	models := h.service.GetSupportedModels()
	resp.SUCCESS(c, models)
}

// SaveGlobalSettings 保存全局配置
func (h *AI3DHandler) SaveConfig(c *gin.Context) {
	var config types.AI3DConfig
	err := c.ShouldBindJSON(&config)
	if err != nil {
		resp.ERROR(c, "参数错误")
		return
	}
	var exist model.Config
	err = h.db.Where("name", types.ConfigKeyAI3D).First(&exist).Error
	if err != nil {
		exist.Name = types.ConfigKeyAI3D
		exist.Value = utils.JsonEncode(config)
		err = h.db.Create(&exist).Error
	} else {
		exist.Value = utils.JsonEncode(config)
		err = h.db.Updates(&exist).Error
	}
	if err != nil {
		resp.ERROR(c, "保存配置失败")
		return
	}

	h.service.UpdateConfig(config)
	h.app.SysConfig.AI3D = config

	resp.SUCCESS(c, "保存成功")
}
