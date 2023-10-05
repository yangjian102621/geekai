package main

import (
	"fmt"
	"net/url"
	"path"
)

func main() {
	imgURL := "https://www.baidu.com/static/upload/2023/10/1696497571220711277.png?ex=6530f4a2&is=651e7fa28hmFd709d069ca05d7855ebdae42e5aa436883a36f9310d546"
	parse, err := url.Parse(imgURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(path.Ext(parse.Path))
}
