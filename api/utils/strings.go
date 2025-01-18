package utils

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/sha3"
	rand2 "math/rand"
)

// RandString generate rand string with specified length
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	data := []byte(str)
	var result []byte
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, data[r.Intn(len(data))])
	}
	return string(result)
}

func RandomNumber(bit int) int {
	minNum := intPow(10, bit-1)
	maxNum := intPow(10, bit) - 1

	rand2.NewSource(time.Now().UnixNano())
	return rand2.Intn(maxNum-minNum+1) + minNum
}

func intPow(x, y int) int {
	result := 1
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}

func Contains(slice []string, item string) bool {
	for _, e := range slice {
		if e == item {
			return true
		}
	}
	return false
}

// Stamp2str 时间戳转字符串
func Stamp2str(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// Str2stamp 字符串转时间戳
func Str2stamp(str string) int64 {
	if len(str) == 0 {
		return 0
	}

	layout := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func GenPassword(pass string, salt string) string {
	data := []byte(pass + salt)
	hash := sha3.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func JsonEncode(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func JsonDecode(src string, dest interface{}) error {
	return json.Unmarshal([]byte(src), dest)
}

func InterfaceToString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return JsonEncode(value)
}

// CutWords 截取前 N 个单词
func CutWords(str string, num int) string {
	// 按空格分割字符串为单词切片
	words := strings.Fields(str)

	// 如果单词数量超过指定数量，则裁剪单词；否则保持原样
	if len(words) > num {
		return strings.Join(words[:num], " ") + " ..."
	} else {
		return str
	}
}

// HasChinese 判断文本是否含有中文
func HasChinese(text string) bool {
	for _, char := range text {
		if unicode.Is(unicode.Scripts["Han"], char) {
			return true
		}
	}
	return false
}

func GenRedeemCode(codeLength int) (string, error) {
	bytes := make([]byte, codeLength/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsValidEmail 检查给定的字符串是否是有效的电子邮件地址
func IsValidEmail(email string) bool {
	// 这个正则表达式匹配大多数常见的邮箱格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidMobile 检查给定的字符串是否是有效的中国大陆手机号
func IsValidMobile(phone string) bool {
	// 支持 13x, 14x, 15x, 16x, 17x, 18x, 19x 开头的号码
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}
