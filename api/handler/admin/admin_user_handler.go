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

type SysUserHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewSysUserHandler(app *core.AppServer, db *gorm.DB) *SysUserHandler {
	h := SysUserHandler{db: db}
	h.App = app
	return &h
}

// List 用户列表
func (h *SysUserHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	username := h.GetTrim(c, "username")

	offset := (page - 1) * pageSize
	var items []model.AdminUser
	var users = make([]vo.AdminUser, 0)
	var total int64

	session := h.db.Session(&gorm.Session{})
	if username != "" {
		session = session.Where("username LIKE ?", "%"+username+"%")
	}

	// 查询total
	session.Model(&model.AdminUser{}).Count(&total)
	res := session.Offset(offset).Limit(pageSize).Find(&items)

	if res.Error == nil {
		for _, item := range items {
			var userVo vo.AdminUser
			err := utils.CopyObject(item, &userVo)
			if err == nil {
				userVo.Id = item.Id
				userVo.CreatedAt = item.CreatedAt.Unix()
				userVo.UpdatedAt = item.UpdatedAt.Unix()
				users = append(users, userVo)
			} else {
				logger.Error(err)
			}
		}
	}
	pageVo := vo.NewPage(total, page, pageSize, users)
	resp.SUCCESS(c, pageVo)
}

// Save 更新或者新增
func (h *SysUserHandler) Save(c *gin.Context) {
	var data struct {
		Id       uint   `json:"id"`
		Password string `json:"password"`
		Username string `json:"username"`
		Status   bool   `json:"status"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 默认id为1是超级管理员
	if data.Id == 1 {
		resp.ERROR(c, "超级管理员不支持更新")
		return
	}

	var user = model.AdminUser{}
	var res *gorm.DB
	var userVo vo.AdminUser
	if data.Id > 0 { // 更新
		user.Id = data.Id
		// 此处需要用 map 更新，用结构体无法更新 0 值
		res = h.db.Model(&user).Updates(map[string]interface{}{
			"username": data.Username,
			"status":   data.Status,
		})
	} else {
		salt := utils.RandString(8)
		u := model.AdminUser{
			Username: data.Username,
			Password: utils.GenPassword(data.Password, salt),
			Salt:     salt,
			Status:   true,
		}
		res = h.db.Create(&u)
		_ = utils.CopyObject(u, &userVo)
		userVo.Id = u.Id
		userVo.CreatedAt = u.CreatedAt.Unix()
		userVo.UpdatedAt = u.UpdatedAt.Unix()
	}
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败")
		return
	}

	resp.SUCCESS(c, userVo)
}

// ResetPass 重置密码
func (h *SysUserHandler) ResetPass(c *gin.Context) {
	var data struct {
		Id       uint
		Password string
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user model.AdminUser
	res := h.db.First(&user, data.Id)
	if res.Error != nil {
		resp.ERROR(c, "No user found")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	user.Password = password
	res = h.db.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c)
	} else {
		resp.SUCCESS(c)
	}
}

// Remove 删除
func (h *SysUserHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		// 默认id为1是超级管理员
		if id == 1 {
			resp.ERROR(c, "超级管理员不能删除")
			return
		}
		res := h.db.Where("id = ?", id).Delete(&model.AdminUser{})
		if res.Error != nil {
			resp.ERROR(c, "删除失败")
			return
		}
	}
	resp.SUCCESS(c)
}
