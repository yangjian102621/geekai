package handler

import (
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/ai3d"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AI3DHandler struct {
	BaseHandler
	service     *ai3d.Service
	userService *service.UserService
}

func NewAI3DHandler(app *core.AppServer, db *gorm.DB, service *ai3d.Service, userService *service.UserService) *AI3DHandler {
	return &AI3DHandler{
		service:     service,
		userService: userService,
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
	}
}

// RegisterRoutes 注册路由
func (h *AI3DHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/ai3d/")

	// 公开接口，不需要授权
	group.GET("configs", h.GetConfigs)

	// 需要用户授权的接口
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.POST("generate", h.Generate)
		group.GET("jobs", h.JobList)
		group.GET("job/:id", h.JobDetail)
		group.DELETE("job/:id", h.DeleteJob)
		group.GET("download/:id", h.Download)
	}
}

// Generate 创建3D生成任务
func (h *AI3DHandler) Generate(c *gin.Context) {
	var request vo.AI3DJobCreate
	if err := c.ShouldBindJSON(&request); err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	// 验证必填参数
	if request.Type == "" || request.Model == "" || request.Power <= 0 {
		resp.ERROR(c, "缺少必要参数")
		return
	}

	// 获取用户ID
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.ERROR(c, "用户未登录")
		return
	}

	// 创建任务
	job, err := h.service.CreateJob(uint(userId), request)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("创建任务失败: %v", err))
		return
	}

	resp.SUCCESS(c, gin.H{
		"job_id":  job.Id,
		"message": "任务创建成功",
	})
}

// JobList 获取任务列表
func (h *AI3DHandler) JobList(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.ERROR(c, "用户未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	jobList, err := h.service.GetJobList(uint(userId), page, pageSize)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("获取任务列表失败: %v", err))
		return
	}

	resp.SUCCESS(c, jobList)
}

// JobDetail 获取任务详情
func (h *AI3DHandler) JobDetail(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.ERROR(c, "用户未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "任务ID格式错误")
		return
	}

	job, err := h.service.GetJobById(uint(id))
	if err != nil {
		resp.ERROR(c, "任务不存在")
		return
	}

	// 检查权限
	if job.UserId != uint(userId) {
		resp.ERROR(c, "无权限访问此任务")
		return
	}

	// 转换为VO
	jobVO := vo.AI3DJob{
		Id:         job.Id,
		UserId:     job.UserId,
		Type:       job.Type,
		Power:      job.Power,
		TaskId:     job.TaskId,
		ImgURL:     job.FileURL,
		PreviewURL: job.PreviewURL,
		Model:      job.Model,
		Status:     job.Status,
		ErrMsg:     job.ErrMsg,
		Params:     job.Params,
		CreatedAt:  job.CreatedAt.Unix(),
		UpdatedAt:  job.UpdatedAt.Unix(),
	}

	resp.SUCCESS(c, jobVO)
}

// DeleteJob 删除任务
func (h *AI3DHandler) DeleteJob(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.ERROR(c, "用户未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "任务ID格式错误")
		return
	}

	err = h.service.DeleteJob(uint(id), uint(userId))
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("删除任务失败: %v", err))
		return
	}

	resp.SUCCESS(c, gin.H{"message": "删除成功"})
}

// Download 下载3D模型
func (h *AI3DHandler) Download(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.ERROR(c, "用户未登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "任务ID格式错误")
		return
	}

	job, err := h.service.GetJobById(uint(id))
	if err != nil {
		resp.ERROR(c, "任务不存在")
		return
	}

	// 检查权限
	if job.UserId != uint(userId) {
		resp.ERROR(c, "无权限访问此任务")
		return
	}

	// 检查任务状态
	if job.Status != types.AI3DJobStatusCompleted {
		resp.ERROR(c, "任务尚未完成")
		return
	}

	if job.FileURL == "" {
		resp.ERROR(c, "模型文件不存在")
		return
	}

	// 重定向到下载链接
	c.Redirect(302, job.FileURL)
}

// GetConfigs 获取3D生成配置
func (h *AI3DHandler) GetConfigs(c *gin.Context) {
	var config model.Config
	err := h.DB.Where("name", types.ConfigKeyAI3D).First(&config).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	var config3d types.AI3DConfig
	err = utils.JsonDecode(config.Value, &config3d)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	models := h.service.GetSupportedModels()
	if len(config3d.Gitee.Models) == 0 {
		config3d.Gitee.Models = models["gitee"]
	}
	if len(config3d.Tencent.Models) == 0 {
		config3d.Tencent.Models = models["tencent"]
	}

	logger.Info("config3d: ", config3d)

	resp.SUCCESS(c, config3d)
}
