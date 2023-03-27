package main

import (
	"fmt"
	"openai/utils"
)

func main() {
	// leveldb 测试
	db, err := utils.NewLevelDB("data")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Put("name", "xiaoming")
	if err != nil {
		panic(err)
	}

	name, err := db.Get("name")
	if err != nil {
		panic(err)
	}

	fmt.Println("name: ", name)
}
