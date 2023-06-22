package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/sha3"
)

// RandString generate rand string with specified length
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	data := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, data[r.Intn(len(data))])
	}
	return string(result)
}

func RandomNumber(bit int) int {
	rand.Seed(time.Now().UnixNano())
	min := 1 // min value
	max := 1 //max value
	for i := 0; i < bit; i++ {
		min = min * 10
		max = max * 10
	}
	max = max * 10
	return rand.Intn(max-min+1) + min
}

func ContainsStr(slice []string, item string) bool {
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
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, str)
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
