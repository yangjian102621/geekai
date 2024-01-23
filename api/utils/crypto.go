package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// AesEncrypt 加密
func AesEncrypt(keyStr string, data []byte) (string, error) {
	//创建加密实例
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	result := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(result, encryptBytes)
	return base64.StdEncoding.EncodeToString(result), nil
}

// AesDecrypt 解密
func AesDecrypt(keyStr string, dataStr string) ([]byte, error) {
	//创建实例
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	result := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(result, data)
	//去除填充
	result, err = pkcs7UnPadding(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty encrypt data")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func Sha256(data string) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, data)
	if err != nil {
		return ""
	}

	hashValue := hash.Sum(nil)
	return fmt.Sprintf("%x", hashValue)
}

func Md5(data string) string {
	md5bs := md5.Sum([]byte(data))
	return hex.EncodeToString(md5bs[:])
}
