package middleware

import (
	"bytes"
	"geekai/utils"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
)

// 静态资源中间件
func StaticMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		url := c.Request.URL.String()
		// 拦截生成缩略图请求
		if strings.HasPrefix(url, "/static/") && strings.Contains(url, "?imageView2") {
			r := strings.SplitAfter(url, "imageView2")
			size := strings.Split(r[1], "/")
			if len(size) != 8 {
				c.String(http.StatusNotFound, "invalid thumb args")
				return
			}
			with := utils.IntValue(size[3], 0)
			height := utils.IntValue(size[5], 0)
			quality := utils.IntValue(size[7], 75)

			// 打开图片文件
			filePath := strings.TrimLeft(c.Request.URL.Path, "/")
			file, err := os.Open(filePath)
			if err != nil {
				c.String(http.StatusNotFound, "Image not found")
				return
			}
			defer file.Close()

			// 解码图片
			img, _, err := image.Decode(file)
			// for .webp image
			if err != nil {
				img, err = webp.Decode(file)
			}
			if err != nil {
				c.String(http.StatusInternalServerError, "Error decoding image")
				return
			}

			var newImg image.Image
			if height == 0 || with == 0 {
				// 固定宽度，高度自适应
				newImg = resize.Resize(uint(with), uint(height), img, resize.Lanczos3)
			} else {
				// 生成缩略图
				newImg = resize.Thumbnail(uint(with), uint(height), img, resize.Lanczos3)
			}
			var buffer bytes.Buffer
			err = jpeg.Encode(&buffer, newImg, &jpeg.Options{Quality: quality})
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			// 设置图片缓存有效期为一年 (365天)
			c.Header("Cache-Control", "max-age=31536000, public")
			// 直接输出图像数据流
			c.Data(http.StatusOK, "image/jpeg", buffer.Bytes())
			c.Abort() // 中断请求

		}
		c.Next()
	}
}
