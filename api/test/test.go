package main

import (
	"geekai/utils"
	"fmt"
)

func main() {
	text := "https://nk.img.r9it.com/chatgpt-plus/1712709360012445.png 请简单描述一下这幅图上的内容 "
	imgURL := utils.ExtractImgURL(text)
	fmt.Println(imgURL)
}
