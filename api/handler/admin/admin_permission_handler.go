package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SysPermissionHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewSysPermissionHandler(app *core.AppServer, db *gorm.DB) *SysPermissionHandler {
	h := SysPermissionHandler{db: db}
	h.App = app
	return &h
}

func (h *SysPermissionHandler) List(c *gin.Context) {
	if err := utils.CheckPermission(c, h.db); err != nil {
		resp.ERROR(c, types.NoPermission)
		return
	}

	var items []model.AdminPermission
	var data = make([]vo.AdminPermission, 0)
	res := h.db.Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "暂无数据")
		return
	}
	for _, item := range items {
		adminPermissionVo := vo.AdminPermission{}
		_ = utils.CopyObject(item, &adminPermissionVo)
		data = append(data, adminPermissionVo)
	}

	data = ArrayToTree(data)
	resp.SUCCESS(c, data)
}

func ArrayToTree(dates []vo.AdminPermission) []vo.AdminPermission {
	group := make(map[int][]vo.AdminPermission, 0)
	for _, node := range dates {
		group[node.Pid] = append(group[node.Pid], node)
	}
	// 初始化递归，从根节点开始构建树
	result := FindSiblings(group[0], group)

	return result
}

func FindSiblings(siblings []vo.AdminPermission, group map[int][]vo.AdminPermission) []vo.AdminPermission {
	result := make([]vo.AdminPermission, 0)
	for _, sibling := range siblings {
		children, ok := group[sibling.Id]
		if ok {
			sibling.Children = FindSiblings(children, group)
		}

		result = append(result, sibling)
	}
	return result
}

func (h *SysPermissionHandler) Save(c *gin.Context) {
	var data struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
		Sort int    `json:"sort"`
		Pid  int    `json:"pid"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var permission = model.AdminPermission{}
	var res *gorm.DB
	if data.Id > 0 { // 更新
		permission.Id = data.Id
		// 此处需要用 map 更新，用结构体无法更新 0 值
		res = h.db.Model(&permission).Updates(map[string]interface{}{
			"name": data.Name,
			"slug": data.Slug,
			"sort": data.Sort,
			"pid":  data.Pid,
		})
	} else {
		p := model.AdminPermission{
			Name: data.Name,
			Slug: data.Slug,
			Sort: data.Sort,
			Pid:  data.Pid,
		}
		res = h.db.Create(&p)
	}
	if res.Error != nil {
		fmt.Println(res.Error)
		resp.ERROR(c, "更新数据库失败")
		return
	}
	resp.SUCCESS(c)
}

func (h *SysPermissionHandler) Remove(c *gin.Context) {
	var data struct {
		Id int
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Id > 0 {
		res := h.db.Where("id = ?", data.Id).Delete(&model.AdminPermission{})
		if res.Error != nil {
			resp.ERROR(c, "删除失败")
			return
		}
	}
	resp.SUCCESS(c)
}
