// socketcount prints the socket summary from sockstat.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/net/sockstat")
	sys.Die(err)
	fmt.Println(strings.TrimSpace(text))
	if text6, err := sys.Read("/proc/net/sockstat6"); err == nil {
		fmt.Println("ipv6:")
		fmt.Println(strings.TrimSpace(text6))
	}
}
