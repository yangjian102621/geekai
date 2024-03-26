package oss

import (
	"chatplus/core/types"
	"chatplus/utils"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MiniOss struct {
	config   *types.MiniOssConfig
	client   *minio.Client
	proxyURL string
}

func NewMiniOss(appConfig *types.AppConfig) (MiniOss, error) {
	config := &appConfig.OSS.Minio
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.AccessSecret, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return MiniOss{}, err
	}
	if config.SubDir == "" {
		config.SubDir = "gpt"
	}
	return MiniOss{config: config, client: minioClient, proxyURL: appConfig.ProxyURL}, nil
}

func (s MiniOss) PutImg(imageURL string, useProxy bool) (string, error) {
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
	filename := fmt.Sprintf("%s/%d%s", s.config.SubDir, time.Now().UnixMicro(), fileExt)
	info, err := s.client.PutObject(
		context.Background(),
		s.config.Bucket,
		filename,
		strings.NewReader(string(imageData)),
		int64(len(imageData)),
		minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s.config.Domain, s.config.Bucket, info.Key), nil
}

func (s MiniOss) PutFile(ctx *gin.Context, name string) (File, error) {
	file, err := ctx.FormFile(name)
	if err != nil {
		return File{}, fmt.Errorf("error with get form: %v", err)
	}
	// Open the uploaded file
	fileReader, err := file.Open()
	if err != nil {
		return File{}, fmt.Errorf("error opening file: %v", err)
	}
	defer fileReader.Close()

	fileExt := utils.GetImgExt(file.Filename)
	filename := fmt.Sprintf("%s/%d%s", s.config.SubDir, time.Now().UnixMicro(), fileExt)
	info, err := s.client.PutObject(ctx, s.config.Bucket, filename, fileReader, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return File{}, fmt.Errorf("error uploading to MinIO: %v", err)
	}

	return File{
		Name:   file.Filename,
		ObjKey: info.Key,
		URL:    fmt.Sprintf("%s/%s/%s", s.config.Domain, s.config.Bucket, info.Key),
		Ext:    fileExt,
		Size:   file.Size,
	}, nil
}

func (s MiniOss) PutBase64(base64Img string) (string, error) {
	imageData, err := base64.StdEncoding.DecodeString(base64Img)
	if err != nil {
		return "", fmt.Errorf("error decoding base64:%v", err)
	}
	objectKey := fmt.Sprintf("%s/%d.png", s.config.SubDir, time.Now().UnixMicro())
	info, err := s.client.PutObject(
		context.Background(),
		s.config.Bucket,
		objectKey,
		strings.NewReader(string(imageData)),
		int64(len(imageData)),
		minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s.config.Domain, s.config.Bucket, info.Key), nil
}

func (s MiniOss) Delete(fileURL string) error {
	var objectKey string
	if strings.HasPrefix(fileURL, "http") {
		filename := filepath.Base(fileURL)
		objectKey = fmt.Sprintf("%s/%s", s.config.SubDir, filename)
	} else {
		objectKey = fileURL
	}
	return s.client.RemoveObject(context.Background(), s.config.Bucket, objectKey, minio.RemoveObjectOptions{})
}

var _ Uploader = MiniOss{}
