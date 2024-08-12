package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

const (
	codeLength = 32 // 兑换码长度
)

var (
	codeMap  = make(map[string]bool)
	mapMutex = &sync.Mutex{}
)

// GenerateUniqueCode 生成唯一兑换码
func GenerateUniqueCode() (string, error) {
	for {
		code, err := generateCode()
		if err != nil {
			return "", err
		}

		mapMutex.Lock()
		if !codeMap[code] {
			codeMap[code] = true
			mapMutex.Unlock()
			return code, nil
		}
		mapMutex.Unlock()
	}
}

// generateCode 生成兑换码
func generateCode() (string, error) {
	bytes := make([]byte, codeLength/2) // 因为 hex 编码会使长度翻倍
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func main() {
	for i := 0; i < 10; i++ {
		code, err := GenerateUniqueCode()
		if err != nil {
			fmt.Println("Error generating code:", err)
			return
		}
		fmt.Println("Generated code:", code)
	}
}
