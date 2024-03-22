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

	fmt.Println(utils.RandString(64))
}
