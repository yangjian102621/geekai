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
	"geekai/store/model"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
	"time"
)

type RedeemHandler struct {
	BaseHandler
	lock sync.Mutex
}

func NewRedeemHandler(app *core.AppServer, db *gorm.DB) *RedeemHandler {
	return &RedeemHandler{BaseHandler: BaseHandler{App: app, DB: db}}
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
	err := tx.Model(&model.User{}).Where("id", userId).UpdateColumn("power", gorm.Expr("power + ?", item.Power)).Error
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

	// 记录算力充值日志
	var user model.User
	err = tx.Where("id", userId).First(&user).Error
	if err != nil {
		tx.Rollback()
		resp.ERROR(c, err.Error())
		return
	}

	h.DB.Create(&model.PowerLog{
		UserId:    userId,
		Username:  user.Username,
		Type:      types.PowerRedeem,
		Amount:    item.Power,
		Balance:   user.Power,
		Mark:      types.PowerAdd,
		Model:     "兑换码",
		Remark:    fmt.Sprintf("兑换码核销，算力：%d，兑换码：%s...", item.Power, item.Code[:10]),
		CreatedAt: time.Now(),
	})
	tx.Commit()
	resp.SUCCESS(c)

}
