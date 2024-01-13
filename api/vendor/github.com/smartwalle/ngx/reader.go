package ngx

import "io"

type Reader interface {
	io.Reader

	Len() int
}

type reader struct {
	Body
	handler  func(total, chunk, finished uint64)
	total    uint64
	finished uint64
}

func NewReader(r Reader, handler func(total, chunk, finished uint64)) *reader {
	return &reader{Body: r, total: uint64(r.Len()), finished: 0, handler: handler}
}

func (this *reader) Read(p []byte) (n int, err error) {
	n, err = this.Body.Read(p)

	if n > 0 {
		var chunk = uint64(n)
		this.finished += chunk
		if this.handler != nil {
			this.handler(this.total, chunk, this.finished)
		}
	}

	return n, err
}
