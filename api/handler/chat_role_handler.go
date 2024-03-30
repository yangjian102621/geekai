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
)

type ChatRoleHandler struct {
	BaseHandler
}

func NewChatRoleHandler(app *core.AppServer, db *gorm.DB) *ChatRoleHandler {
	return &ChatRoleHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// List 获取用户聊天应用列表
func (h *ChatRoleHandler) List(c *gin.Context) {
	all := h.GetBool(c, "all")
	userId := h.GetLoginUserId(c)
	var roles []model.ChatRole
	var roleVos = make([]vo.ChatRole, 0)
	res := h.DB.Where("enable", true).Order("sort_num ASC").Find(&roles)
	if res.Error != nil {
		resp.SUCCESS(c, roleVos)
		return
	}

	// 获取所有角色
	if userId == 0 || all {
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

	var user model.User
	h.DB.First(&user, userId)
	var roleKeys []string
	err := utils.JsonDecode(user.ChatRoles, &roleKeys)
	if err != nil {
		resp.ERROR(c, "角色解析失败！")
		return
	}

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

// UpdateRole 更新用户聊天角色
func (h *ChatRoleHandler) UpdateRole(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	var data struct {
		Keys []string `json:"keys"`
	}
	if err = c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.User{}).Where("id = ?", user.Id).UpdateColumn("chat_roles_json", utils.JsonEncode(data.Keys))
	if res.Error != nil {
		logger.Error("添加应用失败：", err)
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	resp.SUCCESS(c)
}
