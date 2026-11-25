package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const downloadWriteTimeout = 2 * time.Hour

// updateManifest 描述 update 目录中的安装包清单。
// 顶层 version/url 描述 Windows 安装包；mac 下的 arm64 / x86_64 分别描述两个平台。
type updateManifest struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	Mac     map[string]struct {
		Version string `json:"version"`
		URL     string `json:"url"`
	} `json:"mac"`
}

// DownloadHandler 提供官网安装包下载与自动更新清单。
// 安装包和更新清单存放在 static/update，可通过 UPDATE_DIR 指定目录。
type DownloadHandler struct {
	BaseHandler
	updateDir string
}

func NewDownloadHandler(app *core.AppServer) *DownloadHandler {
	dir := strings.TrimSpace(os.Getenv("UPDATE_DIR"))
	if dir == "" {
		dir = filepath.Join("static", "update")
	}
	return &DownloadHandler{
		BaseHandler: BaseHandler{App: app},
		updateDir:   dir,
	}
}

func (h *DownloadHandler) RegisterRoutes(r *core.Router) {
	// 更新接口在登录前可用，旧版调用继续使用原地址。
	for _, prefix := range []string{"/api/update", "/api/v1"} {
		group := r.PublicGroup(prefix)
		group.GET("/latest", h.Latest)
		group.GET("/download", withDownloadWriteDeadline(), h.Download)
	}

	// 安装包与清单直接由后端静态目录提供，nginx 的 /update/ 会转发到这里。
	r.StaticFiles("/update", h.updateDir, withDownloadWriteDeadline(), withUpdateCacheControl())
}

func (h *DownloadHandler) readManifest() (updateManifest, error) {
	data, err := os.ReadFile(filepath.Join(h.updateDir, "latest.json"))
	if err != nil {
		return updateManifest{}, err
	}
	var manifest updateManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return updateManifest{}, err
	}
	return manifest, nil
}

func manifestFileName(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" {
		return "", errors.New("manifest url missing")
	}
	name := filepath.Base(u.Path)
	if name == "" || name == "." || name == "/" {
		return "", errors.New("invalid manifest url")
	}
	return name, nil
}

func (h *DownloadHandler) windowsInstallerFile() (string, error) {
	manifest, err := h.readManifest()
	if err != nil {
		return "", err
	}
	return manifestFileName(manifest.URL)
}

// macInstallerFile 只认清单中的 mac.<arch>.url，不按目录猜测版本。
func (h *DownloadHandler) macInstallerFile(arch string) (string, error) {
	manifest, err := h.readManifest()
	if err != nil {
		return "", err
	}
	entry, ok := manifest.Mac[arch]
	if !ok || entry.URL == "" {
		return "", fmt.Errorf("manifest missing mac.%s", arch)
	}
	return manifestFileName(entry.URL)
}

func (h *DownloadHandler) installerFile(platform string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case "", "windows":
		return h.windowsInstallerFile()
	case "mac-arm64", "macos-arm64", "mac-aarch64":
		return h.macInstallerFile("arm64")
	case "mac-x86_64", "mac-x64", "mac-intel", "macos-x86_64":
		return h.macInstallerFile("x86_64")
	default:
		return "", errors.New("unsupported platform")
	}
}

func (h *DownloadHandler) serveInstaller(c *gin.Context) {
	platform := strings.TrimSpace(c.Query("platform"))
	fileName, err := h.installerFile(platform)
	if err != nil {
		log.Printf("level=ERROR event=download_setup_missing dir=%q platform=%q error=%q", h.updateDir, platform, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "setup_not_found"})
		return
	}
	setupPath := filepath.Join(h.updateDir, fileName)
	if info, statErr := os.Stat(setupPath); statErr != nil || info.IsDir() {
		log.Printf("level=ERROR event=download_setup_missing dir=%q platform=%q file=%q", h.updateDir, platform, fileName)
		c.JSON(http.StatusNotFound, gin.H{"error": "setup_not_found"})
		return
	}
	c.Header("Content-Disposition", `attachment; filename="`+fileName+`"`)
	c.File(setupPath)
}

// Download 官网下载入口：默认 Windows，可用 platform 指定 macOS 架构。
func (h *DownloadHandler) Download(c *gin.Context) {
	h.serveInstaller(c)
}

// Latest 返回当前安装包版本与各平台下载路径。
// 清单中没有配置或文件不存在的平台不会输出，前端据此显示“暂未提供”。
func (h *DownloadHandler) Latest(c *gin.Context) {
	manifest, err := h.readManifest()
	if err != nil {
		log.Printf("level=ERROR event=latest_manifest_invalid error=%q", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "manifest_invalid"})
		return
	}
	downloadURL := manifest.URL
	if u, parseErr := url.Parse(manifest.URL); parseErr == nil && u.Path != "" {
		downloadURL = u.Path
	}
	payload := gin.H{
		"version": manifest.Version,
		"url":     downloadURL,
	}
	platforms := gin.H{}
	for _, arch := range []string{"arm64", "x86_64"} {
		key := "mac-arm64"
		if arch == "x86_64" {
			key = "mac-x64"
		}
		entry, ok := manifest.Mac[arch]
		if !ok || entry.URL == "" {
			continue
		}
		u, parseErr := url.Parse(entry.URL)
		if parseErr != nil || u.Path == "" {
			continue
		}
		if info, statErr := os.Stat(filepath.Join(h.updateDir, filepath.Base(u.Path))); statErr != nil || info.IsDir() {
			continue
		}
		platforms[key] = u.Path
	}
	if len(platforms) > 0 {
		payload["platforms"] = platforms
	}
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.JSON(http.StatusOK, payload)
}

func withDownloadWriteDeadline() gin.HandlerFunc {
	return func(c *gin.Context) {
		controller := http.NewResponseController(c.Writer)
		if err := controller.SetWriteDeadline(time.Now().Add(downloadWriteTimeout)); err != nil && !errors.Is(err, http.ErrNotSupported) {
			log.Printf("level=WARN event=set_download_write_deadline_failed path=%q error=%q", c.Request.URL.Path, err)
		}
		c.Next()
	}
}

// withUpdateCacheControl 为安装包设置缓存：带版本号的文件可长缓存，latest.json 禁止缓存。
func withUpdateCacheControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		switch {
		case strings.HasSuffix(path, "latest.json"):
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		case strings.HasSuffix(path, ".exe"), strings.HasSuffix(path, ".dmg"):
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	}
}
