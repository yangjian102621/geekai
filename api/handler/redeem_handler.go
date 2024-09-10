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
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
	"time"
)

type RedeemHandler struct {
	BaseHandler
	lock        sync.Mutex
	userService *service.UserService
}

func NewRedeemHandler(app *core.AppServer, db *gorm.DB, userService *service.UserService) *RedeemHandler {
	return &RedeemHandler{BaseHandler: BaseHandler{App: app, DB: db}, userService: userService}
}

func (h *RedeemHandler) Verify(c *gin.Context) {
	var data struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	userId := h.GetLoginUserId(c)

	h.lock.Lock()
	defer h.lock.Unlock()

	var item model.Redeem
	res := h.DB.Where("code", data.Code).First(&item)
	if res.Error != nil {
		resp.ERROR(c, "无效的兑换码！")
		return
	}

	if !item.Enabled {
		resp.ERROR(c, "当前兑换码已被禁用！")
		return
	}

	if item.RedeemedAt > 0 {
		resp.ERROR(c, "当前兑换码已使用，请勿重复使用！")
		return
	}

	tx := h.DB.Begin()
	err := h.userService.IncreasePower(int(userId), item.Power, model.PowerLog{
		Type:   types.PowerRedeem,
		Model:  "兑换码",
		Remark: fmt.Sprintf("兑换码核销，算力：%d，兑换码：%s...", item.Power, item.Code[:10]),
	})
	if err != nil {
		tx.Rollback()
		resp.ERROR(c, err.Error())
		return
	}

	// 更新核销状态
	item.RedeemedAt = time.Now().Unix()
	item.UserId = userId
	err = tx.Updates(&item).Error
	if err != nil {
		tx.Rollback()
		resp.ERROR(c, err.Error())
		return
	}

	tx.Commit()
	resp.SUCCESS(c)

}
