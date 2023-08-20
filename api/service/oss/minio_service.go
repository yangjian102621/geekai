package oss

import (
	"chatplus/core/types"
	"chatplus/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"path/filepath"
	"strings"
	"time"
)

type MinioService struct {
	config   *types.MinioConfig
	client   *minio.Client
	proxyURL string
}

func NewMinioService(appConfig *types.AppConfig) (MinioService, error) {
	config := &appConfig.OSS.Minio
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.AccessSecret, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return MinioService{}, err
	}
	return MinioService{config: config, client: minioClient, proxyURL: appConfig.ProxyURL}, nil
}

func (s MinioService) PutImg(imageURL string) (string, error) {
	imageData, err := utils.DownloadImage(imageURL, s.proxyURL)
	if err != nil {
		return "", fmt.Errorf("error with download image: %v", err)
	}
	fileExt := filepath.Ext(filepath.Base(imageURL))
	filename := fmt.Sprintf("%d%s", time.Now().UnixMicro(), fileExt)
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

func (s MinioService) PutFile(ctx *gin.Context, name string) (string, error) {
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
	filename := fmt.Sprintf("%d%s", time.Now().UnixMicro(), fileExt)
	info, err := s.client.PutObject(ctx, s.config.Bucket, filename, fileReader, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("error uploading to MinIO: %v", err)
	}

	return fmt.Sprintf("%s/%s/%s", s.config.Domain, s.config.Bucket, info.Key), nil
}

func (s MinioService) Delete(fileURL string) error {
	objectName := filepath.Base(fileURL)
	return s.client.RemoveObject(context.Background(), s.config.Bucket, objectName, minio.RemoveObjectOptions{})
}

var _ Uploader = MinioService{}
