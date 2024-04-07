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

type ApiKeyHandler struct {
	handler.BaseHandler
}

func NewApiKeyHandler(app *core.AppServer, db *gorm.DB) *ApiKeyHandler {
	return &ApiKeyHandler{BaseHandler: handler.BaseHandler{DB: db, App: app}}
}

func (h *ApiKeyHandler) Save(c *gin.Context) {
	var data struct {
		Id       uint   `json:"id"`
		Platform string `json:"platform"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Value    string `json:"value"`
		ApiURL   string `json:"api_url"`
		Enabled  bool   `json:"enabled"`
		ProxyURL string `json:"proxy_url"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	apiKey := model.ApiKey{}
	if data.Id > 0 {
		h.DB.Find(&apiKey, data.Id)
	}
	apiKey.Platform = data.Platform
	apiKey.Value = data.Value
	apiKey.Type = data.Type
	apiKey.ApiURL = data.ApiURL
	apiKey.Enabled = data.Enabled
	apiKey.ProxyURL = data.ProxyURL
	apiKey.Name = data.Name
	res := h.DB.Save(&apiKey)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	var keyVo vo.ApiKey
	err := utils.CopyObject(apiKey, &keyVo)
	if err != nil {
		resp.ERROR(c, "数据拷贝失败！")
		return
	}
	keyVo.Id = apiKey.Id
	keyVo.CreatedAt = apiKey.CreatedAt.Unix()
	resp.SUCCESS(c, keyVo)
}

func (h *ApiKeyHandler) List(c *gin.Context) {
	if err := utils.CheckPermission(c, h.DB); err != nil {
		resp.NotPermission(c)
		return
	}

	var items []model.ApiKey
	var keys = make([]vo.ApiKey, 0)
	res := h.DB.Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var key vo.ApiKey
			err := utils.CopyObject(item, &key)
			if err == nil {
				key.Id = item.Id
				key.CreatedAt = item.CreatedAt.Unix()
				key.UpdatedAt = item.UpdatedAt.Unix()
				keys = append(keys, key)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, keys)
}

func (h *ApiKeyHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.ApiKey{}).Where("id = ?", data.Id).Update(data.Filed, data.Value)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}
	resp.SUCCESS(c)
}

func (h *ApiKeyHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Where("id", id).Delete(&model.ApiKey{})
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	resp.SUCCESS(c)
}
