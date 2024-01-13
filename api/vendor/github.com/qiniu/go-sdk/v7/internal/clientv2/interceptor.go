package clientv2

import (
	"net/http"
)

const (
	InterceptorPriorityDefault     InterceptorPriority = 100
	InterceptorPriorityRetryHosts  InterceptorPriority = 200
	InterceptorPriorityRetrySimple InterceptorPriority = 300
	InterceptorPrioritySetHeader   InterceptorPriority = 400
	InterceptorPriorityNormal      InterceptorPriority = 500
	InterceptorPriorityAuth        InterceptorPriority = 600
	InterceptorPriorityDebug       InterceptorPriority = 700
)

type InterceptorPriority int

type Interceptor interface {
	// Priority 数字越小优先级越高
	Priority() InterceptorPriority

	// Intercept 拦截处理函数
	Intercept(req *http.Request, handler Handler) (*http.Response, error)
}

type interceptorList []Interceptor

func (l interceptorList) Less(i, j int) bool {
	return l[i].Priority() < l[j].Priority()
}

func (l interceptorList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l interceptorList) Len() int {
	return len(l)
}

type simpleInterceptor struct {
	priority InterceptorPriority
	handler  func(req *http.Request, handler Handler) (*http.Response, error)
}

func NewSimpleInterceptor(interceptorHandler func(req *http.Request, handler Handler) (*http.Response, error)) Interceptor {
	return NewSimpleInterceptorWithPriority(InterceptorPriorityNormal, interceptorHandler)
}

func NewSimpleInterceptorWithPriority(priority InterceptorPriority, interceptorHandler func(req *http.Request, handler Handler) (*http.Response, error)) Interceptor {
	if priority <= 0 {
		priority = InterceptorPriorityNormal
	}

	return &simpleInterceptor{
		priority: priority,
		handler:  interceptorHandler,
	}
}

func (interceptor *simpleInterceptor) Priority() InterceptorPriority {
	return interceptor.priority
}

func (interceptor *simpleInterceptor) Intercept(req *http.Request, handler Handler) (*http.Response, error) {
	if interceptor == nil || interceptor.handler == nil {
		return handler(req)
	}
	return interceptor.handler(req, handler)
}
