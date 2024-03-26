package oss

import (
	"bytes"
	"chatplus/core/types"
	"chatplus/utils"
	"encoding/base64"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
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

	if config.SubDir == "" {
		config.SubDir = "gpt"
	}

	return &AliYunOss{
		config:   config,
		bucket:   bucket,
		proxyURL: appConfig.ProxyURL,
	}, nil

}

func (s AliYunOss) PutFile(ctx *gin.Context, name string) (File, error) {
	// 解析表单
	file, err := ctx.FormFile(name)
	if err != nil {
		return File{}, err
	}
	// 打开上传文件
	src, err := file.Open()
	if err != nil {
		return File{}, err
	}
	defer src.Close()

	fileExt := filepath.Ext(file.Filename)
	objectKey := fmt.Sprintf("%s/%d%s", s.config.SubDir, time.Now().UnixMicro(), fileExt)
	// 上传文件
	err = s.bucket.PutObject(objectKey, src)
	if err != nil {
		return File{}, err
	}

	return File{
		Name:   file.Filename,
		ObjKey: objectKey,
		URL:    fmt.Sprintf("%s/%s", s.config.Domain, objectKey),
		Ext:    fileExt,
		Size:   file.Size,
	}, nil
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
	fileExt := utils.GetImgExt(parse.Path)
	objectKey := fmt.Sprintf("%s/%d%s", s.config.SubDir, time.Now().UnixMicro(), fileExt)
	// 上传文件字节数据
	err = s.bucket.PutObject(objectKey, bytes.NewReader(imageData))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.config.Domain, objectKey), nil
}

func (s AliYunOss) PutBase64(base64Img string) (string, error) {
	imageData, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return "", fmt.Errorf("error decoding base64:%v", err)
	}
	objectKey := fmt.Sprintf("%s/%d.png", s.config.SubDir, time.Now().UnixMicro())
	// 上传文件字节数据
	err = s.bucket.PutObject(objectKey, bytes.NewReader(imageData))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.config.Domain, objectKey), nil
}

func (s AliYunOss) Delete(fileURL string) error {
	var objectKey string
	if strings.HasPrefix(fileURL, "http") {
		filename := filepath.Base(fileURL)
		objectKey = fmt.Sprintf("%s/%s", s.config.SubDir, filename)
	} else {
		objectKey = fileURL
	}
	return s.bucket.DeleteObject(objectKey)
}

var _ Uploader = AliYunOss{}
