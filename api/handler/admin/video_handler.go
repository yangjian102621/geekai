package admin

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
	"geekai/handler"
	"geekai/service"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VideoHandler 管理后台视频生成处理器
type VideoHandler struct {
	handler.BaseHandler
	userService *service.UserService
	uploader    *oss.UploaderManager
}

// NewVideoHandler 创建管理后台视频生成处理器
func NewVideoHandler(app *core.AppServer, db *gorm.DB, userService *service.UserService, manager *oss.UploaderManager) *VideoHandler {
	return &VideoHandler{
		BaseHandler: handler.BaseHandler{App: app, DB: db},
		userService: userService,
		uploader:    manager,
	}
}

// RegisterRoutes 注册视频生成管理后台路由
func (h *VideoHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/video/")
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("config", h.GetConfig)
		rg.POST("config/update", h.UpdateConfig)
		rg.POST("list", h.Videos)
		rg.GET("remove", h.Remove)
	}
}

// GetConfig 获取视频生成配置
func (h *VideoHandler) GetConfig(c *gin.Context) {
	var config model.Config
	err := h.DB.Where("name", types.ConfigKeyVideo).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回空配置
			resp.SUCCESS(c, types.VideoConfig{
				ApiURL:      "",
				ApiKey:      "",
				VideoPowers: make(map[string]types.VideoModelPower),
			})
			return
		}
		resp.ERROR(c, "获取配置失败: "+err.Error())
		return
	}

	var videoConfig types.VideoConfig
	err = utils.JsonDecode(config.Value, &videoConfig)
	if err != nil {
		resp.ERROR(c, "解析配置失败: "+err.Error())
		return
	}

	resp.SUCCESS(c, videoConfig)
}

// UpdateConfig 更新视频生成配置
func (h *VideoHandler) UpdateConfig(c *gin.Context) {
	var req types.VideoConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误")
		return
	}

	// 验证必填字段
	if req.ApiURL == "" {
		resp.ERROR(c, "API地址不能为空")
		return
	}
	if req.ApiKey == "" {
		resp.ERROR(c, "API密钥不能为空")
		return
	}

	// 验证算力配置
	if len(req.VideoPowers) == 0 {
		resp.ERROR(c, "请至少配置一个模型的算力")
		return
	}

	// 新的价格配置方式直接使用 power_config 中的 key（如 "fixed"、"5_720P" 等）
	// 不再区分固定收费和按秒收费，所有价格配置都在 power_config 中
	for key, modelPower := range req.VideoPowers {
		// 验证 provider
		if modelPower.Provider == "" {
			resp.ERROR(c, fmt.Sprintf("模型 %s 的 provider 不能为空", key))
			return
		}

		// 验证 model
		if modelPower.Model == "" {
			resp.ERROR(c, fmt.Sprintf("模型 %s 的 model 不能为空", key))
			return
		}

		// 验证 power_config
		if len(modelPower.PowerConfig) == 0 {
			resp.ERROR(c, fmt.Sprintf("模型 %s 的 power_config 不能为空", key))
			return
		}

		// 验证 power_config 中的值必须大于0
		for configKey, configValue := range modelPower.PowerConfig {
			if configValue <= 0 {
				resp.ERROR(c, fmt.Sprintf("模型 %s 的 power_config.%s 必须大于0", key, configKey))
				return
			}
		}

	}

	// 保存配置
	tx := h.DB.Begin()
	value := utils.JsonEncode(&req)
	var exist model.Config
	tx.Where("name", types.ConfigKeyVideo).First(&exist)

	if exist.Id > 0 {
		exist.Value = value
		err := tx.Updates(&exist).Error
		if err != nil {
			resp.ERROR(c, "更新配置失败: "+err.Error())
			tx.Rollback()
			return
		}
	} else {
		exist.Name = types.ConfigKeyVideo
		exist.Value = value
		err := tx.Create(&exist).Error
		if err != nil {
			resp.ERROR(c, "创建配置失败: "+err.Error())
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	resp.SUCCESS(c, gin.H{"message": "配置更新成功"})
}

type videoQuery struct {
	Type      string   `json:"type"`   // 任务类型 luma, keling
	Status    string   `json:"status"` // 任务状态 pending, in_progress, downloading, success, failed
	Prompt    string   `json:"prompt"`
	CreatedAt []string `json:"created_at"`
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
}

// Videos 视频任务列表
func (h *VideoHandler) Videos(c *gin.Context) {
	var data videoQuery
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.Type != "" {
		session = session.Where("type", data.Type)
	}
	if data.Status != "" {
		session = session.Where("status", data.Status)
	}
	if data.Prompt != "" {
		session = session.Where("prompt LIKE ?", "%"+data.Prompt+"%")
	}
	if len(data.CreatedAt) == 2 {
		session = session.Where("created_at >= ? AND created_at <= ?", data.CreatedAt[0], data.CreatedAt[1])
	}
	var total int64
	session.Model(&model.VideoJob{}).Count(&total)
	var list []model.VideoJob
	var items = make([]vo.VideoJob, 0)
	offset := (data.Page - 1) * data.PageSize
	err := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&list).Error
	if err == nil {
		// 填充数据
		for _, item := range list {
			var job vo.VideoJob
			err = utils.CopyObject(item, &job)
			if err != nil {
				continue
			}
			job.CreatedAt = item.CreatedAt.Unix()
			items = append(items, job)
		}
	}

	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}

func (h *VideoHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	tab := c.Query("tab")

	tx := h.DB.Begin()
	var md, remark, fileURL string
	var power, userId int
	var needRefund bool

	switch tab {
	case "luma", "keling":
		var job model.VideoJob
		if res := h.DB.Where("id", id).First(&job); res.Error != nil {
			resp.ERROR(c, "记录不存在")
			return
		}

		// 删除任务
		tx.Delete(&job)
		md = job.Type
		power = job.Power
		userId = int(job.UserId)
		remark = fmt.Sprintf("视频任务失败，退回算力。任务ID：%d，Err: %s", job.Id, job.ErrMsg)
		needRefund = job.Status != types.VideoStatusSuccess
		fileURL = job.VideoURL
	default:
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if needRefund {
		err := h.userService.IncreasePower(uint(userId), power, model.PowerLog{
			Type:   types.PowerRefund,
			Model:  md,
			Remark: remark,
		})
		if err != nil {
			tx.Rollback()
			resp.ERROR(c, err.Error())
			return
		}
	}
	tx.Commit()
	// remove file
	err := h.uploader.GetUploadHandler().Delete(fileURL)
	if err != nil {
		logger.Error("remove file failed: ", err)
	}

	resp.SUCCESS(c)
}
