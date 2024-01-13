package clientv2

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

// 此处是为了版本兼容，sdk 支持最低版本为 go1.10, go1.13 提供 req.Clone 方法，
// 此处 copy 高版本的 go 标准库方法
func cloneReq(ctx context.Context, r *http.Request) *http.Request {
	if ctx == nil {
		panic("nil context")
	}

	r2 := r.WithContext(ctx)
	if r.Header != nil {
		r2.Header = cloneHeader(r.Header)
	}
	if r.Trailer != nil {
		r2.Trailer = cloneHeader(r.Trailer)
	}
	if s := r.TransferEncoding; s != nil {
		s2 := make([]string, len(s))
		copy(s2, s)
		r2.TransferEncoding = s2
	}
	r2.Form = cloneURLValues(r.Form)
	r2.PostForm = cloneURLValues(r.PostForm)
	r2.MultipartForm = cloneMultipartForm(r.MultipartForm)
	return r2
}

func cloneHeader(h http.Header) http.Header {
	if h == nil {
		return nil
	}

	// Find total number of values.
	nv := 0
	for _, vv := range h {
		nv += len(vv)
	}
	sv := make([]string, nv) // shared backing array for headers' values
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		n := copy(sv, vv)
		h2[k] = sv[:n:n]
		sv = sv[n:]
	}
	return h2
}

func cloneURLValues(v url.Values) url.Values {
	if v == nil {
		return nil
	}

	// http.Header and url.Values have the same representation, so temporarily
	// treat it like http.Header, which does have a clone:
	return url.Values(cloneHeader(http.Header(v)))
}

func cloneMultipartForm(f *multipart.Form) *multipart.Form {
	if f == nil {
		return nil
	}
	f2 := &multipart.Form{
		Value: (map[string][]string)(cloneHeader(http.Header(f.Value))),
	}
	if f.File != nil {
		m := make(map[string][]*multipart.FileHeader)
		for k, vv := range f.File {
			vv2 := make([]*multipart.FileHeader, len(vv))
			for i, v := range vv {
				vv2[i] = cloneMultipartFileHeader(v)
			}
			m[k] = vv2
		}
		f2.File = m
	}
	return f2
}

func cloneMultipartFileHeader(fh *multipart.FileHeader) *multipart.FileHeader {
	if fh == nil {
		return nil
	}
	fh2 := new(multipart.FileHeader)
	*fh2 = *fh
	fh2.Header = textproto.MIMEHeader(cloneHeader(http.Header(fh.Header)))
	return fh2
}
