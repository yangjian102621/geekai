package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatRoleHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewChatRoleHandler(app *core.AppServer, db *gorm.DB) *ChatRoleHandler {
	h := ChatRoleHandler{db: db}
	h.App = app
	return &h
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

	//err := h.service.Create(data)
	//if err != nil {
	//	resp.ERROR(c, "Save failed: "+err.Error())
	//	return
	//}

	resp.SUCCESS(c, data)
}

// Get 获取指定的角色
func (h *ChatRoleHandler) Get(c *gin.Context) {

}

// Update 更新某个聊天角色信息，这里只允许更改名称以及启用和禁用角色操作
func (h *ChatRoleHandler) Update(c *gin.Context) {

}

func (h *ChatRoleHandler) List(c *gin.Context) {
	var items []model.ChatRole
	var roles = make([]vo.ChatRole, 0)
	res := h.db.Where("enable", true).Order("sort ASC").Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "No data found")
		return
	}

	for _, v := range items {
		var role vo.ChatRole
		err := utils.CopyObject(v, &role)
		if err == nil {
			role.Id = v.Id
			role.CreatedAt = v.CreatedAt.Unix()
			role.UpdatedAt = v.UpdatedAt.Unix()
			roles = append(roles, role)
		}
	}

	resp.SUCCESS(c, roles)
}
