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
	"time"
)

type ProductHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewProductHandler(app *core.AppServer, db *gorm.DB) *ProductHandler {
	h := ProductHandler{db: db}
	h.App = app
	return &h
}

func (h *ProductHandler) Save(c *gin.Context) {
	var data struct {
		Id        uint    `json:"id"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Discount  float64 `json:"discount"`
		Enabled   bool    `json:"enabled"`
		Days      int     `json:"days"`
		Calls     int     `json:"calls"`
		ImgCalls  int     `json:"img_calls"`
		CreatedAt int64   `json:"created_at"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	item := model.Product{
		Name:     data.Name,
		Price:    data.Price,
		Discount: data.Discount,
		Days:     data.Days,
		Calls:    data.Calls,
		ImgCalls: data.ImgCalls,
		Enabled:  data.Enabled}
	item.Id = data.Id
	if item.Id > 0 {
		item.CreatedAt = time.Unix(data.CreatedAt, 0)
	}
	res := h.db.Save(&item)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	var itemVo vo.Product
	err := utils.CopyObject(item, &itemVo)
	if err != nil {
		resp.ERROR(c, "数据拷贝失败！")
		return
	}
	itemVo.Id = item.Id
	itemVo.UpdatedAt = item.UpdatedAt.Unix()
	resp.SUCCESS(c, itemVo)
}

// List 模型列表
func (h *ProductHandler) List(c *gin.Context) {
	session := h.db.Session(&gorm.Session{})
	enable := h.GetBool(c, "enable")
	if enable {
		session = session.Where("enabled", enable)
	}
	var items []model.Product
	var list = make([]vo.Product, 0)
	res := session.Order("sort_num ASC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var product vo.Product
			err := utils.CopyObject(item, &product)
			if err == nil {
				product.Id = item.Id
				product.CreatedAt = item.CreatedAt.Unix()
				product.UpdatedAt = item.UpdatedAt.Unix()
				list = append(list, product)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, list)
}

func (h *ProductHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.db.Model(&model.Product{}).Where("id = ?", data.Id).Update("enabled", data.Enabled)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}
	resp.SUCCESS(c)
}

func (h *ProductHandler) Sort(c *gin.Context) {
	var data struct {
		Ids   []uint `json:"ids"`
		Sorts []int  `json:"sorts"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	for index, id := range data.Ids {
		res := h.db.Model(&model.Product{}).Where("id = ?", id).Update("sort_num", data.Sorts[index])
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}

	resp.SUCCESS(c)
}

func (h *ProductHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		res := h.db.Where("id = ?", id).Delete(&model.Product{})
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}
