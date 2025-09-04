package handler

import (
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
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
	service *ai3d.Service
}

func NewAI3DHandler(app *core.AppServer, db *gorm.DB, service *ai3d.Service) *AI3DHandler {
	return &AI3DHandler{
		service: service,
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
		group.GET("job/delete", h.DeleteJob)
	}
}

// Generate 创建3D生成任务
func (h *AI3DHandler) Generate(c *gin.Context) {
	var request vo.AI3DJobParams

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

	// 获取用户ID
	userId := h.GetLoginUserId(c)
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

// DeleteJob 删除任务
func (h *AI3DHandler) DeleteJob(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	id := h.GetInt(c, "id", 0)
	if id == 0 {
		resp.ERROR(c, "任务ID不能为空")
		return
	}

	err := h.service.DeleteUserJob(uint(id), uint(userId))
	if err != nil {
		resp.ERROR(c, "删除任务失败")
		return
	}

	resp.SUCCESS(c, gin.H{"message": "删除成功"})
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
			Status:     types.AI3DJobStatusSuccess,
			ErrMsg:     "",
			Params:     vo.AI3DJobParams{Prompt: "一只可爱的小猫", ImageURL: "", Texture: true, Seed: 42},
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
			Params:     vo.AI3DJobParams{Prompt: "一个现代建筑模型", ImageURL: "", EnablePBR: true},
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
			Params:     vo.AI3DJobParams{Prompt: "一辆跑车模型", ImageURL: "https://example.com/car.jpg", Texture: false},
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
			Params:     vo.AI3DJobParams{Prompt: "一个机器人模型", ImageURL: "https://example.com/robot.jpg", EnablePBR: false},
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
			Status:     types.AI3DJobStatusSuccess,
			ErrMsg:     "",
			Params:     vo.AI3DJobParams{Prompt: "一个复杂的机械装置", ImageURL: "", Texture: true, OctreeResolution: 512},
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
			Params:     vo.AI3DJobParams{Prompt: "一个科幻飞船", ImageURL: "", EnablePBR: true},
			CreatedAt:  1704085200, // 2024-01-01 05:00:00
			UpdatedAt:  1704085200, // 2024-01-01 05:00:00
		},
	}

	// 创建分页响应
	mockResponse := vo.Page{
		Page:     1,
		PageSize: 10,
		Total:    int64(len(mockJobs)),
		Items:    mockJobs,
	}

	resp.SUCCESS(c, mockResponse)
}
