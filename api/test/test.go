package main

import (
	"fmt"
	"reflect"
)

func main() {
	text := 1
	bytes := reflect.ValueOf(text).Bytes()
	fmt.Println(bytes)
}
