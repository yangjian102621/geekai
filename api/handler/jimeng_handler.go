package handler

import (
	"fmt"
	"strconv"
	"time"

	"geekai/core"
	"geekai/core/types"
	"geekai/service/jimeng"
	"geekai/store/model"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JimengHandler 即梦AI处理器
type JimengHandler struct {
	BaseHandler
	jimengService *jimeng.Service
}

// NewJimengHandler 创建即梦AI处理器
func NewJimengHandler(app *core.AppServer, jimengService *jimeng.Service) *JimengHandler {
	return &JimengHandler{
		BaseHandler:   BaseHandler{App: app},
		jimengService: jimengService,
	}
}

// TextToImage 文生图
func (h *JimengHandler) TextToImage(c *gin.Context) {
	var req struct {
		Prompt    string  `json:"prompt" binding:"required"`
		Seed      int64   `json:"seed"`
		Scale     float64 `json:"scale"`
		Width     int     `json:"width"`
		Height    int     `json:"height"`
		UsePreLLM bool    `json:"use_pre_llm"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 获取当前用户
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 检查用户算力
	if user.Power < 20 { // 文生图消耗20算力
		resp.ERROR(c, "算力不足")
		return
	}

	// 设置默认参数
	if req.Scale == 0 {
		req.Scale = 2.5
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

	// 构建任务参数
	params := map[string]interface{}{
		"seed":        req.Seed,
		"scale":       req.Scale,
		"width":       req.Width,
		"height":      req.Height,
		"use_pre_llm": req.UsePreLLM,
	}

	// 创建任务
	taskReq := &jimeng.CreateTaskRequest{
		Type:   model.JimengJobTypeTextToImage,
		Prompt: req.Prompt,
		Params: params,
		ReqKey: model.ReqKeyTextToImage,
		Power:  20,
	}

	job, err := h.jimengService.CreateTask(user.Id, taskReq)
	if err != nil {
		logger.Errorf("create jimeng text to image task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	// 扣除用户算力
	h.subUserPower(user.Id, 20, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "即梦文生图",
		Remark: fmt.Sprintf("任务ID：%d", job.Id),
	})

	resp.SUCCESS(c, job)
}

// ImageToImagePortrait 图生图人像写真
func (h *JimengHandler) ImageToImagePortrait(c *gin.Context) {
	var req struct {
		ImageInput string  `json:"image_input" binding:"required"`
		Prompt     string  `json:"prompt"`
		Width      int     `json:"width"`
		Height     int     `json:"height"`
		Gpen       float64 `json:"gpen"`
		Skin       float64 `json:"skin"`
		SkinUnifi  float64 `json:"skin_unifi"`
		GenMode    string  `json:"gen_mode"`
		Seed       int64   `json:"seed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 检查用户算力
	if user.Power < 30 { // 图生图消耗30算力
		resp.ERROR(c, "算力不足")
		return
	}

	// 设置默认参数
	if req.Width == 0 {
		req.Width = 1328
	}
	if req.Height == 0 {
		req.Height = 1328
	}
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
	if req.Seed == 0 {
		req.Seed = -1
	}
	if req.Prompt == "" {
		req.Prompt = "演唱会现场的合照，闪光灯拍摄"
	}

	// 构建任务参数
	params := map[string]interface{}{
		"image_input": req.ImageInput,
		"width":       req.Width,
		"height":      req.Height,
		"gpen":        req.Gpen,
		"skin":        req.Skin,
		"skin_unifi":  req.SkinUnifi,
		"gen_mode":    req.GenMode,
		"seed":        req.Seed,
	}

	// 创建任务
	taskReq := &jimeng.CreateTaskRequest{
		Type:   model.JimengJobTypeImageToImagePortrait,
		Prompt: req.Prompt,
		Params: params,
		ReqKey: model.ReqKeyImageToImagePortrait,
		Power:  30,
	}

	job, err := h.jimengService.CreateTask(user.Id, taskReq)
	if err != nil {
		logger.Errorf("create jimeng image to image portrait task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	// 扣除用户算力
	h.subUserPower(user.Id, 30, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "即梦图生图",
		Remark: fmt.Sprintf("任务ID：%d", job.Id),
	})

	resp.SUCCESS(c, job)
}

// ImageEdit 图像编辑
func (h *JimengHandler) ImageEdit(c *gin.Context) {
	var req struct {
		ImageUrls        []string `json:"image_urls"`
		BinaryDataBase64 []string `json:"binary_data_base64"`
		Prompt           string   `json:"prompt" binding:"required"`
		Seed             int64    `json:"seed"`
		Scale            float64  `json:"scale"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误: "+err.Error())
		return
	}

	if len(req.ImageUrls) == 0 && len(req.BinaryDataBase64) == 0 {
		resp.ERROR(c, "请提供图片URL或Base64数据")
		return
	}

	// 获取当前用户
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 检查用户算力
	if user.Power < 25 { // 图像编辑消耗25算力
		resp.ERROR(c, "算力不足")
		return
	}

	// 设置默认参数
	if req.Scale == 0 {
		req.Scale = 0.5
	}
	if req.Seed == 0 {
		req.Seed = -1
	}

	// 构建任务参数
	params := map[string]interface{}{
		"seed":  req.Seed,
		"scale": req.Scale,
	}
	if len(req.ImageUrls) > 0 {
		params["image_urls"] = req.ImageUrls
	}
	if len(req.BinaryDataBase64) > 0 {
		params["binary_data_base64"] = req.BinaryDataBase64
	}

	// 创建任务
	taskReq := &jimeng.CreateTaskRequest{
		Type:   model.JimengJobTypeImageEdit,
		Prompt: req.Prompt,
		Params: params,
		ReqKey: model.ReqKeyImageEdit,
		Power:  25,
	}

	job, err := h.jimengService.CreateTask(user.Id, taskReq)
	if err != nil {
		logger.Errorf("create jimeng image edit task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	// 扣除用户算力
	h.subUserPower(user.Id, 25, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "即梦图像编辑",
		Remark: fmt.Sprintf("任务ID：%d", job.Id),
	})

	resp.SUCCESS(c, job)
}

// ImageEffects 图像特效
func (h *JimengHandler) ImageEffects(c *gin.Context) {
	var req struct {
		ImageInput1 string `json:"image_input1" binding:"required"`
		TemplateId  string `json:"template_id" binding:"required"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 检查用户算力
	if user.Power < 15 { // 图像特效消耗15算力
		resp.ERROR(c, "算力不足")
		return
	}

	// 设置默认参数
	if req.Width == 0 {
		req.Width = 1328
	}
	if req.Height == 0 {
		req.Height = 1328
	}

	// 构建任务参数
	params := map[string]interface{}{
		"image_input1": req.ImageInput1,
		"template_id":  req.TemplateId,
		"width":        req.Width,
		"height":       req.Height,
	}

	// 创建任务
	taskReq := &jimeng.CreateTaskRequest{
		Type:   model.JimengJobTypeImageEffects,
		Prompt: "",
		Params: params,
		ReqKey: model.ReqKeyImageEffects,
		Power:  15,
	}

	job, err := h.jimengService.CreateTask(user.Id, taskReq)
	if err != nil {
		logger.Errorf("create jimeng image effects task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	// 扣除用户算力
	h.subUserPower(user.Id, 15, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "即梦图像特效",
		Remark: fmt.Sprintf("任务ID：%d", job.Id),
	})

	resp.SUCCESS(c, job)
}

// TextToVideo 文生视频
func (h *JimengHandler) TextToVideo(c *gin.Context) {
	var req struct {
		Prompt      string `json:"prompt" binding:"required"`
		Seed        int64  `json:"seed"`
		AspectRatio string `json:"aspect_ratio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 检查用户算力
	if user.Power < 100 { // 文生视频消耗100算力
		resp.ERROR(c, "算力不足")
		return
	}

	// 设置默认参数
	if req.Seed == 0 {
		req.Seed = -1
	}
	if req.AspectRatio == "" {
		req.AspectRatio = jimeng.AspectRatio16_9
	}

	// 构建任务参数
	params := map[string]interface{}{
		"seed":         req.Seed,
		"aspect_ratio": req.AspectRatio,
	}

	// 创建任务
	taskReq := &jimeng.CreateTaskRequest{
		Type:   model.JimengJobTypeTextToVideo,
		Prompt: req.Prompt,
		Params: params,
		ReqKey: model.ReqKeyTextToVideo,
		Power:  100,
	}

	job, err := h.jimengService.CreateTask(user.Id, taskReq)
	if err != nil {
		logger.Errorf("create jimeng text to video task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	// 扣除用户算力
	h.subUserPower(user.Id, 100, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "即梦文生视频",
		Remark: fmt.Sprintf("任务ID：%d", job.Id),
	})

	resp.SUCCESS(c, job)
}

// ImageToVideo 图生视频
func (h *JimengHandler) ImageToVideo(c *gin.Context) {
	var req struct {
		ImageUrls        []string `json:"image_urls"`
		BinaryDataBase64 []string `json:"binary_data_base64"`
		Prompt           string   `json:"prompt"`
		Seed             int64    `json:"seed"`
		AspectRatio      string   `json:"aspect_ratio" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误: "+err.Error())
		return
	}

	if len(req.ImageUrls) == 0 && len(req.BinaryDataBase64) == 0 {
		resp.ERROR(c, "请提供图片URL或Base64数据")
		return
	}

	// 获取当前用户
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 检查用户算力
	if user.Power < 120 { // 图生视频消耗120算力
		resp.ERROR(c, "算力不足")
		return
	}

	// 设置默认参数
	if req.Seed == 0 {
		req.Seed = -1
	}

	// 构建任务参数
	params := map[string]interface{}{
		"seed":         req.Seed,
		"aspect_ratio": req.AspectRatio,
	}
	if len(req.ImageUrls) > 0 {
		params["image_urls"] = req.ImageUrls
	}
	if len(req.BinaryDataBase64) > 0 {
		params["binary_data_base64"] = req.BinaryDataBase64
	}

	// 创建任务
	taskReq := &jimeng.CreateTaskRequest{
		Type:   model.JimengJobTypeImageToVideo,
		Prompt: req.Prompt,
		Params: params,
		ReqKey: model.ReqKeyImageToVideo,
		Power:  120,
	}

	job, err := h.jimengService.CreateTask(user.Id, taskReq)
	if err != nil {
		logger.Errorf("create jimeng image to video task failed: %v", err)
		resp.ERROR(c, "创建任务失败")
		return
	}

	// 扣除用户算力
	h.subUserPower(user.Id, 120, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  "即梦图生视频",
		Remark: fmt.Sprintf("任务ID：%d", job.Id),
	})

	resp.SUCCESS(c, job)
}

// Jobs 获取任务列表
func (h *JimengHandler) Jobs(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)

	jobs, total, err := h.jimengService.GetUserJobs(user.Id, page, pageSize)
	if err != nil {
		logger.Errorf("get user jimeng jobs failed: %v", err)
		resp.ERROR(c, "获取任务列表失败")
		return
	}

	resp.SUCCESS(c, gin.H{
		"jobs":      jobs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// PendingCount 获取未完成任务数量
func (h *JimengHandler) PendingCount(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	count, err := h.jimengService.GetPendingTaskCount(user.Id)
	if err != nil {
		logger.Errorf("get pending task count failed: %v", err)
		resp.ERROR(c, "获取待处理任务数量失败")
		return
	}

	resp.SUCCESS(c, gin.H{"count": count})
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

	if err := h.jimengService.DeleteJob(uint(jobId), user.Id); err != nil {
		logger.Errorf("delete jimeng job failed: %v", err)
		resp.ERROR(c, "删除任务失败")
		return
	}

	resp.SUCCESS(c, gin.H{})
}

// Retry 重试任务
func (h *JimengHandler) Retry(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	jobIdStr := c.Param("id")
	jobId, err := strconv.ParseUint(jobIdStr, 10, 32)
	if err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	// 检查任务是否存在且属于当前用户
	job, err := h.jimengService.GetJob(uint(jobId))
	if err != nil {
		resp.ERROR(c, "任务不存在")
		return
	}
	
	if job.UserId != user.Id {
		resp.ERROR(c, "无权限操作")
		return
	}

	// 只有失败的任务才能重试
	if job.Status != model.JimengJobStatusFailed {
		resp.ERROR(c, "只有失败的任务才能重试")
		return
	}

	// 重置任务状态
	if err := h.jimengService.UpdateJobStatus(uint(jobId), model.JimengJobStatusPending, ""); err != nil {
		logger.Errorf("reset job status failed: %v", err)
		resp.ERROR(c, "重置任务状态失败")
		return
	}

	// 重新推送到队列
	task := map[string]interface{}{
		"job_id": jobId,
		"type":   job.Type,
	}
	if err := h.jimengService.PushTaskToQueue(task); err != nil {
		logger.Errorf("push retry task to queue failed: %v", err)
		resp.ERROR(c, "推送重试任务失败")
		return
	}

	resp.SUCCESS(c, gin.H{"message": "重试任务已提交"})
}

// subUserPower 扣除用户算力
func (h *JimengHandler) subUserPower(userId uint, power int, powerLog model.PowerLog) {
	session := h.DB.Session(&gorm.Session{})

	// 更新用户算力
	if err := session.Model(&model.User{}).Where("id = ?", userId).UpdateColumn("power", gorm.Expr("power - ?", power)).Error; err != nil {
		logger.Errorf("update user power failed: %v", err)
		return
	}

	// 记录算力消费日志
	powerLog.UserId = userId
	powerLog.Amount = power
	powerLog.Mark = types.PowerSub
	powerLog.CreatedAt = time.Now()
	if err := session.Create(&powerLog).Error; err != nil {
		logger.Errorf("create power log failed: %v", err)
		return
	}

	session.Commit()
}
