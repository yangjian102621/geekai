package handler

import (
	"chatplus/core"
	"chatplus/service/oss"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadHandler struct {
	BaseHandler
	db              *gorm.DB
	uploaderManager *oss.UploaderManager
}

func NewUploadHandler(app *core.AppServer, db *gorm.DB, manager *oss.UploaderManager) *UploadHandler {
	handler := &UploadHandler{db: db, uploaderManager: manager}
	handler.App = app
	return handler
}

func (h *UploadHandler) Upload(c *gin.Context) {
	fileURL, err := h.uploaderManager.GetUploadHandler().PutFile(c, "file")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, fileURL)
}
