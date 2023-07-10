package main

import (
	"chatplus/core/types"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"context"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/pkoukk/tiktoken-go"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	lMap := types.NewLMap[string, types.ChatSession]()
	lMap.Put("name", types.ChatSession{SessionId: utils.RandString(32)})

	item := lMap.Get("abc")
	fmt.Println(item)
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

func testJson() {

	var role = model.ChatRole{
		Key:      "programmer",
		Name:     "程序员",
		Context:  "[{\"role\":\"user\",\"content\":\"现在开始你扮演一位程序员，你是一名优秀的程序员，具有很强的逻辑思维能力，总能高效的解决问题。你热爱编程，熟悉多种编程语言，尤其精通 Go 语言，注重代码质量，有创新意识，持续学习，良好的沟通协作。\"},{\"role\"\n:\"assistant\",\"content\":\"好的，现在我将扮演一位程序员，非常感谢您对我的评价。作为一名优秀的程序员，我非常热爱编程，并且注重代码质量。我熟悉多种编程语言，尤其是 Go 语言，可以使用它来高效地解决各种问题。\"}]",
		HelloMsg: "Talk is cheap, i will show code!",
		Icon:     "images/avatar/programmer.jpg",
		Enable:   true,
		Sort:     1,
	}
	role.Id = 1
	var v vo.ChatRole

	err := utils.CopyObject(role, &v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", v.Id)

	//var v2 = model.ChatRoles{}
	//err = utils.CopyObject(v, &v2)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("%+v\n", v2.Id)

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
