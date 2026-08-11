package oss

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentOss struct {
	config   types.TencentOssConfig
	client   *cos.Client
	proxyURL string
}

func NewTencentOss(sysConfig *types.SystemConfig, appConfig *types.AppConfig) (*TencentOss, error) {
	s := &TencentOss{
		proxyURL: appConfig.ProxyURL,
	}
	err := s.UpdateConfig(sysConfig.OSS.Tencent)
	if err != nil {
		logger.Warnf("腾讯云COS初始化失败: %v", err)
	}
	return s, nil
}

func (s *TencentOss) UpdateConfig(config types.TencentOssConfig) error {
	if config.Bucket == "" || config.Region == "" || config.SecretId == "" || config.SecretKey == "" {
		// 配置不完整时不初始化客户端
		s.config = config
		return nil
	}

	// 构建 COS 客户端 URL
	cosURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Bucket, config.Region)
	u, err := url.Parse(cosURL)
	if err != nil {
		return fmt.Errorf("error parsing COS URL: %v", err)
	}

	// 创建 COS 客户端
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretId,
			SecretKey: config.SecretKey,
		},
	})

	s.client = client
	s.config = config
	return nil
}

func (s TencentOss) PutFile(ctx *gin.Context, name string) (File, error) {
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
	_, err = s.client.Object.Put(ctx, objectKey, src, nil)
	if err != nil {
		return File{}, err
	}

	// 生成文件 URL
	fileURL := s.generateURL(objectKey)

	return File{
		Name:   file.Filename,
		ObjKey: objectKey,
		URL:    fileURL,
		Ext:    fileExt,
		Size:   file.Size,
	}, nil
}

func (s TencentOss) PutUrlFile(fileURL string, ext string, useProxy bool) (string, error) {
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
	_, err = s.client.Object.Put(context.Background(), objectKey, bytes.NewReader(fileData), nil)
	if err != nil {
		return "", err
	}
	return s.generateURL(objectKey), nil
}

func (s TencentOss) PutBase64(base64Img string) (string, error) {
	imageData, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return "", fmt.Errorf("error decoding base64:%v", err)
	}
	objectKey := fmt.Sprintf("%d.png", time.Now().UnixMicro())
	// 上传文件字节数据
	_, err = s.client.Object.Put(context.Background(), objectKey, bytes.NewReader(imageData), nil)
	if err != nil {
		return "", err
	}
	return s.generateURL(objectKey), nil
}

func (s TencentOss) Delete(fileURL string) error {
	var objectKey string
	if strings.HasPrefix(fileURL, "http") {
		objectKey = filepath.Base(fileURL)
	} else {
		objectKey = fileURL
	}
	_, err := s.client.Object.Delete(context.Background(), objectKey)
	return err
}

// generateURL 生成文件访问 URL
func (s TencentOss) generateURL(objectKey string) string {
	if s.config.Domain != "" {
		// 使用自定义域名
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(s.config.Domain, "/"), objectKey)
	}
	// 使用 COS 默认域名
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", s.config.Bucket, s.config.Region, objectKey)
}

var _ Uploader = TencentOss{}
