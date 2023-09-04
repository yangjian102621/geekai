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

type ChatModelHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewChatModelHandler(app *core.AppServer, db *gorm.DB) *ChatModelHandler {
	h := ChatModelHandler{db: db}
	h.App = app
	return &h
}

// List 模型列表
func (h *ChatModelHandler) List(c *gin.Context) {
	var items []model.ChatModel
	var cms = make([]vo.ChatModel, 0)
	res := h.db.Where("enabled = ?", true).Order("sort_num ASC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var cm vo.ChatModel
			err := utils.CopyObject(item, &cm)
			if err == nil {
				cm.Id = item.Id
				cm.CreatedAt = item.CreatedAt.Unix()
				cm.UpdatedAt = item.UpdatedAt.Unix()
				cms = append(cms, cm)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, cms)
}
