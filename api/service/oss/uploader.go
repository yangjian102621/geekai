package oss

import "github.com/gin-gonic/gin"

type Uploader interface {
	PutFile(ctx *gin.Context) (string, error)
	PutImg(imageURL string) (string, error)
}
