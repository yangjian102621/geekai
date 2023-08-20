package handler

import (
	"chatplus/core"
	"chatplus/utils"
	"chatplus/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadHandler struct {
	BaseHandler
	db *gorm.DB
}

func NewUploadHandler(app *core.AppServer, db *gorm.DB) *UploadHandler {
	handler := &UploadHandler{db: db}
	handler.App = app
	return handler
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("文件上传失败: %s", err.Error()))
		return
	}

	filePath, err := utils.GenUploadPath(h.App.Config.StaticDir, file.Filename)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("文件上传失败: %s", err.Error()))
		return
	}
	// 将文件保存到指定路径
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		resp.ERROR(c, fmt.Sprintf("文件保存失败: %s", err.Error()))
		return
	}

	resp.SUCCESS(c, utils.GenUploadUrl(h.App.Config.StaticDir, h.App.Config.StaticUrl, filePath))
}
