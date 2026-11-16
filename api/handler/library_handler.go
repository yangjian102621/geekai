package handler

import (
	"geekai/core"
	"geekai/service/library"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	app     *core.AppServer
	service *library.Service
}

func NewLibraryHandler(app *core.AppServer, service *library.Service) *LibraryHandler {
	return &LibraryHandler{app: app, service: service}
}

func (h *LibraryHandler) RegisterRoutes() {
	h.app.Engine.GET("/api/v1/library", h.catalog)
}

func (h *LibraryHandler) catalog(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	items, err := h.service.Catalog(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取发布目录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
