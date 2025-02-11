package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	logger2 "geekai/logger"
	"geekai/store/model"
	"geekai/utils"
	"gorm.io/gorm"
	"strings"

	"github.com/gin-gonic/gin"
)

var logger = logger2.GetLogger()

type BaseHandler struct {
	App *core.AppServer
	DB  *gorm.DB
}

func (h *BaseHandler) GetTrim(c *gin.Context, key string) string {
	return strings.TrimSpace(c.Query(key))
}

func (h *BaseHandler) PostInt(c *gin.Context, key string, defaultValue int) int {
	return utils.IntValue(c.PostForm(key), defaultValue)
}

func (h *BaseHandler) GetInt(c *gin.Context, key string, defaultValue int) int {
	return utils.IntValue(c.Query(key), defaultValue)
}

func (h *BaseHandler) GetFloat(c *gin.Context, key string) float64 {
	return utils.FloatValue(c.Query(key))
}
func (h *BaseHandler) PostFloat(c *gin.Context, key string) float64 {
	return utils.FloatValue(c.PostForm(key))
}

func (h *BaseHandler) GetBool(c *gin.Context, key string) bool {
	return utils.BoolValue(c.Query(key))
}
func (h *BaseHandler) PostBool(c *gin.Context, key string) bool {
	return utils.BoolValue(c.PostForm(key))
}
func (h *BaseHandler) GetUserKey(c *gin.Context) string {
	userId, ok := c.Get(types.LoginUserID)
	if !ok {
		return ""
	}
	return fmt.Sprintf("users/%v", userId)
}

func (h *BaseHandler) GetLoginUserId(c *gin.Context) uint {
	userId, ok := c.Get(types.LoginUserID)
	if !ok {
		return 0
	}
	return uint(utils.IntValue(utils.InterfaceToString(userId), 0))
}

func (h *BaseHandler) IsLogin(c *gin.Context) bool {
	return h.GetLoginUserId(c) > 0
}

func (h *BaseHandler) GetLoginUser(c *gin.Context) (model.User, error) {
	value, exists := c.Get(types.LoginUserCache)
	if exists {
		return value.(model.User), nil
	}

	userId, ok := c.Get(types.LoginUserID)
	if !ok {
		return model.User{}, errors.New("user not login")
	}

	var user model.User
	res := h.DB.Where("id", userId).First(&user)
	// 更新缓存
	if res.Error == nil {
		c.Set(types.LoginUserCache, user)
	}
	return user, res.Error
}
