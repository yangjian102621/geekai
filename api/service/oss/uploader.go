package oss

import "github.com/gin-gonic/gin"

type Uploader interface {
	PutFile(ctx *gin.Context, name string) (string, error)
	PutImg(imageURL string, useProxy bool) (string, error)
	Delete(fileURL string) error
}
