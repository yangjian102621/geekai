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

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FunctionHandler struct {
	handler.BaseHandler
}

func NewFunctionHandler(app *core.AppServer, db *gorm.DB) *FunctionHandler {
	return &FunctionHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

func (h *FunctionHandler) Save(c *gin.Context) {
	var data vo.Function
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var f = model.Function{
		Id:          data.Id,
		Name:        data.Name,
		Label:       data.Label,
		Description: data.Description,
		Parameters:  utils.JsonEncode(data.Parameters),
		Action:      data.Action,
		Token:       data.Token,
		Enabled:     data.Enabled,
	}

	res := h.DB.Save(&f)
	if res.Error != nil {
		resp.ERROR(c, "error with save data:"+res.Error.Error())
		return
	}
	data.Id = f.Id
	resp.SUCCESS(c, data)
}

func (h *FunctionHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Model(&model.Function{}).Where("id = ?", data.Id).Update(data.Filed, data.Value).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *FunctionHandler) List(c *gin.Context) {
	var items []model.Function
	res := h.DB.Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "No data found")
		return
	}

	functions := make([]vo.Function, 0)
	for _, v := range items {
		var f vo.Function
		err := utils.CopyObject(v, &f)
		if err != nil {
			continue
		}
		functions = append(functions, f)
	}
	resp.SUCCESS(c, functions)
}

func (h *FunctionHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		err := h.DB.Delete(&model.Function{Id: uint(id)}).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	resp.SUCCESS(c)
}

// GenToken generate function api access token
func (h *FunctionHandler) GenToken(c *gin.Context) {
	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 0,
		"expired": 0,
	})
	tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
	if err != nil {
		logger.Error("error with generate token", err)
		resp.ERROR(c)
		return
	}

	resp.SUCCESS(c, tokenString)
}
