package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/types"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type NetHandler struct {
	BaseHandler
	uploaderManager *oss.UploaderManager
}

func NewNetHandler(app *core.AppServer, db *gorm.DB, manager *oss.UploaderManager) *NetHandler {
	return &NetHandler{BaseHandler: BaseHandler{App: app, DB: db}, uploaderManager: manager}
}

func (h *NetHandler) Upload(c *gin.Context) {
	file, err := h.uploaderManager.GetUploadHandler().PutFile(c, "file")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	logger.Info("upload file: ", file.Name)
	// cut the file name if it's too long
	if len(file.Name) > 100 {
		file.Name = file.Name[:90] + file.Ext
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

func (h *NetHandler) List(c *gin.Context) {
	var data struct {
		Urls     []string `json:"urls,omitempty"`
		Page     int      `json:"page"`
		PageSize int      `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	userId := h.GetLoginUserId(c)
	var items []model.File
	var files = make([]vo.File, 0)
	session := h.DB.Session(&gorm.Session{})
	session = session.Where("user_id = ?", userId)
	if len(data.Urls) > 0 {
		session = session.Where("url IN ?", data.Urls)
	}
	// 统计总数
	var total int64
	session.Model(&model.File{}).Count(&total)

	if data.Page > 0 && data.PageSize > 0 {
		offset := (data.Page - 1) * data.PageSize
		session = session.Offset(offset).Limit(data.PageSize)
	}
	err := session.Order("id desc").Find(&items).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

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

	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, files))
}

// Remove remove files
func (h *NetHandler) Remove(c *gin.Context) {
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

func (h *NetHandler) Download(c *gin.Context) {
	fileUrl := c.Query("url")
	// 使用http工具下载文件
	if fileUrl == "" {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	// 使用http.Get下载文件
	r, err := http.Get(fileUrl)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		resp.ERROR(c, "error status："+r.Status)
		return
	}

	c.Status(http.StatusOK)
	// 将下载的文件内容写入响应
	_, _ = io.Copy(c.Writer, r.Body)
}
