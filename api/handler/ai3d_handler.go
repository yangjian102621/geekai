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
		group.GET("jobs/mock", h.ListMock) // 演示数据接口
		group.GET("job/:id", h.JobDetail)
		group.GET("job/delete", h.DeleteJob)
		group.GET("download/:id", h.Download)
	}
}

// Generate 创建3D生成任务
func (h *AI3DHandler) Generate(c *gin.Context) {
	var request struct {
		// 通用参数
		Type       types.AI3DTaskType `json:"type" binding:"required"`  // API类型 (tencent/gitee)
		Model      string             `json:"model" binding:"required"` // 3D模型类型
		Prompt     string             `json:"prompt"`                   // 文本提示词
		ImageURL   string             `json:"image_url"`                // 输入图片URL
		FileFormat string             `json:"file_format"`              // 输出文件格式
		// 腾讯3d专有参数
		EnablePBR bool `json:"enable_pbr"` // 是否开启PBR材质
		// Gitee3d专有参数
		Texture           bool    `json:"texture"`             // 是否开启纹理
		Seed              int     `json:"seed"`                // 随机种子
		NumInferenceSteps int     `json:"num_inference_steps"` //迭代次数
		GuidanceScale     float64 `json:"guidance_scale"`      //引导系数
		OctreeResolution  int     `json:"octree_resolution"`   // 3D 渲染精度，越高3D 细节越丰富
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	// 提示词和图片不能同时为空
	if request.Prompt == "" && request.ImageURL == "" {
		resp.ERROR(c, "提示词和图片不能同时为空")
		return
	}

	// Gitee 只支持图片
	if request.Type == types.AI3DTaskTypeGitee && request.ImageURL == "" {
		resp.ERROR(c, "Gitee 只支持图生3D")
		return
	}

	logger.Infof("request: %+v", request)

	// // 获取用户ID
	// userId := h.GetLoginUserId(c)
	// // 创建任务
	// job, err := h.service.CreateJob(uint(userId), request)
	// if err != nil {
	// 	resp.ERROR(c, fmt.Sprintf("创建任务失败: %v", err))
	// 	return
	// }

	resp.SUCCESS(c, gin.H{
		"job_id":  0,
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
		FileURL:    job.FileURL,
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
	id := c.Query("id")
	if id == "" {
		resp.ERROR(c, "任务ID不能为空")
		return
	}

	var job model.AI3DJob
	err := h.DB.Where("id = ?", id).Where("user_id = ?", userId).First(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	err = h.DB.Delete(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 失败的任务要退回算力
	if job.Status == types.AI3DJobStatusFailed {
		err = h.userService.IncreasePower(userId, job.Power, model.PowerLog{
			Type:   types.PowerRefund,
			Model:  job.Model,
			Remark: fmt.Sprintf("删除任务，退回%d算力", job.Power),
		})
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
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

// ListMock 返回演示数据
func (h *AI3DHandler) ListMock(c *gin.Context) {
	// 创建各种状态的演示数据
	mockJobs := []vo.AI3DJob{
		{
			Id:         1,
			UserId:     1,
			Type:       "gitee",
			Power:      10,
			TaskId:     "mock_task_1",
			FileURL:    "https://img.r9it.com/R03TQZ7PZ386RGL7PTMNGFOHAJW15WYF.glb",
			PreviewURL: "/static/upload/2025/9/1756873317505073.png",
			Model:      "gitee-3d-v1",
			Status:     types.AI3DJobStatusCompleted,
			ErrMsg:     "",
			Params:     `{"prompt":"一只可爱的小猫","image_url":"","texture":true,"seed":42}`,
			CreatedAt:  1704067200, // 2024-01-01 00:00:00
			UpdatedAt:  1704067800, // 2024-01-01 00:10:00
		},
		{
			Id:         2,
			UserId:     1,
			Type:       "tencent",
			Power:      15,
			TaskId:     "mock_task_2",
			FileURL:    "",
			PreviewURL: "/static/upload/2025/9/1756873317505073.png",
			Model:      "tencent-3d-v2",
			Status:     types.AI3DJobStatusProcessing,
			ErrMsg:     "",
			Params:     `{"prompt":"一个现代建筑模型","image_url":"","enable_pbr":true}`,
			CreatedAt:  1704070800, // 2024-01-01 01:00:00
			UpdatedAt:  1704070800, // 2024-01-01 01:00:00
		},
		{
			Id:         3,
			UserId:     1,
			Type:       "gitee",
			Power:      8,
			TaskId:     "mock_task_3",
			FileURL:    "",
			PreviewURL: "",
			Model:      "gitee-3d-v1",
			Status:     types.AI3DJobStatusPending,
			ErrMsg:     "",
			Params:     `{"prompt":"一辆跑车模型","image_url":"https://example.com/car.jpg","texture":false}`,
			CreatedAt:  1704074400, // 2024-01-01 02:00:00
			UpdatedAt:  1704074400, // 2024-01-01 02:00:00
		},
		{
			Id:         4,
			UserId:     1,
			Type:       "tencent",
			Power:      12,
			TaskId:     "mock_task_4",
			FileURL:    "",
			PreviewURL: "",
			Model:      "tencent-3d-v1",
			Status:     types.AI3DJobStatusFailed,
			ErrMsg:     "模型生成失败：输入图片质量不符合要求",
			Params:     `{"prompt":"一个机器人模型","image_url":"https://example.com/robot.jpg","enable_pbr":false}`,
			CreatedAt:  1704078000, // 2024-01-01 03:00:00
			UpdatedAt:  1704078600, // 2024-01-01 03:10:00
		},
		{
			Id:         5,
			UserId:     1,
			Type:       "gitee",
			Power:      20,
			TaskId:     "mock_task_5",
			FileURL:    "https://ai.gitee.com/a8c1af8e-26e9-4ca6-aa5c-6d4ba86bfdac",
			PreviewURL: "https://ai.gitee.com/a8c1af8e-26e9-4ca6-aa5c-6d4ba86bfdac",
			Model:      "gitee-3d-v2",
			Status:     types.AI3DJobStatusCompleted,
			ErrMsg:     "",
			Params:     `{"prompt":"一个复杂的机械装置","image_url":"","texture":true,"octree_resolution":512}`,
			CreatedAt:  1704081600, // 2024-01-01 04:00:00
			UpdatedAt:  1704082200, // 2024-01-01 04:10:00
		},
		{
			Id:         6,
			UserId:     1,
			Type:       "tencent",
			Power:      18,
			TaskId:     "mock_task_6",
			FileURL:    "",
			PreviewURL: "",
			Model:      "tencent-3d-v2",
			Status:     types.AI3DJobStatusProcessing,
			ErrMsg:     "",
			Params:     `{"prompt":"一个科幻飞船","image_url":"","enable_pbr":true}`,
			CreatedAt:  1704085200, // 2024-01-01 05:00:00
			UpdatedAt:  1704085200, // 2024-01-01 05:00:00
		},
	}

	// 创建分页响应
	mockResponse := vo.ThreeDJobList{
		Page:     1,
		PageSize: 10,
		Total:    len(mockJobs),
		Items:    mockJobs,
	}

	resp.SUCCESS(c, mockResponse)
}
