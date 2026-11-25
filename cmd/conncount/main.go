// conncount counts rows in the TCP and UDP socket tables.
package main

import (
	"fmt"

	"linux-go-utils/internal/sys"
)

func main() {
	tcp, err := sys.TCPSockets()
	if err != nil {
		fmt.Printf("tcp unavailable: %v\n", err)
	} else {
		fmt.Printf("tcp: %d\n", len(tcp))
	}
	for _, item := range []struct {
		name, path string
	}{
		{"udp", "/proc/net/udp"},
		{"udp6", "/proc/net/udp6"},
	} {
		n, err := sys.UDPCount(item.path)
		if err != nil {
			fmt.Printf("%s unavailable: %v\n", item.name, err)
			continue
		}
		fmt.Printf("%s: %d\n", item.name, n)
	}
}
