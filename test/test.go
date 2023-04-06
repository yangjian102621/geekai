package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Unix(1683336167, 0).Format("2006-01-02 15:04:05"))
}
