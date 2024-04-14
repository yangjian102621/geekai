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

// MarkMapHandler 生成思维导图
type MarkMapHandler struct {
	BaseHandler
}

func NewMarkMapHandler(app *core.AppServer, db *gorm.DB) *MarkMapHandler {
	return &MarkMapHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// GetModel get the chat model for generating Markdown text
func (h *MarkMapHandler) GetModel(c *gin.Context) {
	modelId := h.App.SysConfig.XMindModelId
	session := h.DB.Session(&gorm.Session{}).Where("enabled", true)
	if modelId > 0 {
		session = session.Where("id", modelId)
	} else {
		session = session.Where("platform", types.OpenAI)
	}
	var chatModel model.ChatModel
	res := session.First(&chatModel)
	if res.Error != nil {
		resp.ERROR(c, "No available AI model")
		return
	}

	var modelVo vo.ChatModel
	err := utils.CopyObject(chatModel, &modelVo)
	if err != nil {
		resp.ERROR(c, "error with copy object: "+err.Error())
		return
	}

	resp.SUCCESS(c, modelVo)
}
