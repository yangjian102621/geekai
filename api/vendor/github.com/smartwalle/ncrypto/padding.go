package ncrypto

import (
	"bytes"
	"errors"
)

var (
	ErrInvalidPadding = errors.New("invalid padding")
)

type Padding interface {
	Pad(src []byte, blockSize int) ([]byte, error)

	UnPad(src []byte, blockSize int) ([]byte, error)
}

type padding struct {
}

func (padding) Pad(src []byte, blockSize int) ([]byte, error) {
	var pSize = blockSize - len(src)%blockSize
	var pText = bytes.Repeat([]byte{byte(pSize)}, pSize)
	return append(src, pText...), nil
}

func (padding) UnPad(src []byte, blockSize int) ([]byte, error) {
	var srcLen = len(src)
	if srcLen == 0 {
		return nil, ErrInvalidPadding
	}

	var pad = src[srcLen-1]
	if pad == 0 {
		return nil, ErrInvalidPadding
	}

	var pSize = int(pad)
	if pSize > srcLen || pSize > blockSize {
		return nil, ErrInvalidPadding
	}

	var pBytes = src[len(src)-pSize:]
	for i := 0; i < pSize; i++ {
		if pBytes[i] != pad {
			return nil, ErrInvalidPadding
		}
	}

	return src[:(srcLen - pSize)], nil
}

type PKCS5Padding struct {
	padding
}

func (p PKCS5Padding) Pad(src []byte, blockSize int) ([]byte, error) {
	return p.padding.Pad(src, blockSize)
}

func (p PKCS5Padding) UnPad(src []byte, blockSize int) ([]byte, error) {
	return p.padding.UnPad(src, blockSize)
}

type PKCS7Padding struct {
	padding
}

func (p PKCS7Padding) Pad(src []byte, blockSize int) ([]byte, error) {
	return p.padding.Pad(src, blockSize)
}

func (p PKCS7Padding) UnPad(src []byte, blockSize int) ([]byte, error) {
	return p.padding.UnPad(src, blockSize)
}

type ZeroPadding struct {
}

func (p ZeroPadding) Pad(src []byte, blockSize int) ([]byte, error) {
	var pSize = blockSize - len(src)%blockSize
	if pSize == 0 {
		return src, ErrInvalidPadding
	}
	var pText = bytes.Repeat([]byte{0}, pSize)
	return append(src, pText...), nil
}

func (p ZeroPadding) UnPad(src []byte, blockSize int) ([]byte, error) {
	return bytes.TrimFunc(src,
		func(r rune) bool {
			return r == rune(0)
		}), nil
}

type NoPadding struct {
}

func (p NoPadding) Pad(src []byte, blockSize int) ([]byte, error) {
	return src, nil
}

func (p NoPadding) UnPad(src []byte, blockSize int) ([]byte, error) {
	return src, nil
}
