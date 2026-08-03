package sora

import (
	"encoding/json"
	"errors"
	"geekai/service/oss"
	"geekai/store/vo"
	"geekai/utils"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	logger2 "geekai/logger"
)

var logger = logger2.GetLogger()

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

	// 获取 JSON 数据
	resp, err := http.Get(videoDataURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析视频下载地址
	var videoData map[string]any
	err = json.Unmarshal(body, &videoData)
	if err != nil {
		return nil, err
	}

	if v, ok := videoData["url"].(string); ok && v != "" {
		logger.Infof("try to download video: %s", v)
		videoURL, err := s.uploadManager.GetUploadHandler().PutUrlFile(v, ".mp4", true)
		if err != nil {
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
