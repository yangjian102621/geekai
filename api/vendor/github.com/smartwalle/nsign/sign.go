package nsign

import (
	"bytes"
	"crypto"
	"errors"
	"net/url"
	"sync"
)

var ErrVerification = errors.New("verification error")

type Option func(signer *signer)

func WithMethod(method Method) Option {
	return func(signer *signer) {
		if method == nil {
			return
		}
		signer.method = method
	}
}

func WithEncoder(encoder Encoder) Option {
	return func(signer *signer) {
		if encoder == nil {
			return
		}
		signer.encoder = encoder
	}
}

type SignOption func(opt *SignOptions)

type SignOptions struct {
	Prefix  string
	Suffix  string
	Ignores map[string]struct{}
}

func WithPrefix(s string) SignOption {
	return func(opt *SignOptions) {
		opt.Prefix = s
	}
}

func WithSuffix(s string) SignOption {
	return func(opt *SignOptions) {
		opt.Suffix = s
	}
}

func WithIgnore(keys ...string) SignOption {
	return func(opt *SignOptions) {
		if len(keys) > 0 && opt.Ignores == nil {
			opt.Ignores = make(map[string]struct{})
		}
		for _, key := range keys {
			if len(key) > 0 {
				opt.Ignores[key] = struct{}{}
			}
		}
	}
}

type Method interface {
	Sign(data []byte) ([]byte, error)

	Verify(data []byte, signature []byte) error
}

type Signer interface {
	SignValues(values url.Values, opts ...SignOption) ([]byte, error)

	SignBytes(data []byte, opts ...SignOption) ([]byte, error)

	VerifyValues(values url.Values, signature []byte, opts ...SignOption) error

	VerifyBytes(data []byte, signature []byte, opts ...SignOption) error
}

type signer struct {
	pool    *sync.Pool
	method  Method
	encoder Encoder
}

func New(opts ...Option) Signer {
	var s = &signer{}
	s.pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBufferString("")
		},
	}
	s.method = NewHashMethod(crypto.MD5)
	s.encoder = &DefaultEncoder{}

	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	return s
}

func (this *signer) getBuffer() *bytes.Buffer {
	var buffer = this.pool.Get().(*bytes.Buffer)
	buffer.Reset()
	return buffer
}

func (this *signer) putBuffer(buffer *bytes.Buffer) {
	if buffer != nil {
		buffer.Reset()
		this.pool.Put(buffer)
	}
}

func (this *signer) SignValues(values url.Values, opts ...SignOption) ([]byte, error) {
	var buffer = this.getBuffer()
	defer this.putBuffer(buffer)

	var nOptions = &SignOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(nOptions)
		}
	}

	var src, err = this.encoder.EncodeValues(buffer, values, nOptions)
	if err != nil {
		return nil, err
	}
	return this.method.Sign(src)
}

func (this *signer) SignBytes(data []byte, opts ...SignOption) ([]byte, error) {
	var buffer = this.getBuffer()
	defer this.putBuffer(buffer)

	var nOptions = &SignOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(nOptions)
		}
	}

	var src, err = this.encoder.EncodeBytes(buffer, data, nOptions)
	if err != nil {
		return nil, err
	}
	return this.method.Sign(src)
}

func (this *signer) VerifyValues(values url.Values, signature []byte, opts ...SignOption) error {
	var buffer = this.getBuffer()
	defer this.putBuffer(buffer)

	var nOptions = &SignOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(nOptions)
		}
	}

	var src, err = this.encoder.EncodeValues(buffer, values, nOptions)
	if err != nil {
		return err
	}
	return this.method.Verify(src, signature)
}

func (this *signer) VerifyBytes(data []byte, signature []byte, opts ...SignOption) error {
	var buffer = this.getBuffer()
	defer this.putBuffer(buffer)

	var nOptions = &SignOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(nOptions)
		}
	}

	var src, err = this.encoder.EncodeBytes(buffer, data, nOptions)
	if err != nil {
		return err
	}
	return this.method.Verify(src, signature)
}
