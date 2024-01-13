package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"runtime"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/conf"
	"github.com/qiniu/go-sdk/v7/internal/log"
	"github.com/qiniu/go-sdk/v7/reqid"
)

var UserAgent = getUserAgentWithAppName("default")
var DefaultClient = Client{&http.Client{Transport: http.DefaultTransport}}

// 用来打印调试信息
var DebugMode = false
var DeepDebugInfo = false

// --------------------------------------------------------------------

// Client 负责发送HTTP请求到七牛接口服务器
type Client struct {
	*http.Client
}

// TurnOnDebug 开启Debug模式
func TurnOnDebug() {
	DebugMode = true
}

// userApp should be [A-Za-z0-9_\ \-\.]*
func SetAppName(userApp string) error {
	UserAgent = getUserAgentWithAppName(userApp)
	return nil
}

func getUserAgentWithAppName(userApp string) string {
	return fmt.Sprintf("QiniuGo/%s (%s; %s; %s) %s",
		conf.Version, runtime.GOOS, runtime.GOARCH, userApp, runtime.Version())
}

// --------------------------------------------------------------------

func newRequest(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, reqUrl, body)
	if err != nil {
		return
	}

	if headers == nil {
		headers = http.Header{}
	}

	err = addDefaultHeader(headers)
	if err != nil {
		return
	}

	req.Header = headers
	req = req.WithContext(ctx)

	//check access token
	mac, t, ok := auth.CredentialsFromContext(ctx)
	if ok {
		err = mac.AddToken(t, req)
		if err != nil {
			return
		}
	}
	if DebugMode {
		trace := &httptrace.ClientTrace{
			GotConn: func(connInfo httptrace.GotConnInfo) {
				remoteAddr := connInfo.Conn.RemoteAddr()
				log.Debug(fmt.Sprintf("Network: %s, Remote ip:%s, URL: %s", remoteAddr.Network(), remoteAddr.String(), req.URL))
			},
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		bs, bErr := httputil.DumpRequest(req, DeepDebugInfo)
		if bErr != nil {
			err = bErr
			return
		}
		log.Debug(string(bs))
	}
	return
}

func (r Client) DoRequest(ctx context.Context, method, reqUrl string, headers http.Header) (resp *http.Response, err error) {
	req, err := newRequest(ctx, method, reqUrl, headers, nil)
	if err != nil {
		return
	}
	return r.Do(ctx, req)
}

func (r Client) DoRequestWith(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int) (resp *http.Response, err error) {

	req, err := newRequest(ctx, method, reqUrl, headers, body)
	if err != nil {
		return
	}
	req.ContentLength = int64(bodyLength)
	return r.Do(ctx, req)
}

func (r Client) DoRequestWith64(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int64) (resp *http.Response, err error) {

	req, err := newRequest(ctx, method, reqUrl, headers, body)
	if err != nil {
		return
	}
	req.ContentLength = bodyLength
	return r.Do(ctx, req)
}

func (r Client) DoRequestWithBodyGetter(ctx context.Context, method, reqUrl string, headers http.Header, body io.Reader,
	getBody func() (io.ReadCloser, error), bodyLength int64) (resp *http.Response, err error) {

	req, err := newRequest(ctx, method, reqUrl, headers, body)
	if err != nil {
		return
	}
	req.ContentLength = bodyLength
	req.GetBody = getBody
	return r.Do(ctx, req)
}

func (r Client) DoRequestWithForm(ctx context.Context, method, reqUrl string, headers http.Header,
	data map[string][]string) (resp *http.Response, err error) {

	if headers == nil {
		headers = http.Header{}
	}
	headers.Add("Content-Type", "application/x-www-form-urlencoded")

	requestData := url.Values(data).Encode()
	if method == "GET" || method == "HEAD" || method == "DELETE" {
		if strings.ContainsRune(reqUrl, '?') {
			reqUrl += "&"
		} else {
			reqUrl += "?"
		}
		return r.DoRequest(ctx, method, reqUrl+requestData, headers)
	}

	return r.DoRequestWith(ctx, method, reqUrl, headers, strings.NewReader(requestData), len(requestData))
}

func (r Client) DoRequestWithJson(ctx context.Context, method, reqUrl string, headers http.Header,
	data interface{}) (resp *http.Response, err error) {

	reqBody, err := json.Marshal(data)
	if err != nil {
		return
	}

	if headers == nil {
		headers = http.Header{}
	}
	headers.Add("Content-Type", "application/json")
	return r.DoRequestWith(ctx, method, reqUrl, headers, bytes.NewReader(reqBody), len(reqBody))
}

func (r Client) Do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	reqctx := req.Context()

	if reqId, ok := reqid.ReqidFromContext(ctx); ok {
		req.Header.Set("X-Reqid", reqId)
	} else if reqId, ok = reqid.ReqidFromContext(reqctx); ok {
		req.Header.Set("X-Reqid", reqId)
	}

	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", UserAgent)
	}

	resp, err = r.Client.Do(req)
	return
}

// --------------------------------------------------------------------

type ErrorInfo struct {
	Err   string `json:"error,omitempty"`
	Key   string `json:"key,omitempty"`
	Reqid string `json:"reqid,omitempty"`
	Errno int    `json:"errno,omitempty"`
	Code  int    `json:"code"`
}

func (r *ErrorInfo) ErrorDetail() string {

	msg, _ := json.Marshal(r)
	return string(msg)
}

func (r *ErrorInfo) Error() string {

	return r.Err
}

func (r *ErrorInfo) RpcError() (code, errno int, key, err string) {

	return r.Code, r.Errno, r.Key, r.Err
}

func (r *ErrorInfo) HttpCode() int {

	return r.Code
}

// --------------------------------------------------------------------

func parseError(e *ErrorInfo, r io.Reader) {

	body, err1 := ioutil.ReadAll(r)
	if err1 != nil {
		e.Err = err1.Error()
		return
	}

	var ret struct {
		Err   string `json:"error"`
		Key   string `json:"key"`
		Errno int    `json:"errno"`
	}
	if decodeJsonFromData(body, &ret) == nil && ret.Err != "" {
		// qiniu error msg style returns here
		e.Err, e.Key, e.Errno = ret.Err, ret.Key, ret.Errno
		return
	}
	e.Err = string(body)
}

func ResponseError(resp *http.Response) (err error) {

	e := &ErrorInfo{
		Reqid: resp.Header.Get("X-Reqid"),
		Code:  resp.StatusCode,
	}
	if resp.StatusCode > 299 {
		if resp.ContentLength != 0 {
			ct, ok := resp.Header["Content-Type"]
			if ok && strings.HasPrefix(ct[0], "application/json") {
				parseError(e, resp.Body)
			} else {
				bs, rErr := ioutil.ReadAll(resp.Body)
				if rErr != nil {
					err = rErr
				}
				e.Err = strings.TrimRight(string(bs), "\n")
			}
		}
	}
	return e
}

func CallRet(ctx context.Context, ret interface{}, resp *http.Response) (err error) {

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if DebugMode {
		bs, dErr := httputil.DumpResponse(resp, DeepDebugInfo)
		if dErr != nil {
			err = dErr
			return
		}
		log.Debug(string(bs))
	}
	if resp.StatusCode/100 == 2 {
		if ret != nil && resp.ContentLength != 0 {
			err = DecodeJsonFromReader(resp.Body, ret)
			if err != nil {
				return
			}
		}
		return nil
	}
	return ResponseError(resp)
}

func (r Client) CallWithForm(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header,
	param map[string][]string) (err error) {

	resp, err := r.DoRequestWithForm(ctx, method, reqUrl, headers, param)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CallWithJson(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header,
	param interface{}) (err error) {

	resp, err := r.DoRequestWithJson(ctx, method, reqUrl, headers, param)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CallWith(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int) (err error) {

	resp, err := r.DoRequestWith(ctx, method, reqUrl, headers, body, bodyLength)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CallWith64(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int64) (err error) {

	resp, err := r.DoRequestWith64(ctx, method, reqUrl, headers, body, bodyLength)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CallWithBodyGetter(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header, body io.Reader,
	getBody func() (io.ReadCloser, error), bodyLength int64) (err error) {

	resp, err := r.DoRequestWithBodyGetter(ctx, method, reqUrl, headers, body, getBody, bodyLength)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) Call(ctx context.Context, ret interface{}, method, reqUrl string, headers http.Header) (err error) {

	resp, err := r.DoRequestWith(ctx, method, reqUrl, headers, nil, 0)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CredentialedCallWithForm(ctx context.Context, cred *auth.Credentials, tokenType auth.TokenType, ret interface{},
	method, reqUrl string, headers http.Header, param map[string][]string) error {
	ctx = auth.WithCredentialsType(ctx, cred, tokenType)
	return r.CallWithForm(ctx, ret, method, reqUrl, headers, param)
}

func (r Client) CredentialedCallWithJson(ctx context.Context, cred *auth.Credentials, tokenType auth.TokenType, ret interface{},
	method, reqUrl string, headers http.Header, param interface{}) error {
	ctx = auth.WithCredentialsType(ctx, cred, tokenType)
	return r.CallWithJson(ctx, ret, method, reqUrl, headers, param)
}

func (r Client) CredentialedCallWith(ctx context.Context, cred *auth.Credentials, tokenType auth.TokenType, ret interface{},
	method, reqUrl string, headers http.Header, body io.Reader, bodyLength int) error {
	ctx = auth.WithCredentialsType(ctx, cred, tokenType)
	return r.CallWith(ctx, ret, method, reqUrl, headers, body, bodyLength)
}

func (r Client) CredentialedCallWith64(ctx context.Context, cred *auth.Credentials, tokenType auth.TokenType, ret interface{},
	method, reqUrl string, headers http.Header, body io.Reader, bodyLength int64) error {
	ctx = auth.WithCredentialsType(ctx, cred, tokenType)
	return r.CallWith64(ctx, ret, method, reqUrl, headers, body, bodyLength)
}

func (r Client) CredentialedCall(ctx context.Context, cred *auth.Credentials, tokenType auth.TokenType, ret interface{},
	method, reqUrl string, headers http.Header) error {
	ctx = auth.WithCredentialsType(ctx, cred, tokenType)
	return r.Call(ctx, ret, method, reqUrl, headers)
}
