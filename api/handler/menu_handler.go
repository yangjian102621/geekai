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

type MenuHandler struct {
	BaseHandler
}

func NewMenuHandler(app *core.AppServer, db *gorm.DB) *MenuHandler {
	return &MenuHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// List 数据列表
func (h *MenuHandler) List(c *gin.Context) {
	var items []model.Menu
	var list = make([]vo.Menu, 0)
	res := h.DB.Where("enabled", true).Order("sort_num ASC").Find(&items)
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
