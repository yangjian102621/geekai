## Installation
```
go get github.com/go-basic/ipv4
```

## Example
```
package main

import (
	"fmt"
	"github.com/go-basic/ipv4"
)

func main() {
	ip := ipv4.LocalIP()
	fmt.Println(ip)
	ips, _ := ipv4.LocalIPv4s()
	fmt.Println(ips)
}
```