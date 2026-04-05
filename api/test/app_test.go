package test

import (
	"geekai/utils"
	"testing"
)

// TestNewService 测试创建爬虫服务
func TestNewService(t *testing.T) {
	videoURL := `https://p3-aiop-sign.byteimg.com/tos-cn-i-vuqhorh59i/2025072310444223AAB2C93CE2B9BB8573-6843-0~tplv-vuqhorh59i-image.image?rk3s=7f9e702d&x-expires=1753325083&x-signature=%2F5V3H%2FWPQlOej6VtVZyf%2BNJBWok%3D`
	filePath := "test_video.png"
	err := utils.DownloadFile(videoURL, filePath, "")
	if err != nil {
		t.Fatalf("下载视频失败: %v", err)
	}
}
