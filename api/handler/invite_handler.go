package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// InviteHandler 用户邀请
type InviteHandler struct {
	BaseHandler
}

func NewInviteHandler(app *core.AppServer, db *gorm.DB) *InviteHandler {
	return &InviteHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// Code 获取当前用户邀请码
func (h *InviteHandler) Code(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	var inviteCode model.InviteCode
	res := h.DB.Where("user_id = ?", userId).First(&inviteCode)
	// 如果邀请码不存在，则创建一个
	if res.Error != nil {
		code := strings.ToUpper(utils.RandString(8))
		for {
			res = h.DB.Where("code = ?", code).First(&inviteCode)
			if res.Error != nil { // 不存在相同的邀请码则退出
				break
			}
		}
		inviteCode.UserId = userId
		inviteCode.Code = code
		h.DB.Create(&inviteCode)
	}

	var codeVo vo.InviteCode
	err := utils.CopyObject(inviteCode, &codeVo)
	if err != nil {
		resp.ERROR(c, "拷贝对象失败")
		return
	}

	resp.SUCCESS(c, codeVo)
}

// List Log 用户邀请记录
func (h *InviteHandler) List(c *gin.Context) {

	var data struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	userId := h.GetLoginUserId(c)
	session := h.DB.Session(&gorm.Session{}).Where("inviter_id = ?", userId)
	var total int64
	session.Model(&model.InviteLog{}).Count(&total)
	var items []model.InviteLog
	var list = make([]vo.InviteLog, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var v vo.InviteLog
			err := utils.CopyObject(item, &v)
			if err == nil {
				v.Id = item.Id
				v.CreatedAt = item.CreatedAt.Unix()
				list = append(list, v)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}

// Hits 访问邀请码
func (h *InviteHandler) Hits(c *gin.Context) {
	code := c.Query("code")
	h.DB.Model(&model.InviteCode{}).Where("code = ?", code).UpdateColumn("hits", gorm.Expr("hits + ?", 1))
	resp.SUCCESS(c)
}
