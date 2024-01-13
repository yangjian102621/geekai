package ngx

import (
	"io"
)

type writer struct {
	io.Writer
	handler  func(total, chunk, finished uint64)
	finished uint64
	total    uint64
}

func NewWriter(w io.Writer, total uint64, handler func(total, chunk, finished uint64)) *writer {
	return &writer{Writer: w, total: total, handler: handler}
}

func (this *writer) Write(p []byte) (n int, err error) {
	n, err = this.Writer.Write(p)

	if n > 0 {
		var chunk = uint64(n)
		this.finished += chunk
		if this.handler != nil {
			this.handler(this.total, chunk, this.finished)
		}
	}

	return n, err
}
