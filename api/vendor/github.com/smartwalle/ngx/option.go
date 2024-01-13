package ngx

import (
	"net/http"
	"net/url"
)

type Option func(req *Request)

func WithClient(c *http.Client) Option {
	return func(req *Request) {
		req.client = c
	}
}

func WithHeader(header http.Header) Option {
	return func(req *Request) {
		req.header = header
	}
}

func WithParams(params url.Values) Option {
	return func(req *Request) {
		req.params = params
	}
}

func WithQuery(query url.Values) Option {
	return func(req *Request) {
		for key, values := range query {
			for _, value := range values {
				req.query.Add(key, value)
			}
		}
	}
}

func WithBody(body Body) Option {
	return func(req *Request) {
		req.body = body
	}
}

func WithCookies(cookies []*http.Cookie) Option {
	return func(req *Request) {
		req.cookies = cookies
	}
}

// WithReceive 获取从服务端已接收数据大小
func WithReceive(f func(total, chunk, finished uint64)) Option {
	return func(req *Request) {
		req.receive = f
	}
}

// WithSend 获取向服务端已发送数据大小
func WithSend(f func(total, chunk, finished uint64)) Option {
	return func(req *Request) {
		req.send = f
	}
}
