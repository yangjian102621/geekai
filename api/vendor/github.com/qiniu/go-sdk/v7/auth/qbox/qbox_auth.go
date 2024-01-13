package qbox

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"net/http"
)

type Mac = auth.Credentials

// 兼容保留
func NewMac(accessKey, secretKey string) *Mac {
	return auth.New(accessKey, secretKey)
}

// Sign 一般用于下载凭证的签名
func Sign(mac *Mac, data []byte) string {
	return mac.Sign(data)
}

// SignWithData 一般用于上传凭证的签名
func SignWithData(mac *Mac, data []byte) string {
	return mac.SignWithData(data)
}

// VerifyCallback 验证上传回调请求是否来自七牛
func VerifyCallback(mac *Mac, req *http.Request) (bool, error) {
	return mac.VerifyCallback(req)
}
