package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewConfigHandler(app *core.AppServer, db *gorm.DB) *ConfigHandler {
	h := ConfigHandler{db: db}
	h.App = app
	return &h
}

func (h *ConfigHandler) Update(c *gin.Context) {
	var data struct {
		Key    string                 `json:"key"`
		Config map[string]interface{} `json:"config"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	str := utils.JsonEncode(&data.Config)
	config := model.Config{Key: data.Key, Config: str}
	res := h.db.FirstOrCreate(&config, model.Config{Key: data.Key})
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	if config.Id > 0 {
		config.Config = str
		res := h.db.Updates(&config)
		if res.Error != nil {
			resp.ERROR(c, res.Error.Error())
			return
		}

		// update config cache for AppServer
		var cfg model.Config
		h.db.Where("marker", data.Key).First(&cfg)
		var err error
		if data.Key == "system" {
			err = utils.JsonDecode(cfg.Config, &h.App.SysConfig)
		} else if data.Key == "chat" {
			err = utils.JsonDecode(cfg.Config, &h.App.ChatConfig)
		}
		if err != nil {
			resp.ERROR(c, "Failed to update config cache: "+err.Error())
			return
		}
		logger.Infof("Update AppServer's config successfully: %v", config.Config)
	}

	resp.SUCCESS(c, config)
}

// Get 获取指定的系统配置
func (h *ConfigHandler) Get(c *gin.Context) {
	key := c.Query("key")
	var config model.Config
	res := h.db.Where("marker", key).First(&config)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	var m map[string]interface{}
	err := utils.JsonDecode(config.Config, &m)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, m)
}
