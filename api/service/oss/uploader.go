package oss

import "github.com/gin-gonic/gin"

type File struct {
	Size int64
	URL  string
	Ext  string
}
type Uploader interface {
	PutFile(ctx *gin.Context, name string) (File, error)
	PutImg(imageURL string, useProxy bool) (string, error)
	Delete(fileURL string) error
}
