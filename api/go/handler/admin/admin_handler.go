package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	logger2 "chatplus/logger"
	"chatplus/utils"
	"chatplus/utils/resp"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var logger = logger2.GetLogger()

type ManagerHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewAdminHandler(app *core.AppServer, db *gorm.DB) *ManagerHandler {
	h := ManagerHandler{db: db}
	h.App = app
	return &h
}

// Login 登录
func (h *ManagerHandler) Login(c *gin.Context) {
	var data types.Manager
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	manager := h.App.AppConfig.Manager
	if data.Username == manager.Username && data.Password == manager.Password {
		err := utils.SetLoginAdmin(c, manager)
		if err != nil {
			resp.ERROR(c, "Save session failed")
			return
		}
		manager.Password = "" // 清空密码]
		resp.SUCCESS(c, manager)
	} else {
		resp.ERROR(c, "用户名或者密码错误")
	}
}

// Logout 注销
func (h *ManagerHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(types.SessionAdmin)
	err := session.Save()
	if err != nil {
		resp.ERROR(c, "Save session failed")
	} else {
		resp.SUCCESS(c)
	}
}

// Session 会话检测
func (h *ManagerHandler) Session(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get(types.SessionAdmin)
	if admin == nil {
		resp.NotAuth(c)
	} else {
		resp.SUCCESS(c)
	}
}
