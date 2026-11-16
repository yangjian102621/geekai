package admin

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/service/library"
	"geekai/store/model"
	"geekai/utils/resp"
	"net/http"
	"strconv"

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
	group := h.app.Engine.Group("/api/admin/library")
	group.Use(middleware.AdminAuthMiddleware(h.app.Config.AdminSession.SecretKey, h.app.Redis))
	h.RegisterManagementRoutes(group)
}

func (h *LibraryHandler) RegisterManagementRoutes(group *gin.RouterGroup) {
	group.GET("/sources", h.sources)
	group.POST("/sources", h.saveSource)
	group.DELETE("/sources/:id", h.deleteSource)
	group.GET("/sources/:id/collections", h.collections)
	group.GET("/resources", h.resources)
	group.POST("/resources", h.saveResource)
	group.DELETE("/resources/:id", h.deleteResource)
}

func (h *LibraryHandler) sources(c *gin.Context) {
	result, err := h.service.Sources(c.Request.Context())
	if err != nil {
		resp.ERROR(c, "读取资料源失败")
		return
	}
	resp.SUCCESS(c, result)
}

func (h *LibraryHandler) saveSource(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 32768)
	var input library.SourceInput
	if c.ShouldBindJSON(&input) != nil {
		resp.ERROR(c, "参数无效")
		return
	}
	result, err := h.service.SaveSource(c.Request.Context(), input)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, result)
}

func libraryID(c *gin.Context) uint {
	value, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || value == 0 {
		resp.ERROR(c, "ID 无效")
		return 0
	}
	return uint(value)
}

func (h *LibraryHandler) deleteSource(c *gin.Context) {
	id := libraryID(c)
	if id == 0 {
		return
	}
	if err := h.service.DeleteSource(c.Request.Context(), id); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *LibraryHandler) collections(c *gin.Context) {
	id := libraryID(c)
	if id == 0 {
		return
	}
	result, err := h.service.RemoteCollections(c.Request.Context(), id)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, result)
}

func (h *LibraryHandler) resources(c *gin.Context) {
	result, err := h.service.Resources(c.Request.Context())
	if err != nil {
		resp.ERROR(c, "读取发布目录失败")
		return
	}
	resp.SUCCESS(c, result)
}

func (h *LibraryHandler) saveResource(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 32768)
	var input model.LibraryResource
	if c.ShouldBindJSON(&input) != nil {
		resp.ERROR(c, "参数无效")
		return
	}
	result, err := h.service.SaveResource(c.Request.Context(), input)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c, result)
}

func (h *LibraryHandler) deleteResource(c *gin.Context) {
	id := libraryID(c)
	if id == 0 {
		return
	}
	if err := h.service.DeleteResource(c.Request.Context(), id); err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}
