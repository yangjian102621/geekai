package utils

import (
	"context"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"net/http"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	u, err := urlpkg.Parse(url)
	if err != nil {
		return "", err
	}

	base := filepath.Base(u.Path)
	if base == "" || base == "." || base == "/" {
		base = "download"
	}
	dir := os.TempDir()
	filename := filepath.Join(dir, base)

	client := newHTTPClient(60*time.Second, "")
	retries := 2
	var lastErr error
	for attempt := 0; attempt <= retries; attempt++ {
		// 每次重试都重新创建文件，避免写一半残留
		out, err := os.Create(filename)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		if err != nil {
			_ = out.Close()
			return "", err
		}

		resp, err := client.Do(req)
		if err != nil {
			_ = out.Close()
			lastErr = err
			_ = os.Remove(filename)
			if attempt < retries && isRetryableError(err) {
				time.Sleep(retryDelay(attempt))
				continue
			}
			return "", err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			_ = resp.Body.Close()
			_ = out.Close()
			_ = os.Remove(filename)
			lastErr = fmt.Errorf("download file failed: status=%d", resp.StatusCode)
			return "", lastErr
		}

		_, copyErr := io.Copy(out, resp.Body)
		_ = resp.Body.Close()
		_ = out.Close()
		if copyErr != nil {
			lastErr = copyErr
			_ = os.Remove(filename)
			if attempt < retries && isRetryableError(copyErr) {
				time.Sleep(retryDelay(attempt))
				continue
			}
			return "", copyErr
		}

		return filename, nil
	}

	return "", lastErr
}
