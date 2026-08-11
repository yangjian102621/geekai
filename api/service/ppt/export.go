package ppt

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"geekai/core/types"
	"image"
	_ "image/gif"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/ktye/pptx"
	_ "golang.org/x/image/webp"
)

//go:embed embed/minimal.pptx
var minimalPptxTemplate []byte

const (
	exportHTTPTimeout = 60 * time.Second
	// 必须与 embed/minimal.pptx 中 p:sldSz 一致（当前模板为 4:3：10\"×7.5\"）
	slideEmuW pptx.Dimension = 9144000
	slideEmuH pptx.Dimension = 6858000
)

// ExportFormat 导出类型
type ExportFormat string

const (
	ExportFormatPDF  ExportFormat = "pdf"
	ExportFormatPPTX ExportFormat = "pptx"
)

// ExportMimeType 返回 Content-Type
func ExportMimeType(f ExportFormat) string {
	switch f {
	case ExportFormatPDF:
		return "application/pdf"
	case ExportFormatPPTX:
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	default:
		return "application/octet-stream"
	}

}

// ExportFileExt 返回文件扩展名（含点）
func ExportFileExt(f ExportFormat) string {
	switch f {
	case ExportFormatPDF:
		return ".pdf"
	case ExportFormatPPTX:
		return ".pptx"
	default:
		return ""
	}
}

// ParseExportFormat 解析 query format
func ParseExportFormat(s string) (ExportFormat, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "pdf":
		return ExportFormatPDF, true
	case "pptx", "ppt":
		return ExportFormatPPTX, true
	default:
		return "", false
	}
}

// SanitizeExportBaseName 用于下载文件名的主体（不含扩展名）
func SanitizeExportBaseName(title, taskID string) string {
	s := strings.TrimSpace(title)
	repl := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_",
	)
	s = repl.Replace(s)
	var b strings.Builder
	for _, r := range s {
		if r == unicode.ReplacementChar || r < 32 {
			continue
		}
		b.WriteRune(r)
	}
	s = strings.TrimSpace(b.String())
	if s == "" {
		s = strings.TrimSpace(taskID)
	}
	if len([]rune(s)) > 120 {
		rs := []rune(s)
		s = string(rs[:120])
	}
	return s
}

// ContentDispositionAttachment RFC 5987，兼容旧客户端
func ContentDispositionAttachment(filename string) string {
	ascii := filename
	for _, r := range filename {
		if r > 127 || r == '"' || r == '\\' {
			ascii = "export" + strings.ToLower(filepath.Ext(filename))
			if ascii == "export" {
				ascii = "export.bin"
			}
			break
		}
	}
	return fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, ascii, url.PathEscape(filename))
}

// mapLocalUploadFile 将站点相对路径或完整 BaseURL 前缀映射为本地文件路径（local OSS）
func mapLocalUploadFile(raw string, local types.LocalStorageConfig) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || local.BasePath == "" {
		return "", false
	}
	bp := filepath.Clean(local.BasePath)
	bu := strings.TrimSuffix(strings.TrimSpace(local.BaseURL), "/")
	if bu != "" && strings.HasPrefix(raw, bu) {
		suffix := strings.TrimPrefix(strings.TrimPrefix(raw, bu), "/")
		return filepath.Join(bp, suffix), true
	}
	if strings.HasPrefix(raw, local.BaseURL) {
		return filepath.Join(bp, strings.TrimPrefix(raw, local.BaseURL)), true
	}
	return "", false
}

func originFromBaseURL(baseURL string) string {
	u, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

func defaultOriginFromListen(listen string) string {
	listen = strings.TrimSpace(listen)
	if listen == "" {
		return ""
	}
	host, port, err := net.SplitHostPort(listen)
	if err != nil {
		if strings.HasPrefix(listen, ":") {
			return "http://127.0.0.1" + listen
		}
		return ""
	}
	if host == "0.0.0.0" || host == "::" || host == "" {
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, port)
}

// resolveAbsoluteImageURL 将可能为相对路径的地址转为可 HTTP 访问的绝对 URL
func resolveAbsoluteImageURL(raw string, local types.LocalStorageConfig, app *types.AppConfig) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	if u, err := url.Parse(raw); err == nil && u.Scheme != "" && u.Host != "" {
		return raw
	}
	if strings.HasPrefix(raw, "//") {
		return "https:" + raw
	}
	origin := originFromBaseURL(local.BaseURL)
	if origin == "" && app != nil {
		origin = originFromBaseURL(app.StaticUrl)
	}
	if origin == "" && app != nil {
		origin = defaultOriginFromListen(app.Listen)
	}
	if origin == "" {
		return raw
	}
	if strings.HasPrefix(raw, "/") {
		return strings.TrimSuffix(origin, "/") + raw
	}
	return strings.TrimSuffix(origin, "/") + "/" + raw
}

// BuildExportBytes 按幻灯片顺序拉取图片并生成 PDF 或 PPTX（oss/app 用于解析相对路径图片 URL）
func BuildExportBytes(ctx context.Context, slides []SlideData, format ExportFormat, oss types.OSSConfig, app *types.AppConfig) ([]byte, error) {
	raws, err := fetchSlideImages(ctx, slides, oss, app)
	if err != nil {
		return nil, err
	}
	if len(raws) == 0 {
		return nil, fmt.Errorf("没有可导出的幻灯片图片")
	}
	switch format {
	case ExportFormatPDF:
		return buildPDF(raws)
	case ExportFormatPPTX:
		return buildPPTX(raws)
	default:
		return nil, fmt.Errorf("不支持的导出格式")
	}
}

func fetchSlideImages(ctx context.Context, slides []SlideData, oss types.OSSConfig, app *types.AppConfig) ([][]byte, error) {
	cp := append([]SlideData(nil), slides...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].SlideIndex < cp[j].SlideIndex })

	client := &http.Client{Timeout: exportHTTPTimeout}
	var out [][]byte
	for _, s := range cp {
		u := strings.TrimSpace(s.ImageURL)
		if u == "" {
			continue
		}

		var body []byte
		if oss.Active == "local" {
			if fp, ok := mapLocalUploadFile(u, oss.Local); ok {
				b, err := os.ReadFile(fp)
				if err == nil {
					if _, err := decodeImageBytes(b); err == nil {
						body = b
					}
				}
			}
		}
		if len(body) == 0 {
			absURL := resolveAbsoluteImageURL(u, oss.Local, app)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, absURL, nil)
			if err != nil {
				return nil, fmt.Errorf("幻灯片 %d: %w", s.SlideIndex, err)
			}
			resp, err := client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("幻灯片 %d 下载失败: %w", s.SlideIndex, err)
			}
			body, err = io.ReadAll(io.LimitReader(resp.Body, 32<<20))
			_ = resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("幻灯片 %d 读取失败: %w", s.SlideIndex, err)
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return nil, fmt.Errorf("幻灯片 %d 下载失败: HTTP %d", s.SlideIndex, resp.StatusCode)
			}
		}
		if _, err := decodeImageBytes(body); err != nil {
			return nil, fmt.Errorf("幻灯片 %d 不是有效图片: %v", s.SlideIndex, err)
		}
		out = append(out, body)
	}
	return out, nil
}

func decodeImageBytes(b []byte) (image.Image, error) {
	m, _, err := image.Decode(bytes.NewReader(b))
	return m, err
}

// fitSlideBoundsEMU 在幻灯片 EMU 框（与模板 p:sldSz 一致）内按原图比例 contain 居中（不裁切）
func fitSlideBoundsEMU(imgW, imgH int) (x, y, w, h pptx.Dimension) {
	if imgW <= 0 || imgH <= 0 {
		return 0, 0, slideEmuW, slideEmuH
	}
	sw := float64(slideEmuW)
	sh := float64(slideEmuH)
	iw := float64(imgW)
	ih := float64(imgH)
	scale := sw / iw
	if ih*scale > sh {
		scale = sh / ih
	}
	wf := iw * scale
	hf := ih * scale
	w = pptx.Dimension(wf + 0.5)
	h = pptx.Dimension(hf + 0.5)
	x = pptx.Dimension((sw-wf)*0.5 + 0.5)
	y = pptx.Dimension((sh-hf)*0.5 + 0.5)
	return x, y, w, h
}

func buildPDF(images [][]byte) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)
	// 像素 → mm（按 96 DPI），再按页面对比缩放以 contain 放入整页
	const pxPerMM = 96.0 / 25.4

	for i, raw := range images {
		im, err := decodeImageBytes(raw)
		if err != nil {
			return nil, fmt.Errorf("第 %d 页: %w", i+1, err)
		}
		b := im.Bounds()
		pxW := float64(b.Dx())
		pxH := float64(b.Dy())
		if pxW <= 0 || pxH <= 0 {
			return nil, fmt.Errorf("第 %d 页: 图片尺寸无效", i+1)
		}
		pdf.AddPage()
		pageW, pageH := pdf.GetPageSize()
		imgWmm := pxW / pxPerMM
		imgHmm := pxH / pxPerMM
		scale := pageW / imgWmm
		if imgHmm*scale > pageH {
			scale = pageH / imgHmm
		}
		w := imgWmm * scale
		h := imgHmm * scale
		x := (pageW - w) / 2
		y := (pageH - h) / 2

		name := fmt.Sprintf("slide%d", i)
		opt := gofpdf.ImageOptions{ReadDpi: false}
		tp := sniffImageType(raw)
		if tp != "" {
			opt.ImageType = tp
		}
		if pdf.RegisterImageOptionsReader(name, opt, bytes.NewReader(raw)) == nil {
			return nil, fmt.Errorf("第 %d 页: 无法写入 PDF 图片", i+1)
		}
		pdf.ImageOptions(name, x, y, w, h, false, opt, 0, "")
	}
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sniffImageType(b []byte) string {
	if len(b) < 12 {
		return ""
	}
	switch {
	case len(b) >= 2 && b[0] == 0xFF && b[1] == 0xD8:
		return "jpg"
	case len(b) >= 8 && string(b[0:8]) == "\x89PNG\r\n\x1a\n":
		return "png"
	case len(b) >= 6 && string(b[0:6]) == "GIF87a" || string(b[0:6]) == "GIF89a":
		return "gif"
	case len(b) >= 12 && string(b[0:4]) == "RIFF" && string(b[8:12]) == "WEBP":
		return "webp"
	default:
		return ""
	}
}

func buildPPTX(images [][]byte) ([]byte, error) {
	if len(minimalPptxTemplate) == 0 {
		return nil, fmt.Errorf("内置 PPT 模板缺失")
	}
	tmp, err := os.CreateTemp("", "ppt-export-*.pptx")
	if err != nil {
		return nil, err
	}
	path := tmp.Name()
	if _, err := tmp.Write(minimalPptxTemplate); err != nil {
		_ = tmp.Close()
		_ = os.Remove(path)
		return nil, err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(path)
		return nil, err
	}
	defer func() { _ = os.Remove(path) }()

	f, err := pptx.Open(path)
	if err != nil {
		return nil, err
	}

	for _, raw := range images {
		im, err := decodeImageBytes(raw)
		if err != nil {
			f.Abort()
			return nil, err
		}
		b := im.Bounds()
		ex, ey, ew, eh := fitSlideBoundsEMU(b.Dx(), b.Dy())
		slide := pptx.Slide{
			Images: []pptx.Image{
				pptx.NewImage(im, ex, ey, ew, eh),
			},
		}
		if err := f.Add(slide); err != nil {
			f.Abort()
			return nil, err
		}
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}
