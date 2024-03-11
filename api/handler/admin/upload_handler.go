package admin

import (
	"chatplus/core"
	"chatplus/handler"
	"chatplus/service/oss"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type UploadHandler struct {
	handler.BaseHandler
	db              *gorm.DB
	uploaderManager *oss.UploaderManager
}

func NewUploadHandler(app *core.AppServer, db *gorm.DB, manager *oss.UploaderManager) *UploadHandler {
	adminHandler := &UploadHandler{db: db, uploaderManager: manager}
	adminHandler.App = app
	return adminHandler
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := h.uploaderManager.GetUploadHandler().PutFile(c, "file")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	userId := 0
	res := h.db.Create(&model.File{
		UserId:    userId,
		Name:      file.Name,
		ObjKey:    file.ObjKey,
		URL:       file.URL,
		Ext:       file.Ext,
		Size:      file.Size,
		CreatedAt: time.Time{},
	})
	if res.Error != nil || res.RowsAffected == 0 {
		resp.ERROR(c, "error with update database: "+res.Error.Error())
		return
	}

	resp.SUCCESS(c, file)
}

func (h *UploadHandler) List(c *gin.Context) {
	userId := 0
	var items []model.File
	var files = make([]vo.File, 0)
	h.db.Where("user_id = ?", userId).Find(&items)
	if len(items) > 0 {
		for _, v := range items {
			var file vo.File
			err := utils.CopyObject(v, &file)
			if err != nil {
				logger.Error(err)
				continue
			}
			file.CreatedAt = v.CreatedAt.Unix()
			files = append(files, file)
		}
	}

	resp.SUCCESS(c, files)
}

// Remove remove files
func (h *UploadHandler) Remove(c *gin.Context) {
	userId := 0
	id := h.GetInt(c, "id", 0)
	var file model.File
	tx := h.db.Where("user_id = ? AND id = ?", userId, id).First(&file)
	if tx.Error != nil || file.Id == 0 {
		resp.ERROR(c, "file not existed")
		return
	}

	// remove database
	tx = h.db.Model(&model.File{}).Delete("id = ?", id)
	if tx.Error != nil || tx.RowsAffected == 0 {
		resp.ERROR(c, "failed to update database")
		return
	}
	// remove files
	objectKey := file.ObjKey
	if objectKey == "" {
		objectKey = file.URL
	}
	_ = h.uploaderManager.GetUploadHandler().Delete(objectKey)
	resp.SUCCESS(c)
}
