package ncrypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

func HmacMD5(data, key []byte) []byte {
	var h = hmac.New(md5.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func HmacSHA1(data, key []byte) []byte {
	var h = hmac.New(sha1.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func HmacSHA256(data, key []byte) []byte {
	var h = hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func HmacSHA512(data, key []byte) []byte {
	var h = hmac.New(sha512.New, key)
	h.Write(data)
	return h.Sum(nil)
}
