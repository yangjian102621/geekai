package admin

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatAppTypeHandler struct {
	handler.BaseHandler
}

func NewChatAppTypeHandler(app *core.AppServer, db *gorm.DB) *ChatAppTypeHandler {
	return &ChatAppTypeHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// Save 创建或更新App类型
func (h *ChatAppTypeHandler) Save(c *gin.Context) {
	var data struct {
		Id      uint   `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
		Icon    string `json:"icon"`
		SortNum int    `json:"sort_num"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.Id == 0 { // for add
		err := h.DB.Where("name", data.Name).First(&model.AppType{}).Error
		if err == nil {
			resp.ERROR(c, "当前分类已经存在")
			return
		}
		err = h.DB.Create(&model.AppType{
			Name:    data.Name,
			Icon:    data.Icon,
			Enabled: data.Enabled,
			SortNum: data.SortNum,
		}).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	} else { // for update
		err := h.DB.Model(&model.AppType{}).Where("id", data.Id).Updates(map[string]interface{}{
			"name":    data.Name,
			"icon":    data.Icon,
			"enabled": data.Enabled,
		}).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	resp.SUCCESS(c)
}

// List 获取App类型列表
func (h *ChatAppTypeHandler) List(c *gin.Context) {
	var items []model.AppType
	var appTypes = make([]vo.AppType, 0)
	err := h.DB.Order("sort_num ASC").Find(&items).Error
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

// Remove 删除App类型
func (h *ChatAppTypeHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id <= 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	err := h.DB.Where("id", id).Delete(&model.AppType{}).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

// Enable 启用|禁用
func (h *ChatAppTypeHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Model(&model.AppType{}).Where("id", data.Id).UpdateColumn("enabled", data.Enabled).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

// Sort 更新排序
func (h *ChatAppTypeHandler) Sort(c *gin.Context) {
	var data struct {
		Ids   []uint `json:"ids"`
		Sorts []int  `json:"sorts"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	for index, id := range data.Ids {
		err := h.DB.Model(&model.AppType{}).Where("id", id).Update("sort_num", data.Sorts[index]).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	resp.SUCCESS(c)
}
