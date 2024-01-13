package ncrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/smartwalle/ncrypto/pbkdf2"
	"hash"
	"io"
)

const (
	kPKCSSaltLen      = 8
	kPKCSDefaultIter  = 2048
	kPKCSDefaultMagic = "Salted__"
	kEVPMaxIvLen      = 16
)

func RandBytes(length int) (data []byte, err error) {
	data = make([]byte, length)
	if _, err = io.ReadFull(rand.Reader, data); err != nil {
		return nil, err
	}
	return data, err
}

func AESCBCEncryptWithSalt(plaintext, key []byte, iter int, magic string, padding Padding, h func() hash.Hash) ([]byte, error) {
	return AESEncryptWithSalt(plaintext, key, iter, magic, h, padding, AESCBCEncrypt)
}

func AESCBCDecryptWithSalt(data, key []byte, iter int, magic string, padding Padding, h func() hash.Hash) ([]byte, error) {
	return AESDecryptWithSalt(data, key, iter, magic, h, padding, AESCBCDecrypt)
}

func AESEncryptWithSalt(plaintext, key []byte, iter int, magic string, h func() hash.Hash, padding Padding, f func(plaintext, key, iv []byte, padding Padding) ([]byte, error)) ([]byte, error) {
	if iter <= 0 {
		iter = kPKCSDefaultIter
	}

	if h == nil {
		h = md5.New
	}

	var salt, _ = RandBytes(kPKCSSaltLen)
	var sKey = pbkdf2.Key(key, salt, iter, len(key), h)
	var sIV = pbkdf2.Key(sKey, salt, iter, kEVPMaxIvLen, h)

	var ciphertext, err = f(plaintext, sKey, sIV, padding)

	ciphertext = append(salt, ciphertext...)
	ciphertext = append([]byte(magic), ciphertext...)

	return ciphertext, err
}

func AESDecryptWithSalt(ciphertext, key []byte, iterCount int, magic string, h func() hash.Hash, padding Padding, f func(ciphertext, key, iv []byte, padding Padding) ([]byte, error)) ([]byte, error) {
	if iterCount <= 0 {
		iterCount = kPKCSDefaultIter
	}

	if h == nil {
		h = md5.New
	}

	//if len(ciphertext) <= len(magic) + kPKCSSaltLen {
	//	return nil, errors.New("Error")
	//}

	var salt = ciphertext[len(magic) : len(magic)+kPKCSSaltLen]
	var sKey = pbkdf2.Key(key, salt, iterCount, len(key), h)
	var sIV = pbkdf2.Key(sKey, salt, iterCount, kEVPMaxIvLen, h)

	var plaintext, err = f(ciphertext[len(magic)+kPKCSSaltLen:], sKey, sIV, padding)
	return plaintext, err
}

// AESCBCEncrypt 由key的长度决定是128, 192 还是 256
func AESCBCEncrypt(plaintext, key, iv []byte, padding Padding) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	src, err := padding.Pad(plaintext, blockSize)
	if err != nil {
		return nil, err
	}
	var dst = make([]byte, len(src))

	var mode = cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(dst, src)
	return dst, nil
}

func AESCBCDecrypt(ciphertext, key, iv []byte, padding Padding) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	var dst = make([]byte, len(ciphertext))

	var mode = cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, ciphertext)
	return padding.UnPad(dst, block.BlockSize())
}

func AESCFBEncrypt(plaintext, key, iv []byte, padding Padding) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	src, err := padding.Pad(plaintext, blockSize)
	if err != nil {
		return nil, err
	}
	var dst = make([]byte, len(src))

	var mode = cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(dst, src)
	return dst, nil
}

func AESCFBDecrypt(ciphertext, key, iv []byte, padding Padding) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	var dst = make([]byte, len(ciphertext))

	var mode = cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(dst, ciphertext)
	return padding.UnPad(dst, block.BlockSize())
}

func AESGCMEncrypt(plaintext, key, additional []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := RandBytes(mode.NonceSize())
	if err != nil {
		return nil, err
	}
	return mode.Seal(nonce, nonce, plaintext, additional), nil
}

func AESGCMEncryptWithNonce(plaintext, key, nonce, additional []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != mode.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size, must contain %d characters", mode.NonceSize())
	}

	return mode.Seal(nil, nonce, plaintext, additional), nil
}

func AESGCMDecrypt(ciphertext, key, additional []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := mode.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	var nonce []byte
	nonce, ciphertext = ciphertext[:nonceSize], ciphertext[nonceSize:]
	return mode.Open(nil, nonce, ciphertext, additional)
}

func AESGCMDecryptWithNonce(ciphertext, key, nonce, additional []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != mode.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size, must contain %d characters", mode.NonceSize())
	}

	return mode.Open(nil, nonce, ciphertext, additional)
}
