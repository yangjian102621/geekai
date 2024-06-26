package main

import (
	"fmt"
	"geekai/utils"
)

func main() {
	file := "http://nk.img.r9it.com/chatgpt-plus/1719389335351828.xlsx"
	content, err := utils.ReadPdf(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(content)
}
