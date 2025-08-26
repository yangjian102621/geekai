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
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatModelHandler struct {
	handler.BaseHandler
}

func NewChatModelHandler(app *core.AppServer, db *gorm.DB) *ChatModelHandler {
	return &ChatModelHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *ChatModelHandler) RegisterRoutes() {
	group := h.App.Engine.Group("/api/admin/model/")

	// 需要管理员授权的接口
	group.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		group.GET("list", h.List)
		group.POST("save", h.Save)
		group.POST("set", h.Set)
		group.POST("sort", h.Sort)
		group.GET("remove", h.Remove)
		group.POST("batch-remove", h.BatchRemove)
	}
}

func (h *ChatModelHandler) Save(c *gin.Context) {
	var data struct {
		Id          uint              `json:"id"`
		Name        string            `json:"name"`
		Value       string            `json:"value"`
		Enabled     bool              `json:"enabled"`
		SortNum     int               `json:"sort_num"`
		Open        bool              `json:"open"`
		Platform    string            `json:"platform"`
		Power       int               `json:"power"`
		MaxTokens   int               `json:"max_tokens"`  // 最大响应长度
		MaxContext  int               `json:"max_context"` // 最大上下文长度
		Desc        string            `json:"desc"`        //模型描述
		Tag         string            `json:"tag"`         //模型标签
		Temperature float32           `json:"temperature"` // 模型温度
		KeyId       int               `json:"key_id,omitempty"`
		CreatedAt   int64             `json:"created_at"`
		Type        string            `json:"type"`
		Options     map[string]string `json:"options"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	item := model.ChatModel{}
	// 更新
	if data.Id > 0 {
		h.DB.Where("id", data.Id).First(&item)
	}

	item.Name = data.Name
	item.Value = data.Value
	item.Enabled = data.Enabled
	item.Open = data.Open
	item.Power = data.Power
	item.MaxTokens = data.MaxTokens
	item.MaxContext = data.MaxContext
	item.Desc = data.Desc
	item.Tag = data.Tag
	item.Temperature = data.Temperature
	item.KeyId = uint(data.KeyId)
	item.Type = data.Type
	item.Options = utils.JsonEncode(data.Options)
	var res *gorm.DB
	if data.Id > 0 {
		res = h.DB.Save(&item)
	} else {
		res = h.DB.Create(&item)
	}
	if res.Error != nil {
		logger.Error("error with update database：", res.Error)
		resp.ERROR(c, res.Error.Error())
		return
	}

	var itemVo vo.ChatModel
	err := utils.CopyObject(item, &itemVo)
	if err != nil {
		resp.ERROR(c, "数据拷贝失败！")
		return
	}
	itemVo.Id = item.Id
	itemVo.CreatedAt = item.CreatedAt.Unix()
	resp.SUCCESS(c, itemVo)
}

// List 模型列表
func (h *ChatModelHandler) List(c *gin.Context) {
	session := h.DB.Session(&gorm.Session{})
	enable := h.GetBool(c, "enable")
	name := h.GetTrim(c, "name")
	modelType := h.GetTrim(c, "type")
	if enable {
		session = session.Where("enabled", enable)
	}
	if name != "" {
		session = session.Where("name LIKE ?", name+"%")
	}
	if modelType != "" {
		session = session.Where("type", modelType)
	}
	var items []model.ChatModel
	var cms = make([]vo.ChatModel, 0)
	res := session.Order("sort_num ASC").Find(&items)
	if res.Error != nil {
		resp.SUCCESS(c, cms)
		return
	}

	// initialize key name
	keyIds := make([]int, 0)
	for _, v := range items {
		keyIds = append(keyIds, int(v.KeyId))
	}
	var keys []model.ApiKey
	keyMap := make(map[uint]string)
	h.DB.Where("id IN ?", keyIds).Find(&keys)
	for _, v := range keys {
		keyMap[v.Id] = v.Name
	}
	for _, item := range items {
		var cm vo.ChatModel
		err := utils.CopyObject(item, &cm)
		if err == nil {
			cm.Id = item.Id
			cm.CreatedAt = item.CreatedAt.Unix()
			cm.UpdatedAt = item.UpdatedAt.Unix()
			cm.KeyName = keyMap[uint(item.KeyId)]
			cms = append(cms, cm)
		} else {
			logger.Error(err)
		}
	}
	resp.SUCCESS(c, cms)
}

func (h *ChatModelHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Model(&model.ChatModel{}).Where("id = ?", data.Id).Update(data.Filed, data.Value).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *ChatModelHandler) Sort(c *gin.Context) {
	var data struct {
		Ids   []uint `json:"ids"`
		Sorts []int  `json:"sorts"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	for index, id := range data.Ids {
		err := h.DB.Model(&model.ChatModel{}).Where("id = ?", id).Update("sort_num", data.Sorts[index]).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	resp.SUCCESS(c)
}

func (h *ChatModelHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Where("id = ?", id).Delete(&model.ChatModel{}).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

// BatchRemove 批量删除模型
func (h *ChatModelHandler) BatchRemove(c *gin.Context) {
	var data struct {
		Ids []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if len(data.Ids) == 0 {
		resp.ERROR(c, "请选择要删除的模型")
		return
	}

	// 执行批量删除
	err := h.DB.Where("id IN ?", data.Ids).Delete(&model.ChatModel{}).Error
	if err != nil {
		logger.Error("批量删除模型失败：", err)
		resp.ERROR(c, "批量删除失败："+err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"message":       fmt.Sprintf("成功删除 %d 个模型", len(data.Ids)),
		"deleted_count": len(data.Ids),
	})
}
