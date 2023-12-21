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

type FunctionHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewFunctionHandler(app *core.AppServer, db *gorm.DB) *FunctionHandler {
	h := FunctionHandler{db: db}
	h.App = app
	return &h
}

func (h *FunctionHandler) Save(c *gin.Context) {
	var data vo.Function
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var f = model.Function{
		Id:          data.Id,
		Name:        data.Name,
		Description: data.Description,
		Parameters:  utils.JsonEncode(data.Parameters),
		Required:    utils.JsonEncode(data.Required),
		Action:      data.Action,
		Enabled:     false,
	}

	logger.Info(f)
	resp.SUCCESS(c, data)
}

func (h *FunctionHandler) List(c *gin.Context) {

	resp.SUCCESS(c)
}

func (h *FunctionHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		res := h.db.Delete(&model.Function{Id: uint(id)})
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}
