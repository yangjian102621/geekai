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

func ReadFileContent(filePath string, tikaHost string) (string, error) {
	// for remote file, download it first
	if strings.HasPrefix(filePath, "http") {
		file, err := downloadFile(filePath)
		if err != nil {
			return "", err
		}
		filePath = file
	}
	// 创建 Tika 客户端
	client := tika.NewClient(nil, tikaHost)
	// 打开 PDF 文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error with open file: %v", err)
	}
	defer file.Close()

	// 使用 Tika 提取 PDF 文件的文本内容
	content, err := client.Parse(context.TODO(), file)
	if err != nil {
		return "", fmt.Errorf("error with parse file: %v", err)
	}

	ext := filepath.Ext(filePath)
	switch ext {
	case ".doc", ".docx", ".pdf", ".pptx", "ppt":
		return cleanBlankLine(cleanHtml(content, false)), nil
	case ".xls", ".xlsx":
		return cleanBlankLine(cleanHtml(content, true)), nil
	default:
		return cleanBlankLine(content), nil
	}

}

// 清理文本内容
func cleanHtml(html string, keepTable bool) string {
	// 清理 HTML 标签
	var policy *bluemonday.Policy
	if keepTable {
		policy = bluemonday.NewPolicy()
		policy.AllowElements("table", "thead", "tbody", "tfoot", "tr", "td", "th")
	} else {
		policy = bluemonday.StrictPolicy()
	}
	return policy.Sanitize(html)
}

func cleanBlankLine(content string) string {
	lines := strings.Split(content, "\n")
	texts := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 2 {
			continue
		}
		// discard image
		if strings.HasSuffix(line, ".png") ||
			strings.HasSuffix(line, ".jpg") ||
			strings.HasSuffix(line, ".jpeg") {
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
