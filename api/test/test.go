package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "7151109597841850368 一个漂亮的中国女孩，手上拿着一桶爆米花，脸上带着迷人的微笑，电影效果"
	index := strings.Index(str, " ")
	fmt.Println(str[index+1:])
}
