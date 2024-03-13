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

type role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
				var roles []role
				h.db.Raw("SELECT r.id,r.name "+
					"FROM chatgpt_admin_user_roles as ur "+
					"LEFT JOIN chatgpt_admin_roles as r ON ur.role_id = r.id "+
					"WHERE ur.admin_id = ?", item.Id).Scan(&roles)

				userVo.Id = item.Id
				userVo.CreatedAt = item.CreatedAt.Unix()
				userVo.UpdatedAt = item.UpdatedAt.Unix()
				userVo.RoleIds = roles
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
		RoleIds  []int  `json:"role_ids"`
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

	tx := h.db.Begin()

	if data.Id > 0 { // 更新
		user.Id = data.Id
		err := tx.Where("admin_id = ?", user.Id).Delete(&model.AdminUserRole{})
		if err.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "更新数据库失败")
			return
		}
		// 此处需要用 map 更新，用结构体无法更新 0 值
		res = tx.Model(&user).Updates(map[string]interface{}{
			"username": data.Username,
			"status":   data.Status,
		})
	} else {
		salt := utils.RandString(8)

		user.Username = data.Username
		user.Password = utils.GenPassword(data.Password, salt)
		user.Salt = salt
		user.Status = true

		res = tx.Create(&user)
		_ = utils.CopyObject(user, &userVo)
		userVo.Id = user.Id
		userVo.CreatedAt = user.CreatedAt.Unix()
		userVo.UpdatedAt = user.UpdatedAt.Unix()
	}
	if res.Error != nil {
		tx.Rollback()
		resp.ERROR(c, "更新数据库失败")
		return
	}
	// 添加角色
	userRole := make([]model.AdminUserRole, 0)
	if len(data.RoleIds) > 0 {
		for _, roleId := range data.RoleIds {
			userRole = append(userRole, model.AdminUserRole{
				AdminId: user.Id,
				RoleId:  roleId,
			})
		}
		err := tx.CreateInBatches(userRole, len(userRole))
		if err.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "更新数据库失败")
			return
		}
	}
	tx.Commit()

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
	var data struct {
		Id uint
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 默认id为1是超级管理员
	if data.Id == 1 {
		resp.ERROR(c, "超级管理员不能删除")
		return
	}
	if data.Id > 0 {
		tx := h.db.Begin()
		res := tx.Where("id = ?", data.Id).Delete(&model.AdminUser{})
		if res.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "删除失败")
			return
		}
		res2 := tx.Where("admin_id = ?", data.Id).Delete(&model.AdminUserRole{})
		if res2.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "删除失败")
			return
		}
		tx.Commit()
	}
	resp.SUCCESS(c)
}
