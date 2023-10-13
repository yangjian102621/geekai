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
	all := h.GetBool(c, "all")
	var roles []model.ChatRole
	res := h.db.Where("enable", true).Order("sort_num ASC").Find(&roles)
	if res.Error != nil {
		resp.ERROR(c, "No roles found,"+res.Error.Error())
		return
	}

	// 获取所有角色
	if all {
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
		return
	}

	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}
	var roleKeys []string
	err = utils.JsonDecode(user.ChatRoles, &roleKeys)
	if err != nil {
		resp.ERROR(c, "角色解析失败！")
		return
	}
	// 转成 vo
	var roleVos = make([]vo.ChatRole, 0)
	for _, r := range roles {
		if !utils.ContainsStr(roleKeys, r.Key) {
			continue
		}
		var v vo.ChatRole
		err := utils.CopyObject(r, &v)
		if err == nil {
			v.Id = r.Id
			roleVos = append(roleVos, v)
		}
	}
	resp.SUCCESS(c, roleVos)
}
