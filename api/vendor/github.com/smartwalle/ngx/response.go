package ngx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Response struct {
	error error
	*http.Response
	data []byte
}

func (this *Response) Status() string {
	if this.Response != nil {
		return this.Response.Status
	}
	return fmt.Sprintf("%d ServiceUnavailable", http.StatusServiceUnavailable)
}

func (this *Response) StatusCode() int {
	if this.Response != nil {
		return this.Response.StatusCode
	}
	return http.StatusServiceUnavailable
}

func (this *Response) Proto() string {
	if this.Response != nil {
		return this.Response.Proto
	}
	return ""
}

func (this *Response) ProtoMajor() int {
	if this.Response != nil {
		return this.Response.ProtoMajor
	}
	return 1
}

func (this *Response) ProtoMinor() int {
	if this.Response != nil {
		return this.Response.ProtoMinor
	}
	return 0
}

func (this *Response) Header() http.Header {
	if this.Response != nil {
		return this.Response.Header
	}
	return http.Header{}
}

func (this *Response) ContentLength() int64 {
	if this.Response != nil {
		return this.Response.ContentLength
	}
	return 0
}

func (this *Response) TransferEncoding() []string {
	if this.Response != nil {
		return this.Response.TransferEncoding
	}
	return nil
}

func (this *Response) Close() bool {
	if this.Response != nil {
		return this.Response.Close
	}
	return true
}

func (this *Response) Uncompressed() bool {
	if this.Response != nil {
		return this.Response.Uncompressed
	}
	return true
}

func (this *Response) Trailer() http.Header {
	if this.Response != nil {
		return this.Response.Trailer
	}
	return http.Header{}
}

func (this *Response) Request() *http.Request {
	if this.Response != nil {
		return this.Response.Request
	}
	return nil
}

func (this *Response) TLS() *tls.ConnectionState {
	if this.Response != nil {
		return this.Response.TLS
	}
	return nil
}

func (this *Response) Cookies() []*http.Cookie {
	if this.Response != nil {
		return this.Response.Cookies()
	}
	return nil
}

func (this *Response) Location() (*url.URL, error) {
	if this.Response != nil {
		return this.Response.Location()
	}
	return nil, nil
}

func (this *Response) ProtoAtLeast(major, minor int) bool {
	if this.Response != nil {
		return this.Response.ProtoAtLeast(major, minor)
	}
	return false
}

func (this *Response) Write(w io.Writer) error {
	if this.Response != nil {
		return this.Response.Write(w)
	}
	return nil
}

func (this *Response) Error() error {
	return this.error
}

func (this *Response) Reader() io.Reader {
	return bytes.NewReader(this.data)
}

func (this *Response) Bytes() ([]byte, error) {
	return this.data, this.error
}

func (this *Response) MustBytes() []byte {
	return this.data
}

func (this *Response) String() (string, error) {
	return string(this.data), this.error
}

func (this *Response) MustString() string {
	return string(this.data)
}

func (this *Response) UnmarshalJSON(v interface{}) error {
	if this.error != nil {
		return this.error
	}
	return json.Unmarshal(this.data, v)
}
