package ppt

import (
	"encoding/base64"
	"fmt"
	"geekai/core/types"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// PrepareReferenceInputsForImg2Img 将幻灯片参考图转为第三方 API 可消费的输入：本地存储时读文件并转为 data URI（base64），公网 URL 原样传递。
func PrepareReferenceInputsForImg2Img(rawURL string, oss types.OSSConfig, app *types.AppConfig) ([]string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, fmt.Errorf("empty reference image url")
	}
	if strings.HasPrefix(rawURL, "data:") {
		return []string{rawURL}, nil
	}

	if oss.Active == "local" {
		if fp, ok := mapLocalUploadFile(rawURL, oss.Local); ok {
			b, err := os.ReadFile(fp)
			if err != nil {
				return nil, fmt.Errorf("read local reference image: %w", err)
			}
			if _, err := decodeImageBytes(b); err != nil {
				return nil, fmt.Errorf("reference is not a valid image: %w", err)
			}
			mime := http.DetectContentType(b)
			if !strings.HasPrefix(mime, "image/") {
				mime = "image/png"
			}
			dataURI := fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(b))
			return []string{dataURI}, nil
		}
	}

	if u, err := url.Parse(rawURL); err == nil && u.Scheme != "" && u.Host != "" {
		return []string{rawURL}, nil
	}

	abs := resolveAbsoluteImageURL(rawURL, oss.Local, app)
	if strings.TrimSpace(abs) == "" {
		return nil, fmt.Errorf("cannot resolve reference image url")
	}
	return []string{abs}, nil
}
