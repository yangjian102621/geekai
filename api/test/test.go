package main

import (
	"fmt"
	"strconv"
)

func main() {
	value, err := strconv.Atoi("012345")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
