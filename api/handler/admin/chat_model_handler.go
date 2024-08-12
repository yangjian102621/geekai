package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
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

func (h *ChatModelHandler) Save(c *gin.Context) {
	var data struct {
		Id          uint    `json:"id"`
		Name        string  `json:"name"`
		Value       string  `json:"value"`
		Enabled     bool    `json:"enabled"`
		SortNum     int     `json:"sort_num"`
		Open        bool    `json:"open"`
		Platform    string  `json:"platform"`
		Power       int     `json:"power"`
		MaxTokens   int     `json:"max_tokens"`  // 最大响应长度
		MaxContext  int     `json:"max_context"` // 最大上下文长度
		Temperature float32 `json:"temperature"` // 模型温度
		KeyId       int     `json:"key_id,omitempty"`
		CreatedAt   int64   `json:"created_at"`
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
	item.SortNum = data.SortNum
	item.Open = data.Open
	item.Power = data.Power
	item.MaxTokens = data.MaxTokens
	item.MaxContext = data.MaxContext
	item.Temperature = data.Temperature
	item.KeyId = data.KeyId

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
	if enable {
		session = session.Where("enabled", enable)
	}
	if name != "" {
		session = session.Where("name LIKE ?", name+"%")
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
		keyIds = append(keyIds, v.KeyId)
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
