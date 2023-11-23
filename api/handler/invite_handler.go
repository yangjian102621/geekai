package handler

import (
	"chatplus/core"
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
	db *gorm.DB
}

func NewInviteHandler(app *core.AppServer, db *gorm.DB) *InviteHandler {
	h := InviteHandler{db: db}
	h.App = app
	return &h
}

// Code 获取当前用户邀请码
func (h *InviteHandler) Code(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	var inviteCode model.InviteCode
	res := h.db.Where("user_id = ?", userId).First(&inviteCode)
	// 如果邀请码不存在，则创建一个
	if res.Error != nil {
		code := strings.ToUpper(utils.RandString(8))
		for {
			res = h.db.Where("code = ?", code).First(&inviteCode)
			if res.Error != nil { // 不存在相同的邀请码则退出
				break
			}
		}
		inviteCode.UserId = userId
		inviteCode.Code = code
		h.db.Create(&inviteCode)
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

	resp.SUCCESS(c)
}

// Hits 访问邀请码
func (h *InviteHandler) Hits(c *gin.Context) {
	code := c.Query("code")
	h.db.Model(&model.InviteCode{}).Where("code = ?", code).UpdateColumn("hits", gorm.Expr("hits + ?", 1))
	resp.SUCCESS(c)
}
