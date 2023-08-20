package main

import (
	"bufio"
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/service"
	"chatplus/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/pkoukk/tiktoken-go"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	minio()
}

// Http client 取消操作
func testHttpClient(ctx context.Context) {

	req, err := http.NewRequest("GET", "http://localhost:2345", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req = req.WithContext(ctx)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	_, err = io.ReadAll(resp.Body)
	for {
		time.Sleep(time.Second)
		fmt.Println(time.Now())
		select {
		case <-ctx.Done():
			fmt.Println("取消退出")
			return
		default:
			continue
		}
	}

}

func testDate() {
	fmt.Println(time.Unix(1683336167, 0).Format("2006-01-02 15:04:05"))
}

func testIp2Region() {
	dbPath := "res/ip2region.xdb"
	// 1、从 dbPath 加载整个 xdb 到内存
	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		fmt.Printf("failed to load content from `%s`: %s\n", dbPath, err)
		return
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		fmt.Printf("failed to create searcher with content: %s\n", err)
		return
	}

	str, err := searcher.SearchByStr("103.88.46.85")
	fmt.Println(str)
	if err != nil {
		log.Fatal(err)
	}
	arr := strings.Split(str, "|")
	fmt.Println(arr[2], arr[3], arr[4])

}

func calTokens() {
	text := "须知少年凌云志，曾许人间第一流"
	encoding := "cl100k_base"

	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return
	}

	// encode
	token := tke.Encode(text, nil, nil)

	//tokens
	fmt.Println(token)
	// num_tokens
	fmt.Println(len(token))

}

func testAesEncrypt() {
	// 加密
	text := []byte("this is a secret text")
	key := utils.RandString(24)
	encrypt, err := utils.AesEncrypt(key, text)
	if err != nil {
		panic(err)
	}
	fmt.Println("加密密文：", encrypt)
	// 解密
	decrypt, err := utils.AesDecrypt(key, encrypt)
	if err != nil {
		panic(err)
	}
	fmt.Println("解密明文：", string(decrypt))
}

func extractFunction() error {
	open, err := os.Open("res/data.txt")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(open)
	var contents = make([]string, 0)
	var functionCall = false
	var functionName string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if !strings.Contains(line, "data:") {
			continue
		}

		var responseBody = types.ApiResponse{}
		err = json.Unmarshal([]byte(line[6:]), &responseBody)
		if err != nil || len(responseBody.Choices) == 0 { // 数据解析出错
			break
		}

		function := responseBody.Choices[0].Delta.FunctionCall
		if functionCall && function.Name == "" {
			contents = append(contents, function.Arguments)
			continue
		}

		if !utils.IsEmptyValue(function) {
			functionCall = true
			functionName = function.Name
			continue
		}
	}

	fmt.Println("函数名称: ", functionName)
	fmt.Println(strings.Join(contents, ""))
	return err
}

func minio() {
	config := core.NewDefaultConfig()
	config.ProxyURL = "http://localhost:7777"
	config.MinioConfig = types.MinioConfig{
		Endpoint:     "localhost:9010",
		AccessKey:    "ObWIEyXaQUHOYU26L0oI",
		AccessSecret: "AJW3HHhlGrprfPcmiC7jSOSzVCyrlhX4AnOAUzqI",
		Bucket:       "chatgpt-plus",
		UseSSL:       false,
		Domain:       "http://localhost:9010",
	}
	minioService, err := service.NewMinioService(config)
	if err != nil {
		panic(err)
	}

	url, err := minioService.UploadMjImg("https://cdn.discordapp.com/attachments/1139552247693443184/1141619433752768572/lisamiller4099_A_beautiful_fairy_sister_from_Chinese_mythology__3162726e-5ee4-4f60-932b-6b78b375eaef.png")
	if err != nil {
		panic(err)
	}

	fmt.Println(url)
}
