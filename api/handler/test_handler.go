package handler

import (
	"chatplus/store/model"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestHandler struct {
	db *gorm.DB
}

func NewTestHandler(db *gorm.DB) *TestHandler {
	return &TestHandler{db: db}
}

func (h *TestHandler) Test(c *gin.Context) {
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
