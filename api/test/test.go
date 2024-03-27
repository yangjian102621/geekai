package main

import (
	"chatplus/utils"
	"fmt"
)

func main() {
	text := "一只 蜗牛在树干上爬，阳光透过树叶照在蜗牛的背上 --ar 1:1 --iw 0.250000 --v 6"
	fmt.Println(utils.HasChinese(text))
}
