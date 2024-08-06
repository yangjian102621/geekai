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
	"geekai/core/types"
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatRoleHandler struct {
	handler.BaseHandler
}

func NewChatRoleHandler(app *core.AppServer, db *gorm.DB) *ChatRoleHandler {
	return &ChatRoleHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// Save 创建或者更新某个角色
func (h *ChatRoleHandler) Save(c *gin.Context) {
	var data vo.ChatRole
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var role model.ChatRole
	err := utils.CopyObject(data, &role)
	if err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	role.Id = data.Id
	if data.CreatedAt > 0 {
		role.CreatedAt = time.Unix(data.CreatedAt, 0)
	} else {
		err = h.DB.Where("marker", data.Key).First(&role).Error
		if err == nil {
			resp.ERROR(c, fmt.Sprintf("角色 %s 已存在", data.Key))
			return
		}
	}
	err = h.DB.Save(&role).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 填充 ID 数据
	data.Id = role.Id
	data.CreatedAt = role.CreatedAt.Unix()
	resp.SUCCESS(c, data)
}

func (h *ChatRoleHandler) List(c *gin.Context) {
	var items []model.ChatRole
	var roles = make([]vo.ChatRole, 0)
	res := h.DB.Order("sort_num ASC").Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "No data found")
		return
	}

	// initialize model mane for role
	modelIds := make([]int, 0)
	for _, v := range items {
		if v.ModelId > 0 {
			modelIds = append(modelIds, v.ModelId)
		}
	}

	modelNameMap := make(map[int]string)
	if len(modelIds) > 0 {
		var models []model.ChatModel
		tx := h.DB.Where("id IN ?", modelIds).Find(&models)
		if tx.Error == nil {
			for _, m := range models {
				modelNameMap[int(m.Id)] = m.Name
			}
		}
	}

	for _, v := range items {
		var role vo.ChatRole
		err := utils.CopyObject(v, &role)
		if err == nil {
			role.Id = v.Id
			role.CreatedAt = v.CreatedAt.Unix()
			role.UpdatedAt = v.UpdatedAt.Unix()
			role.ModelName = modelNameMap[role.ModelId]
			roles = append(roles, role)
		}
	}

	resp.SUCCESS(c, roles)
}

// Sort 更新角色排序
func (h *ChatRoleHandler) Sort(c *gin.Context) {
	var data struct {
		Ids   []uint `json:"ids"`
		Sorts []int  `json:"sorts"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	for index, id := range data.Ids {
		err := h.DB.Model(&model.ChatRole{}).Where("id = ?", id).Update("sort_num", data.Sorts[index]).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	resp.SUCCESS(c)
}

func (h *ChatRoleHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Model(&model.ChatRole{}).Where("id = ?", data.Id).Update(data.Filed, data.Value).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *ChatRoleHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	res := h.DB.Where("id", id).Delete(&model.ChatRole{})
	if res.Error != nil {
		logger.Error("error with update database：", res.Error)
		resp.ERROR(c, "删除失败！")
		return
	}
	resp.SUCCESS(c)
}
