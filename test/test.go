package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	http.HandleFunc("/cancel", func(w http.ResponseWriter, r *http.Request) {
		cancel()
		_, _ = fmt.Fprintf(w, "请求取消!")
	})

	go func() {
		err := http.ListenAndServe(":9999", nil)
		if err != nil {
			return
		}
	}()

	testHttpClient(ctx)
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

}

func testDate() {
	fmt.Println(time.Unix(1683336167, 0).Format("2006-01-02 15:04:05"))
}
