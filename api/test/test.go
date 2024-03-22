package main

import (
	"chatplus/utils"
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

type Student struct {
	Person
	School string
}

func main() {

	text := "2024-06-01 08:34:46"
	fmt.Println(utils.Str2stamp(text))

	fmt.Println(utils.Stamp2str(utils.Str2stamp(text)))
}
