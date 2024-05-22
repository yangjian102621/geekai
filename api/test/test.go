package main

import (
	"fmt"
	"net/url"
)

func main() {
	text := "https://nk.img.r9it.com/chatgpt-plus/1712709360012445.png"
	parse, _ := url.Parse(text)
	fmt.Println(fmt.Sprintf("%s://%s", parse.Scheme, parse.Host))
}
