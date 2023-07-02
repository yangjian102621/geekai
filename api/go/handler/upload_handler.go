package handler

import (
	"chatplus/core"
	"chatplus/utils/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
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

	filePath, err := h.genFilePath(file.Filename)
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

	resp.SUCCESS(c, h.genFileUrl(filePath))
}

// 生成上传文件路径
func (h *UploadHandler) genFilePath(filename string) (string, error) {
	now := time.Now()
	dir := fmt.Sprintf("%s/upload/%d/%d", h.App.Config.StaticDir, now.Year(), now.Month())
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("创建上传目录失败：%s", err)
		}
	}
	fileExt := filepath.Ext(filename)
	return fmt.Sprintf("%s/%d%s", dir, now.UnixMilli(), fileExt), nil
}

// 生成上传文件 URL
func (h *UploadHandler) genFileUrl(filePath string) string {
	now := time.Now()
	filename := filepath.Base(filePath)
	return fmt.Sprintf("%s/upload/%d/%d/%s", h.App.Config.StaticUrl, now.Year(), now.Month(), filename)
}
