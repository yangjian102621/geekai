package main

import (
	"fmt"
	"os"
)

func main() {
	bytes, err := os.ReadFile("res/text2img.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
