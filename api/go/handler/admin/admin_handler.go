package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	logger2 "chatplus/logger"
	"chatplus/store/model"
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

// Migrate 数据修正
func (h *ManagerHandler) Migrate(c *gin.Context) {
	opt := c.Query("opt")
	switch opt {
	case "user":
		// 将用户订阅角色的数据结构从 map 改成数组
		var users []model.User
		h.db.Find(&users)
		for _, u := range users {
			var m map[string]int
			var roleKeys = make([]string, 0)
			err := utils.JsonDecode(u.ChatRoles, &m)
			if err != nil {
				continue
			}

			for k, _ := range m {
				roleKeys = append(roleKeys, k)
			}
			u.ChatRoles = utils.JsonEncode(roleKeys)
			h.db.Updates(&u)

		}
		break
	case "role":
		// 修改角色图片，改成绝对路径
		var roles []model.ChatRole
		h.db.Find(&roles)
		for _, r := range roles {
			r.Icon = "/" + r.Icon
			h.db.Updates(&r)
		}
		break
	case "history":
		// 修改角色图片，改成绝对路径
		var message []model.HistoryMessage
		h.db.Find(&message)
		for _, r := range message {
			r.Icon = "/" + r.Icon
			h.db.Updates(&r)
		}
		break
	}

	resp.SUCCESS(c, "SUCCESS")
}
