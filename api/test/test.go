package main

import (
	"fmt"
	"net/url"
)

func main() {
	parse, _ := url.Parse("http://localhost:5678/static")

	imgURLPrefix := fmt.Sprintf("%s://%s", parse.Scheme, parse.Host)
	fmt.Println(imgURLPrefix)
}
