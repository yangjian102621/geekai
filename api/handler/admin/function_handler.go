package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FunctionHandler struct {
	handler.BaseHandler
}

func NewFunctionHandler(app *core.AppServer, db *gorm.DB) *FunctionHandler {
	return &FunctionHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
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
		Label:       data.Label,
		Description: data.Description,
		Parameters:  utils.JsonEncode(data.Parameters),
		Action:      data.Action,
		Token:       data.Token,
		Enabled:     data.Enabled,
	}

	res := h.DB.Save(&f)
	if res.Error != nil {
		resp.ERROR(c, "error with save data:"+res.Error.Error())
		return
	}
	data.Id = f.Id
	resp.SUCCESS(c, data)
}

func (h *FunctionHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.Function{}).Where("id = ?", data.Id).Update(data.Filed, data.Value)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}
	resp.SUCCESS(c)
}

func (h *FunctionHandler) List(c *gin.Context) {
	if err := utils.CheckPermission(c, h.DB); err != nil {
		resp.NotPermission(c)
		return
	}

	var items []model.Function
	res := h.DB.Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "No data found")
		return
	}

	functions := make([]vo.Function, 0)
	for _, v := range items {
		var f vo.Function
		err := utils.CopyObject(v, &f)
		if err != nil {
			continue
		}
		functions = append(functions, f)
	}
	resp.SUCCESS(c, functions)
}

func (h *FunctionHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		res := h.DB.Delete(&model.Function{Id: uint(id)})
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}

// GenToken generate function api access token
func (h *FunctionHandler) GenToken(c *gin.Context) {
	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 0,
		"expired": 0,
	})
	tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
	if err != nil {
		logger.Error("error with generate token", err)
		resp.ERROR(c)
		return
	}

	resp.SUCCESS(c, tokenString)
}
