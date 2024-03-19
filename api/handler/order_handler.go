package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	BaseHandler
}

func NewOrderHandler(app *core.AppServer, db *gorm.DB) *OrderHandler {
	return &OrderHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

func (h *OrderHandler) List(c *gin.Context) {
	var data struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	userId := h.GetLoginUserId(c)
	session := h.DB.Session(&gorm.Session{}).Where("user_id = ? AND status = ?", userId, types.OrderPaidSuccess)
	var total int64
	session.Model(&model.Order{}).Count(&total)
	var items []model.Order
	var list = make([]vo.Order, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var order vo.Order
			err := utils.CopyObject(item, &order)
			if err == nil {
				order.Id = item.Id
				order.CreatedAt = item.CreatedAt.Unix()
				order.UpdatedAt = item.UpdatedAt.Unix()
				list = append(list, order)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}
