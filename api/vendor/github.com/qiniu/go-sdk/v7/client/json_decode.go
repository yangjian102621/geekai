package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type jsonDecodeError struct {
	original error
	data     []byte
}

func (e jsonDecodeError) Error() string { return fmt.Sprintf("%s: %s", e.original.Error(), e.data) }

func (e jsonDecodeError) Unwrap() error { return e.original }

func decodeJsonFromData(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return jsonDecodeError{original: err, data: data}
	}
	return nil
}

func DecodeJsonFromReader(reader io.Reader, v interface{}) error {
	buf := new(bytes.Buffer)
	t := io.TeeReader(reader, buf)
	err := json.NewDecoder(t).Decode(v)
	if err != nil {
		return jsonDecodeError{original: err, data: buf.Bytes()}
	}
	return nil
}
