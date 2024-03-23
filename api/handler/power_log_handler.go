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

type PowerLogHandler struct {
	BaseHandler
}

func NewPowerLogHandler(app *core.AppServer, db *gorm.DB) *PowerLogHandler {
	return &PowerLogHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

func (h *PowerLogHandler) List(c *gin.Context) {
	var data struct {
		Model    string   `json:"model"`
		Date     []string `json:"date"`
		Page     int      `json:"page"`
		PageSize int      `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	userId := h.GetLoginUserId(c)
	session = session.Where("user_id", userId)
	if data.Model != "" {
		session = session.Where("model", data.Model)
	}
	if len(data.Date) == 2 {
		start := data.Date[0] + " 00:00:00"
		end := data.Date[1] + " 00:00:00"
		session = session.Where("created_at >= ? AND created_at <= ?", start, end)
	}

	var total int64
	session.Model(&model.PowerLog{}).Count(&total)
	var items []model.PowerLog
	var list = make([]vo.PowerLog, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var log vo.PowerLog
			err := utils.CopyObject(item, &log)
			if err != nil {
				continue
			}
			log.Id = item.Id
			log.CreatedAt = item.CreatedAt.Unix()
			log.TypeStr = item.Type.String()
			list = append(list, log)
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}
