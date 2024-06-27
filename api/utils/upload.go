package utils

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// GenUploadPath 生成上传文件路径
func GenUploadPath(basePath, filename string, isImg bool) (string, error) {
	now := time.Now()
	dir := fmt.Sprintf("%s/%d/%d", basePath, now.Year(), now.Month())
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("error with create upload dir：%v", err)
		}
	}
	var fileExt string
	if isImg {
		fileExt = GetImgExt(filename)
	} else {
		fileExt = filepath.Ext(filename)
	}
	return fmt.Sprintf("%s/%d%s", dir, now.UnixMicro(), fileExt), nil
}

// GenUploadUrl 生成上传文件 URL
func GenUploadUrl(basePath, baseUrl string, filePath string) string {
	return strings.Replace(filePath, basePath, baseUrl, 1)
}

func DownloadFile(fileURL string, filepath string, proxy string) error {
	var client *http.Client
	if proxy == "" {
		client = http.DefaultClient
	} else {
		proxyURL, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetImgExt(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".png"
	}
	return ext
}

func ExtractImgURLs(text string) []string {
	re := regexp.MustCompile(`(http[s]?:\/\/.*?\.(?:png|jpg|jpeg|gif))`)
	matches := re.FindAllStringSubmatch(text, 10)
	urls := make([]string, 0)
	if len(matches) > 0 {
		for _, m := range matches {
			urls = append(urls, m[1])
		}
	}
	return urls
}

func ExtractFileURLs(text string) []string {
	re := regexp.MustCompile(`(http[s]?:\/\/.*?\.(?:docx?|pdf|pptx?|xlsx?|txt))`)
	matches := re.FindAllStringSubmatch(text, 10)
	urls := make([]string, 0)
	if len(matches) > 0 {
		for _, m := range matches {
			urls = append(urls, m[1])
		}
	}
	return urls
}
