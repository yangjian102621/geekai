package utils

import (
	"math/rand"
	"strconv"
	"time"
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

func Long2IP(ipInt int64) string {
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt(ipInt&0xff, 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func ContainsStr(slice []string, item string) bool {
	for _, e := range slice {
		if e == item {
			return true
		}
	}
	return false
}
