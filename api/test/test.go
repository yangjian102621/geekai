package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	r := time.Now()
	f := reflect.ValueOf(r)
	fmt.Println(f.Type().Kind())
}
