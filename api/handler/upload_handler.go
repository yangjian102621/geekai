package handler

import (
	"chatplus/core"
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
	BaseHandler
	uploaderManager *oss.UploaderManager
}

func NewUploadHandler(app *core.AppServer, db *gorm.DB, manager *oss.UploaderManager) *UploadHandler {
	return &UploadHandler{BaseHandler: BaseHandler{App: app, DB: db}, uploaderManager: manager}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := h.uploaderManager.GetUploadHandler().PutFile(c, "file")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	userId := h.GetLoginUserId(c)
	res := h.DB.Create(&model.File{
		UserId:    int(userId),
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
	userId := h.GetLoginUserId(c)
	var items []model.File
	var files = make([]vo.File, 0)
	h.DB.Where("user_id = ?", userId).Find(&items)
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
	userId := h.GetLoginUserId(c)
	id := h.GetInt(c, "id", 0)
	var file model.File
	tx := h.DB.Where("user_id = ? AND id = ?", userId, id).First(&file)
	if tx.Error != nil || file.Id == 0 {
		resp.ERROR(c, "file not existed")
		return
	}

	// remove database
	tx = h.DB.Model(&model.File{}).Delete("id = ?", id)
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
