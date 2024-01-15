package oss

import "github.com/gin-gonic/gin"

type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
	Ext  string `json:"ext"`
}
type Uploader interface {
	PutFile(ctx *gin.Context, name string) (File, error)
	PutImg(imageURL string, useProxy bool) (string, error)
	Delete(fileURL string) error
}
