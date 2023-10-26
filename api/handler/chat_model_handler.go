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
	var chatModels = make([]vo.ChatModel, 0)
	// 只加载用户订阅的 AI 模型
	user, err := utils.GetLoginUser(c, h.db)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	var models []string
	err = utils.JsonDecode(user.ChatModels, &models)
	if err != nil {
		resp.ERROR(c, "当前用户没有订阅任何模型")
		return
	}

	res := h.db.Where("enabled = ?", true).Where("value IN ?", models).Order("sort_num ASC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var cm vo.ChatModel
			err := utils.CopyObject(item, &cm)
			if err == nil {
				cm.Id = item.Id
				cm.CreatedAt = item.CreatedAt.Unix()
				cm.UpdatedAt = item.UpdatedAt.Unix()
				chatModels = append(chatModels, cm)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, chatModels)
}
