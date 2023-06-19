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

type UserHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewUserHandler(app *core.AppServer, db *gorm.DB) *UserHandler {
	h := UserHandler{db: db}
	h.App = app
	return &h
}

// List 用户列表
func (h *UserHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	offset := (page - 1) * pageSize
	var items []model.User
	var users = make([]vo.User, 0)
	var total int64
	h.db.Model(&model.User{}).Count(&total)
	res := h.db.Offset(offset).Limit(pageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var user vo.User
			err := utils.CopyObject(item, &user)
			if err == nil {
				user.Id = item.Id
				user.CreatedAt = item.CreatedAt.Unix()
				user.UpdatedAt = item.UpdatedAt.Unix()
				users = append(users, user)
			} else {
				logger.Error(err)
			}
		}
	}
	pageVo := vo.NewPage(total, page, pageSize, users)
	resp.SUCCESS(c, pageVo)
}

func (h *UserHandler) Update(c *gin.Context) {
	var data struct {
		Id          uint     `json:"id"`
		Nickname    string   `json:"nickname"`
		Calls       int      `json:"calls"`
		ChatRoles   []string `json:"chat_roles"`
		ExpiredTime string   `json:"expired_time"`
		Status      bool     `json:"status"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	var user = model.User{
		Nickname:    data.Nickname,
		Calls:       data.Calls,
		Status:      data.Status,
		ChatRoles:   utils.JsonEncode(data.ChatRoles),
		ExpiredTime: utils.Str2stamp(data.ExpiredTime),
	}
	user.Id = data.Id
	res := h.db.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c)
}

func (h *UserHandler) InitUser(c *gin.Context) {
	var users []model.User
	h.db.Find(&users)
	for _, u := range users {
		var m map[string]int
		var roleKeys = make([]string, 0)
		utils.JsonDecode(u.ChatRoles, &m)
		for k, _ := range m {
			roleKeys = append(roleKeys, k)
		}
		u.ChatRoles = utils.JsonEncode(roleKeys)
		h.db.Updates(&u)

	}
	resp.SUCCESS(c, "SUCCESS")
}
