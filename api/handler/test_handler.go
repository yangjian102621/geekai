package handler

import (
	"chatplus/service"
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestHandler struct {
	db        *gorm.DB
	snowflake *service.Snowflake
}

func NewTestHandler(db *gorm.DB, snowflake *service.Snowflake) *TestHandler {
	return &TestHandler{db: db, snowflake: snowflake}
}

func (h *TestHandler) Test(c *gin.Context) {
	h.initMjTaskId(c)
}

func (h *TestHandler) initUserNickname(c *gin.Context) {
	var users []model.User
	tx := h.db.Find(&users)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	for _, u := range users {
		u.Nickname = fmt.Sprintf("极客学长@%d", utils.RandomNumber(6))
		h.db.Updates(&u)
	}

	resp.SUCCESS(c)
}

func (h *TestHandler) initMjTaskId(c *gin.Context) {
	var jobs []model.MidJourneyJob
	tx := h.db.Find(&jobs)
	if tx.Error != nil {
		resp.ERROR(c, tx.Error.Error())
		return
	}

	for _, job := range jobs {
		id, _ := h.snowflake.Next(true)
		job.TaskId = id
		h.db.Updates(&job)
	}

	resp.SUCCESS(c)
}
