package clientv2

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	RequestMethodGet    = "GET"
	RequestMethodPut    = "PUT"
	RequestMethodPost   = "POST"
	RequestMethodHead   = "HEAD"
	RequestMethodDelete = "DELETE"
)

type RequestBodyCreator func(options *RequestParams) (io.Reader, error)

func RequestBodyCreatorOfJson(object interface{}) RequestBodyCreator {
	body := object
	return func(o *RequestParams) (io.Reader, error) {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		o.Header.Add("Content-Type", "application/json")
		return bytes.NewReader(reqBody), nil
	}
}

func RequestBodyCreatorForm(info map[string][]string) RequestBodyCreator {
	body := FormStringInfo(info)
	return func(o *RequestParams) (io.Reader, error) {
		o.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		return bytes.NewBufferString(body), nil
	}
}

func FormStringInfo(info map[string][]string) string {
	if len(info) == 0 {
		return ""
	}
	return url.Values(info).Encode()
}

type RequestParams struct {
	Context     context.Context
	Method      string
	Url         string
	Header      http.Header
	BodyCreator RequestBodyCreator
}

func (o *RequestParams) init() {
	if o.Context == nil {
		o.Context = context.Background()
	}

	if len(o.Method) == 0 {
		o.Method = RequestMethodGet
	}

	if o.Header == nil {
		o.Header = http.Header{}
	}

	if o.BodyCreator == nil {
		o.BodyCreator = func(options *RequestParams) (io.Reader, error) {
			return nil, nil
		}
	}
}

func NewRequest(options RequestParams) (*http.Request, error) {
	options.init()

	body, cErr := options.BodyCreator(&options)
	if cErr != nil {
		return nil, cErr
	}

	req, err := http.NewRequest(options.Method, options.Url, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(options.Context)
	req.Header = options.Header
	return req, nil
}
