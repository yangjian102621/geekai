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

type SysRoleHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewSysRoleHandler(app *core.AppServer, db *gorm.DB) *SysRoleHandler {
	h := SysRoleHandler{db: db}
	h.App = app
	return &h
}

type permission struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (h *SysRoleHandler) List(c *gin.Context) {
	if err := utils.CheckPermission(c, h.db); err != nil {
		resp.ERROR(c, types.NoPermission)
		return
	}

	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	name := h.GetTrim(c, "name")

	offset := (page - 1) * pageSize
	var items []model.AdminRole
	var data = make([]vo.AdminRole, 0)
	var total int64

	session := h.db.Session(&gorm.Session{})
	if name != "" {
		session = session.Where("name LIKE ?", "%"+name+"%")
	}

	session.Model(&model.AdminRole{}).Count(&total)
	res := session.Offset(offset).Limit(pageSize).Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "暂无数据")
		return
	}
	for _, item := range items {
		adminRoleVo := vo.AdminRole{}
		err := utils.CopyObject(item, &adminRoleVo)
		if err == nil {
			var permissions []permission
			h.db.Raw("SELECT p.id,p.name,p.slug "+
				"FROM chatgpt_admin_role_permissions as rp "+
				"LEFT JOIN chatgpt_admin_permissions as p ON rp.permission_id = p.id "+
				"WHERE rp.role_id = ?", item.Id).Scan(&permissions)

			adminRoleVo.Permissions = permissions
			adminRoleVo.CreatedAt = item.CreatedAt.Format("2006-01-02 15:04:05")
			data = append(data, adminRoleVo)
		}
	}
	pageVo := vo.NewPage(total, page, pageSize, data)
	resp.SUCCESS(c, pageVo)
}

func (h *SysRoleHandler) Save(c *gin.Context) {
	var data struct {
		Id          int
		Name        string
		Description string
		Permissions []int
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var role = model.AdminRole{}
	var res *gorm.DB
	tx := h.db.Begin()
	if data.Id > 0 { // 更新
		role.Id = data.Id
		//删除角色对应的权限
		err := tx.Where("role_id = ?", role.Id).Delete(model.AdminRolePermission{})
		if err.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "更新数据库失败")
			return
		}
		//更新角色名
		res = tx.Model(&role).Updates(map[string]interface{}{
			"name":        data.Name,
			"description": data.Description,
		})
	} else {
		//新建角色
		role.Name = data.Name
		role.Description = data.Description
		res = tx.Create(&role)
	}

	if res.Error != nil {
		tx.Rollback()
		resp.ERROR(c, "更新数据库失败")
		return
	}

	rp := make([]model.AdminRolePermission, 0)
	if len(data.Permissions) > 0 {
		for _, per := range data.Permissions {
			rp = append(rp, model.AdminRolePermission{
				RoleId:       role.Id,
				PermissionId: per,
			})
		}
		res2 := tx.CreateInBatches(rp, len(rp))
		if res2.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "更新数据库失败")
			return
		}
	}

	tx.Commit()

	resp.SUCCESS(c)
}

func (h *SysRoleHandler) Remove(c *gin.Context) {
	var data struct {
		Id int
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Id > 0 {
		tx := h.db.Begin()
		res := tx.Where("id = ?", data.Id).Delete(&model.AdminRole{})
		if res.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "删除失败")
			return
		}
		res = tx.Where("role_id = ?", data.Id).Delete(&model.AdminRolePermission{})
		if res.Error != nil {
			tx.Rollback()
			resp.ERROR(c, "删除失败")
			return
		}
		tx.Commit()
	}
	resp.SUCCESS(c)
}
