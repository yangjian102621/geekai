package oss

import (
	"bytes"
	"chatplus/core/types"
	"chatplus/utils"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"net/url"
	"path/filepath"
	"time"
)

type AliYunOss struct {
	config   *types.AliYunOssConfig
	bucket   *oss.Bucket
	proxyURL string
}

func NewAliYunOss(appConfig *types.AppConfig) (*AliYunOss, error) {
	config := &appConfig.OSS.AliYun
	// 创建 OSS 客户端
	client, err := oss.New(config.Endpoint, config.AccessKey, config.AccessSecret)
	if err != nil {
		return nil, err
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return nil, err
	}

	return &AliYunOss{
		config:   config,
		bucket:   bucket,
		proxyURL: appConfig.ProxyURL,
	}, nil

}

func (s AliYunOss) PutFile(ctx *gin.Context, name string) (string, error) {
	// 解析表单
	file, err := ctx.FormFile(name)
	if err != nil {
		return "", err
	}
	// 打开上传文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileExt := filepath.Ext(file.Filename)
	objectKey := fmt.Sprintf("%d%s", time.Now().UnixMicro(), fileExt)
	// 上传文件
	err = s.bucket.PutObject(objectKey, src)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.%s/%s", s.config.Bucket, s.config.Domain, objectKey), nil
}

func (s AliYunOss) PutImg(imageURL string, useProxy bool) (string, error) {
	var imageData []byte
	var err error
	if useProxy {
		imageData, err = utils.DownloadImage(imageURL, s.proxyURL)
	} else {
		imageData, err = utils.DownloadImage(imageURL, "")
	}
	if err != nil {
		return "", fmt.Errorf("error with download image: %v", err)
	}
	parse, err := url.Parse(imageURL)
	if err != nil {
		return "", fmt.Errorf("error with parse image URL: %v", err)
	}
	fileExt := filepath.Ext(parse.Path)
	objectKey := fmt.Sprintf("%d%s", time.Now().UnixMicro(), fileExt)
	// 上传文件字节数据
	err = s.bucket.PutObject(objectKey, bytes.NewReader(imageData))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s.%s/%s", s.config.Bucket, s.config.Domain, objectKey), nil
}

func (s AliYunOss) Delete(fileURL string) error {
	objectName := filepath.Base(fileURL)
	return s.bucket.DeleteObject(objectName)
}

var _ Uploader = AliYunOss{}
