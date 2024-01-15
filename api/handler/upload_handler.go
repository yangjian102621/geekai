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
	db              *gorm.DB
	uploaderManager *oss.UploaderManager
}

func NewUploadHandler(app *core.AppServer, db *gorm.DB, manager *oss.UploaderManager) *UploadHandler {
	handler := &UploadHandler{db: db, uploaderManager: manager}
	handler.App = app
	return handler
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := h.uploaderManager.GetUploadHandler().PutFile(c, "file")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	userId := h.GetLoginUserId(c)
	res := h.db.Create(&model.File{
		UserId:    userId,
		Name:      file.Name,
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
	h.db.Debug().Where("user_id = ?", userId).Find(&items)
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
