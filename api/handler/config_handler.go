package handler

import (
	"chatplus/core"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	BaseHandler
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB) *ConfigHandler {
	return &ConfigHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

// Get 获取指定的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	key := c.Query("key")
	var config model.Config
	res := h.DB.Where("marker", key).First(&config)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var value map[string]interface{}
	err := utils.JsonDecode(config.Config, &value)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, value)
}
