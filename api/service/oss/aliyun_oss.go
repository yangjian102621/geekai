package oss

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
)

type AliYunOss struct {
	config   types.AliYunOssConfig
	bucket   *oss.Bucket
	proxyURL string
}

func NewAliYunOss(sysConfig *types.SystemConfig, appConfig *types.AppConfig) (*AliYunOss, error) {
	s := &AliYunOss{
		proxyURL: appConfig.ProxyURL,
	}
	if sysConfig.OSS.Active == AliYun {
		err := s.UpdateConfig(sysConfig.OSS.AliYun)
		if err != nil {
			logger.Errorf("阿里云OSS初始化失败: %v", err)
		}
	}
	return s, nil

}

func (s *AliYunOss) UpdateConfig(config types.AliYunOssConfig) error {
	client, err := oss.New(config.Endpoint, config.AccessKey, config.AccessSecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return err
	}
	s.bucket = bucket
	s.config = config
	return nil
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
	objectKey := fmt.Sprintf("%d%s", time.Now().UnixMicro(), fileExt)
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

func (s AliYunOss) PutUrlFile(fileURL string, ext string, useProxy bool) (string, error) {
	var fileData []byte
	var err error
	if useProxy {
		fileData, err = utils.DownloadImage(fileURL, s.proxyURL)
	} else {
		fileData, err = utils.DownloadImage(fileURL, "")
	}
	if err != nil {
		return "", fmt.Errorf("error with download image: %v", err)
	}
	parse, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("error with parse image URL: %v", err)
	}
	if ext == "" {
		ext = filepath.Ext(parse.Path)
	}
	objectKey := fmt.Sprintf("%d%s", time.Now().UnixMicro(), ext)
	// 上传文件字节数据
	err = s.bucket.PutObject(objectKey, bytes.NewReader(fileData))
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
	objectKey := fmt.Sprintf("%d.png", time.Now().UnixMicro())
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
		objectKey = filepath.Base(fileURL)
	} else {
		objectKey = fileURL
	}
	return s.bucket.DeleteObject(objectKey)
}

var _ Uploader = AliYunOss{}
