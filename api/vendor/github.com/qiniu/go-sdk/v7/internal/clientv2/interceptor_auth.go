package clientv2

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"net/http"
)

type AuthConfig struct {
	Credentials auth.Credentials //
	TokenType   auth.TokenType   // 不包含上传
}

type authInterceptor struct {
	config AuthConfig
}

func NewAuthInterceptor(config AuthConfig) Interceptor {
	return &authInterceptor{
		config: config,
	}
}

func (interceptor *authInterceptor) Priority() InterceptorPriority {
	return InterceptorPriorityAuth
}

func (interceptor *authInterceptor) Intercept(req *http.Request, handler Handler) (*http.Response, error) {
	if interceptor == nil || req == nil {
		return handler(req)
	}

	err := interceptor.config.Credentials.AddToken(interceptor.config.TokenType, req)
	if err != nil {
		return nil, err
	}

	return handler(req)
}
