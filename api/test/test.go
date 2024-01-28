package main

import (
	"fmt"
	"net/url"
)

func main() {
	u, err := url.Parse("https://api.chat-plus.net/mj/image/1706368258238514?aaa=bbb")
	fmt.Println(u.Path, u.RawQuery, err)
}
