package oss

import "github.com/gin-gonic/gin"

const Local = "LOCAL"
const Minio = "MINIO"
const QiNiu = "QINIU"
const AliYun = "ALIYUN"

type File struct {
	Name   string `json:"name"`
	ObjKey string `json:"obj_key"`
	Size   int64  `json:"size"`
	URL    string `json:"url"`
	Ext    string `json:"ext"`
}
type Uploader interface {
	PutFile(ctx *gin.Context, name string) (File, error)
	PutImg(imageURL string, useProxy bool) (string, error)
	PutBase64(imageData string) (string, error)
	Delete(fileURL string) error
}
