// dnscheck prints name servers from /etc/resolv.conf.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/etc/resolv.conf")
	sys.Die(err)
	found := false
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "nameserver" {
			found = true
			fmt.Printf("nameserver %s\n", fields[1])
		}
	}
	if !found {
		fmt.Println("no nameserver")
	}
}
