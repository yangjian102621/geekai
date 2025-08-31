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
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiNiuOss struct {
	config    types.QiNiuOssConfig
	mac       *qbox.Mac
	putPolicy storage.PutPolicy
	uploader  *storage.FormUploader
	bucket    *storage.BucketManager
	proxyURL  string
}

func NewQiNiuOss(sysConfig *types.SystemConfig, appConfig *types.AppConfig) *QiNiuOss {
	s := &QiNiuOss{
		proxyURL: appConfig.ProxyURL,
	}
	if sysConfig.OSS.Active == QiNiu {
		s.UpdateConfig(sysConfig.OSS.QiNiu)
	}
	return s
}

func (s *QiNiuOss) UpdateConfig(config types.QiNiuOssConfig) {
	zone, ok := storage.GetRegionByID(storage.RegionID(config.Zone))
	if !ok {
		zone = storage.ZoneHuanan
	}
	storeConfig := storage.Config{Zone: &zone}
	formUploader := storage.NewFormUploader(&storeConfig)
	// generate token
	mac := qbox.NewMac(config.AccessKey, config.AccessSecret)
	putPolicy := storage.PutPolicy{
		Scope: config.Bucket,
	}
	s.config = config
	s.mac = mac
	s.putPolicy = putPolicy
	s.uploader = formUploader
	s.bucket = storage.NewBucketManager(mac, &storeConfig)
}
func (s QiNiuOss) PutFile(ctx *gin.Context, name string) (File, error) {
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
	key := fmt.Sprintf("%d%s", time.Now().UnixMicro(), fileExt)
	// 上传文件
	ret := storage.PutRet{}
	extra := storage.PutExtra{}
	err = s.uploader.Put(ctx, &ret, s.putPolicy.UploadToken(s.mac), key, src, file.Size, &extra)
	if err != nil {
		return File{}, err
	}

	return File{
		Name:   file.Filename,
		ObjKey: key,
		URL:    fmt.Sprintf("%s/%s", s.config.Domain, ret.Key),
		Ext:    fileExt,
		Size:   file.Size,
	}, nil

}

func (s QiNiuOss) PutUrlFile(fileURL string, ext string, useProxy bool) (string, error) {
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
	key := fmt.Sprintf("%d%s", time.Now().UnixMicro(), ext)
	ret := storage.PutRet{}
	extra := storage.PutExtra{}
	// 上传文件字节数据
	err = s.uploader.Put(context.Background(), &ret, s.putPolicy.UploadToken(s.mac), key, bytes.NewReader(fileData), int64(len(fileData)), &extra)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.config.Domain, ret.Key), nil
}

func (s QiNiuOss) PutBase64(base64Img string) (string, error) {
	imageData, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return "", fmt.Errorf("error decoding base64:%v", err)
	}
	objectKey := fmt.Sprintf("%d.png", time.Now().UnixMicro())
	ret := storage.PutRet{}
	extra := storage.PutExtra{}
	// 上传文件字节数据
	err = s.uploader.Put(context.Background(), &ret, s.putPolicy.UploadToken(s.mac), objectKey, bytes.NewReader(imageData), int64(len(imageData)), &extra)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.config.Domain, ret.Key), nil
}

func (s QiNiuOss) Delete(fileURL string) error {
	var objectKey string
	if strings.HasPrefix(fileURL, "http") {
		objectKey = filepath.Base(fileURL)
	} else {
		objectKey = fileURL
	}

	return s.bucket.Delete(s.config.Bucket, objectKey)
}

var _ Uploader = QiNiuOss{}
