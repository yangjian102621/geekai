package handler

import (
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/jimeng"
	"geekai/service/moderation"
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
	jimengService     *jimeng.Service
	userService       *service.UserService
	moderationManager *moderation.ServiceManager
}

// NewJimengHandler 创建即梦AI处理器
func NewJimengHandler(app *core.AppServer, jimengService *jimeng.Service, db *gorm.DB, userService *service.UserService, moderationManager *moderation.ServiceManager) *JimengHandler {
	return &JimengHandler{
		BaseHandler:       BaseHandler{App: app, DB: db},
		jimengService:     jimengService,
		userService:       userService,
		moderationManager: moderationManager,
	}
}

// RegisterRoutes 注册路由，新增统一任务接口
func (h *JimengHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/jimeng/")
	group.GET("power-config", h.GetPowerConfig)

	// 需要用户授权的接口
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.POST("task", h.CreateTask)
		group.POST("jobs", h.Jobs)
		group.GET("remove", h.Remove)
		group.GET("retry", h.Retry)
	}
}

// CreateTask 统一任务创建接口
func (h *JimengHandler) CreateTask(c *gin.Context) {
	var req types.JimengTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 文本审核
	if h.App.SysConfig.Moderation.Enable && req.Prompt != "" {
		moderationResult, err := h.moderationManager.GetService().Moderate(req.Prompt)
		if err != nil {
			logger.Error("failed to moderate content: ", err)
		}
		if moderationResult.Flagged {
			// 记录违规内容
			moderation := model.Moderation{
				UserId: h.GetLoginUserId(c),
				Source: types.ModerationSourceJiMeng,
				Input:  req.Prompt,
				Result: utils.JsonEncode(moderationResult),
			}
			err = h.DB.Create(&moderation).Error
			if err != nil {
				logger.Error("failed to save moderation: ", err)
			}
			resp.ERROR(c, "当前创作内容包含敏感词，请重新输入！")
			return
		}

	}

	if req.Prompt == "" && len(req.ImageUrls) == 0 {
		resp.ERROR(c, "提示词和图片不能同时为空")
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 获取算力消耗
	powerCost, err := h.getTaskPower(req)
	if err != nil {
		resp.ERROR(c, "计算任务消耗积分失败: "+err.Error())
		return
	}

	if user.Power < powerCost {
		resp.ERROR(c, fmt.Sprintf("算力不足，需要%d算力", powerCost))
		return
	}
	req.Power = powerCost

	job, err := h.jimengService.CreateTask(user.Id, &req)
	if err != nil {
		logger.Errorf("create jimeng task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	h.userService.DecreasePower(user.Id, powerCost, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  job.ReqKey,
		Remark: h.getTaskRemark(req, job.Id),
	})

	resp.SUCCESS(c)
}

func (h *JimengHandler) getTaskRemark(req types.JimengTaskRequest, jobId uint) string {
	remark := fmt.Sprintf("即梦任务%s，任务ID：%d", req.ReqKey, jobId)
	perUnit, ok := h.App.SysConfig.Jimeng.Powers[req.ReqKey]
	if !ok || perUnit <= 0 {
		return remark // Fallback if power not found or invalid
	}
	switch req.TaskType {
	case types.JMTaskTypeImage:
		remark = fmt.Sprintf("即梦图片生成，任务ID：%d，%d积分/张", jobId, perUnit)
	case types.JMTaskTypeVideo:
		seconds := 0
		if perUnit > 0 {
			seconds = req.Power / perUnit
		}
		remark = fmt.Sprintf("即梦视频生成，任务ID：%d，%d积分/秒, %d秒", jobId, perUnit, seconds)
	case types.JMTaskTypeVirtualHuman:
		seconds := 0
		if perUnit > 0 {
			seconds = req.Power / perUnit
		}
		remark = fmt.Sprintf("即梦数字人视频生成，任务ID：%d，%d积分/秒, %d秒", jobId, perUnit, seconds)
	case types.JMTaskTypeActionTransfer:
		seconds := 0
		if perUnit > 0 {
			seconds = req.Power / perUnit
		}
		remark = fmt.Sprintf("即梦视频动作迁移，任务ID：%d，%d积分/秒, %d秒", jobId, perUnit, seconds)
	}
	return remark
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
		query = query.Where("type = ?", types.JMTaskTypeImage)
	case "video":
		query = query.Where("type = ?", types.JMTaskTypeVideo)
	case "virtual_human":
		query = query.Where("type = ?", types.JMTaskTypeVirtualHuman)
	case "action_transfer":
		query = query.Where("type = ?", types.JMTaskTypeActionTransfer)
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
	if job.Status == types.JMTaskStatusGenerating || job.Status == types.JMTaskStatusInQueue {
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
	if job.Status == types.JMTaskStatusFailed {
		logger.Infof("delete jimeng job failed, refund power: %d", job.Power)
		err = h.userService.IncreasePower(user.Id, job.Power, model.PowerLog{
			Type:   types.PowerRefund,
			Model:  job.ReqKey,
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
	if job.Status != types.JMTaskStatusFailed {
		resp.ERROR(c, "只有失败的任务才能重试")
		return
	}

	// 重置任务状态
	if err := h.jimengService.UpdateJobStatus(uint(jobId), types.JMTaskStatusInQueue, ""); err != nil {
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

func (h *JimengHandler) getTaskPower(req types.JimengTaskRequest) (int, error) {
	logger.Debugf("getTaskPower req: %+v", req)
	config := h.App.SysConfig.Jimeng
	basePower, ok := config.Powers[req.ReqKey]
	if !ok || basePower <= 0 {
		return 0, errors.New("未配置模型积分或配置不合法")
	}
	switch req.TaskType {
	case types.JMTaskTypeImage:
		return basePower, nil
	case types.JMTaskTypeVideo:
		if req.Duration == 0 {
			return 0, errors.New("视频时长不能为0")
		}
		return basePower * req.Duration, nil
	case types.JMTaskTypeVirtualHuman:
		if req.AudioURL == "" {
			return 0, errors.New("音频URL不能为空")
		}
		audioDuration, err := utils.AudioDurationFromURL(req.AudioURL)
		if err != nil {
			return 0, err
		}
		seconds := int(audioDuration.Seconds())
		if seconds <= 0 {
			return 0, errors.New("音频时长无效")
		}
		return basePower * seconds, nil
	case types.JMTaskTypeActionTransfer:
		if req.VideoURL == "" {
			return 0, errors.New("视频URL不能为空")
		}
		videoDuration, err := utils.VideoDurationMP4FromURL(req.VideoURL)
		if err != nil {
			return 0, err
		}
		seconds := int(videoDuration.Seconds())
		if seconds <= 0 {
			return 0, errors.New("视频时长无效")
		}
		return basePower * seconds, nil
	default:
		return 0, errors.New("任务类型不支持")
	}
}

// GetPowerConfig 获取即梦各任务类型算力消耗配置
func (h *JimengHandler) GetPowerConfig(c *gin.Context) {
	config := h.App.SysConfig.Jimeng
	resp.SUCCESS(c, gin.H{
		"powers": config.Powers,
	})
}
