package handler

import (
	"chatplus/core"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatRoleHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewChatRoleHandler(app *core.AppServer, db *gorm.DB) *ChatRoleHandler {
	handler := &ChatRoleHandler{db: db}
	handler.App = app
	return handler
}

// List get user list
func (h *ChatRoleHandler) List(c *gin.Context) {
	var roles []model.ChatRole
	res := h.db.Find(&roles)
	if res.Error != nil {
		resp.ERROR(c, "No roles found,"+res.Error.Error())
		return
	}
	userId := h.GetInt(c, "user_id", 0)
	if userId > 0 {
		var user model.User
		h.db.First(&user, userId)
		var roleKeys []string
		err := utils.JsonDecode(user.ChatRoles, &roleKeys)
		if err == nil {
			for index, r := range roles {
				if utils.ContainsStr(roleKeys, r.Key) {
					roles = append(roles[:index], roles[index+1:]...)
				}
			}
		}
	}
	// 转成 vo
	var roleVos = make([]vo.ChatRole, 0)
	for _, r := range roles {
		var v vo.ChatRole
		err := utils.CopyObject(r, &v)
		if err == nil {
			v.Id = r.Id
			roleVos = append(roleVos, v)
		}
	}
	resp.SUCCESS(c, roleVos)
}
