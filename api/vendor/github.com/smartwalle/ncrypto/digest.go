package ncrypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

func MD5(value []byte) []byte {
	var m = md5.New()
	m.Write(value)
	return m.Sum(nil)
}

func SHA1(value []byte) []byte {
	var s = sha1.New()
	s.Write(value)
	return s.Sum(nil)
}

func SHA256(value []byte) []byte {
	var s = sha256.New()
	s.Write(value)
	return s.Sum(nil)
}

func SHA512(value []byte) []byte {
	var s = sha512.New()
	s.Write(value)
	return s.Sum(nil)
}
