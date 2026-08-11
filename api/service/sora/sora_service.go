package sora

import (
	"context"
	"encoding/json"
	"errors"
	"geekai/service/oss"
	"geekai/store/vo"
	"geekai/utils"
	"path/filepath"
	"regexp"
	"time"

	"geekai/log"
)

var logger = log.GetLogger()

type SoraService struct {
	uploadManager *oss.UploaderManager
}

func NewSoraService(uploadManager *oss.UploaderManager) *SoraService {
	return &SoraService{
		uploadManager: uploadManager,
	}
}

// 下载视频地址
func (s *SoraService) DownloadVideoURL(text string) (*vo.File, error) {
	videoDataURL, err := s.ExtractVideoURL(text)
	if err != nil {
		return nil, err
	}

	// 用统一的超时/重试策略，避免“偶发 HTTPS 握手超时”直接导致失败
	body, _, err := utils.FetchURLBytes(context.Background(), videoDataURL, "", 30*time.Second, 2, 2<<20)
	if err != nil {
		logger.Errorf("failed to get video data: %v", err)
		return nil, err
	}

	// 解析视频下载地址
	var videoData map[string]any
	err = json.Unmarshal(body, &videoData)
	if err != nil {
		logger.Errorf("failed to unmarshal video data: %v", err)
		return nil, err
	}

	if v, ok := videoData["url"].(string); ok && v != "" {
		logger.Infof("try to download video: %s", v)
		videoURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(v, ".mp4", true)
		if err != nil { // 如果上传失败，则返回原始错误
			return nil, err
		}

		// 获取文件大小
		size, _ := utils.GetFileSize(videoURL)
		name := filepath.Base(videoURL)
		return &vo.File{
			Name:      name,
			ObjKey:    name,
			URL:       videoURL,
			Ext:       ".mp4",
			Size:      size,
			CreatedAt: time.Now().Unix(),
		}, nil
	}

	return nil, errors.New("no video URL found: " + string(body))

}

// 从文本中提取视频URL
func (s *SoraService) ExtractVideoURL(text string) (string, error) {
	// 提取原始 JSON 数据地址
	//[原始数据](https://asyncdata.net/source/task_01k8pye324ef7t6heq6jyaxbfe) 类似这样的链接
	re := regexp.MustCompile(`\[原始数据\]\((https?://.*?)\)`)
	matches := re.FindAllStringSubmatch(text, 1)
	if len(matches) == 0 {
		return "", errors.New("no video URL found")
	}
	return matches[0][1], nil
}
