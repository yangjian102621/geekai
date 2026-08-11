package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/service/moderation"
	"geekai/service/oss"
	"geekai/service/video"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoHandler struct {
	BaseHandler
	videoService      *video.Service
	uploader          *oss.UploaderManager
	userService       *service.UserService
	moderationManager *moderation.ServiceManager
}

func NewVideoHandler(app *core.AppServer, db *gorm.DB, service *video.Service, uploader *oss.UploaderManager, userService *service.UserService, moderationManager *moderation.ServiceManager) *VideoHandler {
	return &VideoHandler{
		BaseHandler: BaseHandler{
			App: app,
			DB:  db,
		},
		videoService:      service,
		uploader:          uploader,
		userService:       userService,
		moderationManager: moderationManager,
	}
}

// RegisterRoutes 注册路由
func (h *VideoHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/video/")

	// 需要用户授权的接口
	group.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		group.POST("create", h.Create)
		group.GET("list", h.List)
		group.GET("remove", h.Remove)
		group.GET("publish", h.Publish)
		group.GET("power-config", h.GetPowerConfig)     // 获取算力配置
		group.GET("power-by-key", h.GetPowerByPriceKey) // 根据 priceKey 获取算力
	}
}

type VideoTaskRequest struct {
	Provider string         `json:"provider"`  // 服务提供商（不带版本号：veo, sora）
	Model    string         `json:"model"`     // 模型标识（带版本号：veo-2.0, sora-2.0）
	Prompt   string         `json:"prompt"`    // 提示词
	Params   map[string]any `json:"params"`    // 模型特定参数
	PriceKey string         `json:"price_key"` // 价格键（如 "fixed", "5_720P" 等）
}

// Create 统一的创建视频任务接口
func (h *VideoHandler) Create(c *gin.Context) {
	var data VideoTaskRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 验证必填字段
	if data.Provider == "" {
		resp.ERROR(c, "provider 不能为空")
		return
	}
	if data.Model == "" {
		resp.ERROR(c, "model 不能为空")
		return
	}
	if data.Prompt == "" {
		resp.ERROR(c, "prompt 不能为空")
		return
	}
	if data.PriceKey == "" {
		resp.ERROR(c, "price_key 不能为空")
		return
	}

	// 文本审查
	if h.App.SysConfig.Moderation.Enable {
		moderationResult, err := h.moderationManager.GetService().Moderate(data.Prompt)
		if err != nil {
			logger.Error("failed to moderate content: ", err)
		}
		if moderationResult.Flagged {
			// 记录违规内容
			moderation := model.Moderation{
				UserId: h.GetLoginUserId(c),
				Source: types.ModerationSourceVideo,
				Input:  data.Prompt,
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

	// 获取用户信息
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	// 计算算力
	power, err := video.CalculatePower(h.DB, data.Model, data.PriceKey)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 检查算力是否充足
	if user.Power < power {
		resp.ERROR(c, "您的算力不足，请充值后再试！")
		return
	}

	// 构建任务
	userId := int(h.GetLoginUserId(c))
	task := types.VideoTask{
		UserId:           userId,
		Type:             data.Provider, // provider 作为 type
		Prompt:           data.Prompt,
		Params:           data.Params,
		TranslateModelId: h.App.SysConfig.Base.AssistantModelId,
	}

	// 插入数据库
	job := model.VideoJob{
		UserId: uint(userId),
		Type:   data.Provider,
		Prompt: data.Prompt,
		Power:  power,
		Params: utils.JsonEncode(task),
	}
	tx := h.DB.Create(&job)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	// 创建任务
	task.Id = job.Id
	h.videoService.PushTask(task)

	// 扣减算力
	err = h.userService.DecreasePower(job.UserId, job.Power, model.PowerLog{
		Type:   types.PowerConsume,
		Model:  data.Provider,
		Remark: fmt.Sprintf("%s 视频生成，任务ID：%d", data.Provider, job.Id),
	})
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{"job_id": job.Id})
}

// GetPowerConfig 获取算力配置
func (h *VideoHandler) GetPowerConfig(c *gin.Context) {
	config, err := video.GetVideoConfig(h.DB)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, config.VideoPowers)
}

// GetPowerByPriceKey 根据 modelKey 和 priceKey 获取算力值
func (h *VideoHandler) GetPowerByPriceKey(c *gin.Context) {
	modelKey := c.Query("model_key")
	priceKey := c.Query("price_key")

	if modelKey == "" {
		resp.ERROR(c, "model_key 不能为空")
		return
	}
	if priceKey == "" {
		resp.ERROR(c, "price_key 不能为空")
		return
	}

	power, err := video.CalculatePower(h.DB, modelKey, priceKey)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{"power": power})
}

func (h *VideoHandler) List(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	t := c.Query("type")
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	all := h.GetBool(c, "all")
	session := h.DB.Session(&gorm.Session{})
	if t != "" {
		session = session.Where("type", t)
	}
	if all {
		session = session.Where("publish", 0).Where("status", types.VideoStatusSuccess)
	} else {
		session = session.Where("user_id", userId)
	}
	// 统计总数
	var total int64
	session.Model(&model.VideoJob{}).Count(&total)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		session = session.Offset(offset).Limit(pageSize)
	}
	var list []model.VideoJob
	err := session.Order("id desc").Find(&list).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 转换为 VO
	items := make([]vo.VideoJob, 0)
	for _, v := range list {
		var item vo.VideoJob
		err = utils.CopyObject(v, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = v.CreatedAt.Unix()
		// 解析任务详情（用于前端展示标签）
		if v.Params != "" {
			task := types.VideoTask{}
			if err := utils.JsonDecode(v.Params, &task); err == nil {
				// 默认从 params map 中提取常用字段
				if paramsMap, ok := task.Params.(map[string]any); ok {
					if item.Params == nil {
						item.Params = make(map[string]any)
					}
					if _, ok := item.Params["task_type"]; !ok {
						if taskType, ok := paramsMap["task_type"]; ok {
							item.Params["task_type"] = taskType
						}
					}
					if _, ok := item.Params["model"]; !ok {
						if modelKey, ok := paramsMap["model"]; ok {
							item.Params["model"] = modelKey
						}
					}
					if _, ok := item.Params["duration"]; !ok {
						if duration, ok := paramsMap["duration"]; ok {
							item.Params["duration"] = duration
						}
					}
					if _, ok := item.Params["size"]; !ok {
						if size, ok := paramsMap["size"]; ok {
							item.Params["size"] = size
						} else if size, ok := paramsMap["resolution"].(string); ok {
							item.Params["size"] = size
						}
					}
					if _, ok := item.Params["mode"]; !ok {
						if mode, ok := paramsMap["mode"]; ok {
							item.Params["mode"] = mode
						}
					}
					if _, ok := item.Params["sound"]; !ok {
						if sound, ok := paramsMap["sound"]; ok {
							item.Params["sound"] = sound
						}
					}
				}
			}
		}

		items = append(items, item)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, items))
}

func (h *VideoHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	var job model.VideoJob
	err := h.DB.Where("id = ?", id).Where("user_id", userId).First(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 只有失败的任务才能删除
	if job.Status != types.VideoStatusFailed {
		resp.ERROR(c, "只有失败的任务才能删除！")
		return
	}

	// 删除任务
	err = h.DB.Delete(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 删除文件
	_ = h.uploader.GetUploadHandler().Delete(job.VideoURL)

	resp.SUCCESS(c)
}

func (h *VideoHandler) Publish(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	userId := h.GetLoginUserId(c)
	publish := h.GetBool(c, "publish")
	var job model.VideoJob
	err := h.DB.Where("id = ?", id).Where("user_id", userId).First(&job).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	err = h.DB.Model(&job).UpdateColumn("publish", publish).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
