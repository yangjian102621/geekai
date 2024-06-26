package utils

import (
	"context"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-tika/tika"
)

func ReadPdf(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "http") {
		file, err := downloadFile(filePath)
		if err != nil {
			return "", err
		}
		filePath = file
	}
	// 创建 Tika 客户端
	client := tika.NewClient(nil, "http://172.22.11.69:9998")
	// 打开 PDF 文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error with open file: %v", err)
	}
	defer file.Close()

	// 使用 Tika 提取 PDF 文件的文本内容
	html, err := client.Parse(context.TODO(), file)
	if err != nil {
		return "", fmt.Errorf("error with parse file: %v", err)
	}

	fmt.Println(html)

	return cleanBlankLine(html), nil

}

// 清理文本内容
func cleanHtml(html string) string {
	// 清理 HTML 标签
	p := bluemonday.StrictPolicy()
	return p.Sanitize(html)
}

func cleanBlankLine(content string) string {
	lines := strings.Split(content, "\n")
	texts := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue
		}
		texts = append(texts, line)
	}

	return strings.Join(texts, "\n")
}

// 下载文件
func downloadFile(url string) (string, error) {
	base := filepath.Base(url)
	dir := os.TempDir()
	filename := filepath.Join(dir, base)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 获取数据
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 写入数据到文件
	_, err = io.Copy(out, resp.Body)
	return filename, err
}
