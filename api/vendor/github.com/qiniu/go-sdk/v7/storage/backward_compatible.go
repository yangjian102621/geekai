// 原来rpc.go包含了客户端的信息，这个部分被调整到了"github.com/qiniu/go-sdk/v7/client"
// 这个文件的内容不应该再被使用
// 客户端应该是所有服务公用的，包括kodo, cdn, dora, atlab等，不应该放在storage下

// 这个文件兼容保留了原来storage暴露出去的变量，函数等
package storage

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/conf"
	"runtime"
)

var DefaultClient = client.DefaultClient
var UserAgent = client.UserAgent

type Client = client.Client
type ErrorInfo = client.ErrorInfo

var ResponseError = client.ResponseError
var CallRet = client.CallRet

// var SetAppName = client.SetAppName
// SetAppName设置的是全局的变量，如果再这个包引入var SetAppName， 那么设置的实际上是
// client包中的UserAgent， 所以为了兼容性重复定义了该函数

// userApp should be [A-Za-z0-9_\ \-\.]*
func SetAppName(userApp string) error {
	UserAgent = fmt.Sprintf(
		"QiniuGo/%s (%s; %s; %s) %s", conf.Version, runtime.GOOS, runtime.GOARCH, userApp, runtime.Version())
	return nil
}
