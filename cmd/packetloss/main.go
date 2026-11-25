// packetloss prints ICMP and IP error counters from /proc/net/snmp.
package main

import (
	"fmt"
	"strings"

	"linux-go-utils/internal/sys"
)

func main() {
	text, err := sys.Read("/proc/net/snmp")
	sys.Die(err)
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "Icmp:") || strings.HasPrefix(line, "Ip:") {
			fmt.Println(line)
		}
	}
}
