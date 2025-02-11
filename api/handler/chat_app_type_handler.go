package handler

import (
	"geekai/core"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatAppTypeHandler struct {
	BaseHandler
}

func NewChatAppTypeHandler(app *core.AppServer, db *gorm.DB) *ChatAppTypeHandler {
	return &ChatAppTypeHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// List 获取App类型列表
func (h *ChatAppTypeHandler) List(c *gin.Context) {
	var items []model.AppType
	var appTypes = make([]vo.AppType, 0)
	err := h.DB.Where("enabled", true).Order("sort_num ASC").Find(&items).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	for _, v := range items {
		var appType vo.AppType
		err = utils.CopyObject(v, &appType)
		if err != nil {
			continue
		}
		appType.Id = v.Id
		appType.CreatedAt = v.CreatedAt.Unix()
		appTypes = append(appTypes, appType)
	}

	resp.SUCCESS(c, appTypes)
}
