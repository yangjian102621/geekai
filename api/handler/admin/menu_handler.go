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

type MenuHandler struct {
	handler.BaseHandler
}

func NewMenuHandler(app *core.AppServer, db *gorm.DB) *MenuHandler {
	return &MenuHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

func (h *MenuHandler) Save(c *gin.Context) {
	var data struct {
		Id      uint   `json:"id"`
		Name    string `json:"name"`
		Icon    string `json:"icon"`
		URL     string `json:"url"`
		SortNum int    `json:"sort_num"`
		Enabled bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Save(&model.Menu{
		Id:      data.Id,
		Name:    data.Name,
		Icon:    data.Icon,
		URL:     data.URL,
		SortNum: data.SortNum,
		Enabled: data.Enabled,
	})
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}
	resp.SUCCESS(c)
}

// List 数据列表
func (h *MenuHandler) List(c *gin.Context) {
	var items []model.Menu
	var list = make([]vo.Menu, 0)
	res := h.DB.Order("sort_num ASC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var product vo.Menu
			err := utils.CopyObject(item, &product)
			if err == nil {
				list = append(list, product)
			}
		}
	}
	resp.SUCCESS(c, list)
}

func (h *MenuHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.Menu{}).Where("id", data.Id).UpdateColumn("enabled", data.Enabled)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}
	resp.SUCCESS(c)
}

func (h *MenuHandler) Sort(c *gin.Context) {
	var data struct {
		Ids   []uint `json:"ids"`
		Sorts []int  `json:"sorts"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	for index, id := range data.Ids {
		res := h.DB.Model(&model.Menu{}).Where("id", id).Update("sort_num", data.Sorts[index])
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}

	resp.SUCCESS(c)
}

func (h *MenuHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		res := h.DB.Where("id", id).Delete(&model.Menu{})
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}
