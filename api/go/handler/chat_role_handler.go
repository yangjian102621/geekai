package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatRoleHandler struct {
	BaseHandler
	service *service.ChatRoleService
}

func NewChatRoleHandler(app *core.AppServer, service *service.ChatRoleService) *ChatRoleHandler {
	handler := &ChatRoleHandler{service: service}
	handler.App = app
	return handler
}

// List get user list
func (h *ChatRoleHandler) List(c *gin.Context) {
	var roles []model.ChatRole
	res := h.service.DB.Find(&roles)
	if res.Error != nil {
		resp.ERROR(c, "No roles found,"+res.Error.Error())
		return
	}
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err == nil && userId > 0 {
		var user model.User
		h.service.DB.First(&user, userId)
		var roleMap map[string]int
		err := utils.JsonDecode(user.ChatRoles, &roleMap)
		if err == nil {
			for index, r := range roles {
				if _, ok := roleMap[r.Key]; !ok {
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

// Add 添加一个聊天角色
func (h *ChatRoleHandler) Add(c *gin.Context) {
	var data vo.ChatRole
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Key == "" || data.Name == "" || data.Icon == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.service.Create(data)
	if err != nil {
		resp.ERROR(c, "Save failed: "+err.Error())
		return
	}

	resp.SUCCESS(c, data)
}

// Get 获取指定的角色
func (h *ChatRoleHandler) Get(c *gin.Context) {

}

// Update 更新某个聊天角色信息，这里只允许更改名称以及启用和禁用角色操作
func (h *ChatRoleHandler) Update(c *gin.Context) {

}
