package service

import (
	"chatplus/core/types"
	"chatplus/utils"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
	"path"
	"strings"
)

type MinioService struct {
	config *types.AppConfig
	client *minio.Client
}

func NewMinioService(config *types.AppConfig) (*MinioService, error) {
	minioClient, err := minio.New(config.MinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioConfig.AccessKey, config.MinioConfig.AccessSecret, ""),
		Secure: config.MinioConfig.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioService{config: config, client: minioClient}, nil
}

func (s *MinioService) UploadMjImg(imageURL string) (string, error) {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return "", err
	}

	filename := path.Base(parsedURL.Path)
	imageBytes, err := utils.DownloadImage(imageURL, s.config.ProxyURL)
	if err != nil {
		return "", err
	}

	info, err := s.client.PutObject(
		context.Background(),
		s.config.MinioConfig.Bucket,
		filename,
		strings.NewReader(string(imageBytes)),
		int64(len(imageBytes)),
		minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s.config.MinioConfig.Domain, s.config.MinioConfig.Bucket, info.Key), nil
}
