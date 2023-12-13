package oss

import (
	"chatplus/core/types"
	"chatplus/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
	"path/filepath"
	"strings"
	"time"
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

func (s MiniOss) PutFile(ctx *gin.Context, name string) (string, error) {
	file, err := ctx.FormFile(name)
	if err != nil {
		return "", fmt.Errorf("error with get form: %v", err)
	}
	// Open the uploaded file
	fileReader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer fileReader.Close()

	fileExt := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%d%s", s.config.SubDir, time.Now().UnixMicro(), fileExt)
	info, err := s.client.PutObject(ctx, s.config.Bucket, filename, fileReader, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("error uploading to MinIO: %v", err)
	}

	return fmt.Sprintf("%s/%s/%s", s.config.Domain, s.config.Bucket, info.Key), nil
}

func (s MiniOss) Delete(fileURL string) error {
	objectName := filepath.Base(fileURL)
	key := fmt.Sprintf("%s/%s", s.config.SubDir, objectName)
	return s.client.RemoveObject(context.Background(), s.config.Bucket, key, minio.RemoveObjectOptions{})
}

var _ Uploader = MiniOss{}
