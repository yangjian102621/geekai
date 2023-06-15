package handler

import (
	"chatplus/core"
	"chatplus/core/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewAdminHandler(config *types.AppConfig, app *core.AppServer, db *gorm.DB) *AdminHandler {
	handler := AdminHandler{db: db}
	handler.app = app
	handler.config = config
	return &handler
}

// Login 登录
func (h *AdminHandler) Login(c *gin.Context) {

}

// Logout 注销
func (h *AdminHandler) Logout(c *gin.Context) {
}
