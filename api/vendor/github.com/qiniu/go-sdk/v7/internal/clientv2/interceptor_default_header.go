package clientv2

import (
	clientV1 "github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/conf"
	"net/http"
	"time"
)

type defaultHeaderInterceptor struct {
}

func newDefaultHeaderInterceptor() Interceptor {
	return &defaultHeaderInterceptor{}
}

func (interceptor *defaultHeaderInterceptor) Priority() InterceptorPriority {
	return InterceptorPrioritySetHeader
}

func (interceptor *defaultHeaderInterceptor) Intercept(req *http.Request, handler Handler) (resp *http.Response, err error) {
	if interceptor == nil || req == nil {
		return handler(req)
	}

	if req.Header == nil {
		req.Header = http.Header{}
	}

	if e := addUseragent(req.Header); e != nil {
		return nil, e
	}

	if e := addXQiniuDate(req.Header); e != nil {
		return nil, e
	}

	return handler(req)
}

func addUseragent(headers http.Header) error {
	headers.Set("User-Agent", clientV1.UserAgent)
	return nil
}

func addXQiniuDate(headers http.Header) error {
	if conf.IsDisableQiniuTimestampSignature() {
		return nil
	}

	timeString := time.Now().UTC().Format("20060102T150405Z")
	headers.Set("X-Qiniu-Date", timeString)
	return nil
}
